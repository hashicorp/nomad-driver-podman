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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"syscall"

	spec "github.com/opencontainers/runtime-spec/specs-go"
)

// ContainerCreate creates a new container
func (c *API) ContainerCreate(ctx context.Context, create SpecGenerator) (ContainerCreateResponse, error) {

	response := ContainerCreateResponse{}

	jsonString, err := json.Marshal(create)
	if err != nil {
		return response, err
	}
	// fmt.Println(string(jsonString))

	res, err := c.Post(ctx, fmt.Sprintf("/%s/libpod/containers/create", PODMAN_API_VERSION), bytes.NewBuffer(jsonString))
	if err != nil {
		return response, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(res.Body)
		return response, fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

type ContainerCreateRequest struct {
	// Name is the name the container will be given.
	// If no name is provided, one will be randomly generated.
	// Optional.
	Name string `json:"name,omitempty"`

	// Command is the container's command.
	// If not given and Image is specified, this will be populated by the
	// image's configuration.
	// Optional.
	Command []string `json:"command,omitempty"`

	// Entrypoint is the container's entrypoint.
	// If not given and Image is specified, this will be populated by the
	// image's configuration.
	// Optional.
	Entrypoint []string `json:"entrypoint,omitempty"`

	// WorkDir is the container's working directory.
	// If unset, the default, /, will be used.
	// Optional.
	WorkDir string `json:"work_dir,omitempty"`
	// Env is a set of environment variables that will be set in the
	// container.
	// Optional.
	Env map[string]string `json:"env,omitempty"`
}

type ContainerCreateResponse struct {
	Id       string
	Warnings []string
}

// -------------------------------------------------------------------------------------------------------
// structs copied from https://github.com/containers/podman/blob/master/pkg/specgen/specgen.go
//
// some unused parts are modified/commented out to not pull
// more dependencies and also to overcome some json unmarshall/version problems
//
// some fields are reordert to make the linter happy (bytes maligned complains)
// -------------------------------------------------------------------------------------------------------

//  LogConfig describes the logging characteristics for a container
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
	// UserNS is the container's user namespace.
	// It defaults to host, indicating that no user namespace will be
	// created.
	// If set to private, IDMappings must be set.
	// Mandatory.
	// UserNS Namespace `json:"userns,omitempty"`

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
	// KeepId indicates a user namespace to keep the owner uid inside
	// of the namespace itself
	KeepID NamespaceMode = "keep-id"
	// KeepId indicates to automatically create a user namespace
	Auto NamespaceMode = "auto"
	// DefaultKernelNamespaces is a comma-separated list of default kernel
	// namespaces.
	DefaultKernelNamespaces = "cgroup,ipc,net,uts"
)

// Namespace describes the namespace
type Namespace struct {
	NSMode NamespaceMode `json:"nsmode,omitempty"`
	Value  string        `json:"string,omitempty"`
}
