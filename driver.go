/*
Copyright 2019 Thomas Weber

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/hashicorp/nomad/nomad/structs"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/client/stats"
	"github.com/hashicorp/nomad/client/taskenv"
	"github.com/hashicorp/nomad/drivers/shared/eventer"
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
	pstructs "github.com/hashicorp/nomad/plugins/shared/structs"
	"github.com/pascomnet/nomad-driver-podman/iopodman"

	shelpers "github.com/hashicorp/nomad/helper/stats"
	"github.com/varlink/go/varlink"
)

const (
	// pluginName is the name of the plugin
	pluginName = "podman"

	// fingerprintPeriod is the interval at which the driver will send fingerprint responses
	fingerprintPeriod = 30 * time.Second

	// taskHandleVersion is the version of task handle which this driver sets
	// and understands how to decode driver state
	taskHandleVersion = 1
)

var (
	// pluginInfo is the response returned for the PluginInfo RPC
	pluginInfo = &base.PluginInfoResponse{
		Type:              base.PluginTypeDriver,
		PluginApiVersions: []string{drivers.ApiVersion010},
		PluginVersion:     "0.0.1-dev",
		Name:              pluginName,
	}

	// capabilities is returned by the Capabilities RPC and indicates what
	// optional features this driver supports
	capabilities = &drivers.Capabilities{
		SendSignals: false,
		Exec:        false,
		FSIsolation: drivers.FSIsolationNone,
	}
)

// Driver is a driver for running podman containers
type Driver struct {
	// eventer is used to handle multiplexing of TaskEvents calls such that an
	// event can be broadcast to all callers
	eventer *eventer.Eventer

	// config is the driver configuration set by the SetConfig RPC
	config *PluginConfig

	// nomadConfig is the client config from nomad
	nomadConfig *base.ClientDriverConfig

	// tasks is the in memory datastore mapping taskIDs to rawExecDriverHandles
	tasks *taskStore

	// ctx is the context for the driver. It is passed to other subsystems to
	// coordinate shutdown
	ctx context.Context

	// signalShutdown is called when the driver is shutting down and cancels the
	// ctx passed to any subsystems
	signalShutdown context.CancelFunc

	// logger will log to the Nomad agent
	logger hclog.Logger
}

// TaskState is the state which is encoded in the handle returned in
// StartTask. This information is needed to rebuild the task state and handler
// during recovery.
type TaskState struct {
	TaskConfig  *drivers.TaskConfig
	ContainerID string
	StartedAt   time.Time
	Net         *drivers.DriverNetwork
}

// NewPodmanDriver returns a new DriverPlugin implementation
func NewPodmanDriver(logger hclog.Logger) drivers.DriverPlugin {
	ctx, cancel := context.WithCancel(context.Background())
	logger = logger.Named(pluginName)
	return &Driver{
		eventer:        eventer.NewEventer(ctx, logger),
		config:         &PluginConfig{},
		tasks:          newTaskStore(),
		ctx:            ctx,
		signalShutdown: cancel,
		logger:         logger,
	}
}

func (d *Driver) PluginInfo() (*base.PluginInfoResponse, error) {
	return pluginInfo, nil
}

func (d *Driver) ConfigSchema() (*hclspec.Spec, error) {
	return configSpec, nil
}

func (d *Driver) SetConfig(cfg *base.Config) error {
	var config PluginConfig
	if len(cfg.PluginConfig) != 0 {
		if err := base.MsgPackDecode(cfg.PluginConfig, &config); err != nil {
			return err
		}
	}

	d.config = &config
	if cfg.AgentConfig != nil {
		d.nomadConfig = cfg.AgentConfig.Driver
	}

	return nil
}

func (d *Driver) Shutdown(ctx context.Context) error {
	d.signalShutdown()
	return nil
}

func (d *Driver) TaskConfigSchema() (*hclspec.Spec, error) {
	return taskConfigSpec, nil
}

func (d *Driver) Capabilities() (*drivers.Capabilities, error) {
	return capabilities, nil
}

func (d *Driver) Fingerprint(ctx context.Context) (<-chan *drivers.Fingerprint, error) {
	err := shelpers.Init()
	if err != nil {
		d.logger.Error("Could not init stats helper", "err", err)
		return nil, err
	}
	ch := make(chan *drivers.Fingerprint)
	go d.handleFingerprint(ctx, ch)
	return ch, nil
}

func (d *Driver) handleFingerprint(ctx context.Context, ch chan<- *drivers.Fingerprint) {
	defer close(ch)
	ticker := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			ticker.Reset(fingerprintPeriod)
			ch <- d.buildFingerprint()
		}
	}
}

func (d *Driver) buildFingerprint() *drivers.Fingerprint {
	var health drivers.HealthState
	var desc string
	attrs := map[string]*pstructs.Attribute{}

	// be negative and guess that we will not be able to get a podman connection
	health = drivers.HealthStateUndetected
	desc = "disabled"

	// try to connect and get version info
	varlinkConnection, err := d.getConnection()
	if err != nil {
		d.logger.Error("Could not connect to podman", "err", err)
	} else {
		defer varlinkConnection.Close()
		//  see if we get the system info from podman
		info, err := iopodman.GetInfo().Call(d.ctx, varlinkConnection)
		if err == nil {
			// yay! we can enable the driver
			health = drivers.HealthStateHealthy
			desc = "ready"
			// TODO: we would have more data here... maybe we can return more details to nomad
			attrs["driver.podman"] = pstructs.NewBoolAttribute(true)
			attrs["driver.podman.version"] = pstructs.NewStringAttribute(info.Podman.Podman_version)
		} else {
			d.logger.Error("Cound not get podman info", "err", err)
		}
	}

	return &drivers.Fingerprint{
		Attributes:        attrs,
		Health:            health,
		HealthDescription: desc,
	}
}

func (d *Driver) RecoverTask(handle *drivers.TaskHandle) error {
	if handle == nil {
		return fmt.Errorf("error: handle cannot be nil")
	}

	if _, ok := d.tasks.Get(handle.Config.ID); ok {
		return nil
	}

	var taskState TaskState
	if err := handle.GetDriverState(&taskState); err != nil {
		return fmt.Errorf("failed to decode task state from handle: %v", err)
	}

	varlinkConnection, err := d.getConnection()
	if err != nil {
		return fmt.Errorf("failed to recover container ref: %v", err)
	}
	defer varlinkConnection.Close()

	// FIXME: duplicated code. share more code with StartTask
	filters := []string{"id=" + taskState.ContainerID}
	psOpts := iopodman.PsOpts{
		Filters: &filters,
	}
	psContainers, err := iopodman.Ps().Call(d.ctx, varlinkConnection, psOpts)
	if err != nil {
		return fmt.Errorf("failt to recover task, could not get container info: %v", err)
	}
	if len(psContainers) != 1 {
		return fmt.Errorf("failt to recover task, problem with Ps()")
	}
	d.logger.Debug("Recovered podman container", "container", taskState.ContainerID, "pid", psContainers[0].PidNum)

	h := &TaskHandle{
		initPid:     int(psContainers[0].PidNum),
		containerID: taskState.ContainerID,
		driver:      d,
		taskConfig:  taskState.TaskConfig,
		procState:   drivers.TaskStateRunning,
		startedAt:   taskState.StartedAt,
		exitResult:  &drivers.ExitResult{},
		logger:      d.logger,

		totalCpuStats:  stats.NewCpuStats(),
		userCpuStats:   stats.NewCpuStats(),
		systemCpuStats: stats.NewCpuStats(),

		removeContainerOnExit: d.config.GC.Container,
	}

	d.tasks.Set(taskState.TaskConfig.ID, h)

	go h.run()

	return nil
}

// BuildContainerName returns the podman container name for a given TaskConfig
func BuildContainerName(cfg *drivers.TaskConfig) string {
	return fmt.Sprintf("%s-%s", cfg.Name, cfg.AllocID)
}

func (d *Driver) StartTask(cfg *drivers.TaskConfig) (*drivers.TaskHandle, *drivers.DriverNetwork, error) {
	if _, ok := d.tasks.Get(cfg.ID); ok {
		return nil, nil, fmt.Errorf("task with ID %q already started", cfg.ID)
	}

	var driverConfig TaskConfig
	if err := cfg.DecodeDriverConfig(&driverConfig); err != nil {
		return nil, nil, fmt.Errorf("failed to decode driver config: %v", err)
	}

	// d.logger.Info("starting podman task", "driver_cfg", hclog.Fmt("%+v", driverConfig))
	handle := drivers.NewTaskHandle(taskHandleVersion)
	handle.Config = cfg

	varlinkConnection, err := d.getConnection()
	if err != nil {
		return nil, nil, fmt.Errorf("failt to start task, could not connect to podman: %v", err)
	}
	defer varlinkConnection.Close()

	if driverConfig.Image == "" {
		return nil, nil, fmt.Errorf("image name required")
	}

	allArgs := []string{driverConfig.Image}
	if driverConfig.Command != "" {
		allArgs = append(allArgs, driverConfig.Command)
	}
	allArgs = append(allArgs, driverConfig.Args...)
	containerName := BuildContainerName(cfg)
	memoryLimit := fmt.Sprintf("%dm", cfg.Resources.NomadResources.Memory.MemoryMB)
	cpuShares := cfg.Resources.LinuxResources.CPUShares
	logOpts := []string{
		fmt.Sprintf("path=%s", cfg.StdoutPath),
	}

	// ensure to include port_map into tasks environment map
	cfg.Env = taskenv.SetPortMapEnvs(cfg.Env, driverConfig.PortMap)

	// convert environment map into a k=v list
	allEnv := cfg.EnvList()

	// ensure to mount nomad alloc dirs into the container
	allVolumes := []string{
		fmt.Sprintf("%s:%s", cfg.TaskDir().SharedAllocDir, cfg.Env[taskenv.AllocDir]),
		fmt.Sprintf("%s:%s", cfg.TaskDir().LocalDir, cfg.Env[taskenv.TaskLocalDir]),
		fmt.Sprintf("%s:%s", cfg.TaskDir().SecretsDir, cfg.Env[taskenv.SecretsDir]),
	}

	if d.config.Volumes.Enabled {
		// add task specific volumes, if enabled
		allVolumes = append(allVolumes, driverConfig.Volumes...)
	}

	// Apply SELinux Label to each volume
	if selinuxLabel := d.config.Volumes.SelinuxLabel; selinuxLabel != "" {
		for i := range allVolumes {
			allVolumes[i] = fmt.Sprintf("%s:%s", allVolumes[i], selinuxLabel)
		}
	}

	createOpts := iopodman.Create{
		Args:       allArgs,
		Env:        &allEnv,
		Name:       &containerName,
		Volume:     &allVolumes,
		Memory:     &memoryLimit,
		MemorySwap: &memoryLimit,
		CpuShares:  &cpuShares,
		LogOpt:     &logOpts,
		Hostname:   &driverConfig.Hostname,
	}

	// Setup port mapping and exposed ports
	if len(cfg.Resources.NomadResources.Networks) == 0 {
		d.logger.Debug("no network interfaces are available")
		if len(driverConfig.PortMap) > 0 {
			return nil, nil, fmt.Errorf("Trying to map ports but no network interface is available")
		}
	} else {
		publishedPorts := []string{}
		network := cfg.Resources.NomadResources.Networks[0]
		allPorts := []structs.Port{}
		allPorts = append(allPorts, network.ReservedPorts...)
		allPorts = append(allPorts, network.DynamicPorts...)

		for _, port := range allPorts {
			hostPort := port.Value
			// By default we will map the allocated port 1:1 to the container
			containerPort := port.Value

			// If the user has mapped a port using port_map we'll change it here
			if mapped, ok := driverConfig.PortMap[port.Label]; ok {
				containerPort = mapped
			}

			d.logger.Debug("Publish port", "ip", network.IP, "hostPort", hostPort, "containerPort", containerPort)
			// we map both udp and tcp ports
			publishedPorts = append(publishedPorts, fmt.Sprintf("%s:%d:%d/tcp", network.IP, hostPort, containerPort))
			publishedPorts = append(publishedPorts, fmt.Sprintf("%s:%d:%d/udp", network.IP, hostPort, containerPort))
		}

		createOpts.Publish = &publishedPorts
	}

	containerID, err := iopodman.CreateContainer().Call(d.ctx, varlinkConnection, createOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failt to start task, could not create container: %v", err)
	}
	d.logger.Info("Created container", "container", containerID)

	cleanup := func() {
		d.logger.Debug("Cleaning up", "container", containerID)
		if _, err := iopodman.RemoveContainer().Call(d.ctx, varlinkConnection, containerID, true, true); err != nil {
			d.logger.Error("failed to clean up from an error in Start", "error", err)
		}
	}

	_, err = iopodman.StartContainer().Call(d.ctx, varlinkConnection, containerID)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failt to start task, could not start container: %v", err)
	}
	d.logger.Info("Started container", "containerId", containerID)

	// FIXME: there is some podman race condition to end up in a deadlock here
	//
	d.logger.Info("inspecting container", "containerId", containerID)
	inspectJSON, err := iopodman.InspectContainer().Call(d.ctx, varlinkConnection, containerID)
	if err != nil {
		d.logger.Error("failt to inspect container", "err", err)
		cleanup()
		return nil, nil, fmt.Errorf("failt to start task, could not inspect container : %v", err)
	}
	var inspectData iopodman.InspectContainerData
	err = json.Unmarshal([]byte(inspectJSON), &inspectData)
	if err != nil {
		cleanup()
		d.logger.Error("failt to unmarshal inspect container", "err", err)
		return nil, nil, fmt.Errorf("failt to start task, could not unmarshal inspect container : %v", err)
	}
	d.logger.Debug("Started podman container", "container", containerID, "pid", inspectData.State.Pid, "ip", inspectData.NetworkSettings.IPAddress)

	net := &drivers.DriverNetwork{
		PortMap:       driverConfig.PortMap,
		IP:            inspectData.NetworkSettings.IPAddress,
		AutoAdvertise: true,
	}

	h := &TaskHandle{
		initPid:     inspectData.State.Pid,
		containerID: containerID,
		driver:      d,
		taskConfig:  cfg,
		procState:   drivers.TaskStateRunning,
		startedAt:   time.Now().Round(time.Millisecond),
		logger:      d.logger,

		totalCpuStats:  stats.NewCpuStats(),
		userCpuStats:   stats.NewCpuStats(),
		systemCpuStats: stats.NewCpuStats(),

		removeContainerOnExit: d.config.GC.Container,
	}

	driverState := TaskState{
		ContainerID: containerID,
		TaskConfig:  cfg,
		StartedAt:   h.startedAt,
		Net:         net,
	}

	if err := handle.SetDriverState(&driverState); err != nil {
		d.logger.Error("failed to start task, error setting driver state", "error", err)
		cleanup()
		return nil, nil, fmt.Errorf("failed to set driver state: %v", err)
	}

	d.tasks.Set(cfg.ID, h)

	go h.run()

	d.logger.Debug("Completely started", "containerID", containerID)

	return handle, net, nil
}

func (d *Driver) WaitTask(ctx context.Context, taskID string) (<-chan *drivers.ExitResult, error) {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	ch := make(chan *drivers.ExitResult)
	go d.handleWait(ctx, handle, ch)

	return ch, nil
}

func (d *Driver) handleWait(ctx context.Context, handle *TaskHandle, ch chan *drivers.ExitResult) {
	defer close(ch)

	//
	// Wait for process completion by polling status from handler.
	// We cannot use the following alternatives:
	//   * Process.Wait() requires LXC container processes to be children
	//     of self process; but LXC runs container in separate PID hierarchy
	//     owned by PID 1.
	//   * lxc.Container.Wait() holds a write lock on container and prevents
	//     any other calls, including stats.
	//
	// Going with simplest approach of polling for handler to mark exit.
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			s := handle.TaskStatus()
			if s.State == drivers.TaskStateExited {
				ch <- handle.exitResult
			}
		}
	}
}

func (d *Driver) StopTask(taskID string, timeout time.Duration, signal string) error {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}
	err := handle.shutdown(timeout)
	return err
}

func (d *Driver) DestroyTask(taskID string, force bool) error {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}

	if handle.IsRunning() && !force {
		return fmt.Errorf("cannot destroy running task")
	}

	if handle.IsRunning() {
		d.logger.Info("Have to destroyTask but container is still running", "containerID", handle.containerID)
		// we can not do anything, so catching the error is useless
		err := handle.shutdown(1 * time.Minute)
		if err != nil {
			d.logger.Warn("failed to stop container during destroy", "error", err)
		}
	}

	if handle.removeContainerOnExit {
		varlinkConnection, err := d.getConnection()
		if err == nil {
			defer varlinkConnection.Close()

			if _, err := iopodman.RemoveContainer().Call(d.ctx, varlinkConnection, handle.containerID, true, true); err == nil {
				d.logger.Debug("Removed container", "container", handle.containerID)
			} else {
				d.logger.Warn("Could not remove container", "container", handle.containerID, "error", err)
			}
		} else {
			d.logger.Warn("Could not remove container, error connecting to podman", "container", handle.containerID, "error", err)
		}
	}

	d.tasks.Delete(taskID)
	return nil
}

func (d *Driver) InspectTask(taskID string) (*drivers.TaskStatus, error) {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	return handle.TaskStatus(), nil
}

func (d *Driver) TaskStats(ctx context.Context, taskID string, interval time.Duration) (<-chan *drivers.TaskResourceUsage, error) {
	d.logger.Debug("TaskStats called")
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	return handle.stats(ctx, interval)
}

func (d *Driver) TaskEvents(ctx context.Context) (<-chan *drivers.TaskEvent, error) {
	return d.eventer.TaskEvents(ctx)
}

func (d *Driver) SignalTask(taskID string, signal string) error {
	return fmt.Errorf("Pdoman driver does not support signals")
}

func (d *Driver) ExecTask(taskID string, cmd []string, timeout time.Duration) (*drivers.ExecTaskResult, error) {
	return nil, fmt.Errorf("Pdoman driver does not support exec")
}

func (d *Driver) getConnection() (*varlink.Connection, error) {
	// FIXME: a parameter for the socket would be nice
	varlinkConnection, err := varlink.NewConnection(d.ctx, "unix://run/podman/io.podman")
	return varlinkConnection, err
}
