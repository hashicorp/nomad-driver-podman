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
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/consul-template/signals"

	"github.com/hashicorp/nomad/nomad/structs"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/iopodman"
	"github.com/hashicorp/nomad-driver-podman/version"
	"github.com/hashicorp/nomad/client/stats"
	"github.com/hashicorp/nomad/client/taskenv"
	"github.com/hashicorp/nomad/drivers/shared/eventer"
	nstructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
	pstructs "github.com/hashicorp/nomad/plugins/shared/structs"

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
		PluginVersion:     version.Version,
		Name:              pluginName,
	}

	// capabilities is returned by the Capabilities RPC and indicates what
	// optional features this driver supports
	capabilities = &drivers.Capabilities{
		SendSignals: true,
		Exec:        false,
		FSIsolation: drivers.FSIsolationImage,
		NetIsolationModes: []drivers.NetIsolationMode{
			drivers.NetIsolationModeGroup,
			drivers.NetIsolationModeHost,
			drivers.NetIsolationModeTask,
		},
		MustInitiateNetwork: false,
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

// PluginInfo returns metadata about the podman driver plugin
func (d *Driver) PluginInfo() (*base.PluginInfoResponse, error) {
	return pluginInfo, nil
}

// ConfigSchema function allows a plugin to tell Nomad the schema for its configuration.
// This configuration is given in a plugin block of the client configuration.
// The schema is defined with the hclspec package.
func (d *Driver) ConfigSchema() (*hclspec.Spec, error) {
	return configSpec, nil
}

// SetConfig function is called when starting the plugin for the first time.
// The Config given has two different configuration fields. The first PluginConfig,
// is an encoded configuration from the plugin block of the client config.
// The second, AgentConfig, is the Nomad agent's configuration which is given to all plugins.
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

	if config.SocketPath != "" {
		d.podmanClient.varlinkSocketPath = config.SocketPath
	} else {
		user, _ := user.Current()
		procFilesystems, err := getProcFilesystems()

		if err != nil {
			return err
		}

		socketPath := guessSocketPath(user, procFilesystems)

		if err != nil {
			return err
		}

		d.podmanClient.varlinkSocketPath = socketPath
	}

	return nil
}

// TaskConfigSchema returns the schema for the driver configuration of the task.
func (d *Driver) TaskConfigSchema() (*hclspec.Spec, error) {
	return taskConfigSpec, nil
}

// Capabilities define what features the driver implements.
func (d *Driver) Capabilities() (*drivers.Capabilities, error) {
	return capabilities, nil
}

// Fingerprint  is called by the client when the plugin is started.
// It allows the driver to indicate its health to the client.
// The channel returned should immediately send an initial Fingerprint,
// then send periodic updates at an interval that is appropriate for the driver
// until the context is canceled.
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
		d.logger.Error("Could not get podman info", "err", err)
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

// RecoverTask detects running tasks when nomad client or task driver is restarted.
// When a driver is restarted it is not expected to persist any internal state to disk.
// To support this, Nomad will attempt to recover a task that was previously started
// if the driver does not recognize the task ID. During task recovery,
// Nomad calls RecoverTask passing the TaskHandle that was returned by the StartTask function.
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
	d.logger.Debug("Checking for recoverable task", "task", handle.Config.Name, "taskid", handle.Config.ID, "container", taskState.ContainerID)

	psInfo, err := d.podmanClient.PsID(taskState.ContainerID)
	if err != nil {
		d.logger.Warn("Recovery lookup failed", "task", handle.Config.ID, "container", taskState.ContainerID, "err", err)
		return nil
	}

	h := &TaskHandle{
		containerID: taskState.ContainerID,
		driver:      d,
		taskConfig:  taskState.TaskConfig,
		procState:   drivers.TaskStateUnknown,
		startedAt:   taskState.StartedAt,
		exitResult:  &drivers.ExitResult{},
		logger:      d.logger.Named("podmanHandle"),

		totalCPUStats:  stats.NewCpuStats(),
		userCPUStats:   stats.NewCpuStats(),
		systemCPUStats: stats.NewCpuStats(),

		removeContainerOnExit: d.config.GC.Container,
	}

	if psInfo.State == "running" {
		d.logger.Info("Recovered a still running container", "container", psInfo.Id)
		h.procState = drivers.TaskStateRunning
	} else if psInfo.State == "exited" {
		// are we allowed to restart a stopped container?
		if d.config.RecoverStopped {
			d.logger.Debug("Found a stopped container, try to start it", "container", psInfo.Id)
			if err = d.podmanClient.StartContainer(psInfo.Id); err != nil {
				d.logger.Warn("Recovery restart failed", "task", handle.Config.ID, "container", taskState.ContainerID, "err", err)
			} else {
				d.logger.Info("Restarted a container during recovery", "container", psInfo.Id)
				h.procState = drivers.TaskStateRunning
			}
		} else {
			// no, let's cleanup here to prepare for a StartTask()
			d.logger.Debug("Found a stopped container, removing it", "container", psInfo.Id)
			if err = d.podmanClient.ForceRemoveContainer(psInfo.Id); err != nil {
				d.logger.Warn("Recovery cleanup failed", "task", handle.Config.ID, "container", psInfo.Id, "err", err)
			}
			h.procState = drivers.TaskStateExited
		}
	} else {
		d.logger.Warn("Recovery restart failed, unknown container state", "state", psInfo.State, "container", taskState.ContainerID)
		h.procState = drivers.TaskStateUnknown
	}

	d.tasks.Set(taskState.TaskConfig.ID, h)

	go h.runContainerMonitor()
	d.logger.Debug("Recovered container handle", "container", taskState.ContainerID)

	return nil
}

// BuildContainerName returns the podman container name for a given TaskConfig
func BuildContainerName(cfg *drivers.TaskConfig) string {
	return fmt.Sprintf("%s-%s", cfg.Name, cfg.AllocID)
}

// StartTask creates and starts a new Container based on the given TaskConfig.
func (d *Driver) StartTask(cfg *drivers.TaskConfig) (*drivers.TaskHandle, *drivers.DriverNetwork, error) {
	if _, ok := d.tasks.Get(cfg.ID); ok {
		return nil, nil, fmt.Errorf("task with ID %q already started", cfg.ID)
	}

	var driverConfig TaskConfig
	if err := cfg.DecodeDriverConfig(&driverConfig); err != nil {
		return nil, nil, fmt.Errorf("failed to decode driver config: %v", err)
	}

	handle := drivers.NewTaskHandle(taskHandleVersion)
	handle.Config = cfg

	if driverConfig.Image == "" {
		return nil, nil, fmt.Errorf("image name required")
	}

	img, err := d.createImage(cfg, &driverConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't create image: %v", err)
	}
	d.logger.Debug("created/pulled image", "img_id", img.ID)

	allArgs := []string{driverConfig.Image}
	if driverConfig.Command != "" {
		allArgs = append(allArgs, driverConfig.Command)
	}
	allArgs = append(allArgs, driverConfig.Args...)

	var entryPoint *string // nil -> image default entryPoint
	if driverConfig.Entrypoint != "" {
		*entryPoint = driverConfig.Entrypoint
	}

	var workingDir *string // nil -> image default workingDir
	if driverConfig.WorkingDir != "" {
		*workingDir = driverConfig.WorkingDir
	}

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

	allVolumes, err := d.containerBinds(cfg, &driverConfig)
	if err != nil {
		return nil, nil, err
	}
	d.logger.Debug("binding volumes", "volumes", allVolumes)

	swap := memoryLimit
	if driverConfig.MemorySwap != "" {
		swap = driverConfig.MemorySwap
	}

	procFilesystems, err := getProcFilesystems()
	swappiness := new(int64)
	if err == nil {
		cgroupv2 := false
		for _, l := range procFilesystems {
			cgroupv2 = cgroupv2 || strings.HasSuffix(l, "cgroup2")
		}
		if !cgroupv2 {
			swappiness = &driverConfig.MemorySwappiness
		}
	}

	// Generate network string
	var network string
	if cfg.NetworkIsolation != nil &&
		cfg.NetworkIsolation.Path != "" {
		network = fmt.Sprintf("ns:%s", cfg.NetworkIsolation.Path)
	} else {
		network = driverConfig.NetworkMode
	}

	createOpts := iopodman.Create{
		Args:              allArgs,
		Entrypoint:        entryPoint,
		WorkDir:           workingDir,
		Env:               &allEnv,
		Name:              &containerName,
		Volume:            &allVolumes,
		Memory:            &memoryLimit,
		CpuShares:         &cpuShares,
		CapAdd:            &driverConfig.CapAdd,
		CapDrop:           &driverConfig.CapDrop,
		Dns:               &driverConfig.Dns,
		LogOpt:            &logOpts,
		Hostname:          &driverConfig.Hostname,
		Init:              &driverConfig.Init,
		InitPath:          &driverConfig.InitPath,
		User:              &cfg.User,
		MemoryReservation: &driverConfig.MemoryReservation,
		MemorySwap:        &swap,
		MemorySwappiness:  swappiness,
		Network:           &network,
		Tmpfs:             &driverConfig.Tmpfs,
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

			// we map both udp and tcp ports
			publishedPorts = append(publishedPorts, fmt.Sprintf("%s:%d:%d/tcp", network.IP, hostPort, containerPort))
			publishedPorts = append(publishedPorts, fmt.Sprintf("%s:%d:%d/udp", network.IP, hostPort, containerPort))
		}

		createOpts.Publish = &publishedPorts
	}

	containerID := ""
	recoverRunningContainer := false
	// check if there is a container with same name
	otherContainerPs, err := d.podmanClient.PsByName(containerName)
	if err == nil {
		// ok, seems we found a container with similar name
		if otherContainerPs.State == "running" {
			// it's still running. So let's use it instead of creating a new one
			d.logger.Info("Detect running container with same name, we reuse it", "task", cfg.ID, "container", otherContainerPs.Id)
			containerID = otherContainerPs.Id
			recoverRunningContainer = true
		} else {
			// let's remove the old, dead container
			d.logger.Info("Detect stopped container with same name, removing it", "task", cfg.ID, "container", otherContainerPs.Id)
			if err = d.podmanClient.ForceRemoveContainer(otherContainerPs.Id); err != nil {
				return nil, nil, nstructs.WrapRecoverable(fmt.Sprintf("failed to remove dead container: %v", err), err)
			}
		}
	}

	if !recoverRunningContainer {
		containerID, err = d.podmanClient.CreateContainer(createOpts)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to start task, could not create container: %v", err)
		}
	}

	cleanup := func() {
		d.logger.Debug("Cleaning up", "container", containerID)
		if err := d.podmanClient.ForceRemoveContainer(containerID); err != nil {
			d.logger.Error("failed to clean up from an error in Start", "error", err)
		}
	}

	if !recoverRunningContainer {
		if err = d.podmanClient.StartContainer(containerID); err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to start task, could not start container: %v", err)
		}
	}

	inspectData, err := d.podmanClient.InspectContainer(containerID)
	if err != nil {
		d.logger.Error("failed to inspect container", "err", err)
		cleanup()
		return nil, nil, fmt.Errorf("failed to start task, could not inspect container : %v", err)
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

		totalCPUStats:  stats.NewCpuStats(),
		userCPUStats:   stats.NewCpuStats(),
		systemCPUStats: stats.NewCpuStats(),

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

	go h.runContainerMonitor()

	d.logger.Info("Completely started container", "taskID", cfg.ID, "container", containerID, "ip", inspectData.NetworkSettings.IPAddress)

	return handle, net, nil
}

// WaitTask function is expected to return a channel that will send an *ExitResult when the task
// exits or close the channel when the context is canceled. It is also expected that calling
// WaitTask on an exited task will immediately send an *ExitResult on the returned channel.
// A call to WaitTask after StopTask is valid and should be handled.
// If WaitTask is called after DestroyTask, it should return drivers.ErrTaskNotFound as no task state should exist after DestroyTask is called.
func (d *Driver) WaitTask(ctx context.Context, taskID string) (<-chan *drivers.ExitResult, error) {
	d.logger.Debug("WaitTask called", "task", taskID)
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}
	ch := make(chan *drivers.ExitResult)
	go handle.runExitWatcher(ctx, ch)
	return ch, nil
}

// StopTask function is expected to stop a running task by sending the given signal to it.
// If the task does not stop during the given timeout, the driver must forcefully kill the task.
// StopTask does not clean up resources of the task or remove it from the driver's internal state.
func (d *Driver) StopTask(taskID string, timeout time.Duration, signal string) error {
	d.logger.Info("Stopping task", "taskID", taskID, "signal", signal)
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}
	// fixme send proper signal to container
	err := d.podmanClient.StopContainer(handle.containerID, int64(timeout.Seconds()))
	if err != nil {
		d.logger.Error("Could not stop/kill container", "containerID", handle.containerID, "err", err)
		return err
	}
	return nil
}

// DestroyTask function cleans up and removes a task that has terminated.
// If force is set to true, the driver must destroy the task even if it is still running.
func (d *Driver) DestroyTask(taskID string, force bool) error {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}

	if handle.isRunning() && !force {
		return fmt.Errorf("cannot destroy running task")
	}

	if handle.isRunning() {
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

// InspectTask function returns detailed status information for the referenced taskID.
func (d *Driver) InspectTask(taskID string) (*drivers.TaskStatus, error) {
	d.logger.Info("InspectTask called")
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	return handle.taskStatus(), nil
}

// TaskStats function returns a channel which the driver should send stats to at the given interval.
// The driver must send stats at the given interval until the given context is canceled or the task terminates.
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

// TaskEvents function allows the driver to publish driver specific events about tasks and
// the Nomad client publishes events associated with an allocation.
func (d *Driver) TaskEvents(ctx context.Context) (<-chan *drivers.TaskEvent, error) {
	return d.eventer.TaskEvents(ctx)
}

// SignalTask function is used by drivers which support sending OS signals (SIGHUP, SIGKILL, SIGUSR1 etc.) to the task.
// It is an optional function and is listed as a capability in the driver Capabilities struct.
func (d *Driver) SignalTask(taskID string, signal string) error {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return drivers.ErrTaskNotFound
	}

	// The given signal will be forwarded to the target taskID.
	// Please checkout https://github.com/hashicorp/consul-template/blob/master/signals/signals_unix.go
	// for a list of supported signals.
	sig, ok := signals.SignalLookup[signal]
	if !ok {
		return fmt.Errorf("Invalid signal: %s", signal)
	}

	return d.podmanClient.SignalContainer(handle.containerID, sig)
}

// ExecTask function is used by the Nomad client to execute commands inside the task execution context.
func (d *Driver) ExecTask(taskID string, cmd []string, timeout time.Duration) (*drivers.ExecTaskResult, error) {
	return nil, fmt.Errorf("Podman driver does not support exec")
}

func (d *Driver) createImage(cfg *drivers.TaskConfig, driverConfig *TaskConfig) (iopodman.InspectImageData, error) {
	img, err := d.podmanClient.InspectImage(driverConfig.Image)
	if err != nil {
		err = d.eventer.EmitEvent(&drivers.TaskEvent{
			TaskID:    cfg.ID,
			AllocID:   cfg.AllocID,
			TaskName:  cfg.Name,
			Timestamp: time.Now(),
			Message:   "Downloading image",
			Annotations: map[string]string{
				"image": driverConfig.Image,
			},
		})
		if err != nil {
			d.logger.Warn("error emitting event", "error", err)
		}

		pullLog, err := d.podmanClient.PullImage(driverConfig.Image)
		if err != nil {
			return iopodman.InspectImageData{}, fmt.Errorf("image %s couldn't be downloaded: %v", driverConfig.Image, err)
		}

		img, err = d.podmanClient.InspectImage(driverConfig.Image)
		if err != nil {
			return iopodman.InspectImageData{}, fmt.Errorf("image %s couldn't be inspected: %v", driverConfig.Image, err)
		}

		err = d.eventer.EmitEvent(&drivers.TaskEvent{
			TaskID:    cfg.ID,
			AllocID:   cfg.AllocID,
			TaskName:  cfg.Name,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("Image downloaded: %s", pullLog),
			Annotations: map[string]string{
				"image": driverConfig.Image,
			},
		})
		if err != nil {
			d.logger.Warn("error emitting event", "error", err)
		}

	}

	d.logger.Debug("Image created", "img_id", img.ID, "config", img.Config)

	return img, nil
}

func (d *Driver) containerBinds(task *drivers.TaskConfig, driverConfig *TaskConfig) ([]string, error) {
	allocDirBind := fmt.Sprintf("%s:%s", task.TaskDir().SharedAllocDir, task.Env[taskenv.AllocDir])
	taskLocalBind := fmt.Sprintf("%s:%s", task.TaskDir().LocalDir, task.Env[taskenv.TaskLocalDir])
	secretDirBind := fmt.Sprintf("%s:%s", task.TaskDir().SecretsDir, task.Env[taskenv.SecretsDir])
	binds := []string{allocDirBind, taskLocalBind, secretDirBind}

	// TODO support volume drivers
	// https://github.com/containers/libpod/pull/4548
	taskLocalBindVolume := false

	for _, userbind := range driverConfig.Volumes {
		src, dst, mode, err := parseVolumeSpec(userbind)
		if err != nil {
			return nil, fmt.Errorf("invalid docker volume %q: %v", userbind, err)
		}

		// Paths inside task dir are always allowed when using the default driver,
		// Relative paths are always allowed as they mount within a container
		// When a VolumeDriver is set, we assume we receive a binding in the format
		// volume-name:container-dest
		// Otherwise, we assume we receive a relative path binding in the format
		// relative/to/task:/also/in/container
		if taskLocalBindVolume {
			src = expandPath(task.TaskDir().Dir, src)
		} else {
			// Resolve dotted path segments
			src = filepath.Clean(src)
		}

		if !d.config.Volumes.Enabled && !isParentPath(task.AllocDir, src) {
			return nil, fmt.Errorf("volumes are not enabled; cannot mount host paths: %+q", userbind)
		}

		bind := src + ":" + dst
		if mode != "" {
			bind += ":" + mode
		}
		binds = append(binds, bind)
	}

	if selinuxLabel := d.config.Volumes.SelinuxLabel; selinuxLabel != "" {
		// Apply SELinux Label to each volume
		for i := range binds {
			binds[i] = fmt.Sprintf("%s:%s", binds[i], selinuxLabel)
		}
	}

	return binds, nil
}

// expandPath returns the absolute path of dir, relative to base if dir is relative path.
// base is expected to be an absolute path
func expandPath(base, dir string) string {
	if filepath.IsAbs(dir) {
		return filepath.Clean(dir)
	}

	return filepath.Clean(filepath.Join(base, dir))
}

func parseVolumeSpec(volBind string) (hostPath string, containerPath string, mode string, err error) {
	// using internal parser to preserve old parsing behavior.  Docker
	// parser has additional validators (e.g. mode validity) and accepts invalid output (per Nomad),
	// e.g. single path entry to be treated as a container path entry with an auto-generated host-path.
	//
	// Reconsider updating to use Docker parser when ready to make incompatible changes.
	parts := strings.Split(volBind, ":")
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("not <src>:<destination> format")
	}

	m := ""
	if len(parts) > 2 {
		m = parts[2]
	}

	return parts[0], parts[1], m, nil
}

// isParentPath returns true if path is a child or a descendant of parent path.
// Both inputs need to be absolute paths.
func isParentPath(parent, path string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && !strings.HasPrefix(rel, "..")
}
