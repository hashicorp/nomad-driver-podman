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
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/nomad/nomad/structs"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/apiclient"
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
	spec "github.com/opencontainers/runtime-spec/specs-go"
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
		FSIsolation: drivers.FSIsolationNone,
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
	podmanClient  *PodmanClient
	podmanClient2 *apiclient.APIClient
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
		podmanClient2: apiclient.NewClient(logger),
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
		d.podmanClient2.SetSocketPath(config.SocketPath)
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
	info, err := d.podmanClient2.SystemInfo(d.ctx)
	if err != nil {
		d.logger.Error("Could not get podman info", "err", err)
	} else {
		// yay! we can enable the driver
		health = drivers.HealthStateHealthy
		desc = "ready"
		// TODO: we would have more data here... maybe we can return more details to nomad
		attrs["driver.podman"] = pstructs.NewBoolAttribute(true)
		attrs["driver.podman.version"] = pstructs.NewStringAttribute(info.Version)
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

	inspectData, err := d.podmanClient2.ContainerInspect(d.ctx, taskState.ContainerID)
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

	if inspectData.State.Running {
		d.logger.Info("Recovered a still running container", "container", inspectData.State.Pid)
		h.procState = drivers.TaskStateRunning
	} else if inspectData.State.Status == "exited" {
		// are we allowed to restart a stopped container?
		if d.config.RecoverStopped {
			d.logger.Debug("Found a stopped container, try to start it", "container", inspectData.State.Pid)
			if err = d.podmanClient2.ContainerStart(d.ctx, inspectData.ID); err != nil {
				d.logger.Warn("Recovery restart failed", "task", handle.Config.ID, "container", taskState.ContainerID, "err", err)
			} else {
				d.logger.Info("Restarted a container during recovery", "container", inspectData.ID)
				h.procState = drivers.TaskStateRunning
			}
		} else {
			// no, let's cleanup here to prepare for a StartTask()
			d.logger.Debug("Found a stopped container, removing it", "container", inspectData.ID)
			if err = d.podmanClient2.ContainerStart(d.ctx, inspectData.ID); err != nil {
				d.logger.Warn("Recovery cleanup failed", "task", handle.Config.ID, "container", inspectData.ID)
			}
			h.procState = drivers.TaskStateExited
		}
	} else {
		d.logger.Warn("Recovery restart failed, unknown container state", "state", inspectData.State.Status, "container", taskState.ContainerID)
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

	createOpts := apiclient.SpecGenerator{}
	createOpts.ContainerBasicConfig.LogConfiguration = &apiclient.LogConfig{}
	allArgs := []string{}
	if driverConfig.Command != "" {
		allArgs = append(allArgs, driverConfig.Command)
	}
	allArgs = append(allArgs, driverConfig.Args...)

	if driverConfig.Entrypoint != "" {
		createOpts.ContainerBasicConfig.Entrypoint = append(createOpts.ContainerBasicConfig.Entrypoint, driverConfig.Entrypoint)
	}

	containerName := BuildContainerName(cfg)
	cpuShares := uint64(cfg.Resources.LinuxResources.CPUShares)

	// ensure to include port_map into tasks environment map
	cfg.Env = taskenv.SetPortMapEnvs(cfg.Env, driverConfig.PortMap)

	procFilesystems, err := getProcFilesystems()
	if err != nil {
		return nil, nil, fmt.Errorf("Couldn't get Proc info: %v", err)
	}

	// -------------------------------------------------------------------------------------------
	// BASIC
	// -------------------------------------------------------------------------------------------
	createOpts.ContainerBasicConfig.Name = containerName
	createOpts.ContainerBasicConfig.Command = allArgs
	createOpts.ContainerBasicConfig.Env = cfg.Env
	createOpts.ContainerBasicConfig.Hostname = driverConfig.Hostname

	createOpts.ContainerBasicConfig.LogConfiguration.Path = cfg.StdoutPath

	// -------------------------------------------------------------------------------------------
	// STORAGE
	// -------------------------------------------------------------------------------------------
	createOpts.ContainerStorageConfig.Init = driverConfig.Init
	createOpts.ContainerStorageConfig.Image = driverConfig.Image
	createOpts.ContainerStorageConfig.InitPath = driverConfig.InitPath
	createOpts.ContainerStorageConfig.WorkDir = driverConfig.WorkingDir
	allMounts, err := d.containerMounts(cfg, &driverConfig)
	if err != nil {
		return nil, nil, err
	}
	createOpts.ContainerStorageConfig.Mounts = allMounts

	// -------------------------------------------------------------------------------------------
	// RESOURCES
	// -------------------------------------------------------------------------------------------
	createOpts.ContainerResourceConfig.ResourceLimits = &spec.LinuxResources{
		Memory: &spec.LinuxMemory{},
		CPU:    &spec.LinuxCPU{},
	}
	if driverConfig.MemoryReservation != "" {
		reservation, err := memoryInBytes(driverConfig.MemoryReservation)
		if err != nil {
			return nil, nil, err
		}
		createOpts.ContainerResourceConfig.ResourceLimits.Memory.Reservation = &reservation
	}

	if cfg.Resources.NomadResources.Memory.MemoryMB > 0 {
		limit := cfg.Resources.NomadResources.Memory.MemoryMB * 1024 * 1024
		createOpts.ContainerResourceConfig.ResourceLimits.Memory.Limit = &limit
	}
	if driverConfig.MemorySwap != "" {
		swap, err := memoryInBytes(driverConfig.MemorySwap)
		if err != nil {
			return nil, nil, err
		}
		createOpts.ContainerResourceConfig.ResourceLimits.Memory.Swap = &swap
	}
	if err == nil {
		cgroupv2 := false
		for _, l := range procFilesystems {
			cgroupv2 = cgroupv2 || strings.HasSuffix(l, "cgroup2")
		}
		if !cgroupv2 {
			swappiness := uint64(driverConfig.MemorySwappiness)
			createOpts.ContainerResourceConfig.ResourceLimits.Memory.Swappiness = &swappiness
		}
	}
	createOpts.ContainerResourceConfig.ResourceLimits.CPU.Shares = &cpuShares

	// -------------------------------------------------------------------------------------------
	// SECURITY
	// -------------------------------------------------------------------------------------------
	createOpts.ContainerSecurityConfig.CapAdd = driverConfig.CapAdd
	createOpts.ContainerSecurityConfig.CapDrop = driverConfig.CapDrop
	createOpts.ContainerSecurityConfig.User = cfg.User

	// -------------------------------------------------------------------------------------------
	// NETWORK
	// -------------------------------------------------------------------------------------------
	for _, strdns := range driverConfig.Dns {
		ipdns := net.ParseIP(strdns)
		if ipdns == nil {
			return nil, nil, fmt.Errorf("Invald dns server address")
		}
		createOpts.ContainerNetworkConfig.DNSServers = append(createOpts.ContainerNetworkConfig.DNSServers, ipdns)
	}
	// Configure network
	if cfg.NetworkIsolation != nil && cfg.NetworkIsolation.Path != "" {
		createOpts.ContainerNetworkConfig.NetNS.NSMode = apiclient.Path
		createOpts.ContainerNetworkConfig.NetNS.Value = cfg.NetworkIsolation.Path
	} else {
		if driverConfig.NetworkMode == "" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = apiclient.Bridge
		} else if driverConfig.NetworkMode == "bridge" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = apiclient.Bridge
		} else if driverConfig.NetworkMode == "host" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = apiclient.Host
		} else if driverConfig.NetworkMode == "none" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = apiclient.NoNetwork
		} else if driverConfig.NetworkMode == "slirp4netns" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = apiclient.Slirp
		} else {
			// FIXME: needs more work, parsing etc.
			return nil, nil, fmt.Errorf("Unknown/Unsupported network mode: %s", driverConfig.NetworkMode)
		}
	}

	// Setup port mapping and exposed ports
	if len(cfg.Resources.NomadResources.Networks) == 0 {
		d.logger.Debug("no network interfaces are available")
		if len(driverConfig.PortMap) > 0 {
			return nil, nil, fmt.Errorf("Trying to map ports but no network interface is available")
		}
	} else {
		publishedPorts := []apiclient.PortMapping{}
		network := cfg.Resources.NomadResources.Networks[0]
		allPorts := []structs.Port{}
		allPorts = append(allPorts, network.ReservedPorts...)
		allPorts = append(allPorts, network.DynamicPorts...)

		for _, port := range allPorts {
			hostPort := uint16(port.Value)
			// By default we will map the allocated port 1:1 to the container
			containerPort := uint16(port.Value)

			// If the user has mapped a port using port_map we'll change it here
			if mapped, ok := driverConfig.PortMap[port.Label]; ok {
				containerPort = uint16(mapped)
			}

			// we map both udp and tcp ports
			publishedPorts = append(publishedPorts, apiclient.PortMapping{
				HostIP:        network.IP,
				HostPort:      hostPort,
				ContainerPort: containerPort,
				Protocol:      "tcp",
			})
			publishedPorts = append(publishedPorts, apiclient.PortMapping{
				HostIP:        network.IP,
				HostPort:      hostPort,
				ContainerPort: containerPort,
				Protocol:      "udp",
			})
		}
		createOpts.ContainerNetworkConfig.PortMappings = publishedPorts
	}
	// -------------------------------------------------------------------------------------------

	containerID := ""
	recoverRunningContainer := false
	// check if there is a container with same name
	otherContainerInspect, err := d.podmanClient2.ContainerInspect(d.ctx, containerName)
	if err == nil {
		// ok, seems we found a container with similar name
		if otherContainerInspect.State.Running {
			// it's still running. So let's use it instead of creating a new one
			d.logger.Info("Detect running container with same name, we reuse it", "task", cfg.ID, "container", otherContainerInspect.ID)
			containerID = otherContainerInspect.ID
			recoverRunningContainer = true
		} else {
			// let's remove the old, dead container
			d.logger.Info("Detect stopped container with same name, removing it", "task", cfg.ID, "container", otherContainerInspect.ID)
			if err = d.podmanClient2.ContainerDelete(d.ctx, otherContainerInspect.ID, true, true); err != nil {
				return nil, nil, nstructs.WrapRecoverable(fmt.Sprintf("failed to remove dead container: %v", err), err)
			}
		}
	}

	if !recoverRunningContainer {
		createResponse, err := d.podmanClient2.ContainerCreate(d.ctx, createOpts)
		for _, w := range createResponse.Warnings {
			d.logger.Warn("Create Warning", "warning", w)
		}
		if err != nil {
			return nil, nil, fmt.Errorf("failed to start task, could not create container: %v", err)
		}
		containerID = createResponse.Id
	}

	cleanup := func() {
		d.logger.Debug("Cleaning up", "container", containerID)
		if err := d.podmanClient2.ContainerDelete(d.ctx, containerID, true, true); err != nil {
			d.logger.Error("failed to clean up from an error in Start", "error", err)
		}
	}

	if !recoverRunningContainer {
		if err = d.podmanClient2.ContainerStart(d.ctx, containerID); err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to start task, could not start container: %v", err)
		}
	}

	inspectData, err := d.podmanClient2.ContainerInspect(d.ctx, containerID)
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

func memoryInBytes(strmem string) (int64, error) {
	l := len(strmem)
	if l < 2 {
		return 0, fmt.Errorf("Invalid memory string: %s", strmem)
	}
	ival, err := strconv.Atoi(strmem[0 : l-1])
	if err != nil {
		return 0, err
	}

	switch strmem[l-1] {
	case 'b':
		return int64(ival), nil
	case 'k':
		return int64(ival) * 1024, nil
	case 'm':
		return int64(ival) * 1024 * 1024, nil
	case 'g':
		return int64(ival) * 1024 * 1024 * 1024, nil
	default:
		return 0, fmt.Errorf("Invalid memory string: %s", strmem)
	}
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
	err := d.podmanClient2.ContainerStop(d.ctx, handle.containerID, int(timeout.Seconds()))
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
		err := d.podmanClient2.ContainerStop(d.ctx, handle.containerID, 60)
		if err != nil {
			d.logger.Warn("failed to stop/kill container during destroy", "error", err)
		}
		// wait a while for stats emitter to collect exit code etc.
		for i := 0; i < 20; i++ {
			if !handle.isRunning() {
				break
			}
			time.Sleep(time.Millisecond * 250)
		}
		if handle.isRunning() {
			d.logger.Warn("stats emitter did not exit while stop/kill container during destroy", "error", err)
		}
	}

	if handle.removeContainerOnExit {
		err := d.podmanClient2.ContainerDelete(d.ctx, handle.containerID, true, true)
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

	return d.podmanClient2.ContainerKill(d.ctx, handle.containerID, signal)
}

// ExecTask function is used by the Nomad client to execute commands inside the task execution context.
func (d *Driver) ExecTask(taskID string, cmd []string, timeout time.Duration) (*drivers.ExecTaskResult, error) {
	return nil, fmt.Errorf("Podman driver does not support exec")
}

func (d *Driver) containerMounts(task *drivers.TaskConfig, driverConfig *TaskConfig) ([]spec.Mount, error) {
	binds := []spec.Mount{}
	binds = append(binds, spec.Mount{Source: task.TaskDir().SharedAllocDir, Destination: task.Env[taskenv.AllocDir], Type: "bind"})
	binds = append(binds, spec.Mount{Source: task.TaskDir().LocalDir, Destination: task.Env[taskenv.TaskLocalDir], Type: "bind"})
	binds = append(binds, spec.Mount{Source: task.TaskDir().SecretsDir, Destination: task.Env[taskenv.SecretsDir], Type: "bind"})

	// TODO support volume drivers
	// https://github.com/containers/libpod/pull/4548
	taskLocalBindVolume := true

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
		bind := spec.Mount{
			Source:      src,
			Destination: dst,
			Type:        "bind",
		}

		if mode != "" {
			bind.Options = append(bind.Options, mode)
		}
		binds = append(binds, bind)
	}

	if selinuxLabel := d.config.Volumes.SelinuxLabel; selinuxLabel != "" {
		// Apply SELinux Label to each volume
		for i := range binds {
			binds[i].Options = append(binds[i].Options, selinuxLabel)
		}
	}

	for _, dst := range driverConfig.Tmpfs {
		bind := spec.Mount{
			Destination: dst,
			Type:        "tmpfs",
		}
		binds = append(binds, bind)
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
