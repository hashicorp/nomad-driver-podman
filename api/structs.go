package api

import (
	"net"
	"syscall"
	"time"

	spec "github.com/opencontainers/runtime-spec/specs-go"
)

// -------------------------------------------------------------------------------------------------------
// structs copied from https://github.com/containers/podman/blob/master/pkg/specgen/specgen.go
//
// some unused parts are modified/commented out to not pull
// more dependencies and also to overcome some json unmarshall/version problems
//
// some fields are reordert to make the linter happy (bytes maligned complains)
// -------------------------------------------------------------------------------------------------------

// LogConfig describes the logging characteristics for a container
type LogConfig struct {
	// LogDriver is the container's log driver.
	// Optional.
	Driver string `json:"driver,omitempty"`
	// LogPath is the path the container's logs will be stored at.
	// Only available if LogDriver is set to "json-file" or "k8s-file".
	// Optional.
	Path string `json:"path,omitempty"`
	// A set of options to accompany the log driver.
	// Optional.
	Options map[string]string `json:"options,omitempty"`
}

// ContainerBasicConfig contains the basic parts of a container.
type ContainerBasicConfig struct {
	// Name is the name the container will be given.
	// If no name is provided, one will be randomly generated.
	// Optional.
	Name string `json:"name,omitempty"`
	// Pod is the ID of the pod the container will join.
	// Optional.
	Pod string `json:"pod,omitempty"`
	// Entrypoint is the container's entrypoint.
	// If not given and Image is specified, this will be populated by the
	// image's configuration.
	// Optional.
	Entrypoint []string `json:"entrypoint,omitempty"`
	// Command is the container's command.
	// If not given and Image is specified, this will be populated by the
	// image's configuration.
	// Optional.
	Command []string `json:"command,omitempty"`
	// Env is a set of environment variables that will be set in the
	// container.
	// Optional.
	Env map[string]string `json:"env,omitempty"`
	// Labels are key-value pairs that are used to add metadata to
	// containers.
	// Optional.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations are key-value options passed into the container runtime
	// that can be used to trigger special behavior.
	// Optional.
	Annotations map[string]string `json:"annotations,omitempty"`
	// StopSignal is the signal that will be used to stop the container.
	// Must be a non-zero integer below SIGRTMAX.
	// If not provided, the default, SIGTERM, will be used.
	// Will conflict with Systemd if Systemd is set to "true" or "always".
	// Optional.
	StopSignal *syscall.Signal `json:"stop_signal,omitempty"`
	// LogConfiguration describes the logging for a container including
	// driver, path, and options.
	// Optional
	LogConfiguration *LogConfig `json:"log_configuration,omitempty"`
	// ConmonPidFile is a path at which a PID file for Conmon will be
	// placed.
	// If not given, a default location will be used.
	// Optional.
	ConmonPidFile string `json:"conmon_pid_file,omitempty"`
	// RestartPolicy is the container's restart policy - an action which
	// will be taken when the container exits.
	// If not given, the default policy, which does nothing, will be used.
	// Optional.
	RestartPolicy string `json:"restart_policy,omitempty"`
	// OCIRuntime is the name of the OCI runtime that will be used to create
	// the container.
	// If not specified, the default will be used.
	// Optional.
	OCIRuntime string `json:"oci_runtime,omitempty"`
	// Systemd is whether the container will be started in systemd mode.
	// Valid options are "true", "false", and "always".
	// "true" enables this mode only if the binary run in the container is
	// /sbin/init or systemd. "always" unconditionally enables systemd mode.
	// "false" unconditionally disables systemd mode.
	// If enabled, mounts and stop signal will be modified.
	// If set to "always" or set to "true" and conditionally triggered,
	// conflicts with StopSignal.
	// If not specified, "false" will be assumed.
	// Optional.
	Systemd string `json:"systemd,omitempty"`
	// Determine how to handle the NOTIFY_SOCKET - do we participate or pass it through
	// "container" - let the OCI runtime deal with it, advertise conmon's MAINPID
	// "conmon-only" - advertise conmon's MAINPID, send READY when started, don't pass to OCI
	// "ignore" - unset NOTIFY_SOCKET
	SdNotifyMode string `json:"sdnotifyMode,omitempty"`
	// Namespace is the libpod namespace the container will be placed in.
	// Optional.
	Namespace string `json:"namespace,omitempty"`

	// PidNS is the container's PID namespace.
	// It defaults to private.
	// Mandatory.
	PidNS Namespace `json:"pidns,omitempty"`

	// UtsNS is the container's UTS namespace.
	// It defaults to private.
	// Must be set to Private to set Hostname.
	// Mandatory.
	UtsNS Namespace `json:"utsns,omitempty"`

	// Hostname is the container's hostname. If not set, the hostname will
	// not be modified (if UtsNS is not private) or will be set to the
	// container ID (if UtsNS is private).
	// Conflicts with UtsNS if UtsNS is not set to private.
	// Optional.
	Hostname string `json:"hostname,omitempty"`
	// Sysctl sets kernel parameters for the container
	Sysctl map[string]string `json:"sysctl,omitempty"`
	// ContainerCreateCommand is the command that was used to create this
	// container.
	// This will be shown in the output of Inspect() on the container, and
	// may also be used by some tools that wish to recreate the container
	// (e.g. `podman generate systemd --new`).
	// Optional.
	ContainerCreateCommand []string `json:"containerCreateCommand,omitempty"`
	// Timezone is the timezone inside the container.
	// Local means it has the same timezone as the host machine
	Timezone string `json:"timezone,omitempty"`
	// PreserveFDs is a number of additional file descriptors (in addition
	// to 0, 1, 2) that will be passed to the executed process. The total FDs
	// passed will be 3 + PreserveFDs.
	// set tags as `json:"-"` for not supported remote
	PreserveFDs uint `json:"-"`
	// StopTimeout is a timeout between the container's stop signal being
	// sent and SIGKILL being sent.
	// If not provided, the default will be used.
	// If 0 is used, stop signal will not be sent, and SIGKILL will be sent
	// instead.
	// Optional.
	StopTimeout *uint `json:"stop_timeout,omitempty"`
	// RestartRetries is the number of attempts that will be made to restart
	// the container.
	// Only available when RestartPolicy is set to "on-failure".
	// Optional.
	RestartRetries *uint `json:"restart_tries,omitempty"`
	// Remove indicates if the container should be removed once it has been started
	// and exits
	Remove bool `json:"remove,omitempty"`
	// Terminal is whether the container will create a PTY.
	// Optional.
	Terminal bool `json:"terminal,omitempty"`
	// Stdin is whether the container will keep its STDIN open.
	Stdin bool `json:"stdin,omitempty"`
}

// ContainerStorageConfig contains information on the storage configuration of a
// container.
type ContainerStorageConfig struct {
	// Image is the image the container will be based on. The image will be
	// used as the container's root filesystem, and its environment vars,
	// volumes, and other configuration will be applied to the container.
	// Conflicts with Rootfs.
	// At least one of Image or Rootfs must be specified.
	Image string `json:"image"`
	// Rootfs is the path to a directory that will be used as the
	// container's root filesystem. No modification will be made to the
	// directory, it will be directly mounted into the container as root.
	// Conflicts with Image.
	// At least one of Image or Rootfs must be specified.
	Rootfs string `json:"rootfs,omitempty"`
	// ImageVolumeMode indicates how image volumes will be created.
	// Supported modes are "ignore" (do not create), "tmpfs" (create as
	// tmpfs), and "anonymous" (create as anonymous volumes).
	// The default if unset is anonymous.
	// Optional.
	ImageVolumeMode string `json:"image_volume_mode,omitempty"`
	// VolumesFrom is a set of containers whose volumes will be added to
	// this container. The name or ID of the container must be provided, and
	// may optionally be followed by a : and then one or more
	// comma-separated options. Valid options are 'ro', 'rw', and 'z'.
	// Options will be used for all volumes sourced from the container.
	VolumesFrom []string `json:"volumes_from,omitempty"`
	// Init specifies that an init binary will be mounted into the
	// container, and will be used as PID1.
	Init bool `json:"init,omitempty"`
	// InitPath specifies the path to the init binary that will be added if
	// Init is specified above. If not specified, the default set in the
	// Libpod config will be used. Ignored if Init above is not set.
	// Optional.
	InitPath string `json:"init_path,omitempty"`
	// Mounts are mounts that will be added to the container.
	// These will supersede Image Volumes and VolumesFrom volumes where
	// there are conflicts.
	// Optional.
	Mounts []spec.Mount `json:"mounts,omitempty"`
	// Volumes are named volumes that will be added to the container.
	// These will supersede Image Volumes and VolumesFrom volumes where
	// there are conflicts.
	// Optional.
	Volumes []*NamedVolume `json:"volumes,omitempty"`
	// Overlay volumes are named volumes that will be added to the container.
	// Optional.
	OverlayVolumes []*OverlayVolume `json:"overlay_volumes,omitempty"`
	// Devices are devices that will be added to the container.
	// Optional.
	Devices []spec.LinuxDevice `json:"devices,omitempty"`
	// IpcNS is the container's IPC namespace.
	// Default is private.
	// Conflicts with ShmSize if not set to private.
	// Mandatory.
	IpcNS Namespace `json:"ipcns,omitempty"`

	// ShmSize is the size of the tmpfs to mount in at /dev/shm, in bytes.
	// Conflicts with ShmSize if IpcNS is not private.
	// Optional.
	ShmSize *int64 `json:"shm_size,omitempty"`
	// WorkDir is the container's working directory.
	// If unset, the default, /, will be used.
	// Optional.
	WorkDir string `json:"work_dir,omitempty"`
	// RootfsPropagation is the rootfs propagation mode for the container.
	// If not set, the default of rslave will be used.
	// Optional.
	RootfsPropagation string `json:"rootfs_propagation,omitempty"`
}

// ContainerSecurityConfig is a container's security features, including
// SELinux, Apparmor, and Seccomp.
type ContainerSecurityConfig struct {
	// User is the user the container will be run as.
	// Can be given as a UID or a username; if a username, it will be
	// resolved within the container, using the container's /etc/passwd.
	// If unset, the container will be run as root.
	// Optional.
	User string `json:"user,omitempty"`
	// UserNS is the container's user namespace.
	// It defaults to host, indicating that no user namespace will be
	// created.
	// If set to private, IDMappings must be set.
	// Mandatory.
	UserNS Namespace `json:"userns,omitempty"`
	// Groups are a list of supplemental groups the container's user will
	// be granted access to.
	// Optional.
	Groups []string `json:"groups,omitempty"`
	// CapAdd are capabilities which will be added to the container.
	// Conflicts with Privileged.
	// Optional.
	CapAdd []string `json:"cap_add,omitempty"`
	// CapDrop are capabilities which will be removed from the container.
	// Conflicts with Privileged.
	// Optional.
	CapDrop []string `json:"cap_drop,omitempty"`
	// SelinuxProcessLabel is the process label the container will use.
	// If SELinux is enabled and this is not specified, a label will be
	// automatically generated if not specified.
	// Optional.
	SelinuxOpts []string `json:"selinux_opts,omitempty"`
	// ApparmorProfile is the name of the Apparmor profile the container
	// will use.
	// Optional.
	ApparmorProfile string `json:"apparmor_profile,omitempty"`
	// SeccompPolicy determines which seccomp profile gets applied
	// the container. valid values: empty,default,image
	SeccompPolicy string `json:"seccomp_policy,omitempty"`
	// SeccompProfilePath is the path to a JSON file containing the
	// container's Seccomp profile.
	// If not specified, no Seccomp profile will be used.
	// Optional.
	SeccompProfilePath string `json:"seccomp_profile_path,omitempty"`
	// Umask is the umask the init process of the container will be run with.
	Umask string `json:"umask,omitempty"`
	// ProcOpts are the options used for the proc mount.
	ProcOpts []string `json:"procfs_opts,omitempty"`
	// Privileged is whether the container is privileged.
	// Privileged does the following:
	// - Adds all devices on the system to the container.
	// - Adds all capabilities to the container.
	// - Disables Seccomp, SELinux, and Apparmor confinement.
	//   (Though SELinux can be manually re-enabled).
	// TODO: this conflicts with things.
	// TODO: this does more.
	Privileged bool `json:"privileged,omitempty"`
	// NoNewPrivileges is whether the container will set the no new
	// privileges flag on create, which disables gaining additional
	// privileges (e.g. via setuid) in the container.
	NoNewPrivileges bool `json:"no_new_privileges,omitempty"`

	// IDMappings are UID and GID mappings that will be used by user
	// namespaces.
	// Required if UserNS is private.
	// IDMappings *storage.IDMappingOptions `json:"idmappings,omitempty"`
	// ReadOnlyFilesystem indicates that everything will be mounted

	// as read-only
	ReadOnlyFilesystem bool `json:"read_only_filesystem,omitempty"`
}

// ContainerCgroupConfig contains configuration information about a container's
// cgroups.
type ContainerCgroupConfig struct {
	// CgroupNS is the container's cgroup namespace.
	// It defaults to private.
	// Mandatory.
	CgroupNS Namespace `json:"cgroupns,omitempty"`

	// CgroupsMode sets a policy for how cgroups will be created in the
	// container, including the ability to disable creation entirely.
	CgroupsMode string `json:"cgroups_mode,omitempty"`
	// CgroupParent is the container's CGroup parent.
	// If not set, the default for the current cgroup driver will be used.
	// Optional.
	CgroupParent string `json:"cgroup_parent,omitempty"`
}

// ContainerNetworkConfig contains information on a container's network
// configuration.
type ContainerNetworkConfig struct {
	// NetNS is the configuration to use for the container's network
	// namespace.
	// Mandatory.
	NetNS Namespace `json:"netns,omitempty"`

	// StaticIP is the a IPv4 address of the container.
	// Only available if NetNS is set to Bridge.
	// Optional.
	StaticIP *net.IP `json:"static_ip,omitempty"`
	// StaticIPv6 is a static IPv6 address to set in the container.
	// Only available if NetNS is set to Bridge.
	// Optional.
	StaticIPv6 *net.IP `json:"static_ipv6,omitempty"`
	// StaticMAC is a static MAC address to set in the container.
	// Only available if NetNS is set to bridge.
	// Optional.
	StaticMAC *net.HardwareAddr `json:"static_mac,omitempty"`
	// PortBindings is a set of ports to map into the container.
	// Only available if NetNS is set to bridge or slirp.
	// Optional.
	PortMappings []PortMapping `json:"portmappings,omitempty"`
	// Expose is a number of ports that will be forwarded to the container
	// if PublishExposedPorts is set.
	// Expose is a map of uint16 (port number) to a string representing
	// protocol. Allowed protocols are "tcp", "udp", and "sctp", or some
	// combination of the three separated by commas.
	// If protocol is set to "" we will assume TCP.
	// Only available if NetNS is set to Bridge or Slirp, and
	// PublishExposedPorts is set.
	// Optional.
	Expose map[uint16]string `json:"expose,omitempty"`
	// CNINetworks is a list of CNI networks to join the container to.
	// If this list is empty, the default CNI network will be joined
	// instead. If at least one entry is present, we will not join the
	// default network (unless it is part of this list).
	// Only available if NetNS is set to bridge.
	// Optional.
	CNINetworks []string `json:"cni_networks,omitempty"`
	// DNSServers is a set of DNS servers that will be used in the
	// container's resolv.conf, replacing the host's DNS Servers which are
	// used by default.
	// Conflicts with UseImageResolvConf.
	// Optional.
	DNSServers []net.IP `json:"dns_server,omitempty"`
	// DNSSearch is a set of DNS search domains that will be used in the
	// container's resolv.conf, replacing the host's DNS search domains
	// which are used by default.
	// Conflicts with UseImageResolvConf.
	// Optional.
	DNSSearch []string `json:"dns_search,omitempty"`
	// DNSOptions is a set of DNS options that will be used in the
	// container's resolv.conf, replacing the host's DNS options which are
	// used by default.
	// Conflicts with UseImageResolvConf.
	// Optional.
	DNSOptions []string `json:"dns_option,omitempty"`
	// HostAdd is a set of hosts which will be added to the container's
	// /etc/hosts file.
	// Conflicts with UseImageHosts.
	// Optional.
	HostAdd []string `json:"hostadd,omitempty"`
	// NetworkOptions are additional options for each network
	// Optional.
	NetworkOptions map[string][]string `json:"network_options,omitempty"`
	// PublishExposedPorts will publish ports specified in the image to
	// random unused ports (guaranteed to be above 1024) on the host.
	// This is based on ports set in Expose below, and any ports specified
	// by the Image (if one is given).
	// Only available if NetNS is set to Bridge or Slirp.
	PublishExposedPorts bool `json:"publish_image_ports,omitempty"`
	// UseImageResolvConf indicates that resolv.conf should not be managed
	// by Podman, but instead sourced from the image.
	// Conflicts with DNSServer, DNSSearch, DNSOption.
	UseImageResolvConf bool `json:"use_image_resolve_conf,omitempty"`
	// UseImageHosts indicates that /etc/hosts should not be managed by
	// Podman, and instead sourced from the image.
	// Conflicts with HostAdd.
	UseImageHosts bool `json:"use_image_hosts,omitempty"`
}

// ContainerResourceConfig contains information on container resource limits.
type ContainerResourceConfig struct {
	// ResourceLimits are resource limits to apply to the container.,
	// Can only be set as root on cgroups v1 systems, but can be set as
	// rootless as well for cgroups v2.
	// Optional.
	ResourceLimits *spec.LinuxResources `json:"resource_limits,omitempty"`
	// Rlimits are POSIX rlimits to apply to the container.
	// Optional.
	Rlimits []spec.POSIXRlimit `json:"r_limits,omitempty"`
	// OOMScoreAdj adjusts the score used by the OOM killer to determine
	// processes to kill for the container's process.
	// Optional.
	OOMScoreAdj *int `json:"oom_score_adj,omitempty"`
	// Weight per cgroup per device, can override BlkioWeight
	WeightDevice map[string]spec.LinuxWeightDevice `json:"weightDevice,omitempty"`
	// IO read rate limit per cgroup per device, bytes per second
	ThrottleReadBpsDevice map[string]spec.LinuxThrottleDevice `json:"throttleReadBpsDevice,omitempty"`
	// IO write rate limit per cgroup per device, bytes per second
	ThrottleWriteBpsDevice map[string]spec.LinuxThrottleDevice `json:"throttleWriteBpsDevice,omitempty"`
	// IO read rate limit per cgroup per device, IO per second
	ThrottleReadIOPSDevice map[string]spec.LinuxThrottleDevice `json:"throttleReadIOPSDevice,omitempty"`
	// IO write rate limit per cgroup per device, IO per second
	ThrottleWriteIOPSDevice map[string]spec.LinuxThrottleDevice `json:"throttleWriteIOPSDevice,omitempty"`
}

// ContainerHealthCheckConfig describes a container healthcheck with attributes
// like command, retries, interval, start period, and timeout.
type ContainerHealthCheckConfig struct {
	// HealthConfig *manifest.Schema2HealthConfig `json:"healthconfig,omitempty"`
}

// SpecGenerator creates an OCI spec and Libpod configuration options to create
// a container based on the given configuration.
// swagger:model SpecGenerator
type SpecGenerator struct {
	ContainerHealthCheckConfig
	ContainerBasicConfig
	ContainerStorageConfig
	ContainerNetworkConfig
	ContainerSecurityConfig
	ContainerResourceConfig
	ContainerCgroupConfig
}

// NamedVolume holds information about a named volume that will be mounted into
// the container.
type NamedVolume struct {
	// Name is the name of the named volume to be mounted. May be empty.
	// If empty, a new named volume with a pseudorandomly generated name
	// will be mounted at the given destination.
	Name string
	// Destination to mount the named volume within the container. Must be
	// an absolute path. Path will be created if it does not exist.
	Dest string
	// Options are options that the named volume will be mounted with.
	Options []string
}

// OverlayVolume holds information about a overlay volume that will be mounted into
// the container.
type OverlayVolume struct {
	// Destination is the absolute path where the mount will be placed in the container.
	Destination string `json:"destination"`
	// Source specifies the source path of the mount.
	Source string `json:"source,omitempty"`
}

// PortMapping is one or more ports that will be mapped into the container.
type PortMapping struct {
	// HostIP is the IP that we will bind to on the host.
	// If unset, assumed to be 0.0.0.0 (all interfaces).
	HostIP string `json:"host_ip,omitempty"`
	// ContainerPort is the port number that will be exposed from the
	// container.
	// Mandatory.
	ContainerPort uint16 `json:"container_port"`
	// HostPort is the port number that will be forwarded from the host into
	// the container.
	// If omitted, a random port on the host (guaranteed to be over 1024)
	// will be assigned.
	HostPort uint16 `json:"host_port,omitempty"`
	// Range is the number of ports that will be forwarded, starting at
	// HostPort and ContainerPort and counting up.
	// This is 1-indexed, so 1 is assumed to be a single port (only the
	// Hostport:Containerport mapping will be added), 2 is two ports (both
	// Hostport:Containerport and Hostport+1:Containerport+1), etc.
	// If unset, assumed to be 1 (a single port).
	// Both hostport + range and containerport + range must be less than
	// 65536.
	Range uint16 `json:"range,omitempty"`
	// Protocol is the protocol forward.
	// Must be either "tcp", "udp", and "sctp", or some combination of these
	// separated by commas.
	// If unset, assumed to be TCP.
	Protocol string `json:"protocol,omitempty"`
}

// taken from https://github.com/containers/podman/blob/master/pkg/specgen/namespaces.go

type NamespaceMode string

const (
	// Default indicates the spec generator should determine
	// a sane default
	Default NamespaceMode = "default"
	// Host means the the namespace is derived from
	// the host
	Host NamespaceMode = "host"
	// Path is the path to a namespace
	Path NamespaceMode = "path"
	// FromContainer means namespace is derived from a
	// different container
	FromContainer NamespaceMode = "container"
	// FromPod indicates the namespace is derived from a pod
	FromPod NamespaceMode = "pod"
	// Private indicates the namespace is private
	Private NamespaceMode = "private"
	// NoNetwork indicates no network namespace should
	// be joined.  loopback should still exists
	NoNetwork NamespaceMode = "none"
	// Bridge indicates that a CNI network stack
	// should be used
	Bridge NamespaceMode = "bridge"
	// Slirp indicates that a slirp4netns network stack should
	// be used
	Slirp NamespaceMode = "slirp4netns"
	// KeepID indicates a user namespace to keep the owner uid inside
	// of the namespace itself
	KeepID NamespaceMode = "keep-id"
	// Auto indicates automatic namespace mode
	Auto NamespaceMode = "auto"
	// DefaultKernelNamespaces is a comma-separated list of default kernel
	// namespaces.
	DefaultKernelNamespaces = "cgroup,ipc,net,uts"
)

// Namespace describes the namespace
type Namespace struct {
	NSMode NamespaceMode `json:"nsmode,omitempty"`
	Value  string        `json:"value,omitempty"`
}

// -------------------------------------------------------------------------------------------------------
// structs copied from https://github.com/containers/podman/blob/master/libpod/define/container_inspect.go
//
// some unused parts are modified/commented out to not pull
// more dependencies and also to overcome some json unmarshall/version problems
//
// some fields are reordert to make the linter happy (bytes maligned complains)
// -------------------------------------------------------------------------------------------------------

// InspectContainerConfig holds further data about how a container was initially
// configured.
type InspectContainerConfig struct {
	// Container hostname
	Hostname string `json:"Hostname"`
	// Container domain name - unused at present
	DomainName string `json:"Domainname"`
	// User the container was launched with
	User string `json:"User"`
	// Container environment variables
	Env []string `json:"Env"`
	// Container command
	Cmd []string `json:"Cmd"`
	// Container image
	Image string `json:"Image"`
	// Unused, at present. I've never seen this field populated.
	Volumes map[string]struct{} `json:"Volumes"`
	// Container working directory
	WorkingDir string `json:"WorkingDir"`
	// Container entrypoint
	Entrypoint string `json:"Entrypoint"`
	// On-build arguments - presently unused. More of Buildah's domain.
	OnBuild *string `json:"OnBuild"`
	// Container labels
	Labels map[string]string `json:"Labels"`
	// Container annotations
	Annotations map[string]string `json:"Annotations"`
	// Container stop signal
	StopSignal uint `json:"StopSignal"`
	// Configured healthcheck for the container
	//Healthcheck *manifest.Schema2HealthConfig `json:"Healthcheck,omitempty"`
	// CreateCommand is the full command plus arguments of the process the
	// container has been created with.
	CreateCommand []string `json:"CreateCommand,omitempty"`
	// Timezone is the timezone inside the container.
	// Local means it has the same timezone as the host machine
	Timezone string `json:"Timezone,omitempty"`
	// Umask is the umask inside the container.
	Umask string `json:"Umask,omitempty"`
	// SystemdMode is whether the container is running in systemd mode. In
	// systemd mode, the container configuration is customized to optimize
	// running systemd in the container.
	SystemdMode bool `json:"SystemdMode,omitempty"`
	// Unused, at present
	AttachStdin bool `json:"AttachStdin"`
	// Unused, at present
	AttachStdout bool `json:"AttachStdout"`
	// Unused, at present
	AttachStderr bool `json:"AttachStderr"`
	// Whether the container creates a TTY
	Tty bool `json:"Tty"`
	// Whether the container leaves STDIN open
	OpenStdin bool `json:"OpenStdin"`
	// Whether STDIN is only left open once.
	// Presently not supported by Podman, unused.
	StdinOnce bool `json:"StdinOnce"`
}

// InspectRestartPolicy holds information about the container's restart policy.
type InspectRestartPolicy struct {
	// Name contains the container's restart policy.
	// Allowable values are "no" or "" (take no action),
	// "on-failure" (restart on non-zero exit code, with an optional max
	// retry count), and "always" (always restart on container stop, unless
	// explicitly requested by API).
	// Note that this is NOT actually a name of any sort - the poor naming
	// is for Docker compatibility.
	Name string `json:"Name"`
	// MaximumRetryCount is the maximum number of retries allowed if the
	// "on-failure" restart policy is in use. Not used if "on-failure" is
	// not set.
	MaximumRetryCount uint `json:"MaximumRetryCount"`
}

// InspectLogConfig holds information about a container's configured log driver
// and is presently unused. It is retained for Docker compatibility.
type InspectLogConfig struct {
	Type   string            `json:"Type"`
	Config map[string]string `json:"Config"` //idk type, TODO
}

// InspectBlkioWeightDevice holds information about the relative weight
// of an individual device node. Weights are used in the I/O scheduler to give
// relative priority to some accesses.
type InspectBlkioWeightDevice struct {
	// Path is the path to the device this applies to.
	Path string `json:"Path"`
	// Weight is the relative weight the scheduler will use when scheduling
	// I/O.
	Weight uint16 `json:"Weight"`
}

// InspectBlkioThrottleDevice holds information about a speed cap for a device
// node. This cap applies to a specific operation (read, write, etc) on the given
// node.
type InspectBlkioThrottleDevice struct {
	// Path is the path to the device this applies to.
	Path string `json:"Path"`
	// Rate is the maximum rate. It is in either bytes per second or iops
	// per second, determined by where it is used - documentation will
	// indicate which is appropriate.
	Rate uint64 `json:"Rate"`
}

// InspectUlimit is a ulimit that will be applied to the container.
type InspectUlimit struct {
	// Name is the name (type) of the ulimit.
	Name string `json:"Name"`
	// Soft is the soft limit that will be applied.
	Soft uint64 `json:"Soft"`
	// Hard is the hard limit that will be applied.
	Hard uint64 `json:"Hard"`
}

// InspectDevice is a single device that will be mounted into the container.
type InspectDevice struct {
	// PathOnHost is the path of the device on the host.
	PathOnHost string `json:"PathOnHost"`
	// PathInContainer is the path of the device within the container.
	PathInContainer string `json:"PathInContainer"`
	// CgroupPermissions is the permissions of the mounted device.
	// Presently not populated.
	// TODO.
	CgroupPermissions string `json:"CgroupPermissions"`
}

// InspectHostPort provides information on a port on the host that a container's
// port is bound to.
type InspectHostPort struct {
	// IP on the host we are bound to. "" if not specified (binding to all
	// IPs).
	HostIP string `json:"HostIp"`
	// Port on the host we are bound to. No special formatting - just an
	// integer stuffed into a string.
	HostPort string `json:"HostPort"`
}

// InspectMount provides a record of a single mount in a container. It contains
// fields for both named and normal volumes. Only user-specified volumes will be
// included, and tmpfs volumes are not included even if the user specified them.
type InspectMount struct {
	// Whether the mount is a volume or bind mount. Allowed values are
	// "volume" and "bind".
	Type string `json:"Type"`
	// The name of the volume. Empty for bind mounts.
	Name string `json:"Name,omitempty"`
	// The source directory for the volume.
	Source string `json:"Source"`
	// The destination directory for the volume. Specified as a path within
	// the container, as it would be passed into the OCI runtime.
	Destination string `json:"Destination"`
	// The driver used for the named volume. Empty for bind mounts.
	Driver string `json:"Driver"`
	// Contains SELinux :z/:Z mount options. Unclear what, if anything, else
	// goes in here.
	Mode string `json:"Mode"`
	// All remaining mount options. Additional data, not present in the
	// original output.
	Options []string `json:"Options"`
	// Whether the volume is read-write
	RW bool `json:"RW"`
	// Mount propagation for the mount. Can be empty if not specified, but
	// is always printed - no omitempty.
	Propagation string `json:"Propagation"`
}

// InspectContainerState provides a detailed record of a container's current
// state. It is returned as part of InspectContainerData.
// As with InspectContainerData, many portions of this struct are matched to
// Docker, but here we see more fields that are unused (nonsensical in the
// context of Libpod).
type InspectContainerState struct {
	OciVersion  string             `json:"OciVersion"`
	Status      string             `json:"Status"`
	Running     bool               `json:"Running"`
	Paused      bool               `json:"Paused"`
	Restarting  bool               `json:"Restarting"` // TODO
	OOMKilled   bool               `json:"OOMKilled"`
	Dead        bool               `json:"Dead"`
	Pid         int                `json:"Pid"`
	ConmonPid   int                `json:"ConmonPid,omitempty"`
	ExitCode    int32              `json:"ExitCode"`
	Error       string             `json:"Error"` // TODO
	StartedAt   time.Time          `json:"StartedAt"`
	FinishedAt  time.Time          `json:"FinishedAt"`
	Healthcheck HealthCheckResults `json:"Healthcheck,omitempty"`
}

// HealthCheckResults describes the results/logs from a healthcheck
type HealthCheckResults struct {
	// Status healthy or unhealthy
	Status string `json:"Status"`
	// FailingStreak is the number of consecutive failed healthchecks
	FailingStreak int `json:"FailingStreak"`
	// Log describes healthcheck attempts and results
	Log []HealthCheckLog `json:"Log"`
}

// HealthCheckLog describes the results of a single healthcheck
type HealthCheckLog struct {
	// Start time as string
	Start string `json:"Start"`
	// End time as a string
	End string `json:"End"`
	// Exitcode is 0 or 1
	ExitCode int `json:"ExitCode"`
	// Output is the stdout/stderr from the healthcheck command
	Output string `json:"Output"`
}

// InspectContainerHostConfig holds information used when the container was
// created.
// It's very much a Docker-specific struct, retained (mostly) as-is for
// compatibility. We fill individual fields as best as we can, inferring as much
// as possible from the spec and container config.
// Some things cannot be inferred. These will be populated by spec annotations
// (if available).
// Field names are fixed for compatibility and cannot be changed.
// As such, silence lint warnings about them.
// nolint
type InspectContainerHostConfig struct {
	// Binds contains an array of user-added mounts.
	// Both volume mounts and named volumes are included.
	// Tmpfs mounts are NOT included.
	// In 'docker inspect' this is separated into 'Binds' and 'Mounts' based
	// on how a mount was added. We do not make this distinction and do not
	// include a Mounts field in inspect.
	// Format: <src>:<destination>[:<comma-separated options>]
	Binds []string `json:"Binds"`
	// CgroupMode is the configuration of the container's cgroup namespace.
	// Populated as follows:
	// private - a cgroup namespace has been created
	// host - No cgroup namespace created
	// container:<id> - Using another container's cgroup namespace
	// ns:<path> - A path to a cgroup namespace has been specified
	CgroupMode string `json:"CgroupMode"`
	// ContainerIDFile is a file created during container creation to hold
	// the ID of the created container.
	// This is not handled within libpod and is stored in an annotation.
	ContainerIDFile string `json:"ContainerIDFile"`
	// LogConfig contains information on the container's logging backend
	LogConfig *InspectLogConfig `json:"LogConfig"`
	// NetworkMode is the configuration of the container's network
	// namespace.
	// Populated as follows:
	// default - A network namespace is being created and configured via CNI
	// none - A network namespace is being created, not configured via CNI
	// host - No network namespace created
	// container:<id> - Using another container's network namespace
	// ns:<path> - A path to a network namespace has been specified
	NetworkMode string `json:"NetworkMode"`
	// PortBindings contains the container's port bindings.
	// It is formatted as map[string][]InspectHostPort.
	// The string key here is formatted as <integer port number>/<protocol>
	// and represents the container port. A single container port may be
	// bound to multiple host ports (on different IPs).
	PortBindings map[string][]InspectHostPort `json:"PortBindings"`
	// RestartPolicy contains the container's restart policy.
	RestartPolicy *InspectRestartPolicy `json:"RestartPolicy"`
	// AutoRemove is whether the container will be automatically removed on
	// exiting.
	// It is not handled directly within libpod and is stored in an
	// annotation.
	AutoRemove bool `json:"AutoRemove"`
	// VolumeDriver is presently unused and is retained for Docker
	// compatibility.
	VolumeDriver string `json:"VolumeDriver"`
	// VolumesFrom is a list of containers which this container uses volumes
	// from. This is not handled directly within libpod and is stored in an
	// annotation.
	// It is formatted as an array of container names and IDs.
	VolumesFrom []string `json:"VolumesFrom"`
	// CapAdd is a list of capabilities added to the container.
	// It is not directly stored by Libpod, and instead computed from the
	// capabilities listed in the container's spec, compared against a set
	// of default capabilities.
	CapAdd []string `json:"CapAdd"`
	// CapDrop is a list of capabilities removed from the container.
	// It is not directly stored by libpod, and instead computed from the
	// capabilities listed in the container's spec, compared against a set
	// of default capabilities.
	CapDrop []string `json:"CapDrop"`
	// SelinuxProcessLabel is the process label the container will use.
	// If SELinux is enabled and this is not specified, a label will be
	// automatically generated if not specified.
	// Optional.
	SelinuxOpts []string `json:"SelinuxOpts"`
	// Dns is a list of DNS nameservers that will be added to the
	// container's resolv.conf
	Dns []string `json:"Dns"`
	// DnsOptions is a list of DNS options that will be set in the
	// container's resolv.conf
	DnsOptions []string `json:"DnsOptions"`
	// DnsSearch is a list of DNS search domains that will be set in the
	// container's resolv.conf
	DnsSearch []string `json:"DnsSearch"`
	// ExtraHosts contains hosts that will be aded to the container's
	// /etc/hosts.
	ExtraHosts []string `json:"ExtraHosts"`
	// GroupAdd contains groups that the user inside the container will be
	// added to.
	GroupAdd []string `json:"GroupAdd"`
	// IpcMode represents the configuration of the container's IPC
	// namespace.
	// Populated as follows:
	// "" (empty string) - Default, an IPC namespace will be created
	// host - No IPC namespace created
	// container:<id> - Using another container's IPC namespace
	// ns:<path> - A path to an IPC namespace has been specified
	IpcMode string `json:"IpcMode"`
	// Cgroup contains the container's cgroup. It is presently not
	// populated.
	// TODO.
	Cgroup string `json:"Cgroup"`
	// Cgroups contains the container's CGroup mode.
	// Allowed values are "default" (container is creating CGroups) and
	// "disabled" (container is not creating CGroups).
	// This is Libpod-specific and not included in `docker inspect`.
	Cgroups string `json:"Cgroups"`
	// Links is unused, and provided purely for Docker compatibility.
	Links []string `json:"Links"`
	// OOMScoreAdj is an adjustment that will be made to the container's OOM
	// score.
	OomScoreAdj int `json:"OomScoreAdj"`
	// PidMode represents the configuration of the container's PID
	// namespace.
	// Populated as follows:
	// "" (empty string) - Default, a PID namespace will be created
	// host - No PID namespace created
	// container:<id> - Using another container's PID namespace
	// ns:<path> - A path to a PID namespace has been specified
	PidMode string `json:"PidMode"`
	// Privileged indicates whether the container is running with elevated
	// privileges.
	// This has a very specific meaning in the Docker sense, so it's very
	// difficult to decode from the spec and config, and so is stored as an
	// annotation.
	Privileged bool `json:"Privileged"`
	// PublishAllPorts indicates whether image ports are being published.
	// This is not directly stored in libpod and is saved as an annotation.
	PublishAllPorts bool `json:"PublishAllPorts"`
	// ReadonlyRootfs is whether the container will be mounted read-only.
	ReadonlyRootfs bool `json:"ReadonlyRootfs"`
	// SecurityOpt is a list of security-related options that are set in the
	// container.
	SecurityOpt []string `json:"SecurityOpt"`
	// Tmpfs is a list of tmpfs filesystems that will be mounted into the
	// container.
	// It is a map of destination path to options for the mount.
	Tmpfs map[string]string `json:"Tmpfs"`
	// UTSMode represents the configuration of the container's UID
	// namespace.
	// Populated as follows:
	// "" (empty string) - Default, a UTS namespace will be created
	// host - no UTS namespace created
	// container:<id> - Using another container's UTS namespace
	// ns:<path> - A path to a UTS namespace has been specified
	UTSMode string `json:"UTSMode"`
	// UsernsMode represents the configuration of the container's user
	// namespace.
	// When running rootless, a user namespace is created outside of libpod
	// to allow some privileged operations. This will not be reflected here.
	// Populated as follows:
	// "" (empty string) - No user namespace will be created
	// private - The container will be run in a user namespace
	// container:<id> - Using another container's user namespace
	// ns:<path> - A path to a user namespace has been specified
	// TODO Rootless has an additional 'keep-id' option, presently not
	// reflected here.
	UsernsMode string `json:"UsernsMode"`
	// ShmSize is the size of the container's SHM device.
	ShmSize int64 `json:"ShmSize"`
	// Runtime is provided purely for Docker compatibility.
	// It is set unconditionally to "oci" as Podman does not presently
	// support non-OCI runtimes.
	Runtime string `json:"Runtime"`
	// ConsoleSize is an array of 2 integers showing the size of the
	// container's console.
	// It is only set if the container is creating a terminal.
	// TODO.
	ConsoleSize []uint `json:"ConsoleSize"`
	// Isolation is presently unused and provided solely for Docker
	// compatibility.
	Isolation string `json:"Isolation"`
	// CpuShares indicates the CPU resources allocated to the container.
	// It is a relative weight in the scheduler for assigning CPU time
	// versus other CGroups.
	CpuShares uint64 `json:"CpuShares"`
	// Memory indicates the memory resources allocated to the container.
	// This is the limit (in bytes) of RAM the container may use.
	Memory int64 `json:"Memory"`
	// NanoCpus indicates number of CPUs allocated to the container.
	// It is an integer where one full CPU is indicated by 1000000000 (one
	// billion).
	// Thus, 2.5 CPUs (fractional portions of CPUs are allowed) would be
	// 2500000000 (2.5 billion).
	// In 'docker inspect' this is set exclusively of two further options in
	// the output (CpuPeriod and CpuQuota) which are both used to implement
	// this functionality.
	// We can't distinguish here, so if CpuQuota is set to the default of
	// 100000, we will set both CpuQuota, CpuPeriod, and NanoCpus. If
	// CpuQuota is not the default, we will not set NanoCpus.
	NanoCpus int64 `json:"NanoCpus"`
	// CgroupParent is the CGroup parent of the container.
	// Only set if not default.
	CgroupParent string `json:"CgroupParent"`
	// BlkioWeight indicates the I/O resources allocated to the container.
	// It is a relative weight in the scheduler for assigning I/O time
	// versus other CGroups.
	BlkioWeight uint16 `json:"BlkioWeight"`
	// BlkioWeightDevice is an array of I/O resource priorities for
	// individual device nodes.
	// Unfortunately, the spec only stores the device's Major/Minor numbers
	// and not the path, which is used here.
	// Fortunately, the kernel provides an interface for retrieving the path
	// of a given node by major:minor at /sys/dev/. However, the exact path
	// in use may not be what was used in the original CLI invocation -
	// though it is guaranteed that the device node will be the same, and
	// using the given path will be functionally identical.
	BlkioWeightDevice []InspectBlkioWeightDevice `json:"BlkioWeightDevice"`
	// BlkioDeviceReadBps is an array of I/O throttle parameters for
	// individual device nodes.
	// This specifically sets read rate cap in bytes per second for device
	// nodes.
	// As with BlkioWeightDevice, we pull the path from /sys/dev, and we
	// don't guarantee the path will be identical to the original (though
	// the node will be).
	BlkioDeviceReadBps []InspectBlkioThrottleDevice `json:"BlkioDeviceReadBps"`
	// BlkioDeviceWriteBps is an array of I/O throttle parameters for
	// individual device nodes.
	// this specifically sets write rate cap in bytes per second for device
	// nodes.
	// as with BlkioWeightDevice, we pull the path from /sys/dev, and we
	// don't guarantee the path will be identical to the original (though
	// the node will be).
	BlkioDeviceWriteBps []InspectBlkioThrottleDevice `json:"BlkioDeviceWriteBps"`
	// BlkioDeviceReadIOps is an array of I/O throttle parameters for
	// individual device nodes.
	// This specifically sets the read rate cap in iops per second for
	// device nodes.
	// As with BlkioWeightDevice, we pull the path from /sys/dev, and we
	// don't guarantee the path will be identical to the original (though
	// the node will be).
	BlkioDeviceReadIOps []InspectBlkioThrottleDevice `json:"BlkioDeviceReadIOps"`
	// BlkioDeviceWriteIOps is an array of I/O throttle parameters for
	// individual device nodes.
	// This specifically sets the write rate cap in iops per second for
	// device nodes.
	// As with BlkioWeightDevice, we pull the path from /sys/dev, and we
	// don't guarantee the path will be identical to the original (though
	// the node will be).
	BlkioDeviceWriteIOps []InspectBlkioThrottleDevice `json:"BlkioDeviceWriteIOps"`
	// CpuPeriod is the length of a CPU period in microseconds.
	// It relates directly to CpuQuota.
	CpuPeriod uint64 `json:"CpuPeriod"`
	// CpuPeriod is the amount of time (in microseconds) that a container
	// can use the CPU in every CpuPeriod.
	CpuQuota int64 `json:"CpuQuota"`
	// CpuRealtimePeriod is the length of time (in microseconds) of the CPU
	// realtime period. If set to 0, no time will be allocated to realtime
	// tasks.
	CpuRealtimePeriod uint64 `json:"CpuRealtimePeriod"`
	// CpuRealtimeRuntime is the length of time (in microseconds) allocated
	// for realtime tasks within every CpuRealtimePeriod.
	CpuRealtimeRuntime int64 `json:"CpuRealtimeRuntime"`
	// CpusetCpus is the is the set of CPUs that the container will execute
	// on. Formatted as `0-3` or `0,2`. Default (if unset) is all CPUs.
	CpusetCpus string `json:"CpusetCpus"`
	// CpusetMems is the set of memory nodes the container will use.
	// Formatted as `0-3` or `0,2`. Default (if unset) is all memory nodes.
	CpusetMems string `json:"CpusetMems"`
	// Devices is a list of device nodes that will be added to the
	// container.
	// These are stored in the OCI spec only as type, major, minor while we
	// display the host path. We convert this with /sys/dev, but we cannot
	// guarantee that the host path will be identical - only that the actual
	// device will be.
	Devices []InspectDevice `json:"Devices"`
	// DiskQuota is the maximum amount of disk space the container may use
	// (in bytes).
	// Presently not populated.
	// TODO.
	DiskQuota uint64 `json:"DiskQuota"`
	// KernelMemory is the maximum amount of memory the kernel will devote
	// to the container.
	KernelMemory int64 `json:"KernelMemory"`
	// MemoryReservation is the reservation (soft limit) of memory available
	// to the container. Soft limits are warnings only and can be exceeded.
	MemoryReservation int64 `json:"MemoryReservation"`
	// MemorySwap is the total limit for all memory available to the
	// container, including swap. 0 indicates that there is no limit to the
	// amount of memory available.
	MemorySwap int64 `json:"MemorySwap"`
	// MemorySwappiness is the willingness of the kernel to page container
	// memory to swap. It is an integer from 0 to 100, with low numbers
	// being more likely to be put into swap.
	// -1, the default, will not set swappiness and use the system defaults.
	MemorySwappiness int64 `json:"MemorySwappiness"`
	// OomKillDisable indicates whether the kernel OOM killer is disabled
	// for the container.
	OomKillDisable bool `json:"OomKillDisable"`
	// Init indicates whether the container has an init mounted into it.
	Init bool `json:"Init,omitempty"`
	// PidsLimit is the maximum number of PIDs what may be created within
	// the container. 0, the default, indicates no limit.
	PidsLimit int64 `json:"PidsLimit"`
	// Ulimits is a set of ulimits that will be set within the container.
	Ulimits []InspectUlimit `json:"Ulimits"`
	// CpuCount is Windows-only and not presently implemented.
	CpuCount uint64 `json:"CpuCount"`
	// CpuPercent is Windows-only and not presently implemented.
	CpuPercent uint64 `json:"CpuPercent"`
	// IOMaximumIOps is Windows-only and not presently implemented.
	IOMaximumIOps uint64 `json:"IOMaximumIOps"`
	// IOMaximumBandwidth is Windows-only and not presently implemented.
	IOMaximumBandwidth uint64 `json:"IOMaximumBandwidth"`
}

// InspectBasicNetworkConfig holds basic configuration information (e.g. IP
// addresses, MAC address, subnet masks, etc) that are common for all networks
// (both additional and main).
type InspectBasicNetworkConfig struct {
	// EndpointID is unused, maintained exclusively for compatibility.
	EndpointID string `json:"EndpointID"`
	// Gateway is the IP address of the gateway this network will use.
	Gateway string `json:"Gateway"`
	// IPAddress is the IP address for this network.
	IPAddress string `json:"IPAddress"`
	// IPPrefixLen is the length of the subnet mask of this network.
	IPPrefixLen int `json:"IPPrefixLen"`
	// SecondaryIPAddresses is a list of extra IP Addresses that the
	// container has been assigned in this network.
	SecondaryIPAddresses []string `json:"SecondaryIPAddresses,omitempty"`
	// IPv6Gateway is the IPv6 gateway this network will use.
	IPv6Gateway string `json:"IPv6Gateway"`
	// GlobalIPv6Address is the global-scope IPv6 Address for this network.
	GlobalIPv6Address string `json:"GlobalIPv6Address"`
	// GlobalIPv6PrefixLen is the length of the subnet mask of this network.
	GlobalIPv6PrefixLen int `json:"GlobalIPv6PrefixLen"`
	// SecondaryIPv6Addresses is a list of extra IPv6 Addresses that the
	// container has been assigned in this networ.
	SecondaryIPv6Addresses []string `json:"SecondaryIPv6Addresses,omitempty"`
	// MacAddress is the MAC address for the interface in this network.
	MacAddress string `json:"MacAddress"`
	// AdditionalMacAddresses is a set of additional MAC Addresses beyond
	// the first. CNI may configure more than one interface for a single
	// network, which can cause this.
	AdditionalMacAddresses []string `json:"AdditionalMACAddresses,omitempty"`
}

// InspectAdditionalNetwork holds information about non-default CNI networks the
// container has been connected to.
// As with InspectNetworkSettings, many fields are unused and maintained only
// for compatibility with Docker.
type InspectAdditionalNetwork struct {
	InspectBasicNetworkConfig

	// Name of the network we're connecting to.
	NetworkID string `json:"NetworkID,omitempty"`
	// DriverOpts is presently unused and maintained exclusively for
	// compatibility.
	DriverOpts map[string]string `json:"DriverOpts"`
	// IPAMConfig is presently unused and maintained exclusively for
	// compatibility.
	IPAMConfig map[string]string `json:"IPAMConfig"`
	// Links is presently unused and maintained exclusively for
	// compatibility.
	Links []string `json:"Links"`
}

// InspectNetworkSettings holds information about the network settings of the
// container.
// Many fields are maintained only for compatibility with `docker inspect` and
// are unused within Libpod.
type InspectNetworkSettings struct {
	InspectBasicNetworkConfig

	Bridge                 string                       `json:"Bridge"`
	SandboxID              string                       `json:"SandboxID"`
	HairpinMode            bool                         `json:"HairpinMode"`
	LinkLocalIPv6Address   string                       `json:"LinkLocalIPv6Address"`
	LinkLocalIPv6PrefixLen int                          `json:"LinkLocalIPv6PrefixLen"`
	Ports                  map[string][]InspectHostPort `json:"Ports"`
	SandboxKey             string                       `json:"SandboxKey"`
	// Networks contains information on non-default CNI networks this
	// container has joined.
	// It is a map of network name to network information.
	Networks map[string]*InspectAdditionalNetwork `json:"Networks,omitempty"`
}

// InspectContainerData provides a detailed record of a container's configuration
// and state as viewed by Libpod.
// Large portions of this structure are defined such that the output is
// compatible with `docker inspect` JSON, but additional fields have been added
// as required to share information not in the original output.
type InspectContainerData struct {
	State           *InspectContainerState      `json:"State"`
	Mounts          []InspectMount              `json:"Mounts"`
	NetworkSettings *InspectNetworkSettings     `json:"NetworkSettings"` //TODO
	Config          *InspectContainerConfig     `json:"Config"`
	HostConfig      *InspectContainerHostConfig `json:"HostConfig"`
	ID              string                      `json:"Id"`
	// FIXME can not parse date/time: "Created": "2020-07-05 11:32:38.541987006 -0400 -0400",
	//Created         time.Time              `json:"Created"`
	Path            string   `json:"Path"`
	Args            []string `json:"Args"`
	Image           string   `json:"Image"`
	ImageName       string   `json:"ImageName"`
	Rootfs          string   `json:"Rootfs"`
	Pod             string   `json:"Pod"`
	ResolvConfPath  string   `json:"ResolvConfPath"`
	HostnamePath    string   `json:"HostnamePath"`
	HostsPath       string   `json:"HostsPath"`
	StaticDir       string   `json:"StaticDir"`
	OCIConfigPath   string   `json:"OCIConfigPath,omitempty"`
	OCIRuntime      string   `json:"OCIRuntime,omitempty"`
	LogPath         string   `json:"LogPath"`
	LogTag          string   `json:"LogTag"`
	ConmonPidFile   string   `json:"ConmonPidFile"`
	Name            string   `json:"Name"`
	Driver          string   `json:"Driver"`
	MountLabel      string   `json:"MountLabel"`
	ProcessLabel    string   `json:"ProcessLabel"`
	AppArmorProfile string   `json:"AppArmorProfile"`
	EffectiveCaps   []string `json:"EffectiveCaps"`
	BoundingCaps    []string `json:"BoundingCaps"`
	ExecIDs         []string `json:"ExecIDs"`
	Dependencies    []string `json:"Dependencies"`
	ExitCommand     []string `json:"ExitCommand"`
	Namespace       string   `json:"Namespace"`
	//GraphDriver     *driver.Data                `json:"GraphDriver"`
	SizeRw       *int64 `json:"SizeRw,omitempty"`
	SizeRootFs   int64  `json:"SizeRootFs,omitempty"`
	RestartCount int32  `json:"RestartCount"`
	IsInfra      bool   `json:"IsInfra"`
}

// InspectExecSession contains information about a given exec session.
type InspectExecSession struct {
	// ProcessConfig contains information about the exec session's process.
	// ProcessConfig *InspectExecProcess `json:"ProcessConfig"`
	// ContainerID is the ID of the container this exec session is attached
	// to.
	ContainerID string `json:"ContainerID"`
	// DetachKeys are the detach keys used by the exec session.
	// If set to "" the default keys are being used.
	// Will show "<none>" if no detach keys are set.
	DetachKeys string `json:"DetachKeys"`
	// ID is the ID of the exec session.
	ID string `json:"ID"`
	// ExitCode is the exit code of the exec session. Will be set to 0 if
	// the exec session has not yet exited.
	ExitCode int `json:"ExitCode"`
	// Pid is the PID of the exec session's process.
	// Will be set to 0 if the exec session is not running.
	Pid int `json:"Pid"`
	// CanRemove is legacy and used purely for compatibility reasons.
	// Will always be set to true, unless the exec session is running.
	CanRemove bool `json:"CanRemove"`
	// OpenStderr is whether the container's STDERR stream will be attached.
	// Always set to true if the exec session created a TTY.
	OpenStderr bool `json:"OpenStderr"`
	// OpenStdin is whether the container's STDIN stream will be attached
	// to.
	OpenStdin bool `json:"OpenStdin"`
	// OpenStdout is whether the container's STDOUT stream will be attached.
	// Always set to true if the exec session created a TTY.
	OpenStdout bool `json:"OpenStdout"`
	// Running is whether the exec session is running.
	Running bool `json:"Running"`
}

// InspectExecProcess contains information about the process in a given exec
// session.
type InspectExecProcess struct {
	// Arguments are the arguments to the entrypoint command of the exec
	// session.
	Arguments []string `json:"arguments"`
	// Entrypoint is the entrypoint for the exec session (the command that
	// will be executed in the container).
	// FIXME: was string instead of []string ??
	Entrypoint []string `json:"entrypoint"`
	// Privileged is whether the exec session will be started with elevated
	// privileges.
	Privileged bool `json:"privileged"`
	// Tty is whether the exec session created a terminal.
	Tty bool `json:"tty"`
	// User is the user the exec session was started as.
	User string `json:"user"`
}

// ImagePullReport is the response from pulling one or more images.
type ImagePullReport struct {
	// Stream used to provide output from c/image
	Stream string `json:"stream,omitempty"`
	// Error contains text of errors from c/image
	Error string `json:"error,omitempty"`
	// Images contains the ID's of the images pulled
	Images []string `json:"images,omitempty"`
	// ID contains image id (retained for backwards compatibility)
	ID string `json:"id,omitempty"`
}

// -------------------------------------------------------------------------------------------------------
// structs loosly copied from https://github.com/containers/podman/blob/master/pkg/api/handlers/compat/types.go
//         and https://github.com/moby/moby/blob/master/api/types/stats.go
// some unused parts are modified/commented out to not pull
// more dependencies and also to overcome some json unmarshall/version problems
//
// some fields are reordert to make the linter happy (bytes maligned complains)
// -------------------------------------------------------------------------------------------------------
//

// CPUStats aggregates and wraps all CPU related info of container
type CPUStats struct {
	// CPU Usage. Linux and Windows.
	CPUUsage CPUUsage `json:"cpu_usage"`

	// System Usage. Linux only.
	SystemUsage uint64 `json:"system_cpu_usage,omitempty"`

	// Online CPUs. Linux only.
	OnlineCPUs uint32 `json:"online_cpus,omitempty"`

	// Usage of CPU in %. Linux only.
	CPU float64 `json:"cpu"`

	// Throttling Data. Linux only.
	ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
}

// Stats is Ultimate struct aggregating all types of stats of one container
type Stats struct {
	// Common stats
	Read    time.Time `json:"read"`
	PreRead time.Time `json:"preread"`

	// Linux specific stats, not populated on Windows.
	PidsStats  PidsStats  `json:"pids_stats,omitempty"`
	BlkioStats BlkioStats `json:"blkio_stats,omitempty"`

	// Windows specific stats, not populated on Linux.
	// NumProcs     uint32              `json:"num_procs"`
	// StorageStats docker.StorageStats `json:"storage_stats,omitempty"`

	// Shared stats
	CPUStats    CPUStats    `json:"cpu_stats,omitempty"`
	PreCPUStats CPUStats    `json:"precpu_stats,omitempty"` // "Pre"="Previous"
	MemoryStats MemoryStats `json:"memory_stats,omitempty"`
}

// ThrottlingData stores CPU throttling stats of one running container.
// Not used on Windows.
type ThrottlingData struct {
	// Number of periods with throttling active
	Periods uint64 `json:"periods"`
	// Number of periods when the container hits its throttling limit.
	ThrottledPeriods uint64 `json:"throttled_periods"`
	// Aggregate time the container was throttled for in nanoseconds.
	ThrottledTime uint64 `json:"throttled_time"`
}

// CPUUsage stores All CPU stats aggregated since container inception.
type CPUUsage struct {
	// Total CPU time consumed.
	// Units: nanoseconds (Linux)
	// Units: 100's of nanoseconds (Windows)
	TotalUsage uint64 `json:"total_usage"`

	// Total CPU time consumed per core (Linux). Not used on Windows.
	// Units: nanoseconds.
	PercpuUsage []uint64 `json:"percpu_usage,omitempty"`

	// Time spent by tasks of the cgroup in kernel mode (Linux).
	// Time spent by all container processes in kernel mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers.
	UsageInKernelmode uint64 `json:"usage_in_kernelmode"`

	// Time spent by tasks of the cgroup in user mode (Linux).
	// Time spent by all container processes in user mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers
	UsageInUsermode uint64 `json:"usage_in_usermode"`
}

// PidsStats contains the stats of a container's pids
type PidsStats struct {
	// Current is the number of pids in the cgroup
	Current uint64 `json:"current,omitempty"`
	// Limit is the hard limit on the number of pids in the cgroup.
	// A "Limit" of 0 means that there is no limit.
	Limit uint64 `json:"limit,omitempty"`
}

// BlkioStatEntry is one small entity to store a piece of Blkio stats
// Not used on Windows.
type BlkioStatEntry struct {
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Op    string `json:"op"`
	Value uint64 `json:"value"`
}

// BlkioStats stores All IO service stats for data read and write.
// This is a Linux specific structure as the differences between expressing
// block I/O on Windows and Linux are sufficiently significant to make
// little sense attempting to morph into a combined structure.
type BlkioStats struct {
	// number of bytes transferred to and from the block device
	IoServiceBytesRecursive []BlkioStatEntry `json:"io_service_bytes_recursive"`
	IoServicedRecursive     []BlkioStatEntry `json:"io_serviced_recursive"`
	IoQueuedRecursive       []BlkioStatEntry `json:"io_queue_recursive"`
	IoServiceTimeRecursive  []BlkioStatEntry `json:"io_service_time_recursive"`
	IoWaitTimeRecursive     []BlkioStatEntry `json:"io_wait_time_recursive"`
	IoMergedRecursive       []BlkioStatEntry `json:"io_merged_recursive"`
	IoTimeRecursive         []BlkioStatEntry `json:"io_time_recursive"`
	SectorsRecursive        []BlkioStatEntry `json:"sectors_recursive"`
}

// MemoryStats aggregates all memory stats since container inception on Linux.
// Windows returns stats for commit and private working set only.
type MemoryStats struct {
	// Linux Memory Stats

	// current res_counter usage for memory
	Usage uint64 `json:"usage,omitempty"`
	// maximum usage ever recorded.
	MaxUsage uint64 `json:"max_usage,omitempty"`
	// TODO(vishh): Export these as stronger types.
	// all the stats exported via memory.stat.
	Stats map[string]uint64 `json:"stats,omitempty"`
	// number of times memory usage hits limits.
	Failcnt uint64 `json:"failcnt,omitempty"`
	Limit   uint64 `json:"limit,omitempty"`

	// Windows Memory Stats
	// See https://technet.microsoft.com/en-us/magazine/ff382715.aspx

	// committed bytes
	Commit uint64 `json:"commitbytes,omitempty"`
	// peak committed bytes
	CommitPeak uint64 `json:"commitpeakbytes,omitempty"`
	// private working set
	PrivateWorkingSet uint64 `json:"privateworkingset,omitempty"`
}

// -------------------------------------------------------------------------------------------------------
// structs copied from https://github.com/containers/podman/blob/master/libpod/define/info.go
//
// some unused parts are modified/commented out to not pull more dependencies
//
// some fields are reordert to make the linter happy (bytes maligned complains)
// -------------------------------------------------------------------------------------------------------

// Info is the overall struct that describes the host system
// running libpod/podman
type Info struct {
	Host       *HostInfo              `json:"host"`
	Store      *StoreInfo             `json:"store"`
	Registries map[string]interface{} `json:"registries"`
	Version    Version                `json:"version"`
}

// HostInfo describes the libpod host
type SecurityInfo struct {
	DefaultCapabilities string `json:"capabilities"`
	AppArmorEnabled     bool   `json:"apparmorEnabled"`
	Rootless            bool   `json:"rootless"`
	SECCOMPEnabled      bool   `json:"seccompEnabled"`
	SELinuxEnabled      bool   `json:"selinuxEnabled"`
}

// HostInfo describes the libpod host
type HostInfo struct {
	Conmon       *ConmonInfo      `json:"conmon"`
	Distribution DistributionInfo `json:"distribution"`
	//IDMappings     IDMappings             `json:"idMappings,omitempty"`
	OCIRuntime     *OCIRuntimeInfo        `json:"ociRuntime"`
	RemoteSocket   *RemoteSocket          `json:"remoteSocket,omitempty"`
	Security       SecurityInfo           `json:"security"`
	Slirp4NetNS    SlirpInfo              `json:"slirp4netns,omitempty"`
	RuntimeInfo    map[string]interface{} `json:"runtimeInfo,omitempty"`
	Arch           string                 `json:"arch"`
	BuildahVersion string                 `json:"buildahVersion"`
	CgroupManager  string                 `json:"cgroupManager"`
	CGroupsVersion string                 `json:"cgroupVersion"`
	EventLogger    string                 `json:"eventLogger"`
	Hostname       string                 `json:"hostname"`
	Kernel         string                 `json:"kernel"`
	OS             string                 `json:"os"`
	Uptime         string                 `json:"uptime"`
	Linkmode       string                 `json:"linkmode"`
	MemFree        int64                  `json:"memFree"`
	MemTotal       int64                  `json:"memTotal"`
	SwapFree       int64                  `json:"swapFree"`
	SwapTotal      int64                  `json:"swapTotal"`
	CPUs           int                    `json:"cpus"`
}

// RemoteSocket describes information about the API socket
type RemoteSocket struct {
	Path   string `json:"path,omitempty"`
	Exists bool   `json:"exists,omitempty"`
}

// SlirpInfo describes the slirp executable that
// is being being used.
type SlirpInfo struct {
	Executable string `json:"executable"`
	Package    string `json:"package"`
	Version    string `json:"version"`
}

// IDMappings describe the GID and UID mappings
// type IDMappings struct {
// 	GIDMap []idtools.IDMap `json:"gidmap"`
// 	UIDMap []idtools.IDMap `json:"uidmap"`
// }

// DistributionInfo describes the host distribution
// for libpod
type DistributionInfo struct {
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
}

// ConmonInfo describes the conmon executable being used
type ConmonInfo struct {
	Package string `json:"package"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

// OCIRuntimeInfo describes the runtime (crun or runc) being
// used with podman
type OCIRuntimeInfo struct {
	Name    string `json:"name"`
	Package string `json:"package"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

// StoreInfo describes the container storage and its
// attributes
type StoreInfo struct {
	ConfigFile      string                 `json:"configFile"`
	ContainerStore  ContainerStore         `json:"containerStore"`
	GraphDriverName string                 `json:"graphDriverName"`
	GraphOptions    map[string]interface{} `json:"graphOptions"`
	GraphRoot       string                 `json:"graphRoot"`
	GraphStatus     map[string]string      `json:"graphStatus"`
	ImageStore      ImageStore             `json:"imageStore"`
	RunRoot         string                 `json:"runRoot"`
	VolumePath      string                 `json:"volumePath"`
}

// ImageStore describes the image store.  Right now only the number
// of images present
type ImageStore struct {
	Number int `json:"number"`
}

// ContainerStore describes the quantity of containers in the
// store by status
type ContainerStore struct {
	Number  int `json:"number"`
	Paused  int `json:"paused"`
	Running int `json:"running"`
	Stopped int `json:"stopped"`
}

// Version is an output struct for API
type Version struct {
	APIVersion string
	Version    string
	GoVersion  string
	GitCommit  string
	BuiltTime  string
	Built      int64
	OsArch     string
}
