package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/client/stats"
	"github.com/hashicorp/nomad/drivers/shared/eventer"
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
	pstructs "github.com/hashicorp/nomad/plugins/shared/structs"
	"github.com/pascomnet/nomad-driver-podman/iopodman"

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

	// configSpec is the hcl specification returned by the ConfigSchema RPC
	configSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"enabled": hclspec.NewDefault(
			hclspec.NewAttr("enabled", "bool", false),
			hclspec.NewLiteral("true"),
		),
	})

	// taskConfigSpec is the hcl specification for the driver config section of
	// a task within a job. It is returned in the TaskConfigSchema RPC
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"image": hclspec.NewAttr("image", "string", true),
		// "args":    hclspec.NewAttr("args", "list(string)", false),
	})

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
	config *Config

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

// Config is the driver configuration set by the SetConfig RPC call
type Config struct {
	Enabled bool `codec:"enabled"`
}

// TaskConfig is the driver configuration of a task within a job
type TaskConfig struct {
	Image string `codec:"image"`
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
		config:         &Config{},
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
	var config Config
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

	// try to connect if the driver is enabled
	if d.config.Enabled {
		varlinkConnection, err := d.getConnection()
		if err != nil {
			d.logger.Error("Could not connect to podman", "err", err)
		} else {
			defer varlinkConnection.Close()
			//  see if we get the system info from podman
			info, err := iopodman.GetInfo().Call(varlinkConnection)
			if err == nil {
				// yay! we can enable the driver
				health = drivers.HealthStateHealthy
				desc = "ready"
				// FIXME: we would have more data here...
				attrs["driver.podman"] = pstructs.NewBoolAttribute(true)
				attrs["driver.podman.version"] = pstructs.NewStringAttribute(info.Podman.Podman_version)
			} else {
				d.logger.Error("Cound not get podman info", "err", err)
			}
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

	// FIXME: use inspect... share code with StartTask
	filters := []string{"id=" + taskState.ContainerID}
	psOpts := iopodman.PsOpts{
		Filters: &filters,
	}
	psContainers, err := iopodman.Ps().Call(varlinkConnection, psOpts)
	if err != nil {
		return fmt.Errorf("failt to recover task, could not get container info: %v", err)
	}
	if len(psContainers) != 1 {
		return fmt.Errorf("failt to recover task, problem with Ps()")
	}
	// pid := strconv.FormatInt(psContainers[0].PidNum, 10)
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
	}

	d.tasks.Set(taskState.TaskConfig.ID, h)

	go h.run()

	return nil
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

	args := []string{driverConfig.Image}
	containerName := fmt.Sprintf("%s-%s", cfg.Name, cfg.AllocID)
	truep := true

	createOpts := iopodman.Create{
		Args:   args,
		Name:   &containerName,
		Detach: &truep,
	}

	containerID, err := iopodman.CreateContainer().Call(varlinkConnection, createOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failt to start task, could not create container: %v", err)
	}
	d.logger.Info("Created container", "container", containerID)

	_, err = iopodman.StartContainer().Call(varlinkConnection, containerID)
	if err != nil {
		return nil, nil, fmt.Errorf("failt to start task, could not start container: %v", err)
	}
	d.logger.Info("Started container", "containerId", containerID)

	inspectJson, err := iopodman.InspectContainer().Call(varlinkConnection, containerID)
	if err != nil {
		d.logger.Error("failt to inspect container", "err", err)
		return nil, nil, fmt.Errorf("failt to start task, could not inspect container : %v", err)
	}
	var inspectData iopodman.InspectContainerData
	err = json.Unmarshal([]byte(inspectJson), &inspectData)
	if err != nil {
		d.logger.Error("failt to unmarshal inspect container", "err", err)
		return nil, nil, fmt.Errorf("failt to start task, could not unmarshal inspect container : %v", err)
	}
	// pid := strconv.FormatInt(psContainers[0].PidNum, 10)
	d.logger.Debug("Started podman container", "container", containerID, "pid", inspectData.State.Pid, "ip", inspectData.NetworkSettings.IPAddress)

	net := &drivers.DriverNetwork{
		// 	// PortMap:       driverConfig.PortMap,
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
	}

	driverState := TaskState{
		ContainerID: containerID,
		TaskConfig:  cfg,
		StartedAt:   h.startedAt,
		Net:         net,
	}

	if err := handle.SetDriverState(&driverState); err != nil {
		d.logger.Error("failed to start task, error setting driver state", "error", err)
		//cleanup()
		return nil, nil, fmt.Errorf("failed to set driver state: %v", err)
	}

	d.tasks.Set(cfg.ID, h)

	go h.run()

	d.logger.Debug("DONE STARTING")

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

	varlinkConnection, err := d.getConnection()
	if err != nil {
		return fmt.Errorf("executor Shutdown failed, could not get podman connection: %v", err)

	}
	defer varlinkConnection.Close()
	d.logger.Debug("Stopping podman container", "container", handle.containerID)
	if _, err := iopodman.StopContainer().Call(varlinkConnection, handle.containerID, int64(timeout)); err != nil {
		return fmt.Errorf("executor Shutdown failed: %v", err)
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
		// grace period is chosen arbitrary here
		if err := handle.shutdown(1 * time.Minute); err != nil {
			handle.logger.Error("failed to destroy executor", "err", err)
		}
	}

	// FIXME: implement!
	// if d.config.DestroyContainers {
	// 	handle.logger.Debug("Destroying container", "container", handle.container.Name())
	// 	// delete the container itself
	// 	if err := handle.container.Destroy(); err != nil {
	// 		handle.logger.Error("failed to destroy lxc container", "err", err)
	// 	}
	// }

	// finally cleanup task map
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
	return fmt.Errorf("LXC driver does not support signals")
}

func (d *Driver) ExecTask(taskID string, cmd []string, timeout time.Duration) (*drivers.ExecTaskResult, error) {
	return nil, fmt.Errorf("LXC driver does not support exec")
}

func (d *Driver) getConnection() (*varlink.Connection, error) {
	// FIXME: a parameter for the socket would be nice
	varlinkConnection, err := varlink.NewConnection("unix://run/podman/io.podman")
	return varlinkConnection, err
}
