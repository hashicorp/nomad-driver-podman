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
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/nomad/nomad/structs"

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

	// podmanClient encapsulates podman remote calls
	podmanClient *PodmanClient
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
	return &Driver{
		eventer:        eventer.NewEventer(ctx, logger),
		config:         &PluginConfig{},
		tasks:          newTaskStore(),
		ctx:            ctx,
		signalShutdown: cancel,
		logger:         logger.Named(pluginName),
		podmanClient: &PodmanClient{
			ctx:    ctx,
			logger: logger.Named("podmanClient"),
		},
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
	info, err := d.podmanClient.GetInfo()
	if err != nil {
		d.logger.Error("Cound not get podman info", "err", err)
	} else {
		// yay! we can enable the driver
		health = drivers.HealthStateHealthy
		desc = "ready"
		// TODO: we would have more data here... maybe we can return more details to nomad
		attrs["driver.podman"] = pstructs.NewBoolAttribute(true)
		attrs["driver.podman.version"] = pstructs.NewStringAttribute(info.Podman.Podman_version)
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

	// FIXME: duplicated code. share more code with StartTask
	_, err := d.podmanClient.PsID(taskState.ContainerID)
	if err != nil {
		return fmt.Errorf("failt to recover task, could not get container info: %v", err)
	}

	h := &TaskHandle{
		containerID: taskState.ContainerID,
		driver:      d,
		taskConfig:  taskState.TaskConfig,
		procState:   drivers.TaskStateRunning,
		startedAt:   taskState.StartedAt,
		exitResult:  &drivers.ExitResult{},
		logger:      d.logger.Named("podmanHandle"),

		totalCpuStats:  stats.NewCpuStats(),
		userCpuStats:   stats.NewCpuStats(),
		systemCpuStats: stats.NewCpuStats(),

		removeContainerOnExit: d.config.GC.Container,
	}

	d.tasks.Set(taskState.TaskConfig.ID, h)

	go h.MonitorContainer()
	d.logger.Debug("Recovered podman container", "container", taskState.ContainerID)

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

	swap := memoryLimit
	if driverConfig.MemorySwap != "" {
		swap = driverConfig.MemorySwap
	}

	createOpts := iopodman.Create{
		Args:              allArgs,
		Env:               &allEnv,
		Name:              &containerName,
		Volume:            &allVolumes,
		Memory:            &memoryLimit,
		CpuShares:         &cpuShares,
		LogOpt:            &logOpts,
		Hostname:          &driverConfig.Hostname,
		Init:              &driverConfig.Init,
		InitPath:          &driverConfig.InitPath,
		User:              &cfg.User,
		MemoryReservation: &driverConfig.MemoryReservation,
		MemorySwap:        &swap,
		MemorySwappiness:  &driverConfig.MemorySwappiness,
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

	containerID, err := d.podmanClient.CreateContainer(createOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failt to start task, could not create container: %v", err)
	}

	cleanup := func() {
		d.logger.Debug("Cleaning up", "container", containerID)
		if err := d.podmanClient.ForceRemoveContainer(containerID); err != nil {
			d.logger.Error("failed to clean up from an error in Start", "error", err)
		}
	}

	err = d.podmanClient.StartContainer(containerID)
	if err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("failt to start task, could not start container: %v", err)
	}

	inspectData, err := d.podmanClient.InspectContainer(containerID)
	if err != nil {
		d.logger.Error("failt to inspect container", "err", err)
		cleanup()
		return nil, nil, fmt.Errorf("failt to start task, could not inspect container : %v", err)
	}

	net := &drivers.DriverNetwork{
		PortMap:       driverConfig.PortMap,
		IP:            inspectData.NetworkSettings.IPAddress,
		AutoAdvertise: true,
	}

	h := &TaskHandle{
		containerID: containerID,
		driver:      d,
		taskConfig:  cfg,
		procState:   drivers.TaskStateRunning,
		exitResult:  &drivers.ExitResult{},
		startedAt:   time.Now().Round(time.Millisecond),
		logger:      d.logger.Named("podmanHandle"),

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

	go h.MonitorContainer()

	d.logger.Debug("Completely started container", "taskID", cfg.ID, "container", containerID, "ip", inspectData.NetworkSettings.IPAddress)

	return handle, net, nil
}

func (d *Driver) WaitTask(ctx context.Context, taskID string) (<-chan *drivers.ExitResult, error) {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	ch := make(chan *drivers.ExitResult)

	// check if the container is already dead
	s := handle.TaskStatus()
	if s.State == drivers.TaskStateExited {
		d.logger.Warn("Not setting exitChannel for a stopped container", "container", handle.containerID)
		// tell nomad about this
		ch <- handle.exitResult
		close(ch)
		// and bail out
		return nil, errors.New("Container is not running")
	}

	// otherwise forward the exit channel into the handle.
	// it's used when the stats collector detects a stopped channel
	d.logger.Debug("Setting exitChannel for Handle", "container", handle.containerID)
	handle.exitChannel = ch
	return ch, nil
}

func (d *Driver) StopTask(taskID string, timeout time.Duration, signal string) error {
	d.logger.Info("Stopping task", "taskID", taskID)
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}
	err := d.podmanClient.StopContainer(handle.containerID, int64(timeout.Seconds()))
	if err != nil {
		d.logger.Error("Could not stop/kill container", "containerID", handle.containerID, "err", err)
		return err
	}
	return nil
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
		err := d.podmanClient.StopContainer(handle.containerID, int64(60))
		if err != nil {
			d.logger.Warn("failed to stop/kill container during destroy", "error", err)
		}
	}

	if handle.removeContainerOnExit {
		err := d.podmanClient.ForceRemoveContainer(handle.containerID)
		if err != nil {
			d.logger.Warn("Could not remove container", "container", handle.containerID, "error", err)
		}
	}

	d.tasks.Delete(taskID)
	return nil
}

func (d *Driver) InspectTask(taskID string) (*drivers.TaskStatus, error) {
	d.logger.Info("InspectTask called")
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	return handle.TaskStatus(), nil
}

func (d *Driver) TaskStats(ctx context.Context, taskID string, interval time.Duration) (<-chan *drivers.TaskResourceUsage, error) {
	d.logger.Debug("TaskStats called", "taskID", taskID)
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}
	statsChannel := make(chan *drivers.TaskResourceUsage)
	go handle.runStatsEmitter(ctx, statsChannel, interval)
	return statsChannel, nil
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
