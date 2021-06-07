package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/armon/circbuf"
	"github.com/hashicorp/nomad/nomad/structs"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
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

	"github.com/containers/image/v5/docker"
	dockerArchive "github.com/containers/image/v5/docker/archive"
	ociArchive "github.com/containers/image/v5/oci/archive"
	"github.com/containers/image/v5/pkg/shortnames"
	"github.com/containers/image/v5/types"
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
		Exec:        true,
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
	podman *api.API

	// SystemInfo collected at first fingerprint query
	systemInfo api.Info
	// Queried from systemInfo: is podman running on a cgroupv2 system?
	cgroupV2 bool
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
	var pluginConfig PluginConfig
	if len(cfg.PluginConfig) != 0 {
		if err := base.MsgPackDecode(cfg.PluginConfig, &pluginConfig); err != nil {
			return err
		}
	}

	d.config = &pluginConfig
	if cfg.AgentConfig != nil {
		d.nomadConfig = cfg.AgentConfig.Driver
	}

	clientConfig := api.DefaultClientConfig()
	if pluginConfig.SocketPath != "" {
		clientConfig.SocketPath = pluginConfig.SocketPath
	}

	d.podman = api.NewClient(d.logger, clientConfig)
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
	info, err := d.podman.SystemInfo(d.ctx)
	if err != nil {
		d.logger.Error("Could not get podman info", "err", err)
	} else {
		// yay! we can enable the driver
		health = drivers.HealthStateHealthy
		desc = "ready"
		attrs["driver.podman"] = pstructs.NewBoolAttribute(true)
		attrs["driver.podman.version"] = pstructs.NewStringAttribute(info.Version.Version)
		attrs["driver.podman.rootless"] = pstructs.NewBoolAttribute(info.Host.Security.Rootless)
		attrs["driver.podman.cgroupVersion"] = pstructs.NewStringAttribute(info.Host.CGroupsVersion)
		if d.systemInfo.Version.Version == "" {
			// keep first received systemInfo in driver struct
			// it is used to toggle cgroup v1/v2, rootless/rootful behavior
			d.systemInfo = info
			d.cgroupV2 = info.Host.CGroupsVersion == "v2"
		}
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

	inspectData, err := d.podman.ContainerInspect(d.ctx, taskState.ContainerID)
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
			if err = d.podman.ContainerStart(d.ctx, inspectData.ID); err != nil {
				d.logger.Warn("Recovery restart failed", "task", handle.Config.ID, "container", taskState.ContainerID, "err", err)
			} else {
				d.logger.Info("Restarted a container during recovery", "container", inspectData.ID)
				h.procState = drivers.TaskStateRunning
			}
		} else {
			// no, let's cleanup here to prepare for a StartTask()
			d.logger.Debug("Found a stopped container, removing it", "container", inspectData.ID)
			if err = d.podman.ContainerStart(d.ctx, inspectData.ID); err != nil {
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
	return BuildContainerNameForTask(cfg.Name, cfg)
}

// BuildContainerName returns the podman container name for a specific Task in our group
func BuildContainerNameForTask(taskName string, cfg *drivers.TaskConfig) string {
	return fmt.Sprintf("%s-%s", taskName, cfg.AllocID)
}

// StartTask creates and starts a new Container based on the given TaskConfig.
func (d *Driver) StartTask(cfg *drivers.TaskConfig) (*drivers.TaskHandle, *drivers.DriverNetwork, error) {
	rootless := d.systemInfo.Host.Security.Rootless

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

	createOpts := api.SpecGenerator{}
	createOpts.ContainerBasicConfig.LogConfiguration = &api.LogConfig{}
	allArgs := []string{}
	if driverConfig.Command != "" {
		allArgs = append(allArgs, driverConfig.Command)
	}
	allArgs = append(allArgs, driverConfig.Args...)

	if driverConfig.Entrypoint != "" {
		createOpts.ContainerBasicConfig.Entrypoint = append(createOpts.ContainerBasicConfig.Entrypoint, driverConfig.Entrypoint)
	}

	containerName := BuildContainerName(cfg)

	// ensure to include port_map into tasks environment map
	cfg.Env = taskenv.SetPortMapEnvs(cfg.Env, driverConfig.PortMap)

	// Basic config options
	createOpts.ContainerBasicConfig.Name = containerName
	createOpts.ContainerBasicConfig.Command = allArgs
	createOpts.ContainerBasicConfig.Env = cfg.Env
	createOpts.ContainerBasicConfig.Hostname = driverConfig.Hostname
	createOpts.ContainerBasicConfig.Sysctl = driverConfig.Sysctl
	createOpts.ContainerBasicConfig.Terminal = driverConfig.Tty
	createOpts.ContainerBasicConfig.Labels = driverConfig.Labels

	createOpts.ContainerBasicConfig.LogConfiguration.Path = cfg.StdoutPath

	// Storage config options
	createOpts.ContainerStorageConfig.Init = driverConfig.Init
	createOpts.ContainerStorageConfig.Image = driverConfig.Image
	createOpts.ContainerStorageConfig.InitPath = driverConfig.InitPath
	createOpts.ContainerStorageConfig.WorkDir = driverConfig.WorkingDir
	allMounts, err := d.containerMounts(cfg, &driverConfig)
	if err != nil {
		return nil, nil, err
	}
	createOpts.ContainerStorageConfig.Mounts = allMounts

	// Resources config options
	createOpts.ContainerResourceConfig.ResourceLimits = &spec.LinuxResources{
		Memory: &spec.LinuxMemory{},
		CPU:    &spec.LinuxCPU{},
	}

	hard, soft, err := memoryLimits(cfg.Resources.NomadResources.Memory, driverConfig.MemoryReservation)
	if err != nil {
		return nil, nil, err
	}
	createOpts.ContainerResourceConfig.ResourceLimits.Memory.Reservation = soft
	createOpts.ContainerResourceConfig.ResourceLimits.Memory.Limit = hard

	if driverConfig.MemorySwap != "" {
		swap, err := memoryInBytes(driverConfig.MemorySwap)
		if err != nil {
			return nil, nil, err
		}
		createOpts.ContainerResourceConfig.ResourceLimits.Memory.Swap = &swap
	}
	if !d.cgroupV2 {
		swappiness := uint64(driverConfig.MemorySwappiness)
		createOpts.ContainerResourceConfig.ResourceLimits.Memory.Swappiness = &swappiness
	}
	// FIXME: can fail for nonRoot due to missing cpu limit delegation permissions,
	//        see https://github.com/containers/podman/blob/master/troubleshooting.md
	if !rootless {
		cpuShares := uint64(cfg.Resources.LinuxResources.CPUShares)
		createOpts.ContainerResourceConfig.ResourceLimits.CPU.Shares = &cpuShares
	}

	// Security config options
	createOpts.ContainerSecurityConfig.CapAdd = driverConfig.CapAdd
	createOpts.ContainerSecurityConfig.CapDrop = driverConfig.CapDrop
	createOpts.ContainerSecurityConfig.User = cfg.User

	// Network config options
	if cfg.DNS != nil {
		for _, strdns := range cfg.DNS.Servers {
			ipdns := net.ParseIP(strdns)
			if ipdns == nil {
				return nil, nil, fmt.Errorf("Invald dns server address")
			}
			createOpts.ContainerNetworkConfig.DNSServers = append(createOpts.ContainerNetworkConfig.DNSServers, ipdns)
		}
		createOpts.ContainerNetworkConfig.DNSSearch = append(createOpts.ContainerNetworkConfig.DNSSearch, cfg.DNS.Searches...)
		createOpts.ContainerNetworkConfig.DNSOptions = append(createOpts.ContainerNetworkConfig.DNSOptions, cfg.DNS.Options...)
	}
	// Configure network
	if cfg.NetworkIsolation != nil && cfg.NetworkIsolation.Path != "" {
		createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Path
		createOpts.ContainerNetworkConfig.NetNS.Value = cfg.NetworkIsolation.Path
	} else {
		if driverConfig.NetworkMode == "" {
			if !rootless {
				// should we join the group shared network namespace?
				if cfg.NetworkIsolation != nil && cfg.NetworkIsolation.Mode == drivers.NetIsolationModeGroup {
					// yes, join the group ns namespace
					createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Path
					createOpts.ContainerNetworkConfig.NetNS.Value = cfg.NetworkIsolation.Path
				} else {
					// no, simply attach a rootful container to the default podman bridge
					createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Bridge
				}
			} else {
				// slirp4netns is default for rootless podman
				createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Slirp
			}
		} else if driverConfig.NetworkMode == "bridge" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Bridge
		} else if driverConfig.NetworkMode == "host" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Host
		} else if driverConfig.NetworkMode == "none" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.NoNetwork
		} else if driverConfig.NetworkMode == "slirp4netns" {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Slirp
		} else if strings.HasPrefix(driverConfig.NetworkMode, "container:") {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.FromContainer
			createOpts.ContainerNetworkConfig.NetNS.Value = strings.TrimPrefix(driverConfig.NetworkMode, "container:")
		} else if strings.HasPrefix(driverConfig.NetworkMode, "ns:") {
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Path
			createOpts.ContainerNetworkConfig.NetNS.Value = strings.TrimPrefix(driverConfig.NetworkMode, "ns:")
		} else if strings.HasPrefix(driverConfig.NetworkMode, "task:") {
			otherTaskName := strings.TrimPrefix(driverConfig.NetworkMode, "task:")
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.FromContainer
			createOpts.ContainerNetworkConfig.NetNS.Value = BuildContainerNameForTask(otherTaskName, cfg)
		} else {
			return nil, nil, fmt.Errorf("Unknown/Unsupported network mode: %s", driverConfig.NetworkMode)
		}
	}

	portMappings, err := d.portMappings(cfg, driverConfig)
	if err != nil {
		return nil, nil, err
	}
	createOpts.ContainerNetworkConfig.PortMappings = portMappings

	containerID := ""
	recoverRunningContainer := false
	// check if there is a container with same name
	otherContainerInspect, err := d.podman.ContainerInspect(d.ctx, containerName)
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
			if err = d.podman.ContainerDelete(d.ctx, otherContainerInspect.ID, true, true); err != nil {
				return nil, nil, nstructs.WrapRecoverable(fmt.Sprintf("failed to remove dead container: %v", err), err)
			}
		}
	}

	if !recoverRunningContainer {
		imageID, err := d.createImage(createOpts.Image, &driverConfig.Auth, driverConfig.ForcePull)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create image: %s: %w", createOpts.Image, err)
		}
		createOpts.Image = imageID

		createResponse, err := d.podman.ContainerCreate(d.ctx, createOpts)
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
		if err := d.podman.ContainerDelete(d.ctx, containerID, true, true); err != nil {
			d.logger.Error("failed to clean up from an error in Start", "error", err)
		}
	}

	if !recoverRunningContainer {
		if err = d.podman.ContainerStart(d.ctx, containerID); err != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to start task, could not start container: %v", err)
		}
	}

	inspectData, err := d.podman.ContainerInspect(d.ctx, containerID)
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

func memoryLimits(r drivers.MemoryResources, reservation string) (hard, soft *int64, err error) {
	memoryMax := r.MemoryMaxMB * 1024 * 1024
	memory := r.MemoryMB * 1024 * 1024

	var reserved *int64
	if reservation != "" {
		reservation, err := memoryInBytes(reservation)
		if err != nil {
			return nil, nil, err
		}
		reserved = &reservation
	}

	if memoryMax > 0 {
		if reserved != nil && *reserved < memory {
			memory = *reserved
		}
		return &memoryMax, &memory, nil
	}

	if memory > 0 {
		return &memory, reserved, nil
	}

	// We may never actually be here
	return nil, reserved, nil
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

// Creates the requested image if missing from storage
// returns the 64-byte image ID as an unique image identifier
func (d *Driver) createImage(image string, auth *AuthConfig, forcePull bool) (string, error) {
	var imageID string
	imageName := image
	// If it is a shortname, we should not have to worry
	// Let podman deal with it according to user configuration
	if !shortnames.IsShortName(image) {
		imageRef, err := parseImage(image)
		if err != nil {
			return imageID, fmt.Errorf("invalid image reference %s: %w", image, err)
		}
		switch transport := imageRef.Transport().Name(); transport {
		case "docker":
			imageName = imageRef.DockerReference().String()
		case "oci-archive", "docker-archive":
			// For archive transports, we cannot ask for a pull or
			// check for existence in the API without image plumbing.
			// Load the images instead
			archiveData := imageRef.StringWithinTransport()
			path := strings.Split(archiveData, ":")[0]
			d.logger.Debug("Load image archive", "path", path)
			imageName, err = d.podman.ImageLoad(d.ctx, path)
			if err != nil {
				return imageID, fmt.Errorf("error while loading image: %w", err)
			}
		}
	}

	imageID, err := d.podman.ImageInspectID(d.ctx, imageName)
	if err != nil {
		// If ImageInspectID errors, continue the operation and try
		// to pull the image instead
		d.logger.Warn("Unable to check for local image", "image", imageName, "err", err)
	}
	if !forcePull && imageID != "" {
		d.logger.Info("Found imageID", imageID, "for image", imageName, "in local storage")
		return imageID, nil
	}

	d.logger.Debug("Pull image", "image", imageName)
	imageAuth := api.ImageAuthConfig{
		Username: auth.Username,
		Password: auth.Password,
	}
	if imageID, err = d.podman.ImagePull(d.ctx, imageName, imageAuth); err != nil {
		return imageID, fmt.Errorf("failed to start task, unable to pull image %s : %w", imageName, err)
	}
	d.logger.Debug("Pulled image ID", "imageID", imageID)
	return imageID, nil
}

func parseImage(image string) (types.ImageReference, error) {
	var transport, name string
	parts := strings.SplitN(image, ":", 2)
	// In case the transport is missing, assume docker://
	if len(parts) == 1 {
		transport = "docker"
		name = "//" + image
	} else {
		transport = parts[0]
		name = parts[1]
	}

	switch transport {
	case "docker":
		return docker.ParseReference(name)
	case "oci-archive":
		return ociArchive.ParseReference(name)
	case "docker-archive":
		return dockerArchive.ParseReference(name)
	default:
		// We could have both an unknown/malformed transport
		// or an image:tag with default "docker://" transport omitted
		ref, err := docker.ParseReference("//" + image)
		if err == nil {
			return ref, nil
		}
		return nil, fmt.Errorf("unsupported transport %s", transport)
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
	err := d.podman.ContainerStop(d.ctx, handle.containerID, int(timeout.Seconds()))
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
		d.logger.Debug("Have to destroyTask but container is still running", "containerID", handle.containerID)
		// we can not do anything, so catching the error is useless
		err := d.podman.ContainerStop(d.ctx, handle.containerID, 60)
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
		err := d.podman.ContainerDelete(d.ctx, handle.containerID, true, true)
		if err != nil {
			d.logger.Warn("Could not remove container", "container", handle.containerID, "error", err)
		}
	}

	d.tasks.Delete(taskID)
	return nil
}

// InspectTask function returns detailed status information for the referenced taskID.
func (d *Driver) InspectTask(taskID string) (*drivers.TaskStatus, error) {
	d.logger.Debug("InspectTask called")
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

	return d.podman.ContainerKill(d.ctx, handle.containerID, signal)
}

// ExecTask function is used by the Nomad client to execute scripted health checks inside the task execution context.
func (d *Driver) ExecTask(taskID string, cmd []string, timeout time.Duration) (*drivers.ExecTaskResult, error) {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}
	createRequest := api.ExecConfig{
		Command:      cmd,
		Tty:          false,
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
	}
	ctx, cancel := context.WithTimeout(d.ctx, timeout)
	defer cancel()
	sessionId, err := d.podman.ExecCreate(ctx, handle.containerID, createRequest)
	if err != nil {
		d.logger.Error("Unable to create ExecTask session", "err", err)
		return nil, err
	}
	stdout, err := circbuf.NewBuffer(int64(drivers.CheckBufSize))
	if err != nil {
		d.logger.Error("ExecTask session failed, unable to allocate stdout buffer", "sessionId", sessionId, "err", err)
		return nil, err
	}
	stderr, err := circbuf.NewBuffer(int64(drivers.CheckBufSize))
	if err != nil {
		d.logger.Error("ExecTask session failed, unable to allocate stderr buffer", "sessionId", sessionId, "err", err)
		return nil, err
	}
	startRequest := api.ExecStartRequest{
		Tty:          false,
		AttachInput:  false,
		AttachOutput: true,
		Stdout:       stdout,
		AttachError:  true,
		Stderr:       stderr,
	}
	err = d.podman.ExecStart(ctx, sessionId, startRequest)
	if err != nil {
		d.logger.Error("ExecTask session returned with error", "sessionId", sessionId, "err", err)
		return nil, err
	}

	inspectData, err := d.podman.ExecInspect(ctx, sessionId)
	if err != nil {
		d.logger.Error("Unable to inspect finished ExecTask session", "sessionId", sessionId, "err", err)
		return nil, err
	}
	execResult := &drivers.ExecTaskResult{
		ExitResult: &drivers.ExitResult{
			ExitCode: inspectData.ExitCode,
		},
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}
	d.logger.Trace("ExecTask result", "code", execResult.ExitResult.ExitCode, "out", string(execResult.Stdout), "err", string(execResult.Stderr))

	return execResult, nil
}

// ExecTask function is used by the Nomad client to execute commands inside the task execution context.
// i.E. nomad alloc exec ....
func (d *Driver) ExecTaskStreaming(ctx context.Context, taskID string, execOptions *drivers.ExecOptions) (*drivers.ExitResult, error) {
	handle, ok := d.tasks.Get(taskID)
	if !ok {
		return nil, drivers.ErrTaskNotFound
	}

	createRequest := api.ExecConfig{
		Command:      execOptions.Command,
		Tty:          execOptions.Tty,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	sessionId, err := d.podman.ExecCreate(ctx, handle.containerID, createRequest)
	if err != nil {
		d.logger.Error("Unable to create exec session", "err", err)
		return nil, err
	}

	startRequest := api.ExecStartRequest{
		Tty:          execOptions.Tty,
		AttachInput:  createRequest.AttachStdin,
		Stdin:        execOptions.Stdin,
		AttachOutput: createRequest.AttachStdout,
		Stdout:       execOptions.Stdout,
		AttachError:  createRequest.AttachStderr,
		Stderr:       execOptions.Stderr,
		ResizeCh:     execOptions.ResizeCh,
	}
	err = d.podman.ExecStart(ctx, sessionId, startRequest)
	if err != nil {
		d.logger.Error("Exec session returned with error", "sessionId", sessionId, "err", err)
		return nil, err
	}

	inspectData, err := d.podman.ExecInspect(ctx, sessionId)
	if err != nil {
		d.logger.Error("Unable to inspect finished exec session", "sessionId", sessionId, "err", err)
		return nil, err
	}
	exitResult := drivers.ExitResult{
		ExitCode: inspectData.ExitCode,
		Err:      err,
	}
	return &exitResult, nil
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

func (d *Driver) portMappings(taskCfg *drivers.TaskConfig, driverCfg TaskConfig) ([]api.PortMapping, error) {
	if taskCfg.Resources.Ports != nil && len(driverCfg.Ports) > 0 && len(driverCfg.PortMap) > 0 {
		return nil, errors.New("Invalid port declaration; use of port_map and ports")
	}

	if len(driverCfg.PortMap) > 0 && len(taskCfg.Resources.NomadResources.Networks) == 0 {
		return nil, fmt.Errorf("Trying to map ports but no network interface is available")
	}

	if taskCfg.Resources.Ports == nil && len(driverCfg.Ports) > 0 {
		return nil, errors.New("No ports defined in network stanza")
	}

	var publishedPorts []api.PortMapping
	if len(driverCfg.Ports) > 0 {
		for _, port := range driverCfg.Ports {
			mapping, ok := taskCfg.Resources.Ports.Get(port)
			if !ok {
				return nil, fmt.Errorf("Port %q not found, check network stanza", port)
			}
			to := mapping.To
			if to == 0 {
				to = mapping.Value
			}
			publishedPorts = append(publishedPorts, api.PortMapping{
				HostIP:        mapping.HostIP,
				HostPort:      uint16(mapping.Value),
				ContainerPort: uint16(to),
				Protocol:      "tcp",
			})
			publishedPorts = append(publishedPorts, api.PortMapping{
				HostIP:        mapping.HostIP,
				HostPort:      uint16(mapping.Value),
				ContainerPort: uint16(to),
				Protocol:      "udp",
			})
		}
	} else if len(driverCfg.PortMap) > 0 {
		// DEPRECATED: This style of PortMapping was Deprecated in Nomad 0.12
		network := taskCfg.Resources.NomadResources.Networks[0]
		allPorts := []structs.Port{}
		allPorts = append(allPorts, network.ReservedPorts...)
		allPorts = append(allPorts, network.DynamicPorts...)

		for _, port := range allPorts {
			hostPort := uint16(port.Value)
			// By default we will map the allocated port 1:1 to the container
			containerPort := uint16(port.Value)

			// If the user has mapped a port using port_map we'll change it here
			if mapped, ok := driverCfg.PortMap[port.Label]; ok {
				containerPort = uint16(mapped)
			}

			publishedPorts = append(publishedPorts, api.PortMapping{
				HostIP:        network.IP,
				HostPort:      hostPort,
				ContainerPort: containerPort,
				Protocol:      "tcp",
			})
			publishedPorts = append(publishedPorts, api.PortMapping{
				HostIP:        network.IP,
				HostPort:      hostPort,
				ContainerPort: containerPort,
				Protocol:      "udp",
			})
		}
	}
	return publishedPorts, nil
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
