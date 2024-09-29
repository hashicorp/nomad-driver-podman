// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
	"os"

	"github.com/armon/circbuf"
	"github.com/containers/image/v5/docker"
	dockerArchive "github.com/containers/image/v5/docker/archive"
	ociArchive "github.com/containers/image/v5/oci/archive"
	"github.com/containers/image/v5/pkg/shortnames"
	"github.com/containers/image/v5/types"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad-driver-podman/registry"
	"github.com/hashicorp/nomad-driver-podman/version"
	"github.com/hashicorp/nomad/client/lib/cpustats"
	"github.com/hashicorp/nomad/client/taskenv"
	"github.com/hashicorp/nomad/drivers/shared/eventer"
	nstructs "github.com/hashicorp/nomad/nomad/structs"
	"github.com/hashicorp/nomad/plugins/base"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/hashicorp/nomad/plugins/drivers/fsisolation"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
	pstructs "github.com/hashicorp/nomad/plugins/shared/structs"
	spec "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/ryanuber/go-glob"
	"golang.org/x/sync/singleflight"
)

const (
	// pluginName is the name of the plugin
	pluginName = "podman"

	// fingerprintPeriod is the interval at which the driver will send fingerprint responses
	fingerprintPeriod = 30 * time.Second

	// taskHandleVersion is the version of task handle which this driver sets
	// and understands how to decode driver state
	taskHandleVersion = 1

	LOG_DRIVER_NOMAD    = "nomad"
	LOG_DRIVER_JOURNALD = "journald"

	labelAllocID       = "com.hashicorp.nomad.alloc_id"
	labelJobName       = "com.hashicorp.nomad.job_name"
	labelJobID         = "com.hashicorp.nomad.job_id"
	labelTaskGroupName = "com.hashicorp.nomad.task_group_name"
	labelTaskName      = "com.hashicorp.nomad.task_name"
	labelNamespace     = "com.hashicorp.nomad.namespace"
	labelNodeName      = "com.hashicorp.nomad.node_name"
	labelNodeID        = "com.hashicorp.nomad.node_id"
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
		FSIsolation: fsisolation.Image,
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

	// compute contains information about the available cpu compute
	compute cpustats.Compute

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

	// Podman clients
	podmanClients map[string]*api.API

	// For any call where it's unspecified/unknown which podman should be used
	defaultPodman *api.API

	// singleflight group to prevent parallel image downloads
	pullGroup singleflight.Group
}

// TaskState is the state which is encoded in the handle returned in
// StartTask. This information is needed to rebuild the task state and handler
// during recovery.
type TaskState struct {
	TaskConfig  *drivers.TaskConfig
	ContainerID string
	StartedAt   time.Time
	Net         *drivers.DriverNetwork
	LogStreamer bool
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

	var timeout time.Duration
	if pluginConfig.ClientHttpTimeout != "" {
		t, err := time.ParseDuration(pluginConfig.ClientHttpTimeout)
		if err != nil {
			return err
		}
		timeout = t
	}

	if len(d.config.Socket) > 0 && d.config.SocketPath != "" {
		return fmt.Errorf("error: can't define socket blocks and socket_path, they're mutually exclusive.")
	} else if len(d.config.Socket) > 0 {
		d.podmanClients = d.makePodmanClients(d.config.Socket, timeout)
	} else if d.config.SocketPath != "" {
		d.podmanClients = make(map[string]*api.API)
		d.podmanClients["default"] = d.newPodmanClient(timeout, d.config.SocketPath, true)
		d.defaultPodman = d.podmanClients["default"]
	} else {
		d.podmanClients = make(map[string]*api.API)
		uid := os.Getuid()
		if uid == 0 {
			d.podmanClients["default"] = d.newPodmanClient(timeout, "unix:///run/podman/podman.sock", true)
		} else {
			d.podmanClients["default"] = d.newPodmanClient(timeout, fmt.Sprintf("unix:///run/user/%d/podman/podman.sock", uid), true)
		}
		d.defaultPodman = d.podmanClients["default"]
	}
	d.compute = cfg.AgentConfig.Compute()

	return nil
}

func (d *Driver) makePodmanClients(sockets []PluginSocketConfig, timeout time.Duration) map[string]*api.API {
	podmanClients := make(map[string]*api.API)
	var firstEntry *api.API
	foundDefaultPodman := false
	for i, sock := range sockets {
		var podmanClient *api.API
		if sock.Name == "" {
			sock.Name = "default"
		}
		sock.Name = cleanUpSocketName(sock.Name)
		if sock.Name == "default" && !foundDefaultPodman {
			foundDefaultPodman = true
			podmanClient = d.newPodmanClient(timeout, sock.SocketPath, true)
			d.defaultPodman = podmanClient
		} else {
			podmanClient = d.newPodmanClient(timeout, sock.SocketPath, false)
		}
		// Case when there are two socket blocks with the same name or
		// if there's one with name="default" and one without name.
		if _, ok := podmanClients[sock.Name]; ok {
			d.logger.Error(fmt.Sprintf("There is already a socket with the name: %s ", sock.Name))
		} else {
			podmanClients[sock.Name] = podmanClient
		}

		if i == 0 {
			firstEntry = podmanClient
		}
	}

	// If no socket was default, the first entry becomes the default
	if foundDefaultPodman == false {
		firstEntry.SetClientAsDefault(true)
		d.defaultPodman = firstEntry
	}
	return podmanClients
}

// We need to make a "clean" name that can be used safely in logging, journald, attributes in the UI, ...
func cleanUpSocketName(name string) string {
	var result strings.Builder
    for i := 0; i < len(name); i++ {
        b := name[i]
        if ! (('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z') || ('0' <= b && b <= '9')) {
			result.WriteByte('_')
		} else {
            result.WriteByte(b)
        }
    }
    return result.String()
}

func (d *Driver) getPodmanClient(clientName string) (*api.API, error) {
	p, ok := d.podmanClients[clientName]
	if ok {
		return p, nil
	}
	return nil, fmt.Errorf("podman client with name %s was not found, check your podman driver config", clientName)
}

// newPodmanClient returns Podman client configured with the provided timeout.
// This method must be called after the driver configuration has been loaded.
func (d *Driver) newPodmanClient(timeout time.Duration, socketPath string, defaultPodman bool) *api.API {
	clientConfig := api.DefaultClientConfig()
	clientConfig.HttpTimeout = timeout
	clientConfig.SocketPath = socketPath
	clientConfig.DefaultPodman = defaultPodman

	return api.NewClient(d.logger, clientConfig)
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
	// emit warnings about known bad configs
	d.config.LogWarnings(d.logger)
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
	attrs := map[string]*pstructs.Attribute{}
	allClientsAreHealthy := true
	unhealthyClients := []string{}
	allClientsAreUnhealthy := true

	for name, podmanClient := range d.podmanClients {
		// Ping podman api
		apiVersion, err := podmanClient.Ping(d.ctx)
		attrPrefix := fmt.Sprintf("driver.podman.%s", name)

		if err != nil || apiVersion == "" {
			d.logger.Error("Could not get podman version", "error", err)
			attrs[fmt.Sprintf("%s.health", attrPrefix)] = pstructs.NewStringAttribute("unhealthy")
			allClientsAreHealthy = false
			unhealthyClients = append(unhealthyClients, name)
			continue
		}

		info, err := podmanClient.SystemInfo(d.ctx)
		if err != nil {
			d.logger.Error("Could not get podman info", "error", err)
			attrs[fmt.Sprintf("%s.health", attrPrefix)] = pstructs.NewStringAttribute("unhealthy")
			allClientsAreHealthy = false
			unhealthyClients = append(unhealthyClients, name)
			continue
		}
		allClientsAreUnhealthy = false

		podmanClient.SetRootless(info.Host.Security.Rootless)
		podmanClient.SetCgroupV2(info.Host.CGroupsVersion == "v2")
		podmanClient.SetCgroupMgr(info.Host.CgroupManager)
		podmanClient.SetAppArmor(info.Host.Security.AppArmorEnabled)

		attrs[fmt.Sprintf("%s.appArmor", attrPrefix)] = pstructs.NewBoolAttribute(info.Host.Security.AppArmorEnabled)
		attrs[fmt.Sprintf("%s.capabilities", attrPrefix)] = pstructs.NewStringAttribute(info.Host.Security.DefaultCapabilities)
		attrs[fmt.Sprintf("%s.cgroupVersion", attrPrefix)] = pstructs.NewStringAttribute(info.Host.CGroupsVersion)
		attrs[fmt.Sprintf("%s.defaultPodman", attrPrefix)] = pstructs.NewBoolAttribute(podmanClient.IsDefaultClient())
		attrs[fmt.Sprintf("%s.health", attrPrefix)] = pstructs.NewStringAttribute("healthy")
		attrs[fmt.Sprintf("%s.ociRuntime", attrPrefix)] = pstructs.NewStringAttribute(info.Host.OCIRuntime.Path)
		attrs[fmt.Sprintf("%s.rootless", attrPrefix)] = pstructs.NewBoolAttribute(info.Host.Security.Rootless)
		attrs[fmt.Sprintf("%s.seccompEnabled", attrPrefix)] = pstructs.NewBoolAttribute(info.Host.Security.SECCOMPEnabled)
		attrs[fmt.Sprintf("%s.selinuxEnabled", attrPrefix)] = pstructs.NewBoolAttribute(info.Host.Security.SELinuxEnabled)
		attrs[fmt.Sprintf("%s.socket", attrPrefix)] = pstructs.NewStringAttribute(info.Host.RemoteSocket.Path)
		attrs[fmt.Sprintf("%s.version", attrPrefix)] = pstructs.NewStringAttribute(apiVersion)
	}

	if allClientsAreHealthy {
		attrs["driver.podman"] = pstructs.NewBoolAttribute(true)
		return &drivers.Fingerprint{
			Attributes:        attrs,
			Health:            drivers.HealthStateHealthy,
			HealthDescription: "All Podman sockets are responding.",
		}
	} else if allClientsAreUnhealthy {
		attrs["driver.podman"] = pstructs.NewBoolAttribute(false)
		return &drivers.Fingerprint{
			Attributes:        attrs,
			Health:            drivers.HealthStateUndetected,
			HealthDescription: "Cannot connect to any Podman socket.",
		}
	} else {
		attrs["driver.podman"] = pstructs.NewBoolAttribute(true)
		slices.Sort(unhealthyClients) // If not sorted, we generate a new fingerprint log entry every time the order changes
		return &drivers.Fingerprint{
			Attributes:        attrs,
			Health:            drivers.HealthStateUnhealthy,
			HealthDescription: fmt.Sprintf("Cannot fingerprint certain podman sockets: %s", strings.Join(unhealthyClients[:], ", ")),
		}


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
		return fmt.Errorf("failed to decode task state from handle: %w", err)
	}
	d.logger.Debug("Checking for recoverable task", "task", handle.Config.Name, "taskid", handle.Config.ID, "container", taskState.ContainerID)

	var inspectData api.InspectContainerData
	// We need to parse the task config for our driver to be able to find the Socket field (*drivers.TaskConfig itself is a generic task struct from the nomad repo)
	var podmanTaskConfig TaskConfig
	if err := taskState.TaskConfig.DecodeDriverConfig(&podmanTaskConfig); err != nil {
		return fmt.Errorf("error: cannot decode task config")
	}
	podmanTaskSocketName := podmanTaskConfig.Socket
	if podmanTaskConfig.Socket == "" {
		podmanTaskSocketName = "default"
	}
	taskPodmanClient, err := d.getPodmanClient(podmanTaskSocketName)
	if err == nil {
		inspectData, err = taskPodmanClient.ContainerInspect(d.ctx, taskState.ContainerID)
		if errors.Is(err, api.ContainerNotFound) {
			d.logger.Debug("Recovery lookup found no container", "task", handle.Config.ID, "container", taskState.ContainerID, "error", err)
			return err
		} else if err != nil {
			d.logger.Warn("Recovery lookup failed", "task", handle.Config.ID, "container", taskState.ContainerID, "error", err)
			return err
		}
	} else {
		d.logger.Warn("Did not find podman client for this task", "task", handle.Config.ID, "container", taskState.ContainerID, "podmanClient", podmanTaskSocketName, "error", err)
		return fmt.Errorf("error: cannot find the podman socket for this task. Socket might have been removed from driver config but still referenced in Task. podmanClient=%s", podmanTaskSocketName)
	}

	h := &TaskHandle{
		containerID:           taskState.ContainerID,
		driver:                d,
		podmanClient:          taskPodmanClient,
		taskConfig:            taskState.TaskConfig,
		procState:             drivers.TaskStateUnknown,
		startedAt:             taskState.StartedAt,
		exitResult:            &drivers.ExitResult{},
		logger:                d.logger.Named(fmt.Sprintf("podman.%s", podmanTaskSocketName)),
		logPointer:            time.Now(), // do not rewind log to the startetAt date.
		logStreamer:           taskState.LogStreamer,
		collectionInterval:    time.Second,
		totalCPUStats:         cpustats.New(d.compute),
		userCPUStats:          cpustats.New(d.compute),
		systemCPUStats:        cpustats.New(d.compute),
		removeContainerOnExit: d.config.GC.Container,
	}

	switch {
	case inspectData.State.Running:
		d.logger.Info("Recovered a still running container", "container", inspectData.State.Pid)
		h.procState = drivers.TaskStateRunning
	case inspectData.State.Status == "exited":
		// Are we allowed to restart a stopped container?
		if d.config.RecoverStopped {
			d.logger.Debug("Found a stopped container, try to start it", "container", inspectData.State.Pid, "podman client", podmanTaskSocketName)
			if err = taskPodmanClient.ContainerStart(d.ctx, inspectData.ID); err != nil {
				d.logger.Warn("Recovery restart failed", "task", handle.Config.ID, "container", taskState.ContainerID, "podman client", podmanTaskSocketName, "error", err)
			} else {
				d.logger.Info("Restarted a container during recovery", "container", inspectData.ID, "podman client", podmanTaskSocketName)
				h.procState = drivers.TaskStateRunning
			}
		} else {
			// No, let's cleanup here to prepare for a StartTask()
			d.logger.Debug("Found a stopped container, removing it", "container", inspectData.ID, "podman client", podmanTaskSocketName)
			if err = taskPodmanClient.ContainerDelete(d.ctx, inspectData.ID, true, true); err != nil {
				d.logger.Warn("Recovery cleanup failed", "task", handle.Config.ID, "container", inspectData.ID, "podman client", podmanTaskSocketName)
			}
			h.procState = drivers.TaskStateExited
		}
	default:
		d.logger.Warn("Recovery restart failed, unknown container state", "state", inspectData.State.Status, "container", taskState.ContainerID, "podman client", podmanTaskSocketName)
		h.procState = drivers.TaskStateUnknown
	}

	d.tasks.Set(handle.Config.ID, h)

	go h.runContainerMonitor()
	d.logger.Debug("Recovered container handle", "container", taskState.ContainerID, "podman client", podmanTaskSocketName)

	return nil
}

// BuildContainerName returns the podman container name for a given TaskConfig
func BuildContainerName(cfg *drivers.TaskConfig) string {
	return BuildContainerNameForTask(cfg.Name, cfg)
}

// BuildContainerNameForTask returns the podman container name for a specific Task in our group
func BuildContainerNameForTask(taskName string, cfg *drivers.TaskConfig) string {
	return fmt.Sprintf("%s-%s", taskName, cfg.AllocID)
}

// StartTask creates and starts a new Container based on the given TaskConfig.
func (d *Driver) StartTask(cfg *drivers.TaskConfig) (*drivers.TaskHandle, *drivers.DriverNetwork, error) {
	if _, ok := d.tasks.Get(cfg.ID); ok {
		return nil, nil, fmt.Errorf("task with ID %q already started", cfg.ID)
	}

	var podmanTaskConfig TaskConfig
	if err := cfg.DecodeDriverConfig(&podmanTaskConfig); err != nil {
		return nil, nil, fmt.Errorf("failed to decode driver config: %w", err)
	}

	handle := drivers.NewTaskHandle(taskHandleVersion)
	handle.Config = cfg

	if podmanTaskConfig.Image == "" {
		return nil, nil, fmt.Errorf("image name required")
	}

	var podmanClient *api.API
	podmanTaskSocketName := podmanTaskConfig.Socket
	if podmanTaskConfig.Socket == "" {
		podmanTaskSocketName = "default"
		podmanClient = d.defaultPodman
	} else {
		var err error
		podmanClient, err = d.getPodmanClient(podmanTaskSocketName)
		if err != nil {
			return nil, nil, fmt.Errorf("podman client with name %s not found, check your podman driver config", podmanTaskSocketName)
		}
	}
	rootless := podmanClient.IsRootless()

	createOpts := api.SpecGenerator{}
	createOpts.ContainerBasicConfig.LogConfiguration = &api.LogConfig{}
	var allArgs []string
	if podmanTaskConfig.Command != "" {
		allArgs = append(allArgs, podmanTaskConfig.Command)
	}
	allArgs = append(allArgs, podmanTaskConfig.Args...)

	// Parse entrypoint.
	switch v := podmanTaskConfig.Entrypoint.(type) {
	case string:
		// Check for a string type to maintain backwards compatibility.
		d.logger.Warn("Defining the entrypoint as a string has been deprecated, use a list of strings instead.")
		createOpts.ContainerBasicConfig.Entrypoint = append(createOpts.ContainerBasicConfig.Entrypoint, v)
	case []interface{}:
		entrypoint := make([]string, len(v))
		for i, e := range v {
			entrypoint[i] = fmt.Sprintf("%v", e)
		}
		createOpts.ContainerBasicConfig.Entrypoint = entrypoint
	case nil:
	default:
		return nil, nil, fmt.Errorf("invalid entrypoint type %T", podmanTaskConfig.Entrypoint)
	}

	containerName := BuildContainerName(cfg)

	// ensure to include port_map into tasks environment map
	cfg.Env = taskenv.SetPortMapEnvs(cfg.Env, podmanTaskConfig.PortMap)

	if len(podmanTaskConfig.Labels) > 0 {
		createOpts.ContainerBasicConfig.Labels = podmanTaskConfig.Labels
	}

	labels := make(map[string]string, len(podmanTaskConfig.Labels)+1)
	for k, v := range podmanTaskConfig.Labels {
		labels[k] = v
	}

	if len(d.config.ExtraLabels) > 0 {
		// main mandatory label
		labels[labelAllocID] = cfg.AllocID
	}

	// optional labels, as configured in plugin configuration
	for _, configurationExtraLabel := range d.config.ExtraLabels {
		if glob.Glob(configurationExtraLabel, "job_name") {
			labels[labelJobName] = cfg.JobName
		}
		if glob.Glob(configurationExtraLabel, "job_id") {
			labels[labelJobID] = cfg.JobID
		}
		if glob.Glob(configurationExtraLabel, "task_group_name") {
			labels[labelTaskGroupName] = cfg.TaskGroupName
		}
		if glob.Glob(configurationExtraLabel, "task_name") {
			labels[labelTaskName] = cfg.Name
		}
		if glob.Glob(configurationExtraLabel, "namespace") {
			labels[labelNamespace] = cfg.Namespace
		}
		if glob.Glob(configurationExtraLabel, "node_name") {
			labels[labelNodeName] = cfg.NodeName
		}
		if glob.Glob(configurationExtraLabel, "node_id") {
			labels[labelNodeID] = cfg.NodeID
		}
	}

	podmanTaskConfig.Labels = labels
	d.logger.Debug("applied labels on the container", "labels", podmanTaskConfig.Labels)

	// Basic config options
	createOpts.ContainerBasicConfig.Name = containerName
	createOpts.ContainerBasicConfig.Command = allArgs
	createOpts.ContainerBasicConfig.Env = cfg.Env
	createOpts.ContainerBasicConfig.Hostname = podmanTaskConfig.Hostname
	createOpts.ContainerBasicConfig.Sysctl = podmanTaskConfig.Sysctl
	createOpts.ContainerBasicConfig.Terminal = podmanTaskConfig.Tty
	createOpts.ContainerBasicConfig.Labels = podmanTaskConfig.Labels

	// Logging
	if !podmanTaskConfig.Logging.Empty() {
		switch podmanTaskConfig.Logging.Driver {
		case "", LOG_DRIVER_NOMAD:
			if !d.config.DisableLogCollection {
				createOpts.LogConfiguration.Driver = "k8s-file"
				createOpts.ContainerBasicConfig.LogConfiguration.Path = cfg.StdoutPath
			}
		case LOG_DRIVER_JOURNALD:
			createOpts.LogConfiguration.Driver = "journald"
		default:
			return nil, nil, fmt.Errorf("Invalid task logging.driver option")
		}
		createOpts.ContainerBasicConfig.LogConfiguration.Options = podmanTaskConfig.Logging.Options
	} else {
		d.logger.Trace("no podman log driver provided, defaulting to plugin config")
		switch d.config.Logging.Driver {
		case "", LOG_DRIVER_NOMAD:
			if !d.config.DisableLogCollection {
				createOpts.LogConfiguration.Driver = "k8s-file"
				createOpts.ContainerBasicConfig.LogConfiguration.Path = cfg.StdoutPath
			}
		case LOG_DRIVER_JOURNALD:
			createOpts.LogConfiguration.Driver = "journald"
		default:
			return nil, nil, fmt.Errorf("Invalid plugin logging.driver option")
		}
		createOpts.ContainerBasicConfig.LogConfiguration.Options = d.config.Logging.Options
	}

	// Storage config options
	createOpts.ContainerStorageConfig.Init = podmanTaskConfig.Init
	createOpts.ContainerStorageConfig.Image = podmanTaskConfig.Image
	createOpts.ContainerStorageConfig.InitPath = podmanTaskConfig.InitPath
	createOpts.ContainerStorageConfig.WorkDir = podmanTaskConfig.WorkingDir
	allMounts, err := d.containerMounts(cfg, &podmanTaskConfig)
	if err != nil {
		return nil, nil, err
	}
	createOpts.ContainerStorageConfig.Mounts = allMounts
	createOpts.ContainerStorageConfig.Devices = make([]spec.LinuxDevice, len(podmanTaskConfig.Devices))
	for idx, device := range podmanTaskConfig.Devices {
		createOpts.ContainerStorageConfig.Devices[idx] = spec.LinuxDevice{Path: device}
	}

	// Set the nomad slice as cgroup parent
	if podmanClient.IsCgroupV2() && podmanClient.GetCgroupMgr() == "systemd" {
		createOpts.ContainerCgroupConfig.CgroupParent = "nomad.slice"
	}

	// Resources config options
	createOpts.ContainerResourceConfig.ResourceLimits = &spec.LinuxResources{
		Memory: &spec.LinuxMemory{},
		CPU:    &spec.LinuxCPU{},
	}

	err = setCPUResources(podmanTaskConfig, cfg.Resources.LinuxResources, createOpts.ContainerResourceConfig.ResourceLimits.CPU)
	if err != nil {
		return nil, nil, err
	}

	hard, soft, err := memoryLimits(cfg.Resources.NomadResources.Memory, podmanTaskConfig.MemoryReservation)
	if err != nil {
		return nil, nil, err
	}
	createOpts.ContainerResourceConfig.ResourceLimits.Memory.Reservation = soft
	createOpts.ContainerResourceConfig.ResourceLimits.Memory.Limit = hard
	// set PidsLimit only if configured.
	if podmanTaskConfig.PidsLimit > 0 {
		createOpts.ContainerResourceConfig.ResourceLimits.Pids = &spec.LinuxPids{
			Limit: podmanTaskConfig.PidsLimit,
		}
	}

	if podmanTaskConfig.MemorySwap != "" {
		swap, memErr := memoryInBytes(podmanTaskConfig.MemorySwap)
		if memErr != nil {
			return nil, nil, memErr
		}
		createOpts.ContainerResourceConfig.ResourceLimits.Memory.Swap = &swap
	}
	if !podmanClient.IsCgroupV2() {
		swappiness := uint64(podmanTaskConfig.MemorySwappiness)
		createOpts.ContainerResourceConfig.ResourceLimits.Memory.Swappiness = &swappiness
	}
	// FIXME: can fail for nonRoot due to missing cpu limit delegation permissions,
	//        see https://github.com/containers/podman/blob/master/troubleshooting.md
	if !rootless {
		cpuShares := uint64(cfg.Resources.LinuxResources.CPUShares)
		createOpts.ContainerResourceConfig.ResourceLimits.CPU.Shares = &cpuShares
	}

	ulimits, ulimitErr := sliceMergeUlimit(podmanTaskConfig.Ulimit)
	if ulimitErr != nil {
		return nil, nil, fmt.Errorf("failed to parse ulimit configuration: %w", ulimitErr)
	}
	createOpts.ContainerResourceConfig.Rlimits = ulimits

	// Security config options
	createOpts.ContainerSecurityConfig.CapAdd = podmanTaskConfig.CapAdd
	createOpts.ContainerSecurityConfig.CapDrop = podmanTaskConfig.CapDrop
	createOpts.ContainerSecurityConfig.SelinuxOpts = podmanTaskConfig.SelinuxOpts
	createOpts.ContainerSecurityConfig.User = cfg.User
	createOpts.ContainerSecurityConfig.Privileged = podmanTaskConfig.Privileged
	createOpts.ContainerSecurityConfig.ReadOnlyFilesystem = podmanTaskConfig.ReadOnlyRootfs
	createOpts.ContainerSecurityConfig.ApparmorProfile = podmanTaskConfig.ApparmorProfile

	// Populate --userns mode only if configured
	if podmanTaskConfig.UserNS != "" {
		userns := strings.SplitN(podmanTaskConfig.UserNS, ":", 2)
		mode := api.NamespaceMode(userns[0])
		// Populate value only if specified
		if len(userns) > 1 {
			createOpts.ContainerSecurityConfig.UserNS = api.Namespace{NSMode: mode, Value: userns[1]}
		} else {
			createOpts.ContainerSecurityConfig.UserNS = api.Namespace{NSMode: mode}
		}
	}

	// populate shm_size if configured
	if podmanTaskConfig.ShmSize != "" {
		shmsize, memErr := memoryInBytes(podmanTaskConfig.ShmSize)
		if memErr != nil {
			return nil, nil, memErr
		}
		createOpts.ContainerStorageConfig.ShmSize = &shmsize
	}

	// Network config options
	if cfg.DNS != nil {
		for _, strdns := range cfg.DNS.Servers {
			ipdns := net.ParseIP(strdns)
			if ipdns == nil {
				return nil, nil, fmt.Errorf("Invalid dns server address, dns=%v", strdns)
			}
			createOpts.ContainerNetworkConfig.DNSServers = append(createOpts.ContainerNetworkConfig.DNSServers, ipdns)
		}
		createOpts.ContainerNetworkConfig.DNSSearch = append(createOpts.ContainerNetworkConfig.DNSSearch, cfg.DNS.Searches...)
		createOpts.ContainerNetworkConfig.DNSOptions = append(createOpts.ContainerNetworkConfig.DNSOptions, cfg.DNS.Options...)
	} else if len(d.config.DNSServers) > 0 {
		// no task DNS specific, try load default DNS from plugin config
		for _, strdns := range d.config.DNSServers {
			ipdns := net.ParseIP(strdns)
			if ipdns == nil {
				return nil, nil, fmt.Errorf("Invalid dns server address from plugin config, dns=%v", strdns)
			}
			createOpts.ContainerNetworkConfig.DNSServers = append(createOpts.ContainerNetworkConfig.DNSServers, ipdns)
		}
	}
	// Configure network
	if cfg.NetworkIsolation != nil && cfg.NetworkIsolation.Path != "" {
		createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Path
		createOpts.ContainerNetworkConfig.NetNS.Value = cfg.NetworkIsolation.Path
	} else {
		switch {
		case podmanTaskConfig.NetworkMode == "":
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
		case podmanTaskConfig.NetworkMode == "bridge":
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Bridge
		case podmanTaskConfig.NetworkMode == "host":
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Host
		case podmanTaskConfig.NetworkMode == "none":
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.NoNetwork
		case podmanTaskConfig.NetworkMode == "slirp4netns":
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Slirp
		case strings.HasPrefix(podmanTaskConfig.NetworkMode, "container:"):
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.FromContainer
			createOpts.ContainerNetworkConfig.NetNS.Value = strings.TrimPrefix(podmanTaskConfig.NetworkMode, "container:")
		case strings.HasPrefix(podmanTaskConfig.NetworkMode, "ns:"):
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.Path
			createOpts.ContainerNetworkConfig.NetNS.Value = strings.TrimPrefix(podmanTaskConfig.NetworkMode, "ns:")
		case strings.HasPrefix(podmanTaskConfig.NetworkMode, "task:"):
			otherTaskName := strings.TrimPrefix(podmanTaskConfig.NetworkMode, "task:")
			createOpts.ContainerNetworkConfig.NetNS.NSMode = api.FromContainer
			createOpts.ContainerNetworkConfig.NetNS.Value = BuildContainerNameForTask(otherTaskName, cfg)
		default:
			return nil, nil, fmt.Errorf("Unknown/Unsupported network mode: %s", podmanTaskConfig.NetworkMode)
		}
	}

	// carefully add extra hosts (--add-host)
	if extraHostsErr := setExtraHosts(podmanTaskConfig.ExtraHosts, &createOpts); extraHostsErr != nil {
		return nil, nil, extraHostsErr
	}

	portMappings, err := d.portMappings(cfg, podmanTaskConfig)
	if err != nil {
		return nil, nil, err
	}
	createOpts.ContainerNetworkConfig.PortMappings = portMappings

	containerID := ""
	recoverRunningContainer := false
	// check if there is a container with same name
	otherContainerInspect, err := podmanClient.ContainerInspect(d.ctx, containerName)
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
			if err = podmanClient.ContainerDelete(d.ctx, otherContainerInspect.ID, true, true); err != nil {
				return nil, nil, nstructs.WrapRecoverable(fmt.Sprintf("failed to remove dead container: %v", err), err)
			}
		}
	} else if !errors.Is(err, api.ContainerNotFound) {
		return nil, nil, fmt.Errorf("failed to inspect container: %s: %w", containerName, err)
	}

	if !recoverRunningContainer {
		imagePullTimeout, parseErr := time.ParseDuration(podmanTaskConfig.ImagePullTimeout)
		if parseErr != nil {
			return nil, nil, fmt.Errorf("failed to parse image_pull_timeout: %w", parseErr)
		}

		imageID, createErr := d.createImage(
			createOpts.Image,
			&podmanTaskConfig.Auth,
			podmanTaskConfig.AuthSoftFail,
			podmanTaskConfig.ForcePull,
			podmanClient,
			imagePullTimeout,
			cfg,
		)
		if createErr != nil {
			return nil, nil, fmt.Errorf("failed to create image: %s: %w", createOpts.Image, createErr)
		}
		createOpts.Image = imageID

		createResponse, createErr := podmanClient.ContainerCreate(d.ctx, createOpts)
		for _, w := range createResponse.Warnings {
			d.logger.Warn("Create Warning", "warning", w)
		}
		if createErr != nil {
			return nil, nil, fmt.Errorf("failed to start task, could not create container: %w", createErr)
		}
		containerID = createResponse.Id
	}

	cleanup := func() {
		d.logger.Debug("Cleaning up", "container", containerID)
		if cleanupErr := podmanClient.ContainerDelete(d.ctx, containerID, true, true); cleanupErr != nil {
			d.logger.Error("failed to clean up from an error in Start", "error", cleanupErr)
		}
	}

	h := &TaskHandle{
		containerID:           containerID,
		driver:                d,
		podmanClient:          podmanClient,
		taskConfig:            cfg,
		procState:             drivers.TaskStateRunning,
		exitResult:            &drivers.ExitResult{},
		startedAt:             time.Now(),
		logger:                d.logger.Named(fmt.Sprintf("podman.%s", podmanTaskSocketName)),
		logStreamer:           podmanTaskConfig.Logging.Driver == LOG_DRIVER_JOURNALD,
		logPointer:            time.Now(),
		collectionInterval:    time.Second,
		totalCPUStats:         cpustats.New(d.compute),
		userCPUStats:          cpustats.New(d.compute),
		systemCPUStats:        cpustats.New(d.compute),
		removeContainerOnExit: d.config.GC.Container,
	}

	if !recoverRunningContainer {
		if startErr := podmanClient.ContainerStart(d.ctx, containerID); startErr != nil {
			cleanup()
			return nil, nil, fmt.Errorf("failed to start task, could not start container: %w", startErr)
		}
	}

	inspectData, err := podmanClient.ContainerInspect(d.ctx, containerID)
	if err != nil {
		d.logger.Error("failed to inspect container", "error", err)
		cleanup()
		return nil, nil, fmt.Errorf("failed to start task, could not inspect container : %w", err)
	}

	driverNet := &drivers.DriverNetwork{
		PortMap:       podmanTaskConfig.PortMap,
		IP:            inspectData.NetworkSettings.IPAddress,
		AutoAdvertise: true,
	}

	driverState := TaskState{
		ContainerID: containerID,
		TaskConfig:  cfg,
		LogStreamer: h.logStreamer,
		StartedAt:   h.startedAt,
		Net:         driverNet,
	}

	if err := handle.SetDriverState(&driverState); err != nil {
		d.logger.Error("failed to start task, error setting driver state", "error", err)
		cleanup()
		return nil, nil, fmt.Errorf("failed to set driver state: %w", err)
	}

	d.tasks.Set(cfg.ID, h)

	go h.runContainerMonitor()

	d.logger.Info("Completely started container", "taskID", cfg.ID, "container", containerID, "ip", inspectData.NetworkSettings.IPAddress)

	return handle, driverNet, nil
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

func setCPUResources(cfg TaskConfig, systemResources *drivers.LinuxResources, taskCPU *spec.LinuxCPU) error {
	// always assign cpuset
	taskCPU.Cpus = systemResources.CpusetCpus

	// only set bandwidth if hard limit is enabled
	// currently only docker and podman drivers provide this ability
	if !cfg.CPUHardLimit {
		return nil
	}

	period := cfg.CPUCFSPeriod
	if period > 1000000 {
		return fmt.Errorf("invalid value for cpu_cfs_period, %d is bigger than 1000000", period)
	}
	if period == 0 {
		period = 100000 // matches cgroup default
	}
	if period < 1000 {
		period = 1000
	}

	numCores := runtime.NumCPU()
	quota := int64(systemResources.PercentTicks*float64(period)) * int64(numCores)
	if quota < 1000 {
		quota = 1000
	}

	taskCPU.Period = &period
	taskCPU.Quota = &quota

	return nil
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

func sliceMergeUlimit(ulimitsRaw map[string]string) ([]spec.POSIXRlimit, error) {
	var ulimits []spec.POSIXRlimit

	for name, ulimitRaw := range ulimitsRaw {
		if len(ulimitRaw) == 0 {
			return []spec.POSIXRlimit{}, fmt.Errorf("Malformed ulimit specification %v: %q, cannot be empty", name, ulimitRaw)
		}
		// hard limit is optional
		if !strings.Contains(ulimitRaw, ":") {
			ulimitRaw = ulimitRaw + ":" + ulimitRaw
		}

		split := strings.SplitN(ulimitRaw, ":", 2)
		if len(split) < 2 {
			return []spec.POSIXRlimit{}, fmt.Errorf("Malformed ulimit specification %v: %v", name, ulimitRaw)
		}
		soft, err := strconv.Atoi(split[0])
		if err != nil {
			return []spec.POSIXRlimit{}, fmt.Errorf("Malformed soft ulimit %v: %v", name, ulimitRaw)
		}
		hard, err := strconv.Atoi(split[1])
		if err != nil {
			return []spec.POSIXRlimit{}, fmt.Errorf("Malformed hard ulimit %v: %v", name, ulimitRaw)
		}

		ulimit := spec.POSIXRlimit{
			Type: name,
			Soft: uint64(soft),
			Hard: uint64(hard),
		}
		ulimits = append(ulimits, ulimit)
	}
	return ulimits, nil
}

// Creates the requested image if missing from storage
// returns the 64-byte image ID as an unique image identifier
func (d *Driver) createImage(
	image string,
	auth *TaskAuthConfig,
	authSoftFail bool,
	forcePull bool,
	podmanClient *api.API,
	imagePullTimeout time.Duration,
	cfg *drivers.TaskConfig,
) (string, error) {
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
			_ = d.eventer.EmitEvent(&drivers.TaskEvent{
				TaskID:    cfg.ID,
				TaskName:  cfg.Name,
				AllocID:   cfg.AllocID,
				Timestamp: time.Now(),
				Message:   "Loading image " + path,
			})
			imageName, err = podmanClient.ImageLoad(d.ctx, path)
			if err != nil {
				return imageID, fmt.Errorf("error while loading image: %w", err)
			}
		}
	}

	imageID, err := podmanClient.ImageInspectID(d.ctx, imageName)
	if err != nil && !errors.Is(err, api.ImageNotFound) {
		// If ImageInspectID errors, continue the operation and try
		// to pull the image instead
		d.logger.Warn("Unable to check for local image", "image", imageName, "error", err)
	}
	if !forcePull && imageID != "" {
		d.logger.Debug("Found imageID", imageID, "for image", imageName, "in local storage")
		return imageID, nil
	}

	d.logger.Info("Pulling image", "image", imageName)

	_ = d.eventer.EmitEvent(&drivers.TaskEvent{
		TaskID:    cfg.ID,
		TaskName:  cfg.Name,
		AllocID:   cfg.AllocID,
		Timestamp: time.Now(),
		Message:   "Pulling image " + imageName,
	})

	pc := &registry.PullConfig{
		Image:     imageName,
		TLSVerify: auth.TLSVerify,
		RegistryConfig: &registry.RegistryAuthConfig{
			Username: auth.Username,
			Password: auth.Password,
		},
		CredentialsFile:   d.config.Auth.FileConfig,
		CredentialsHelper: d.config.Auth.Helper,
		AuthSoftFail:      authSoftFail,
	}

	result, err, _ := d.pullGroup.Do(imageName, func() (interface{}, error) {

		ctx, cancel := context.WithTimeout(context.Background(), imagePullTimeout)
		defer cancel()

		if imageID, err = podmanClient.ImagePull(ctx, pc); err != nil {
			return imageID, fmt.Errorf("failed to start task, unable to pull image %s : %w", imageName, err)
		}
		return imageID, nil
	})
	if err != nil {
		return "", err
	}
	imageID = result.(string)
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
	// only start logstreamer if we have to...
	// start to stream logs if journald log driver is configured and LogCollection is not disabled
	if handle.logStreamer && !d.config.DisableLogCollection {
		go handle.runLogStreamer(ctx)
	}

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
	err := handle.podmanClient.ContainerStop(d.ctx, handle.containerID, int(timeout.Seconds()), true)
	switch {
	case err == nil:
		return nil
	case errors.Is(err, api.ContainerNotFound):
		d.logger.Debug("Container not found while we wanted to stop it", "task", taskID, "container", handle.containerID, "error", err)
		return nil
	default:
		d.logger.Error("Could not stop/kill container", "containerID", handle.containerID, "error", err)
		return err
	}
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
		err := handle.podmanClient.ContainerStop(d.ctx, handle.containerID, 60, true)
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
		err := handle.podmanClient.ContainerDelete(d.ctx, handle.containerID, true, true)
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

	return handle.podmanClient.ContainerKill(d.ctx, handle.containerID, signal)
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
	sessionId, err := handle.podmanClient.ExecCreate(ctx, handle.containerID, createRequest)
	if err != nil {
		d.logger.Error("Unable to create ExecTask session", "error", err)
		return nil, err
	}
	stdout, err := circbuf.NewBuffer(int64(drivers.CheckBufSize))
	if err != nil {
		d.logger.Error("ExecTask session failed, unable to allocate stdout buffer", "sessionId", sessionId, "error", err)
		return nil, err
	}
	stderr, err := circbuf.NewBuffer(int64(drivers.CheckBufSize))
	if err != nil {
		d.logger.Error("ExecTask session failed, unable to allocate stderr buffer", "sessionId", sessionId, "error", err)
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
	err = handle.podmanClient.ExecStart(ctx, sessionId, startRequest)
	if err != nil {
		d.logger.Error("ExecTask session returned with error", "sessionId", sessionId, "error", err)
		return nil, err
	}

	inspectData, err := handle.podmanClient.ExecInspect(ctx, sessionId)
	if err != nil {
		d.logger.Error("Unable to inspect finished ExecTask session", "sessionId", sessionId, "error", err)
		return nil, err
	}
	execResult := &drivers.ExecTaskResult{
		ExitResult: &drivers.ExitResult{
			ExitCode: inspectData.ExitCode,
		},
		Stdout: stdout.Bytes(),
		Stderr: stderr.Bytes(),
	}
	d.logger.Trace("ExecTask result", "code", execResult.ExitResult.ExitCode, "out", string(execResult.Stdout), "error", string(execResult.Stderr))

	return execResult, nil
}

// ExecTaskStreaming function is used by the Nomad client to execute commands inside the task execution context.
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

	sessionId, err := handle.podmanClient.ExecCreate(ctx, handle.containerID, createRequest)
	if err != nil {
		d.logger.Error("Unable to create exec session", "error", err)
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
	err = handle.podmanClient.ExecStart(ctx, sessionId, startRequest)
	if err != nil {
		d.logger.Error("Exec session returned with error", "sessionId", sessionId, "error", err)
		return nil, err
	}

	inspectData, err := handle.podmanClient.ExecInspect(ctx, sessionId)
	if err != nil {
		d.logger.Error("Unable to inspect finished exec session", "sessionId", sessionId, "error", err)
		return nil, err
	}
	exitResult := drivers.ExitResult{
		ExitCode: inspectData.ExitCode,
		Err:      err,
	}
	return &exitResult, nil
}

func getSElinuxVolumeLabel(vc VolumeConfig, mc *drivers.MountConfig) string {
	if mc.SELinuxLabel != vc.SelinuxLabel && mc.SELinuxLabel != "" {
		return mc.SELinuxLabel
	}

	return vc.SelinuxLabel
}

func (d *Driver) containerMounts(task *drivers.TaskConfig, driverConfig *TaskConfig) ([]spec.Mount, error) {
	var binds []spec.Mount
	binds = append(binds, spec.Mount{Source: task.TaskDir().SharedAllocDir, Destination: task.Env[taskenv.AllocDir], Type: "bind"})
	binds = append(binds, spec.Mount{Source: task.TaskDir().LocalDir, Destination: task.Env[taskenv.TaskLocalDir], Type: "bind"})
	binds = append(binds, spec.Mount{Source: task.TaskDir().SecretsDir, Destination: task.Env[taskenv.SecretsDir], Type: "bind"})

	// TODO support volume drivers
	// https://github.com/containers/libpod/pull/4548
	taskLocalBindVolume := true

	if selinuxLabel := d.config.Volumes.SelinuxLabel; selinuxLabel != "" && !driverConfig.Privileged {
		// Apply SELinux Label to each volume
		for i := range binds {
			binds[i].Options = append(binds[i].Options, selinuxLabel)
		}
	}

	for _, userbind := range driverConfig.Volumes {
		src, dst, mode, err := parseVolumeSpec(userbind)
		if err != nil {
			return nil, fmt.Errorf("invalid docker volume %q: %w", userbind, err)
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
			bind.Options = append(bind.Options, strings.Split(mode, ",")...)
		}
		binds = append(binds, bind)
	}

	// create binds for host volumes, CSI plugins, and CSI volumes
	for _, m := range task.Mounts {
		bind := spec.Mount{
			Type:        "bind",
			Destination: m.TaskPath,
			Source:      m.HostPath,
		}
		if m.Readonly {
			bind.Options = append(bind.Options, "ro")
		}

		switch m.PropagationMode {
		case nstructs.VolumeMountPropagationPrivate:
			bind.Options = append(bind.Options, "rprivate")
		case nstructs.VolumeMountPropagationHostToTask:
			bind.Options = append(bind.Options, "rslave")
		case nstructs.VolumeMountPropagationBidirectional:
			bind.Options = append(bind.Options, "rshared")
			// If PropagationMode is something else or unset, Podman defaults to rprivate
		}

		selinuxLabel := getSElinuxVolumeLabel(d.config.Volumes, m)
		if selinuxLabel != "" && !driverConfig.Privileged {
			bind.Options = append(bind.Options, selinuxLabel)
		}

		binds = append(binds, bind)
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
		var allPorts []nstructs.Port
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

func setExtraHosts(hosts []string, createOpts *api.SpecGenerator) error {
	for _, host := range hosts {
		// confirm this is a valid address:port, otherwise podman replies with
		// an indecipherable slice index out of bounds panic
		tokens := strings.SplitN(host, ":", 2)
		if len(tokens) != 2 || tokens[0] == "" || tokens[1] == "" {
			return fmt.Errorf("cannot use %q as extra_hosts (must be host:ip)", host)
		}
	}
	createOpts.ContainerNetworkConfig.HostAdd = slices.Clone(hosts)
	return nil
}
