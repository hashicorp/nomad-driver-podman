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

package iopodman

import (
	"time"
)

// this partly replicates https://github.com/containers/libpod/blob/master/libpod/container_inspect.go
// but we do not want to depend on (huge) libpod directly
// and varlink api lacks support for this structures
// so we resorted to unmarshal only some parts of the inspected container

type InspectContainerData struct {
	ID              string                 `json:"Id"`
	Created         time.Time              `json:"Created"`
	Path            string                 `json:"Path"`
	Args            []string               `json:"Args"`
	State           *InspectContainerState `json:"State"`
	ImageID         string                 `json:"Image"`
	ImageName       string                 `json:"ImageName"`
	Rootfs          string                 `json:"Rootfs"`
	ResolvConfPath  string                 `json:"ResolvConfPath"`
	HostnamePath    string                 `json:"HostnamePath"`
	HostsPath       string                 `json:"HostsPath"`
	StaticDir       string                 `json:"StaticDir"`
	OCIConfigPath   string                 `json:"OCIConfigPath,omitempty"`
	OCIRuntime      string                 `json:"OCIRuntime,omitempty"`
	LogPath         string                 `json:"LogPath"`
	ConmonPidFile   string                 `json:"ConmonPidFile"`
	Name            string                 `json:"Name"`
	RestartCount    int32                  `json:"RestartCount"`
	Driver          string                 `json:"Driver"`
	MountLabel      string                 `json:"MountLabel"`
	ProcessLabel    string                 `json:"ProcessLabel"`
	AppArmorProfile string                 `json:"AppArmorProfile"`
	EffectiveCaps   []string               `json:"EffectiveCaps"`
	BoundingCaps    []string               `json:"BoundingCaps"`
	ExecIDs         []string               `json:"ExecIDs"`
	// GraphDriver     *driver.Data            `json:"GraphDriver"`
	SizeRw          int64                       `json:"SizeRw,omitempty"`
	SizeRootFs      int64                       `json:"SizeRootFs,omitempty"`
	Mounts          []InspectMount              `json:"Mounts"`
	Dependencies    []string                    `json:"Dependencies"`
	NetworkSettings *InspectNetworkSettings     `json:"NetworkSettings"` //TODO
	ExitCommand     []string                    `json:"ExitCommand"`
	Namespace       string                      `json:"Namespace"`
	IsInfra         bool                        `json:"IsInfra"`
	Config          *InspectContainerConfig     `json:"Config"`
	HostConfig      *InspectContainerHostConfig `json:"HostConfig"`
}

// InspectContainerConfig holds further data about how a container was initially
// configured.
type InspectContainerConfig struct {
	// Container hostname
	Hostname string `json:"Hostname"`
	// Container domain name - unused at present
	DomainName string `json:"Domainname"`
	// User the container was launched with
	User string `json:"User"`
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
	// Healthcheck *manifest.Schema2HealthConfig `json:"Healthcheck,omitempty"`
}

type InspectImageData struct {
	ID              string                 `json:"Id"`
	Created         time.Time              `json:"Created"`
	Config          *InspectImageConfig    `json:"Config"`
}

type InspectImageConfig struct {
	Env         []string              `json:"Env"`
	Entrypoint  []string              `json:"Entrypoint"`
	Cmd         []string              `json:"Cmd"`
	WorkingDir  string                `json:"WorkingDir"`
}

// InspectMount provides a record of a single mount in a container. It contains
// fields for both named and normal volumes. Only user-specified volumes will be
// included, and tmpfs volumes are not included even if the user specified them.
type InspectMount struct {
	// Whether the mount is a volume or bind mount. Allowed values are
	// "volume" and "bind".
	Type string `json:"Type"`
	// The name of the volume. Empty for bind mounts.
	Name string `json:"Name,omptempty"`
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

// InspectContainerState provides a detailed record of a container's current
// state. It is returned as part of InspectContainerData.
// As with InspectContainerData, many portions of this struct are matched to
// Docker, but here we see more fields that are unused (nonsensical in the
// context of Libpod).
type InspectContainerState struct {
	OciVersion string    `json:"OciVersion"`
	Status     string    `json:"Status"`
	Running    bool      `json:"Running"`
	Paused     bool      `json:"Paused"`
	Restarting bool      `json:"Restarting"` // TODO
	OOMKilled  bool      `json:"OOMKilled"`
	Dead       bool      `json:"Dead"`
	Pid        int       `json:"Pid"`
	ExitCode   int32     `json:"ExitCode"`
	Error      string    `json:"Error"` // TODO
	StartedAt  time.Time `json:"StartedAt"`
	FinishedAt time.Time `json:"FinishedAt"`
	// Healthcheck HealthCheckResults `json:"Healthcheck,omitempty"`
}

// InspectNetworkSettings holds information about the network settings of the
// container.
// Many fields are maintained only for compatibility with `docker inspect` and
// are unused within Libpod.
type InspectNetworkSettings struct {
	Bridge                 string `json:"Bridge"`
	SandboxID              string `json:"SandboxID"`
	HairpinMode            bool   `json:"HairpinMode"`
	LinkLocalIPv6Address   string `json:"LinkLocalIPv6Address"`
	LinkLocalIPv6PrefixLen int    `json:"LinkLocalIPv6PrefixLen"`
	// Ports                  []ocicni.PortMapping `json:"Ports"`
	SandboxKey             string   `json:"SandboxKey"`
	SecondaryIPAddresses   []string `json:"SecondaryIPAddresses"`
	SecondaryIPv6Addresses []string `json:"SecondaryIPv6Addresses"`
	EndpointID             string   `json:"EndpointID"`
	Gateway                string   `json:"Gateway"`
	GlobalIPv6Address      string   `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen    int      `json:"GlobalIPv6PrefixLen"`
	IPAddress              string   `json:"IPAddress"`
	IPPrefixLen            int      `json:"IPPrefixLen"`
	IPv6Gateway            string   `json:"IPv6Gateway"`
	MacAddress             string   `json:"MacAddress"`
}

// InspectContainerHostConfig holds information used when the container was created.
type InspectContainerHostConfig struct {
	// Binds contains an array of user-added mounts.
	// Both volume mounts and named volumes are included.
	// Tmpfs mounts are NOT included.
	// In 'docker inspect' this is separated into 'Binds' and 'Mounts' based
	// on how a mount was added. We do not make this distinction and do not
	// include a Mounts field in inspect.
	// Format: <src>:<destination>[:<comma-separated options>]
	Binds []string `json:"Binds"`
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
	// Tmpfs is a list of tmpfs filesystems that will be mounted into the
	// container.
	// It is a map of destination path to options for the mount.
	Tmpfs map[string]string `json:"Tmpfs"`
	// Memory indicates the memory resources allocated to the container.
	// This is the limit (in bytes) of RAM the container may use.
	Memory int64 `json:"Memory"`
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
}
