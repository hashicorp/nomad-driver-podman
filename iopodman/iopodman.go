// Generated with github.com/varlink/go/cmd/varlink-go-interface-generator

// Podman Service Interface and API description.  The master version of this document can be found
// in the [API.md](https://github.com/containers/libpod/blob/master/API.md) file in the upstream libpod repository.
package iopodman

import (
	"encoding/json"
	"github.com/varlink/go/varlink"
)

// Generated type declarations

type Volume struct {
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels"`
	MountPoint string            `json:"mountPoint"`
	Driver     string            `json:"driver"`
	Options    map[string]string `json:"options"`
	Scope      string            `json:"scope"`
}

type NotImplemented struct {
	Comment string `json:"comment"`
}

type StringResponse struct {
	Message string `json:"message"`
}

type LogLine struct {
	Device       string `json:"device"`
	ParseLogType string `json:"parseLogType"`
	Time         string `json:"time"`
	Msg          string `json:"msg"`
	Cid          string `json:"cid"`
}

// ContainerChanges describes the return struct for ListContainerChanges
type ContainerChanges struct {
	Changed []string `json:"changed"`
	Added   []string `json:"added"`
	Deleted []string `json:"deleted"`
}

type ImageSaveOptions struct {
	Name       string   `json:"name"`
	Format     string   `json:"format"`
	Output     string   `json:"output"`
	OutputType string   `json:"outputType"`
	MoreTags   []string `json:"moreTags"`
	Quiet      bool     `json:"quiet"`
	Compress   bool     `json:"compress"`
}

type VolumeCreateOpts struct {
	VolumeName string            `json:"volumeName"`
	Driver     string            `json:"driver"`
	Labels     map[string]string `json:"labels"`
	Options    map[string]string `json:"options"`
}

type VolumeRemoveOpts struct {
	Volumes []string `json:"volumes"`
	All     bool     `json:"all"`
	Force   bool     `json:"force"`
}

type Image struct {
	Id          string            `json:"id"`
	Digest      string            `json:"digest"`
	ParentId    string            `json:"parentId"`
	RepoTags    []string          `json:"repoTags"`
	RepoDigests []string          `json:"repoDigests"`
	Created     string            `json:"created"`
	Size        int64             `json:"size"`
	VirtualSize int64             `json:"virtualSize"`
	Containers  int64             `json:"containers"`
	Labels      map[string]string `json:"labels"`
	IsParent    bool              `json:"isParent"`
	TopLayer    string            `json:"topLayer"`
}

// ImageHistory describes the returned structure from ImageHistory.
type ImageHistory struct {
	Id        string   `json:"id"`
	Created   string   `json:"created"`
	CreatedBy string   `json:"createdBy"`
	Tags      []string `json:"tags"`
	Size      int64    `json:"size"`
	Comment   string   `json:"comment"`
}

// Represents a single search result from SearchImages
type ImageSearchResult struct {
	Description  string `json:"description"`
	Is_official  bool   `json:"is_official"`
	Is_automated bool   `json:"is_automated"`
	Registry     string `json:"registry"`
	Name         string `json:"name"`
	Star_count   int64  `json:"star_count"`
}

type ImageSearchFilter struct {
	Is_official  *bool `json:"is_official,omitempty"`
	Is_automated *bool `json:"is_automated,omitempty"`
	Star_count   int64 `json:"star_count"`
}

type KubePodService struct {
	Pod     string `json:"pod"`
	Service string `json:"service"`
}

type Container struct {
	Id               string                  `json:"id"`
	Image            string                  `json:"image"`
	Imageid          string                  `json:"imageid"`
	Command          []string                `json:"command"`
	Createdat        string                  `json:"createdat"`
	Runningfor       string                  `json:"runningfor"`
	Status           string                  `json:"status"`
	Ports            []ContainerPortMappings `json:"ports"`
	Rootfssize       int64                   `json:"rootfssize"`
	Rwsize           int64                   `json:"rwsize"`
	Names            string                  `json:"names"`
	Labels           map[string]string       `json:"labels"`
	Mounts           []ContainerMount        `json:"mounts"`
	Containerrunning bool                    `json:"containerrunning"`
	Namespaces       ContainerNameSpace      `json:"namespaces"`
}

// ContainerStats is the return struct for the stats of a container
type ContainerStats struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Cpu          float64 `json:"cpu"`
	Cpu_nano     int64   `json:"cpu_nano"`
	System_nano  int64   `json:"system_nano"`
	Mem_usage    int64   `json:"mem_usage"`
	Mem_limit    int64   `json:"mem_limit"`
	Mem_perc     float64 `json:"mem_perc"`
	Net_input    int64   `json:"net_input"`
	Net_output   int64   `json:"net_output"`
	Block_output int64   `json:"block_output"`
	Block_input  int64   `json:"block_input"`
	Pids         int64   `json:"pids"`
}

type PsOpts struct {
	All     bool      `json:"all"`
	Filters *[]string `json:"filters,omitempty"`
	Last    *int64    `json:"last,omitempty"`
	Latest  *bool     `json:"latest,omitempty"`
	NoTrunc *bool     `json:"noTrunc,omitempty"`
	Pod     *bool     `json:"pod,omitempty"`
	Quiet   *bool     `json:"quiet,omitempty"`
	Sort    *string   `json:"sort,omitempty"`
	Sync    *bool     `json:"sync,omitempty"`
}

type PsContainer struct {
	Id         string            `json:"id"`
	Image      string            `json:"image"`
	Command    string            `json:"command"`
	Created    string            `json:"created"`
	Ports      string            `json:"ports"`
	Names      string            `json:"names"`
	IsInfra    bool              `json:"isInfra"`
	Status     string            `json:"status"`
	State      string            `json:"state"`
	PidNum     int64             `json:"pidNum"`
	RootFsSize int64             `json:"rootFsSize"`
	RwSize     int64             `json:"rwSize"`
	Pod        string            `json:"pod"`
	CreatedAt  string            `json:"createdAt"`
	ExitedAt   string            `json:"exitedAt"`
	StartedAt  string            `json:"startedAt"`
	Labels     map[string]string `json:"labels"`
	NsPid      string            `json:"nsPid"`
	Cgroup     string            `json:"cgroup"`
	Ipc        string            `json:"ipc"`
	Mnt        string            `json:"mnt"`
	Net        string            `json:"net"`
	PidNs      string            `json:"pidNs"`
	User       string            `json:"user"`
	Uts        string            `json:"uts"`
	Mounts     string            `json:"mounts"`
}

// ContainerMount describes the struct for mounts in a container
type ContainerMount struct {
	Destination string   `json:"destination"`
	Type        string   `json:"type"`
	Source      string   `json:"source"`
	Options     []string `json:"options"`
}

// ContainerPortMappings describes the struct for portmappings in an existing container
type ContainerPortMappings struct {
	Host_port      string `json:"host_port"`
	Host_ip        string `json:"host_ip"`
	Protocol       string `json:"protocol"`
	Container_port string `json:"container_port"`
}

// ContainerNamespace describes the namespace structure for an existing container
type ContainerNameSpace struct {
	User   string `json:"user"`
	Uts    string `json:"uts"`
	Pidns  string `json:"pidns"`
	Pid    string `json:"pid"`
	Cgroup string `json:"cgroup"`
	Net    string `json:"net"`
	Mnt    string `json:"mnt"`
	Ipc    string `json:"ipc"`
}

// InfoDistribution describes the host's distribution
type InfoDistribution struct {
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
}

// InfoHost describes the host stats portion of PodmanInfo
type InfoHost struct {
	Buildah_version string           `json:"buildah_version"`
	Distribution    InfoDistribution `json:"distribution"`
	Mem_free        int64            `json:"mem_free"`
	Mem_total       int64            `json:"mem_total"`
	Swap_free       int64            `json:"swap_free"`
	Swap_total      int64            `json:"swap_total"`
	Arch            string           `json:"arch"`
	Cpus            int64            `json:"cpus"`
	Hostname        string           `json:"hostname"`
	Kernel          string           `json:"kernel"`
	Os              string           `json:"os"`
	Uptime          string           `json:"uptime"`
}

// InfoGraphStatus describes the detailed status of the storage driver
type InfoGraphStatus struct {
	Backing_filesystem  string `json:"backing_filesystem"`
	Native_overlay_diff string `json:"native_overlay_diff"`
	Supports_d_type     string `json:"supports_d_type"`
}

// InfoStore describes the host's storage informatoin
type InfoStore struct {
	Containers           int64           `json:"containers"`
	Images               int64           `json:"images"`
	Graph_driver_name    string          `json:"graph_driver_name"`
	Graph_driver_options string          `json:"graph_driver_options"`
	Graph_root           string          `json:"graph_root"`
	Graph_status         InfoGraphStatus `json:"graph_status"`
	Run_root             string          `json:"run_root"`
}

// InfoPodman provides details on the podman binary
type InfoPodmanBinary struct {
	Compiler       string `json:"compiler"`
	Go_version     string `json:"go_version"`
	Podman_version string `json:"podman_version"`
	Git_commit     string `json:"git_commit"`
}

// PodmanInfo describes the Podman host and build
type PodmanInfo struct {
	Host                InfoHost         `json:"host"`
	Registries          []string         `json:"registries"`
	Insecure_registries []string         `json:"insecure_registries"`
	Store               InfoStore        `json:"store"`
	Podman              InfoPodmanBinary `json:"podman"`
}

// Sockets describes sockets location for a container
type Sockets struct {
	Container_id   string `json:"container_id"`
	Io_socket      string `json:"io_socket"`
	Control_socket string `json:"control_socket"`
}

// Create is an input structure for creating containers.
type Create struct {
	Args                   []string  `json:"args"`
	AddHost                *[]string `json:"addHost,omitempty"`
	Annotation             *[]string `json:"annotation,omitempty"`
	Attach                 *[]string `json:"attach,omitempty"`
	BlkioWeight            *string   `json:"blkioWeight,omitempty"`
	BlkioWeightDevice      *[]string `json:"blkioWeightDevice,omitempty"`
	CapAdd                 *[]string `json:"capAdd,omitempty"`
	CapDrop                *[]string `json:"capDrop,omitempty"`
	CgroupParent           *string   `json:"cgroupParent,omitempty"`
	CidFile                *string   `json:"cidFile,omitempty"`
	ConmonPidfile          *string   `json:"conmonPidfile,omitempty"`
	Command                *[]string `json:"command,omitempty"`
	CpuPeriod              *int64    `json:"cpuPeriod,omitempty"`
	CpuQuota               *int64    `json:"cpuQuota,omitempty"`
	CpuRtPeriod            *int64    `json:"cpuRtPeriod,omitempty"`
	CpuRtRuntime           *int64    `json:"cpuRtRuntime,omitempty"`
	CpuShares              *int64    `json:"cpuShares,omitempty"`
	Cpus                   *float64  `json:"cpus,omitempty"`
	CpuSetCpus             *string   `json:"cpuSetCpus,omitempty"`
	CpuSetMems             *string   `json:"cpuSetMems,omitempty"`
	Detach                 *bool     `json:"detach,omitempty"`
	DetachKeys             *string   `json:"detachKeys,omitempty"`
	Device                 *[]string `json:"device,omitempty"`
	DeviceReadBps          *[]string `json:"deviceReadBps,omitempty"`
	DeviceReadIops         *[]string `json:"deviceReadIops,omitempty"`
	DeviceWriteBps         *[]string `json:"deviceWriteBps,omitempty"`
	DeviceWriteIops        *[]string `json:"deviceWriteIops,omitempty"`
	Dns                    *[]string `json:"dns,omitempty"`
	DnsOpt                 *[]string `json:"dnsOpt,omitempty"`
	DnsSearch              *[]string `json:"dnsSearch,omitempty"`
	DnsServers             *[]string `json:"dnsServers,omitempty"`
	Entrypoint             *string   `json:"entrypoint,omitempty"`
	Env                    *[]string `json:"env,omitempty"`
	EnvFile                *[]string `json:"envFile,omitempty"`
	Expose                 *[]string `json:"expose,omitempty"`
	Gidmap                 *[]string `json:"gidmap,omitempty"`
	Groupadd               *[]string `json:"groupadd,omitempty"`
	HealthcheckCommand     *string   `json:"healthcheckCommand,omitempty"`
	HealthcheckInterval    *string   `json:"healthcheckInterval,omitempty"`
	HealthcheckRetries     *int64    `json:"healthcheckRetries,omitempty"`
	HealthcheckStartPeriod *string   `json:"healthcheckStartPeriod,omitempty"`
	HealthcheckTimeout     *string   `json:"healthcheckTimeout,omitempty"`
	Hostname               *string   `json:"hostname,omitempty"`
	ImageVolume            *string   `json:"imageVolume,omitempty"`
	Init                   *bool     `json:"init,omitempty"`
	InitPath               *string   `json:"initPath,omitempty"`
	Interactive            *bool     `json:"interactive,omitempty"`
	Ip                     *string   `json:"ip,omitempty"`
	Ipc                    *string   `json:"ipc,omitempty"`
	KernelMemory           *string   `json:"kernelMemory,omitempty"`
	Label                  *[]string `json:"label,omitempty"`
	LabelFile              *[]string `json:"labelFile,omitempty"`
	LogDriver              *string   `json:"logDriver,omitempty"`
	LogOpt                 *[]string `json:"logOpt,omitempty"`
	MacAddress             *string   `json:"macAddress,omitempty"`
	Memory                 *string   `json:"memory,omitempty"`
	MemoryReservation      *string   `json:"memoryReservation,omitempty"`
	MemorySwap             *string   `json:"memorySwap,omitempty"`
	MemorySwappiness       *int64    `json:"memorySwappiness,omitempty"`
	Name                   *string   `json:"name,omitempty"`
	Net                    *string   `json:"net,omitempty"`
	Network                *string   `json:"network,omitempty"`
	NoHosts                *bool     `json:"noHosts,omitempty"`
	OomKillDisable         *bool     `json:"oomKillDisable,omitempty"`
	OomScoreAdj            *int64    `json:"oomScoreAdj,omitempty"`
	Pid                    *string   `json:"pid,omitempty"`
	PidsLimit              *int64    `json:"pidsLimit,omitempty"`
	Pod                    *string   `json:"pod,omitempty"`
	Privileged             *bool     `json:"privileged,omitempty"`
	Publish                *[]string `json:"publish,omitempty"`
	PublishAll             *bool     `json:"publishAll,omitempty"`
	Quiet                  *bool     `json:"quiet,omitempty"`
	Readonly               *bool     `json:"readonly,omitempty"`
	Readonlytmpfs          *bool     `json:"readonlytmpfs,omitempty"`
	Restart                *string   `json:"restart,omitempty"`
	Rm                     *bool     `json:"rm,omitempty"`
	Rootfs                 *bool     `json:"rootfs,omitempty"`
	SecurityOpt            *[]string `json:"securityOpt,omitempty"`
	ShmSize                *string   `json:"shmSize,omitempty"`
	StopSignal             *string   `json:"stopSignal,omitempty"`
	StopTimeout            *int64    `json:"stopTimeout,omitempty"`
	StorageOpt             *[]string `json:"storageOpt,omitempty"`
	Subuidname             *string   `json:"subuidname,omitempty"`
	Subgidname             *string   `json:"subgidname,omitempty"`
	Sysctl                 *[]string `json:"sysctl,omitempty"`
	Systemd                *bool     `json:"systemd,omitempty"`
	Tmpfs                  *[]string `json:"tmpfs,omitempty"`
	Tty                    *bool     `json:"tty,omitempty"`
	Uidmap                 *[]string `json:"uidmap,omitempty"`
	Ulimit                 *[]string `json:"ulimit,omitempty"`
	User                   *string   `json:"user,omitempty"`
	Userns                 *string   `json:"userns,omitempty"`
	Uts                    *string   `json:"uts,omitempty"`
	Mount                  *[]string `json:"mount,omitempty"`
	Volume                 *[]string `json:"volume,omitempty"`
	VolumesFrom            *[]string `json:"volumesFrom,omitempty"`
	WorkDir                *string   `json:"workDir,omitempty"`
}

// BuildOptions are are used to describe describe physical attributes of the build
type BuildOptions struct {
	AddHosts     []string `json:"addHosts"`
	CgroupParent string   `json:"cgroupParent"`
	CpuPeriod    int64    `json:"cpuPeriod"`
	CpuQuota     int64    `json:"cpuQuota"`
	CpuShares    int64    `json:"cpuShares"`
	CpusetCpus   string   `json:"cpusetCpus"`
	CpusetMems   string   `json:"cpusetMems"`
	Memory       int64    `json:"memory"`
	MemorySwap   int64    `json:"memorySwap"`
	ShmSize      string   `json:"shmSize"`
	Ulimit       []string `json:"ulimit"`
	Volume       []string `json:"volume"`
}

// BuildInfo is used to describe user input for building images
type BuildInfo struct {
	AdditionalTags          []string          `json:"additionalTags"`
	Annotations             []string          `json:"annotations"`
	BuildArgs               map[string]string `json:"buildArgs"`
	BuildOptions            BuildOptions      `json:"buildOptions"`
	CniConfigDir            string            `json:"cniConfigDir"`
	CniPluginDir            string            `json:"cniPluginDir"`
	Compression             string            `json:"compression"`
	ContextDir              string            `json:"contextDir"`
	DefaultsMountFilePath   string            `json:"defaultsMountFilePath"`
	Dockerfiles             []string          `json:"dockerfiles"`
	Err                     string            `json:"err"`
	ForceRmIntermediateCtrs bool              `json:"forceRmIntermediateCtrs"`
	Iidfile                 string            `json:"iidfile"`
	Label                   []string          `json:"label"`
	Layers                  bool              `json:"layers"`
	Nocache                 bool              `json:"nocache"`
	Out                     string            `json:"out"`
	Output                  string            `json:"output"`
	OutputFormat            string            `json:"outputFormat"`
	PullPolicy              string            `json:"pullPolicy"`
	Quiet                   bool              `json:"quiet"`
	RemoteIntermediateCtrs  bool              `json:"remoteIntermediateCtrs"`
	ReportWriter            string            `json:"reportWriter"`
	RuntimeArgs             []string          `json:"runtimeArgs"`
	Squash                  bool              `json:"squash"`
}

// MoreResponse is a struct for when responses from varlink requires longer output
type MoreResponse struct {
	Logs []string `json:"logs"`
	Id   string   `json:"id"`
}

// ListPodContainerInfo is a returned struct for describing containers
// in a pod.
type ListPodContainerInfo struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Status string `json:"status"`
}

// PodCreate is an input structure for creating pods.
// It emulates options to podman pod create. The infraCommand and
// infraImage options are currently NotSupported.
type PodCreate struct {
	Name         string            `json:"name"`
	CgroupParent string            `json:"cgroupParent"`
	Labels       map[string]string `json:"labels"`
	Share        []string          `json:"share"`
	Infra        bool              `json:"infra"`
	InfraCommand string            `json:"infraCommand"`
	InfraImage   string            `json:"infraImage"`
	Publish      []string          `json:"publish"`
}

// ListPodData is the returned struct for an individual pod
type ListPodData struct {
	Id                 string                 `json:"id"`
	Name               string                 `json:"name"`
	Createdat          string                 `json:"createdat"`
	Cgroup             string                 `json:"cgroup"`
	Status             string                 `json:"status"`
	Labels             map[string]string      `json:"labels"`
	Numberofcontainers string                 `json:"numberofcontainers"`
	Containersinfo     []ListPodContainerInfo `json:"containersinfo"`
}

type PodContainerErrorData struct {
	Containerid string `json:"containerid"`
	Reason      string `json:"reason"`
}

// Runlabel describes the required input for container runlabel
type Runlabel struct {
	Image     string            `json:"image"`
	Authfile  string            `json:"authfile"`
	Display   bool              `json:"display"`
	Name      string            `json:"name"`
	Pull      bool              `json:"pull"`
	Label     string            `json:"label"`
	ExtraArgs []string          `json:"extraArgs"`
	Opts      map[string]string `json:"opts"`
}

// Event describes a libpod struct
type Event struct {
	Id     string `json:"id"`
	Image  string `json:"image"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Time   string `json:"time"`
	Type   string `json:"type"`
}

type DiffInfo struct {
	Path       string `json:"path"`
	ChangeType string `json:"changeType"`
}

// ImageNotFound means the image could not be found by the provided name or ID in local storage.
type ImageNotFound struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

func (e ImageNotFound) Error() string {
	return "io.podman.ImageNotFound"
}

// ContainerNotFound means the container could not be found by the provided name or ID in local storage.
type ContainerNotFound struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

func (e ContainerNotFound) Error() string {
	return "io.podman.ContainerNotFound"
}

// NoContainerRunning means none of the containers requested are running in a command that requires a running container.
type NoContainerRunning struct{}

func (e NoContainerRunning) Error() string {
	return "io.podman.NoContainerRunning"
}

// PodNotFound means the pod could not be found by the provided name or ID in local storage.
type PodNotFound struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func (e PodNotFound) Error() string {
	return "io.podman.PodNotFound"
}

// VolumeNotFound means the volume could not be found by the name or ID in local storage.
type VolumeNotFound struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

func (e VolumeNotFound) Error() string {
	return "io.podman.VolumeNotFound"
}

// PodContainerError means a container associated with a pod failed to perform an operation. It contains
// a container ID of the container that failed.
type PodContainerError struct {
	Podname string                  `json:"podname"`
	Errors  []PodContainerErrorData `json:"errors"`
}

func (e PodContainerError) Error() string {
	return "io.podman.PodContainerError"
}

// NoContainersInPod means a pod has no containers on which to perform the operation. It contains
// the pod ID.
type NoContainersInPod struct {
	Name string `json:"name"`
}

func (e NoContainersInPod) Error() string {
	return "io.podman.NoContainersInPod"
}

// InvalidState indicates that a container or pod was in an improper state for the requested operation
type InvalidState struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

func (e InvalidState) Error() string {
	return "io.podman.InvalidState"
}

// ErrorOccurred is a generic error for an error that occurs during the execution.  The actual error message
// is includes as part of the error's text.
type ErrorOccurred struct {
	Reason string `json:"reason"`
}

func (e ErrorOccurred) Error() string {
	return "io.podman.ErrorOccurred"
}

// RuntimeErrors generally means a runtime could not be found or gotten.
type RuntimeError struct {
	Reason string `json:"reason"`
}

func (e RuntimeError) Error() string {
	return "io.podman.RuntimeError"
}

// The Podman endpoint requires that you use a streaming connection.
type WantsMoreRequired struct {
	Reason string `json:"reason"`
}

func (e WantsMoreRequired) Error() string {
	return "io.podman.WantsMoreRequired"
}

// Container is already stopped
type ErrCtrStopped struct {
	Id string `json:"id"`
}

func (e ErrCtrStopped) Error() string {
	return "io.podman.ErrCtrStopped"
}

func Dispatch_Error(err error) error {
	if e, ok := err.(*varlink.Error); ok {
		switch e.Name {
		case "io.podman.ImageNotFound":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param ImageNotFound
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.ContainerNotFound":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param ContainerNotFound
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.NoContainerRunning":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param NoContainerRunning
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.PodNotFound":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param PodNotFound
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.VolumeNotFound":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param VolumeNotFound
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.PodContainerError":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param PodContainerError
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.NoContainersInPod":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param NoContainersInPod
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.InvalidState":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param InvalidState
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.ErrorOccurred":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param ErrorOccurred
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.RuntimeError":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param RuntimeError
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.WantsMoreRequired":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param WantsMoreRequired
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		case "io.podman.ErrCtrStopped":
			errorRawParameters := e.Parameters.(*json.RawMessage)
			if errorRawParameters == nil {
				return e
			}
			var param ErrCtrStopped
			err := json.Unmarshal(*errorRawParameters, &param)
			if err != nil {
				return e
			}
			return &param
		}
	}
	return err
}

// Generated client method calls

// GetVersion returns version and build information of the podman service
type GetVersion_methods struct{}

func GetVersion() GetVersion_methods { return GetVersion_methods{} }

func (m GetVersion_methods) Call(c *varlink.Connection) (version_out_ string, go_version_out_ string, git_commit_out_ string, built_out_ string, os_arch_out_ string, remote_api_version_out_ int64, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	version_out_, go_version_out_, git_commit_out_, built_out_, os_arch_out_, remote_api_version_out_, _, err_ = receive()
	return
}

func (m GetVersion_methods) Send(c *varlink.Connection, flags uint64) (func() (string, string, string, string, string, int64, uint64, error), error) {
	receive, err := c.Send("io.podman.GetVersion", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (version_out_ string, go_version_out_ string, git_commit_out_ string, built_out_ string, os_arch_out_ string, remote_api_version_out_ int64, flags uint64, err error) {
		var out struct {
			Version            string `json:"version"`
			Go_version         string `json:"go_version"`
			Git_commit         string `json:"git_commit"`
			Built              string `json:"built"`
			Os_arch            string `json:"os_arch"`
			Remote_api_version int64  `json:"remote_api_version"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		version_out_ = out.Version
		go_version_out_ = out.Go_version
		git_commit_out_ = out.Git_commit
		built_out_ = out.Built
		os_arch_out_ = out.Os_arch
		remote_api_version_out_ = out.Remote_api_version
		return
	}, nil
}

// GetInfo returns a [PodmanInfo](#PodmanInfo) struct that describes podman and its host such as storage stats,
// build information of Podman, and system-wide registries.
type GetInfo_methods struct{}

func GetInfo() GetInfo_methods { return GetInfo_methods{} }

func (m GetInfo_methods) Call(c *varlink.Connection) (info_out_ PodmanInfo, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	info_out_, _, err_ = receive()
	return
}

func (m GetInfo_methods) Send(c *varlink.Connection, flags uint64) (func() (PodmanInfo, uint64, error), error) {
	receive, err := c.Send("io.podman.GetInfo", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (info_out_ PodmanInfo, flags uint64, err error) {
		var out struct {
			Info PodmanInfo `json:"info"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		info_out_ = out.Info
		return
	}, nil
}

// ListContainers returns information about all containers.
// See also [GetContainer](#GetContainer).
type ListContainers_methods struct{}

func ListContainers() ListContainers_methods { return ListContainers_methods{} }

func (m ListContainers_methods) Call(c *varlink.Connection) (containers_out_ []Container, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	containers_out_, _, err_ = receive()
	return
}

func (m ListContainers_methods) Send(c *varlink.Connection, flags uint64) (func() ([]Container, uint64, error), error) {
	receive, err := c.Send("io.podman.ListContainers", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (containers_out_ []Container, flags uint64, err error) {
		var out struct {
			Containers []Container `json:"containers"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		containers_out_ = []Container(out.Containers)
		return
	}, nil
}

type Ps_methods struct{}

func Ps() Ps_methods { return Ps_methods{} }

func (m Ps_methods) Call(c *varlink.Connection, opts_in_ PsOpts) (containers_out_ []PsContainer, err_ error) {
	receive, err_ := m.Send(c, 0, opts_in_)
	if err_ != nil {
		return
	}
	containers_out_, _, err_ = receive()
	return
}

func (m Ps_methods) Send(c *varlink.Connection, flags uint64, opts_in_ PsOpts) (func() ([]PsContainer, uint64, error), error) {
	var in struct {
		Opts PsOpts `json:"opts"`
	}
	in.Opts = opts_in_
	receive, err := c.Send("io.podman.Ps", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (containers_out_ []PsContainer, flags uint64, err error) {
		var out struct {
			Containers []PsContainer `json:"containers"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		containers_out_ = []PsContainer(out.Containers)
		return
	}, nil
}

type GetContainersByStatus_methods struct{}

func GetContainersByStatus() GetContainersByStatus_methods { return GetContainersByStatus_methods{} }

func (m GetContainersByStatus_methods) Call(c *varlink.Connection, status_in_ []string) (containerS_out_ []Container, err_ error) {
	receive, err_ := m.Send(c, 0, status_in_)
	if err_ != nil {
		return
	}
	containerS_out_, _, err_ = receive()
	return
}

func (m GetContainersByStatus_methods) Send(c *varlink.Connection, flags uint64, status_in_ []string) (func() ([]Container, uint64, error), error) {
	var in struct {
		Status []string `json:"status"`
	}
	in.Status = []string(status_in_)
	receive, err := c.Send("io.podman.GetContainersByStatus", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (containerS_out_ []Container, flags uint64, err error) {
		var out struct {
			ContainerS []Container `json:"containerS"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		containerS_out_ = []Container(out.ContainerS)
		return
	}, nil
}

type Top_methods struct{}

func Top() Top_methods { return Top_methods{} }

func (m Top_methods) Call(c *varlink.Connection, nameOrID_in_ string, descriptors_in_ []string) (top_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, nameOrID_in_, descriptors_in_)
	if err_ != nil {
		return
	}
	top_out_, _, err_ = receive()
	return
}

func (m Top_methods) Send(c *varlink.Connection, flags uint64, nameOrID_in_ string, descriptors_in_ []string) (func() ([]string, uint64, error), error) {
	var in struct {
		NameOrID    string   `json:"nameOrID"`
		Descriptors []string `json:"descriptors"`
	}
	in.NameOrID = nameOrID_in_
	in.Descriptors = []string(descriptors_in_)
	receive, err := c.Send("io.podman.Top", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (top_out_ []string, flags uint64, err error) {
		var out struct {
			Top []string `json:"top"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		top_out_ = []string(out.Top)
		return
	}, nil
}

// GetContainer returns information about a single container.  If a container
// with the given id doesn't exist, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.  See also [ListContainers](ListContainers) and
// [InspectContainer](#InspectContainer).
type GetContainer_methods struct{}

func GetContainer() GetContainer_methods { return GetContainer_methods{} }

func (m GetContainer_methods) Call(c *varlink.Connection, id_in_ string) (container_out_ Container, err_ error) {
	receive, err_ := m.Send(c, 0, id_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m GetContainer_methods) Send(c *varlink.Connection, flags uint64, id_in_ string) (func() (Container, uint64, error), error) {
	var in struct {
		Id string `json:"id"`
	}
	in.Id = id_in_
	receive, err := c.Send("io.podman.GetContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ Container, flags uint64, err error) {
		var out struct {
			Container Container `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// GetContainersByContext allows you to get a list of container ids depending on all, latest, or a list of
// container names.  The definition of latest container means the latest by creation date.  In a multi-
// user environment, results might differ from what you expect.
type GetContainersByContext_methods struct{}

func GetContainersByContext() GetContainersByContext_methods { return GetContainersByContext_methods{} }

func (m GetContainersByContext_methods) Call(c *varlink.Connection, all_in_ bool, latest_in_ bool, args_in_ []string) (containers_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, all_in_, latest_in_, args_in_)
	if err_ != nil {
		return
	}
	containers_out_, _, err_ = receive()
	return
}

func (m GetContainersByContext_methods) Send(c *varlink.Connection, flags uint64, all_in_ bool, latest_in_ bool, args_in_ []string) (func() ([]string, uint64, error), error) {
	var in struct {
		All    bool     `json:"all"`
		Latest bool     `json:"latest"`
		Args   []string `json:"args"`
	}
	in.All = all_in_
	in.Latest = latest_in_
	in.Args = []string(args_in_)
	receive, err := c.Send("io.podman.GetContainersByContext", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (containers_out_ []string, flags uint64, err error) {
		var out struct {
			Containers []string `json:"containers"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		containers_out_ = []string(out.Containers)
		return
	}, nil
}

// CreateContainer creates a new container from an image.  It uses a [Create](#Create) type for input.
type CreateContainer_methods struct{}

func CreateContainer() CreateContainer_methods { return CreateContainer_methods{} }

func (m CreateContainer_methods) Call(c *varlink.Connection, create_in_ Create) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, create_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m CreateContainer_methods) Send(c *varlink.Connection, flags uint64, create_in_ Create) (func() (string, uint64, error), error) {
	var in struct {
		Create Create `json:"create"`
	}
	in.Create = create_in_
	receive, err := c.Send("io.podman.CreateContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// InspectContainer data takes a name or ID of a container returns the inspection
// data in string format.  You can then serialize the string into JSON.  A [ContainerNotFound](#ContainerNotFound)
// error will be returned if the container cannot be found. See also [InspectImage](#InspectImage).
type InspectContainer_methods struct{}

func InspectContainer() InspectContainer_methods { return InspectContainer_methods{} }

func (m InspectContainer_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m InspectContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.InspectContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// ListContainerProcesses takes a name or ID of a container and returns the processes
// running inside the container as array of strings.  It will accept an array of string
// arguments that represent ps options.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ListContainerProcesses '{"name": "135d71b9495f", "opts": []}'
// {
//   "container": [
//     "  UID   PID  PPID  C STIME TTY          TIME CMD",
//     "    0 21220 21210  0 09:05 pts/0    00:00:00 /bin/sh",
//     "    0 21232 21220  0 09:05 pts/0    00:00:00 top",
//     "    0 21284 21220  0 09:05 pts/0    00:00:00 vi /etc/hosts"
//   ]
// }
// ~~~
type ListContainerProcesses_methods struct{}

func ListContainerProcesses() ListContainerProcesses_methods { return ListContainerProcesses_methods{} }

func (m ListContainerProcesses_methods) Call(c *varlink.Connection, name_in_ string, opts_in_ []string) (container_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, opts_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m ListContainerProcesses_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, opts_in_ []string) (func() ([]string, uint64, error), error) {
	var in struct {
		Name string   `json:"name"`
		Opts []string `json:"opts"`
	}
	in.Name = name_in_
	in.Opts = []string(opts_in_)
	receive, err := c.Send("io.podman.ListContainerProcesses", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ []string, flags uint64, err error) {
		var out struct {
			Container []string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = []string(out.Container)
		return
	}, nil
}

// GetContainerLogs takes a name or ID of a container and returns the logs of that container.
// If the container cannot be found, a [ContainerNotFound](#ContainerNotFound) error will be returned.
// The container logs are returned as an array of strings.  GetContainerLogs will honor the streaming
// capability of varlink if the client invokes it.
type GetContainerLogs_methods struct{}

func GetContainerLogs() GetContainerLogs_methods { return GetContainerLogs_methods{} }

func (m GetContainerLogs_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m GetContainerLogs_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() ([]string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.GetContainerLogs", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ []string, flags uint64, err error) {
		var out struct {
			Container []string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = []string(out.Container)
		return
	}, nil
}

type GetContainersLogs_methods struct{}

func GetContainersLogs() GetContainersLogs_methods { return GetContainersLogs_methods{} }

func (m GetContainersLogs_methods) Call(c *varlink.Connection, names_in_ []string, follow_in_ bool, latest_in_ bool, since_in_ string, tail_in_ int64, timestamps_in_ bool) (log_out_ LogLine, err_ error) {
	receive, err_ := m.Send(c, 0, names_in_, follow_in_, latest_in_, since_in_, tail_in_, timestamps_in_)
	if err_ != nil {
		return
	}
	log_out_, _, err_ = receive()
	return
}

func (m GetContainersLogs_methods) Send(c *varlink.Connection, flags uint64, names_in_ []string, follow_in_ bool, latest_in_ bool, since_in_ string, tail_in_ int64, timestamps_in_ bool) (func() (LogLine, uint64, error), error) {
	var in struct {
		Names      []string `json:"names"`
		Follow     bool     `json:"follow"`
		Latest     bool     `json:"latest"`
		Since      string   `json:"since"`
		Tail       int64    `json:"tail"`
		Timestamps bool     `json:"timestamps"`
	}
	in.Names = []string(names_in_)
	in.Follow = follow_in_
	in.Latest = latest_in_
	in.Since = since_in_
	in.Tail = tail_in_
	in.Timestamps = timestamps_in_
	receive, err := c.Send("io.podman.GetContainersLogs", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (log_out_ LogLine, flags uint64, err error) {
		var out struct {
			Log LogLine `json:"log"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		log_out_ = out.Log
		return
	}, nil
}

// ListContainerChanges takes a name or ID of a container and returns changes between the container and
// its base image. It returns a struct of changed, deleted, and added path names.
type ListContainerChanges_methods struct{}

func ListContainerChanges() ListContainerChanges_methods { return ListContainerChanges_methods{} }

func (m ListContainerChanges_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ ContainerChanges, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m ListContainerChanges_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (ContainerChanges, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.ListContainerChanges", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ ContainerChanges, flags uint64, err error) {
		var out struct {
			Container ContainerChanges `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// ExportContainer creates an image from a container.  It takes the name or ID of a container and a
// path representing the target tarfile.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.
// The return value is the written tarfile.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ExportContainer '{"name": "flamboyant_payne", "path": "/tmp/payne.tar" }'
// {
//   "tarfile": "/tmp/payne.tar"
// }
// ~~~
type ExportContainer_methods struct{}

func ExportContainer() ExportContainer_methods { return ExportContainer_methods{} }

func (m ExportContainer_methods) Call(c *varlink.Connection, name_in_ string, path_in_ string) (tarfile_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, path_in_)
	if err_ != nil {
		return
	}
	tarfile_out_, _, err_ = receive()
	return
}

func (m ExportContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, path_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	in.Name = name_in_
	in.Path = path_in_
	receive, err := c.Send("io.podman.ExportContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (tarfile_out_ string, flags uint64, err error) {
		var out struct {
			Tarfile string `json:"tarfile"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		tarfile_out_ = out.Tarfile
		return
	}, nil
}

// GetContainerStats takes the name or ID of a container and returns a single ContainerStats structure which
// contains attributes like memory and cpu usage.  If the container cannot be found, a
// [ContainerNotFound](#ContainerNotFound) error will be returned. If the container is not running, a [NoContainerRunning](#NoContainerRunning)
// error will be returned
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.GetContainerStats '{"name": "c33e4164f384"}'
// {
//   "container": {
//     "block_input": 0,
//     "block_output": 0,
//     "cpu": 2.571123918839990154678e-08,
//     "cpu_nano": 49037378,
//     "id": "c33e4164f384aa9d979072a63319d66b74fd7a128be71fa68ede24f33ec6cfee",
//     "mem_limit": 33080606720,
//     "mem_perc": 2.166828456524753747370e-03,
//     "mem_usage": 716800,
//     "name": "competent_wozniak",
//     "net_input": 768,
//     "net_output": 5910,
//     "pids": 1,
//     "system_nano": 10000000
//   }
// }
// ~~~
type GetContainerStats_methods struct{}

func GetContainerStats() GetContainerStats_methods { return GetContainerStats_methods{} }

func (m GetContainerStats_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ ContainerStats, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m GetContainerStats_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (ContainerStats, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.GetContainerStats", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ ContainerStats, flags uint64, err error) {
		var out struct {
			Container ContainerStats `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// GetContainerStatsWithHistory takes a previous set of container statistics and uses libpod functions
// to calculate the containers statistics based on current and previous measurements.
type GetContainerStatsWithHistory_methods struct{}

func GetContainerStatsWithHistory() GetContainerStatsWithHistory_methods {
	return GetContainerStatsWithHistory_methods{}
}

func (m GetContainerStatsWithHistory_methods) Call(c *varlink.Connection, previousStats_in_ ContainerStats) (container_out_ ContainerStats, err_ error) {
	receive, err_ := m.Send(c, 0, previousStats_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m GetContainerStatsWithHistory_methods) Send(c *varlink.Connection, flags uint64, previousStats_in_ ContainerStats) (func() (ContainerStats, uint64, error), error) {
	var in struct {
		PreviousStats ContainerStats `json:"previousStats"`
	}
	in.PreviousStats = previousStats_in_
	receive, err := c.Send("io.podman.GetContainerStatsWithHistory", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ ContainerStats, flags uint64, err error) {
		var out struct {
			Container ContainerStats `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// StartContainer starts a created or stopped container. It takes the name or ID of container.  It returns
// the container ID once started.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.  See also [CreateContainer](#CreateContainer).
type StartContainer_methods struct{}

func StartContainer() StartContainer_methods { return StartContainer_methods{} }

func (m StartContainer_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m StartContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.StartContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// StopContainer stops a container given a timeout.  It takes the name or ID of a container as well as a
// timeout value.  The timeout value the time before a forcible stop to the container is applied.  It
// returns the container ID once stopped. If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned instead. See also [KillContainer](KillContainer).
// #### Error
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.StopContainer '{"name": "135d71b9495f", "timeout": 5}'
// {
//   "container": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
// }
// ~~~
type StopContainer_methods struct{}

func StopContainer() StopContainer_methods { return StopContainer_methods{} }

func (m StopContainer_methods) Call(c *varlink.Connection, name_in_ string, timeout_in_ int64) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, timeout_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m StopContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, timeout_in_ int64) (func() (string, uint64, error), error) {
	var in struct {
		Name    string `json:"name"`
		Timeout int64  `json:"timeout"`
	}
	in.Name = name_in_
	in.Timeout = timeout_in_
	receive, err := c.Send("io.podman.StopContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// InitContainer initializes the given container. It accepts a container name or
// ID, and will initialize the container matching that ID if possible, and error
// if not. Containers can only be initialized when they are in the Created or
// Exited states. Initialization prepares a container to be started, but does not
// start the container. It is intended to be used to debug a container's state
// prior to starting it.
type InitContainer_methods struct{}

func InitContainer() InitContainer_methods { return InitContainer_methods{} }

func (m InitContainer_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m InitContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.InitContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// RestartContainer will restart a running container given a container name or ID and timeout value. The timeout
// value is the time before a forcible stop is used to stop the container.  If the container cannot be found by
// name or ID, a [ContainerNotFound](#ContainerNotFound)  error will be returned; otherwise, the ID of the
// container will be returned.
type RestartContainer_methods struct{}

func RestartContainer() RestartContainer_methods { return RestartContainer_methods{} }

func (m RestartContainer_methods) Call(c *varlink.Connection, name_in_ string, timeout_in_ int64) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, timeout_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m RestartContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, timeout_in_ int64) (func() (string, uint64, error), error) {
	var in struct {
		Name    string `json:"name"`
		Timeout int64  `json:"timeout"`
	}
	in.Name = name_in_
	in.Timeout = timeout_in_
	receive, err := c.Send("io.podman.RestartContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// KillContainer takes the name or ID of a container as well as a signal to be applied to the container.  Once the
// container has been killed, the container's ID is returned.  If the container cannot be found, a
// [ContainerNotFound](#ContainerNotFound) error is returned. See also [StopContainer](StopContainer).
type KillContainer_methods struct{}

func KillContainer() KillContainer_methods { return KillContainer_methods{} }

func (m KillContainer_methods) Call(c *varlink.Connection, name_in_ string, signal_in_ int64) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, signal_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m KillContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, signal_in_ int64) (func() (string, uint64, error), error) {
	var in struct {
		Name   string `json:"name"`
		Signal int64  `json:"signal"`
	}
	in.Name = name_in_
	in.Signal = signal_in_
	receive, err := c.Send("io.podman.KillContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// PauseContainer takes the name or ID of container and pauses it.  If the container cannot be found,
// a [ContainerNotFound](#ContainerNotFound) error will be returned; otherwise the ID of the container is returned.
// See also [UnpauseContainer](#UnpauseContainer).
type PauseContainer_methods struct{}

func PauseContainer() PauseContainer_methods { return PauseContainer_methods{} }

func (m PauseContainer_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m PauseContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.PauseContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// UnpauseContainer takes the name or ID of container and unpauses a paused container.  If the container cannot be
// found, a [ContainerNotFound](#ContainerNotFound) error will be returned; otherwise the ID of the container is returned.
// See also [PauseContainer](#PauseContainer).
type UnpauseContainer_methods struct{}

func UnpauseContainer() UnpauseContainer_methods { return UnpauseContainer_methods{} }

func (m UnpauseContainer_methods) Call(c *varlink.Connection, name_in_ string) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m UnpauseContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.UnpauseContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// Attach takes the name or ID of a container and sets up the ability to remotely attach to its console. The start
// bool is whether you wish to start the container in question first.
type Attach_methods struct{}

func Attach() Attach_methods { return Attach_methods{} }

func (m Attach_methods) Call(c *varlink.Connection, name_in_ string, detachKeys_in_ string, start_in_ bool) (err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, detachKeys_in_, start_in_)
	if err_ != nil {
		return
	}
	_, err_ = receive()
	return
}

func (m Attach_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, detachKeys_in_ string, start_in_ bool) (func() (uint64, error), error) {
	var in struct {
		Name       string `json:"name"`
		DetachKeys string `json:"detachKeys"`
		Start      bool   `json:"start"`
	}
	in.Name = name_in_
	in.DetachKeys = detachKeys_in_
	in.Start = start_in_
	receive, err := c.Send("io.podman.Attach", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (flags uint64, err error) {
		flags, err = receive(nil)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		return
	}, nil
}

type AttachControl_methods struct{}

func AttachControl() AttachControl_methods { return AttachControl_methods{} }

func (m AttachControl_methods) Call(c *varlink.Connection, name_in_ string) (err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	_, err_ = receive()
	return
}

func (m AttachControl_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.AttachControl", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (flags uint64, err error) {
		flags, err = receive(nil)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		return
	}, nil
}

// GetAttachSockets takes the name or ID of an existing container.  It returns file paths for two sockets needed
// to properly communicate with a container.  The first is the actual I/O socket that the container uses.  The
// second is a "control" socket where things like resizing the TTY events are sent. If the container cannot be
// found, a [ContainerNotFound](#ContainerNotFound) error will be returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/io.podman/io.podman.GetAttachSockets '{"name": "b7624e775431219161"}'
// {
//   "sockets": {
//     "container_id": "b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca",
//     "control_socket": "/var/lib/containers/storage/overlay-containers/b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca/userdata/ctl",
//     "io_socket": "/var/run/libpod/socket/b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca/attach"
//   }
// }
// ~~~
type GetAttachSockets_methods struct{}

func GetAttachSockets() GetAttachSockets_methods { return GetAttachSockets_methods{} }

func (m GetAttachSockets_methods) Call(c *varlink.Connection, name_in_ string) (sockets_out_ Sockets, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	sockets_out_, _, err_ = receive()
	return
}

func (m GetAttachSockets_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (Sockets, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.GetAttachSockets", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (sockets_out_ Sockets, flags uint64, err error) {
		var out struct {
			Sockets Sockets `json:"sockets"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		sockets_out_ = out.Sockets
		return
	}, nil
}

// WaitContainer takes the name or ID of a container and waits the given interval in milliseconds until the container
// stops.  Upon stopping, the return code of the container is returned. If the container container cannot be found by ID
// or name, a [ContainerNotFound](#ContainerNotFound) error is returned.
type WaitContainer_methods struct{}

func WaitContainer() WaitContainer_methods { return WaitContainer_methods{} }

func (m WaitContainer_methods) Call(c *varlink.Connection, name_in_ string, interval_in_ int64) (exitcode_out_ int64, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, interval_in_)
	if err_ != nil {
		return
	}
	exitcode_out_, _, err_ = receive()
	return
}

func (m WaitContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, interval_in_ int64) (func() (int64, uint64, error), error) {
	var in struct {
		Name     string `json:"name"`
		Interval int64  `json:"interval"`
	}
	in.Name = name_in_
	in.Interval = interval_in_
	receive, err := c.Send("io.podman.WaitContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (exitcode_out_ int64, flags uint64, err error) {
		var out struct {
			Exitcode int64 `json:"exitcode"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		exitcode_out_ = out.Exitcode
		return
	}, nil
}

// RemoveContainer requires the name or ID of container as well a boolean representing whether a running container can be stopped and removed, and a boolean
// indicating whether to remove builtin volumes. Upon successful removal of the
// container, its ID is returned.  If the
// container cannot be found by name or ID, a [ContainerNotFound](#ContainerNotFound) error will be returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.RemoveContainer '{"name": "62f4fd98cb57"}'
// {
//   "container": "62f4fd98cb57f529831e8f90610e54bba74bd6f02920ffb485e15376ed365c20"
// }
// ~~~
type RemoveContainer_methods struct{}

func RemoveContainer() RemoveContainer_methods { return RemoveContainer_methods{} }

func (m RemoveContainer_methods) Call(c *varlink.Connection, name_in_ string, force_in_ bool, removeVolumes_in_ bool) (container_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, force_in_, removeVolumes_in_)
	if err_ != nil {
		return
	}
	container_out_, _, err_ = receive()
	return
}

func (m RemoveContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, force_in_ bool, removeVolumes_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Name          string `json:"name"`
		Force         bool   `json:"force"`
		RemoveVolumes bool   `json:"removeVolumes"`
	}
	in.Name = name_in_
	in.Force = force_in_
	in.RemoveVolumes = removeVolumes_in_
	receive, err := c.Send("io.podman.RemoveContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (container_out_ string, flags uint64, err error) {
		var out struct {
			Container string `json:"container"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		container_out_ = out.Container
		return
	}, nil
}

// DeleteStoppedContainers will delete all containers that are not running. It will return a list the deleted
// container IDs.  See also [RemoveContainer](RemoveContainer).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.DeleteStoppedContainers
// {
//   "containers": [
//     "451410b931d00def8aa9b4f8084e4d4a39e5e04ea61f358cf53a5cf95afcdcee",
//     "8b60f754a3e01389494a9581ade97d35c2765b6e2f19acd2d3040c82a32d1bc0",
//     "cf2e99d4d3cad6073df199ed32bbe64b124f3e1aba6d78821aa8460e70d30084",
//     "db901a329587312366e5ecff583d08f0875b4b79294322df67d90fc6eed08fc1"
//   ]
// }
// ~~~
type DeleteStoppedContainers_methods struct{}

func DeleteStoppedContainers() DeleteStoppedContainers_methods {
	return DeleteStoppedContainers_methods{}
}

func (m DeleteStoppedContainers_methods) Call(c *varlink.Connection) (containers_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	containers_out_, _, err_ = receive()
	return
}

func (m DeleteStoppedContainers_methods) Send(c *varlink.Connection, flags uint64) (func() ([]string, uint64, error), error) {
	receive, err := c.Send("io.podman.DeleteStoppedContainers", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (containers_out_ []string, flags uint64, err error) {
		var out struct {
			Containers []string `json:"containers"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		containers_out_ = []string(out.Containers)
		return
	}, nil
}

// ListImages returns information about the images that are currently in storage.
// See also [InspectImage](#InspectImage).
type ListImages_methods struct{}

func ListImages() ListImages_methods { return ListImages_methods{} }

func (m ListImages_methods) Call(c *varlink.Connection) (images_out_ []Image, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	images_out_, _, err_ = receive()
	return
}

func (m ListImages_methods) Send(c *varlink.Connection, flags uint64) (func() ([]Image, uint64, error), error) {
	receive, err := c.Send("io.podman.ListImages", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (images_out_ []Image, flags uint64, err error) {
		var out struct {
			Images []Image `json:"images"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		images_out_ = []Image(out.Images)
		return
	}, nil
}

// GetImage returns information about a single image in storage.
// If the image caGetImage returns be found, [ImageNotFound](#ImageNotFound) will be returned.
type GetImage_methods struct{}

func GetImage() GetImage_methods { return GetImage_methods{} }

func (m GetImage_methods) Call(c *varlink.Connection, id_in_ string) (image_out_ Image, err_ error) {
	receive, err_ := m.Send(c, 0, id_in_)
	if err_ != nil {
		return
	}
	image_out_, _, err_ = receive()
	return
}

func (m GetImage_methods) Send(c *varlink.Connection, flags uint64, id_in_ string) (func() (Image, uint64, error), error) {
	var in struct {
		Id string `json:"id"`
	}
	in.Id = id_in_
	receive, err := c.Send("io.podman.GetImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (image_out_ Image, flags uint64, err error) {
		var out struct {
			Image Image `json:"image"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		image_out_ = out.Image
		return
	}, nil
}

// BuildImage takes a [BuildInfo](#BuildInfo) structure and builds an image.  At a minimum, you must provide the
// 'dockerfile' and 'tags' options in the BuildInfo structure. It will return a [MoreResponse](#MoreResponse) structure
// that contains the build logs and resulting image ID.
type BuildImage_methods struct{}

func BuildImage() BuildImage_methods { return BuildImage_methods{} }

func (m BuildImage_methods) Call(c *varlink.Connection, build_in_ BuildInfo) (image_out_ MoreResponse, err_ error) {
	receive, err_ := m.Send(c, 0, build_in_)
	if err_ != nil {
		return
	}
	image_out_, _, err_ = receive()
	return
}

func (m BuildImage_methods) Send(c *varlink.Connection, flags uint64, build_in_ BuildInfo) (func() (MoreResponse, uint64, error), error) {
	var in struct {
		Build BuildInfo `json:"build"`
	}
	in.Build = build_in_
	receive, err := c.Send("io.podman.BuildImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (image_out_ MoreResponse, flags uint64, err error) {
		var out struct {
			Image MoreResponse `json:"image"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		image_out_ = out.Image
		return
	}, nil
}

// InspectImage takes the name or ID of an image and returns a string representation of data associated with the
// mage.  You must serialize the string into JSON to use it further.  An [ImageNotFound](#ImageNotFound) error will
// be returned if the image cannot be found.
type InspectImage_methods struct{}

func InspectImage() InspectImage_methods { return InspectImage_methods{} }

func (m InspectImage_methods) Call(c *varlink.Connection, name_in_ string) (image_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	image_out_, _, err_ = receive()
	return
}

func (m InspectImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.InspectImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (image_out_ string, flags uint64, err error) {
		var out struct {
			Image string `json:"image"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		image_out_ = out.Image
		return
	}, nil
}

// HistoryImage takes the name or ID of an image and returns information about its history and layers.  The returned
// history is in the form of an array of ImageHistory structures.  If the image cannot be found, an
// [ImageNotFound](#ImageNotFound) error is returned.
type HistoryImage_methods struct{}

func HistoryImage() HistoryImage_methods { return HistoryImage_methods{} }

func (m HistoryImage_methods) Call(c *varlink.Connection, name_in_ string) (history_out_ []ImageHistory, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	history_out_, _, err_ = receive()
	return
}

func (m HistoryImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() ([]ImageHistory, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.HistoryImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (history_out_ []ImageHistory, flags uint64, err error) {
		var out struct {
			History []ImageHistory `json:"history"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		history_out_ = []ImageHistory(out.History)
		return
	}, nil
}

// PushImage takes two input arguments: the name or ID of an image, the fully-qualified destination name of the image,
// It will return an [ImageNotFound](#ImageNotFound) error if
// the image cannot be found in local storage; otherwise it will return a [MoreResponse](#MoreResponse)
type PushImage_methods struct{}

func PushImage() PushImage_methods { return PushImage_methods{} }

func (m PushImage_methods) Call(c *varlink.Connection, name_in_ string, tag_in_ string, compress_in_ bool, format_in_ string, removeSignatures_in_ bool, signBy_in_ string) (reply_out_ MoreResponse, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, tag_in_, compress_in_, format_in_, removeSignatures_in_, signBy_in_)
	if err_ != nil {
		return
	}
	reply_out_, _, err_ = receive()
	return
}

func (m PushImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, tag_in_ string, compress_in_ bool, format_in_ string, removeSignatures_in_ bool, signBy_in_ string) (func() (MoreResponse, uint64, error), error) {
	var in struct {
		Name             string `json:"name"`
		Tag              string `json:"tag"`
		Compress         bool   `json:"compress"`
		Format           string `json:"format"`
		RemoveSignatures bool   `json:"removeSignatures"`
		SignBy           string `json:"signBy"`
	}
	in.Name = name_in_
	in.Tag = tag_in_
	in.Compress = compress_in_
	in.Format = format_in_
	in.RemoveSignatures = removeSignatures_in_
	in.SignBy = signBy_in_
	receive, err := c.Send("io.podman.PushImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (reply_out_ MoreResponse, flags uint64, err error) {
		var out struct {
			Reply MoreResponse `json:"reply"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		reply_out_ = out.Reply
		return
	}, nil
}

// TagImage takes the name or ID of an image in local storage as well as the desired tag name.  If the image cannot
// be found, an [ImageNotFound](#ImageNotFound) error will be returned; otherwise, the ID of the image is returned on success.
type TagImage_methods struct{}

func TagImage() TagImage_methods { return TagImage_methods{} }

func (m TagImage_methods) Call(c *varlink.Connection, name_in_ string, tagged_in_ string) (image_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, tagged_in_)
	if err_ != nil {
		return
	}
	image_out_, _, err_ = receive()
	return
}

func (m TagImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, tagged_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name   string `json:"name"`
		Tagged string `json:"tagged"`
	}
	in.Name = name_in_
	in.Tagged = tagged_in_
	receive, err := c.Send("io.podman.TagImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (image_out_ string, flags uint64, err error) {
		var out struct {
			Image string `json:"image"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		image_out_ = out.Image
		return
	}, nil
}

// RemoveImage takes the name or ID of an image as well as a boolean that determines if containers using that image
// should be deleted.  If the image cannot be found, an [ImageNotFound](#ImageNotFound) error will be returned.  The
// ID of the removed image is returned when complete.  See also [DeleteUnusedImages](DeleteUnusedImages).
// #### Example
// ~~~
// varlink call -m unix:/run/podman/io.podman/io.podman.RemoveImage '{"name": "registry.fedoraproject.org/fedora", "force": true}'
// {
//   "image": "426866d6fa419873f97e5cbd320eeb22778244c1dfffa01c944db3114f55772e"
// }
// ~~~
type RemoveImage_methods struct{}

func RemoveImage() RemoveImage_methods { return RemoveImage_methods{} }

func (m RemoveImage_methods) Call(c *varlink.Connection, name_in_ string, force_in_ bool) (image_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, force_in_)
	if err_ != nil {
		return
	}
	image_out_, _, err_ = receive()
	return
}

func (m RemoveImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, force_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Name  string `json:"name"`
		Force bool   `json:"force"`
	}
	in.Name = name_in_
	in.Force = force_in_
	receive, err := c.Send("io.podman.RemoveImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (image_out_ string, flags uint64, err error) {
		var out struct {
			Image string `json:"image"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		image_out_ = out.Image
		return
	}, nil
}

// SearchImages searches available registries for images that contain the
// contents of "query" in their name. If "limit" is given, limits the amount of
// search results per registry.
type SearchImages_methods struct{}

func SearchImages() SearchImages_methods { return SearchImages_methods{} }

func (m SearchImages_methods) Call(c *varlink.Connection, query_in_ string, limit_in_ *int64, filter_in_ ImageSearchFilter) (results_out_ []ImageSearchResult, err_ error) {
	receive, err_ := m.Send(c, 0, query_in_, limit_in_, filter_in_)
	if err_ != nil {
		return
	}
	results_out_, _, err_ = receive()
	return
}

func (m SearchImages_methods) Send(c *varlink.Connection, flags uint64, query_in_ string, limit_in_ *int64, filter_in_ ImageSearchFilter) (func() ([]ImageSearchResult, uint64, error), error) {
	var in struct {
		Query  string            `json:"query"`
		Limit  *int64            `json:"limit,omitempty"`
		Filter ImageSearchFilter `json:"filter"`
	}
	in.Query = query_in_
	in.Limit = limit_in_
	in.Filter = filter_in_
	receive, err := c.Send("io.podman.SearchImages", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (results_out_ []ImageSearchResult, flags uint64, err error) {
		var out struct {
			Results []ImageSearchResult `json:"results"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		results_out_ = []ImageSearchResult(out.Results)
		return
	}, nil
}

// DeleteUnusedImages deletes any images not associated with a container.  The IDs of the deleted images are returned
// in a string array.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.DeleteUnusedImages
// {
//   "images": [
//     "166ea6588079559c724c15223f52927f514f73dd5c5cf2ae2d143e3b2e6e9b52",
//     "da86e6ba6ca197bf6bc5e9d900febd906b133eaa4750e6bed647b0fbe50ed43e",
//     "3ef70f7291f47dfe2b82931a993e16f5a44a0e7a68034c3e0e086d77f5829adc",
//     "59788edf1f3e78cd0ebe6ce1446e9d10788225db3dedcfd1a59f764bad2b2690"
//   ]
// }
// ~~~
type DeleteUnusedImages_methods struct{}

func DeleteUnusedImages() DeleteUnusedImages_methods { return DeleteUnusedImages_methods{} }

func (m DeleteUnusedImages_methods) Call(c *varlink.Connection) (images_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	images_out_, _, err_ = receive()
	return
}

func (m DeleteUnusedImages_methods) Send(c *varlink.Connection, flags uint64) (func() ([]string, uint64, error), error) {
	receive, err := c.Send("io.podman.DeleteUnusedImages", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (images_out_ []string, flags uint64, err error) {
		var out struct {
			Images []string `json:"images"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		images_out_ = []string(out.Images)
		return
	}, nil
}

// Commit, creates an image from an existing container. It requires the name or
// ID of the container as well as the resulting image name.  Optionally, you can define an author and message
// to be added to the resulting image.  You can also define changes to the resulting image for the following
// attributes: _CMD, ENTRYPOINT, ENV, EXPOSE, LABEL, ONBUILD, STOPSIGNAL, USER, VOLUME, and WORKDIR_.  To pause the
// container while it is being committed, pass a _true_ bool for the pause argument.  If the container cannot
// be found by the ID or name provided, a (ContainerNotFound)[#ContainerNotFound] error will be returned; otherwise,
// the resulting image's ID will be returned as a string inside a MoreResponse.
type Commit_methods struct{}

func Commit() Commit_methods { return Commit_methods{} }

func (m Commit_methods) Call(c *varlink.Connection, name_in_ string, image_name_in_ string, changes_in_ []string, author_in_ string, message_in_ string, pause_in_ bool, manifestType_in_ string) (reply_out_ MoreResponse, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, image_name_in_, changes_in_, author_in_, message_in_, pause_in_, manifestType_in_)
	if err_ != nil {
		return
	}
	reply_out_, _, err_ = receive()
	return
}

func (m Commit_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, image_name_in_ string, changes_in_ []string, author_in_ string, message_in_ string, pause_in_ bool, manifestType_in_ string) (func() (MoreResponse, uint64, error), error) {
	var in struct {
		Name         string   `json:"name"`
		Image_name   string   `json:"image_name"`
		Changes      []string `json:"changes"`
		Author       string   `json:"author"`
		Message      string   `json:"message"`
		Pause        bool     `json:"pause"`
		ManifestType string   `json:"manifestType"`
	}
	in.Name = name_in_
	in.Image_name = image_name_in_
	in.Changes = []string(changes_in_)
	in.Author = author_in_
	in.Message = message_in_
	in.Pause = pause_in_
	in.ManifestType = manifestType_in_
	receive, err := c.Send("io.podman.Commit", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (reply_out_ MoreResponse, flags uint64, err error) {
		var out struct {
			Reply MoreResponse `json:"reply"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		reply_out_ = out.Reply
		return
	}, nil
}

// ImportImage imports an image from a source (like tarball) into local storage.  The image can have additional
// descriptions added to it using the message and changes options. See also [ExportImage](ExportImage).
type ImportImage_methods struct{}

func ImportImage() ImportImage_methods { return ImportImage_methods{} }

func (m ImportImage_methods) Call(c *varlink.Connection, source_in_ string, reference_in_ string, message_in_ string, changes_in_ []string, delete_in_ bool) (image_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, source_in_, reference_in_, message_in_, changes_in_, delete_in_)
	if err_ != nil {
		return
	}
	image_out_, _, err_ = receive()
	return
}

func (m ImportImage_methods) Send(c *varlink.Connection, flags uint64, source_in_ string, reference_in_ string, message_in_ string, changes_in_ []string, delete_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Source    string   `json:"source"`
		Reference string   `json:"reference"`
		Message   string   `json:"message"`
		Changes   []string `json:"changes"`
		Delete    bool     `json:"delete"`
	}
	in.Source = source_in_
	in.Reference = reference_in_
	in.Message = message_in_
	in.Changes = []string(changes_in_)
	in.Delete = delete_in_
	receive, err := c.Send("io.podman.ImportImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (image_out_ string, flags uint64, err error) {
		var out struct {
			Image string `json:"image"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		image_out_ = out.Image
		return
	}, nil
}

// ExportImage takes the name or ID of an image and exports it to a destination like a tarball.  There is also
// a boolean option to force compression.  It also takes in a string array of tags to be able to save multiple
// tags of the same image to a tarball (each tag should be of the form <image>:<tag>).  Upon completion, the ID
// of the image is returned. If the image cannot be found in local storage, an [ImageNotFound](#ImageNotFound)
// error will be returned. See also [ImportImage](ImportImage).
type ExportImage_methods struct{}

func ExportImage() ExportImage_methods { return ExportImage_methods{} }

func (m ExportImage_methods) Call(c *varlink.Connection, name_in_ string, destination_in_ string, compress_in_ bool, tags_in_ []string) (image_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, destination_in_, compress_in_, tags_in_)
	if err_ != nil {
		return
	}
	image_out_, _, err_ = receive()
	return
}

func (m ExportImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, destination_in_ string, compress_in_ bool, tags_in_ []string) (func() (string, uint64, error), error) {
	var in struct {
		Name        string   `json:"name"`
		Destination string   `json:"destination"`
		Compress    bool     `json:"compress"`
		Tags        []string `json:"tags"`
	}
	in.Name = name_in_
	in.Destination = destination_in_
	in.Compress = compress_in_
	in.Tags = []string(tags_in_)
	receive, err := c.Send("io.podman.ExportImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (image_out_ string, flags uint64, err error) {
		var out struct {
			Image string `json:"image"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		image_out_ = out.Image
		return
	}, nil
}

// PullImage pulls an image from a repository to local storage.  After a successful pull, the image id and logs
// are returned as a [MoreResponse](#MoreResponse).  This connection also will handle a WantsMores request to send
// status as it occurs.
type PullImage_methods struct{}

func PullImage() PullImage_methods { return PullImage_methods{} }

func (m PullImage_methods) Call(c *varlink.Connection, name_in_ string) (reply_out_ MoreResponse, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	reply_out_, _, err_ = receive()
	return
}

func (m PullImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (MoreResponse, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.PullImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (reply_out_ MoreResponse, flags uint64, err error) {
		var out struct {
			Reply MoreResponse `json:"reply"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		reply_out_ = out.Reply
		return
	}, nil
}

// CreatePod creates a new empty pod.  It uses a [PodCreate](#PodCreate) type for input.
// On success, the ID of the newly created pod will be returned.
// #### Example
// ~~~
// $ varlink call unix:/run/podman/io.podman/io.podman.CreatePod '{"create": {"name": "test"}}'
// {
//   "pod": "b05dee7bd4ccfee688099fe1588a7a898d6ddd6897de9251d4671c9b0feacb2a"
// }
// # $ varlink call unix:/run/podman/io.podman/io.podman.CreatePod '{"create": {"infra": true, "share": ["ipc", "net", "uts"]}}'
// {
//   "pod": "d7697449a8035f613c1a8891286502aca68fff7d5d49a85279b3bda229af3b28"
// }
// ~~~
type CreatePod_methods struct{}

func CreatePod() CreatePod_methods { return CreatePod_methods{} }

func (m CreatePod_methods) Call(c *varlink.Connection, create_in_ PodCreate) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, create_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m CreatePod_methods) Send(c *varlink.Connection, flags uint64, create_in_ PodCreate) (func() (string, uint64, error), error) {
	var in struct {
		Create PodCreate `json:"create"`
	}
	in.Create = create_in_
	receive, err := c.Send("io.podman.CreatePod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// ListPods returns a list of pods in no particular order.  They are
// returned as an array of ListPodData structs.  See also [GetPod](#GetPod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ListPods
// {
//   "pods": [
//     {
//       "cgroup": "machine.slice",
//       "containersinfo": [
//         {
//           "id": "00c130a45de0411f109f1a0cfea2e298df71db20fa939de5cab8b2160a36be45",
//           "name": "1840835294cf-infra",
//           "status": "running"
//         },
//         {
//           "id": "49a5cce72093a5ca47c6de86f10ad7bb36391e2d89cef765f807e460865a0ec6",
//           "name": "upbeat_murdock",
//           "status": "running"
//         }
//       ],
//       "createdat": "2018-12-07 13:10:15.014139258 -0600 CST",
//       "id": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f",
//       "name": "foobar",
//       "numberofcontainers": "2",
//       "status": "Running"
//     },
//     {
//       "cgroup": "machine.slice",
//       "containersinfo": [
//         {
//           "id": "1ca4b7bbba14a75ba00072d4b705c77f3df87db0109afaa44d50cb37c04a477e",
//           "name": "784306f655c6-infra",
//           "status": "running"
//         }
//       ],
//       "createdat": "2018-12-07 13:09:57.105112457 -0600 CST",
//       "id": "784306f655c6200aea321dd430ba685e9b2cc1f7d7528a72f3ff74ffb29485a2",
//       "name": "nostalgic_pike",
//       "numberofcontainers": "1",
//       "status": "Running"
//     }
//   ]
// }
// ~~~
type ListPods_methods struct{}

func ListPods() ListPods_methods { return ListPods_methods{} }

func (m ListPods_methods) Call(c *varlink.Connection) (pods_out_ []ListPodData, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	pods_out_, _, err_ = receive()
	return
}

func (m ListPods_methods) Send(c *varlink.Connection, flags uint64) (func() ([]ListPodData, uint64, error), error) {
	receive, err := c.Send("io.podman.ListPods", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (pods_out_ []ListPodData, flags uint64, err error) {
		var out struct {
			Pods []ListPodData `json:"pods"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pods_out_ = []ListPodData(out.Pods)
		return
	}, nil
}

// GetPod takes a name or ID of a pod and returns single [ListPodData](#ListPodData)
// structure.  A [PodNotFound](#PodNotFound) error will be returned if the pod cannot be found.
// See also [ListPods](ListPods).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.GetPod '{"name": "foobar"}'
// {
//   "pod": {
//     "cgroup": "machine.slice",
//     "containersinfo": [
//       {
//         "id": "00c130a45de0411f109f1a0cfea2e298df71db20fa939de5cab8b2160a36be45",
//         "name": "1840835294cf-infra",
//         "status": "running"
//       },
//       {
//         "id": "49a5cce72093a5ca47c6de86f10ad7bb36391e2d89cef765f807e460865a0ec6",
//         "name": "upbeat_murdock",
//         "status": "running"
//       }
//     ],
//     "createdat": "2018-12-07 13:10:15.014139258 -0600 CST",
//     "id": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f",
//     "name": "foobar",
//     "numberofcontainers": "2",
//     "status": "Running"
//   }
// }
// ~~~
type GetPod_methods struct{}

func GetPod() GetPod_methods { return GetPod_methods{} }

func (m GetPod_methods) Call(c *varlink.Connection, name_in_ string) (pod_out_ ListPodData, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m GetPod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (ListPodData, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.GetPod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ ListPodData, flags uint64, err error) {
		var out struct {
			Pod ListPodData `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// InspectPod takes the name or ID of an image and returns a string representation of data associated with the
// pod.  You must serialize the string into JSON to use it further.  A [PodNotFound](#PodNotFound) error will
// be returned if the pod cannot be found.
type InspectPod_methods struct{}

func InspectPod() InspectPod_methods { return InspectPod_methods{} }

func (m InspectPod_methods) Call(c *varlink.Connection, name_in_ string) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m InspectPod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.InspectPod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// StartPod starts containers in a pod.  It takes the name or ID of pod.  If the pod cannot be found, a [PodNotFound](#PodNotFound)
// error will be returned.  Containers in a pod are started independently. If there is an error starting one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was started with no errors, the pod ID is returned.
// See also [CreatePod](#CreatePod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.StartPod '{"name": "135d71b9495f"}'
// {
//   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6",
// }
// ~~~
type StartPod_methods struct{}

func StartPod() StartPod_methods { return StartPod_methods{} }

func (m StartPod_methods) Call(c *varlink.Connection, name_in_ string) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m StartPod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.StartPod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// StopPod stops containers in a pod.  It takes the name or ID of a pod and a timeout.
// If the pod cannot be found, a [PodNotFound](#PodNotFound) error will be returned instead.
// Containers in a pod are stopped independently. If there is an error stopping one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was stopped with no errors, the pod ID is returned.
// See also [KillPod](KillPod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.StopPod '{"name": "135d71b9495f"}'
// {
//   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
// }
// ~~~
type StopPod_methods struct{}

func StopPod() StopPod_methods { return StopPod_methods{} }

func (m StopPod_methods) Call(c *varlink.Connection, name_in_ string, timeout_in_ int64) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, timeout_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m StopPod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, timeout_in_ int64) (func() (string, uint64, error), error) {
	var in struct {
		Name    string `json:"name"`
		Timeout int64  `json:"timeout"`
	}
	in.Name = name_in_
	in.Timeout = timeout_in_
	receive, err := c.Send("io.podman.StopPod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// RestartPod will restart containers in a pod given a pod name or ID. Containers in
// the pod that are running will be stopped, then all stopped containers will be run.
// If the pod cannot be found by name or ID, a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are restarted independently. If there is an error restarting one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was restarted with no errors, the pod ID is returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.RestartPod '{"name": "135d71b9495f"}'
// {
//   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
// }
// ~~~
type RestartPod_methods struct{}

func RestartPod() RestartPod_methods { return RestartPod_methods{} }

func (m RestartPod_methods) Call(c *varlink.Connection, name_in_ string) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m RestartPod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.RestartPod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// KillPod takes the name or ID of a pod as well as a signal to be applied to the pod.  If the pod cannot be found, a
// [PodNotFound](#PodNotFound) error is returned.
// Containers in a pod are killed independently. If there is an error killing one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was killed with no errors, the pod ID is returned.
// See also [StopPod](StopPod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.KillPod '{"name": "foobar", "signal": 15}'
// {
//   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
// }
// ~~~
type KillPod_methods struct{}

func KillPod() KillPod_methods { return KillPod_methods{} }

func (m KillPod_methods) Call(c *varlink.Connection, name_in_ string, signal_in_ int64) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, signal_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m KillPod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, signal_in_ int64) (func() (string, uint64, error), error) {
	var in struct {
		Name   string `json:"name"`
		Signal int64  `json:"signal"`
	}
	in.Name = name_in_
	in.Signal = signal_in_
	receive, err := c.Send("io.podman.KillPod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// PausePod takes the name or ID of a pod and pauses the running containers associated with it.  If the pod cannot be found,
// a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are paused independently. If there is an error pausing one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was paused with no errors, the pod ID is returned.
// See also [UnpausePod](#UnpausePod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.PausePod '{"name": "foobar"}'
// {
//   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
// }
// ~~~
type PausePod_methods struct{}

func PausePod() PausePod_methods { return PausePod_methods{} }

func (m PausePod_methods) Call(c *varlink.Connection, name_in_ string) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m PausePod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.PausePod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// UnpausePod takes the name or ID of a pod and unpauses the paused containers associated with it.  If the pod cannot be
// found, a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are unpaused independently. If there is an error unpausing one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was unpaused with no errors, the pod ID is returned.
// See also [PausePod](#PausePod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.UnpausePod '{"name": "foobar"}'
// {
//   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
// }
// ~~~
type UnpausePod_methods struct{}

func UnpausePod() UnpausePod_methods { return UnpausePod_methods{} }

func (m UnpausePod_methods) Call(c *varlink.Connection, name_in_ string) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m UnpausePod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.UnpausePod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// RemovePod takes the name or ID of a pod as well a boolean representing whether a running
// container in the pod can be stopped and removed.  If a pod has containers associated with it, and force is not true,
// an error will occur.
// If the pod cannot be found by name or ID, a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are removed independently. If there is an error removing any container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was removed with no errors, the pod ID is returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.RemovePod '{"name": "62f4fd98cb57", "force": "true"}'
// {
//   "pod": "62f4fd98cb57f529831e8f90610e54bba74bd6f02920ffb485e15376ed365c20"
// }
// ~~~
type RemovePod_methods struct{}

func RemovePod() RemovePod_methods { return RemovePod_methods{} }

func (m RemovePod_methods) Call(c *varlink.Connection, name_in_ string, force_in_ bool) (pod_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, force_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m RemovePod_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, force_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Name  string `json:"name"`
		Force bool   `json:"force"`
	}
	in.Name = name_in_
	in.Force = force_in_
	receive, err := c.Send("io.podman.RemovePod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, flags uint64, err error) {
		var out struct {
			Pod string `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

type TopPod_methods struct{}

func TopPod() TopPod_methods { return TopPod_methods{} }

func (m TopPod_methods) Call(c *varlink.Connection, pod_in_ string, latest_in_ bool, descriptors_in_ []string) (stats_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, pod_in_, latest_in_, descriptors_in_)
	if err_ != nil {
		return
	}
	stats_out_, _, err_ = receive()
	return
}

func (m TopPod_methods) Send(c *varlink.Connection, flags uint64, pod_in_ string, latest_in_ bool, descriptors_in_ []string) (func() ([]string, uint64, error), error) {
	var in struct {
		Pod         string   `json:"pod"`
		Latest      bool     `json:"latest"`
		Descriptors []string `json:"descriptors"`
	}
	in.Pod = pod_in_
	in.Latest = latest_in_
	in.Descriptors = []string(descriptors_in_)
	receive, err := c.Send("io.podman.TopPod", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (stats_out_ []string, flags uint64, err error) {
		var out struct {
			Stats []string `json:"stats"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		stats_out_ = []string(out.Stats)
		return
	}, nil
}

// GetPodStats takes the name or ID of a pod and returns a pod name and slice of ContainerStats structure which
// contains attributes like memory and cpu usage.  If the pod cannot be found, a [PodNotFound](#PodNotFound)
// error will be returned.  If the pod has no running containers associated with it, a [NoContainerRunning](#NoContainerRunning)
// error will be returned.
// #### Example
// ~~~
// $ varlink call unix:/run/podman/io.podman/io.podman.GetPodStats '{"name": "7f62b508b6f12b11d8fe02e"}'
// {
//   "containers": [
//     {
//       "block_input": 0,
//       "block_output": 0,
//       "cpu": 2.833470544016107524276e-08,
//       "cpu_nano": 54363072,
//       "id": "a64b51f805121fe2c5a3dc5112eb61d6ed139e3d1c99110360d08b58d48e4a93",
//       "mem_limit": 12276146176,
//       "mem_perc": 7.974359265237864966003e-03,
//       "mem_usage": 978944,
//       "name": "quirky_heisenberg",
//       "net_input": 866,
//       "net_output": 7388,
//       "pids": 1,
//       "system_nano": 20000000
//     }
//   ],
//   "pod": "7f62b508b6f12b11d8fe02e0db4de6b9e43a7d7699b33a4fc0d574f6e82b4ebd"
// }
// ~~~
type GetPodStats_methods struct{}

func GetPodStats() GetPodStats_methods { return GetPodStats_methods{} }

func (m GetPodStats_methods) Call(c *varlink.Connection, name_in_ string) (pod_out_ string, containers_out_ []ContainerStats, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	pod_out_, containers_out_, _, err_ = receive()
	return
}

func (m GetPodStats_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, []ContainerStats, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.GetPodStats", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ string, containers_out_ []ContainerStats, flags uint64, err error) {
		var out struct {
			Pod        string           `json:"pod"`
			Containers []ContainerStats `json:"containers"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		containers_out_ = []ContainerStats(out.Containers)
		return
	}, nil
}

// GetPodsByStatus searches for pods whose status is included in statuses
type GetPodsByStatus_methods struct{}

func GetPodsByStatus() GetPodsByStatus_methods { return GetPodsByStatus_methods{} }

func (m GetPodsByStatus_methods) Call(c *varlink.Connection, statuses_in_ []string) (pods_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, statuses_in_)
	if err_ != nil {
		return
	}
	pods_out_, _, err_ = receive()
	return
}

func (m GetPodsByStatus_methods) Send(c *varlink.Connection, flags uint64, statuses_in_ []string) (func() ([]string, uint64, error), error) {
	var in struct {
		Statuses []string `json:"statuses"`
	}
	in.Statuses = []string(statuses_in_)
	receive, err := c.Send("io.podman.GetPodsByStatus", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pods_out_ []string, flags uint64, err error) {
		var out struct {
			Pods []string `json:"pods"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pods_out_ = []string(out.Pods)
		return
	}, nil
}

// ImageExists talks a full or partial image ID or name and returns an int as to whether
// the image exists in local storage. An int result of 0 means the image does exist in
// local storage; whereas 1 indicates the image does not exists in local storage.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ImageExists '{"name": "imageddoesntexist"}'
// {
//   "exists": 1
// }
// ~~~
type ImageExists_methods struct{}

func ImageExists() ImageExists_methods { return ImageExists_methods{} }

func (m ImageExists_methods) Call(c *varlink.Connection, name_in_ string) (exists_out_ int64, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	exists_out_, _, err_ = receive()
	return
}

func (m ImageExists_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (int64, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.ImageExists", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (exists_out_ int64, flags uint64, err error) {
		var out struct {
			Exists int64 `json:"exists"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		exists_out_ = out.Exists
		return
	}, nil
}

// ContainerExists takes a full or partial container ID or name and returns an int as to
// whether the container exists in local storage.  A result of 0 means the container does
// exists; whereas a result of 1 means it could not be found.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ContainerExists '{"name": "flamboyant_payne"}'{
//   "exists": 0
// }
// ~~~
type ContainerExists_methods struct{}

func ContainerExists() ContainerExists_methods { return ContainerExists_methods{} }

func (m ContainerExists_methods) Call(c *varlink.Connection, name_in_ string) (exists_out_ int64, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	exists_out_, _, err_ = receive()
	return
}

func (m ContainerExists_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (int64, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.ContainerExists", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (exists_out_ int64, flags uint64, err error) {
		var out struct {
			Exists int64 `json:"exists"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		exists_out_ = out.Exists
		return
	}, nil
}

// ContainerCheckPoint performs a checkpopint on a container by its name or full/partial container
// ID.  On successful checkpoint, the id of the checkpointed container is returned.
type ContainerCheckpoint_methods struct{}

func ContainerCheckpoint() ContainerCheckpoint_methods { return ContainerCheckpoint_methods{} }

func (m ContainerCheckpoint_methods) Call(c *varlink.Connection, name_in_ string, keep_in_ bool, leaveRunning_in_ bool, tcpEstablished_in_ bool) (id_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, keep_in_, leaveRunning_in_, tcpEstablished_in_)
	if err_ != nil {
		return
	}
	id_out_, _, err_ = receive()
	return
}

func (m ContainerCheckpoint_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, keep_in_ bool, leaveRunning_in_ bool, tcpEstablished_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Name           string `json:"name"`
		Keep           bool   `json:"keep"`
		LeaveRunning   bool   `json:"leaveRunning"`
		TcpEstablished bool   `json:"tcpEstablished"`
	}
	in.Name = name_in_
	in.Keep = keep_in_
	in.LeaveRunning = leaveRunning_in_
	in.TcpEstablished = tcpEstablished_in_
	receive, err := c.Send("io.podman.ContainerCheckpoint", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (id_out_ string, flags uint64, err error) {
		var out struct {
			Id string `json:"id"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		id_out_ = out.Id
		return
	}, nil
}

// ContainerRestore restores a container that has been checkpointed.  The container to be restored can
// be identified by its name or full/partial container ID.  A successful restore will result in the return
// of the container's ID.
type ContainerRestore_methods struct{}

func ContainerRestore() ContainerRestore_methods { return ContainerRestore_methods{} }

func (m ContainerRestore_methods) Call(c *varlink.Connection, name_in_ string, keep_in_ bool, tcpEstablished_in_ bool) (id_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, keep_in_, tcpEstablished_in_)
	if err_ != nil {
		return
	}
	id_out_, _, err_ = receive()
	return
}

func (m ContainerRestore_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, keep_in_ bool, tcpEstablished_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Name           string `json:"name"`
		Keep           bool   `json:"keep"`
		TcpEstablished bool   `json:"tcpEstablished"`
	}
	in.Name = name_in_
	in.Keep = keep_in_
	in.TcpEstablished = tcpEstablished_in_
	receive, err := c.Send("io.podman.ContainerRestore", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (id_out_ string, flags uint64, err error) {
		var out struct {
			Id string `json:"id"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		id_out_ = out.Id
		return
	}, nil
}

// ContainerRunlabel runs executes a command as described by a given container image label.
type ContainerRunlabel_methods struct{}

func ContainerRunlabel() ContainerRunlabel_methods { return ContainerRunlabel_methods{} }

func (m ContainerRunlabel_methods) Call(c *varlink.Connection, runlabel_in_ Runlabel) (err_ error) {
	receive, err_ := m.Send(c, 0, runlabel_in_)
	if err_ != nil {
		return
	}
	_, err_ = receive()
	return
}

func (m ContainerRunlabel_methods) Send(c *varlink.Connection, flags uint64, runlabel_in_ Runlabel) (func() (uint64, error), error) {
	var in struct {
		Runlabel Runlabel `json:"runlabel"`
	}
	in.Runlabel = runlabel_in_
	receive, err := c.Send("io.podman.ContainerRunlabel", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (flags uint64, err error) {
		flags, err = receive(nil)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		return
	}, nil
}

// ListContainerMounts gathers all the mounted container mount points and returns them as an array
// of strings
// #### Example
// ~~~
// $ varlink call unix:/run/podman/io.podman/io.podman.ListContainerMounts
// {
//   "mounts": {
//     "04e4c255269ed2545e7f8bd1395a75f7949c50c223415c00c1d54bfa20f3b3d9": "/var/lib/containers/storage/overlay/a078925828f57e20467ca31cfca8a849210d21ec7e5757332b72b6924f441c17/merged",
//     "1d58c319f9e881a644a5122ff84419dccf6d138f744469281446ab243ef38924": "/var/lib/containers/storage/overlay/948fcf93f8cb932f0f03fd52e3180a58627d547192ffe3b88e0013b98ddcd0d2/merged"
//   }
// }
// ~~~
type ListContainerMounts_methods struct{}

func ListContainerMounts() ListContainerMounts_methods { return ListContainerMounts_methods{} }

func (m ListContainerMounts_methods) Call(c *varlink.Connection) (mounts_out_ map[string]string, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	mounts_out_, _, err_ = receive()
	return
}

func (m ListContainerMounts_methods) Send(c *varlink.Connection, flags uint64) (func() (map[string]string, uint64, error), error) {
	receive, err := c.Send("io.podman.ListContainerMounts", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (mounts_out_ map[string]string, flags uint64, err error) {
		var out struct {
			Mounts map[string]string `json:"mounts"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		mounts_out_ = map[string]string(out.Mounts)
		return
	}, nil
}

// MountContainer mounts a container by name or full/partial ID.  Upon a successful mount, the destination
// mount is returned as a string.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.MountContainer '{"name": "jolly_shannon"}'{
//   "path": "/var/lib/containers/storage/overlay/419eeb04e783ea159149ced67d9fcfc15211084d65e894792a96bedfae0470ca/merged"
// }
// ~~~
type MountContainer_methods struct{}

func MountContainer() MountContainer_methods { return MountContainer_methods{} }

func (m MountContainer_methods) Call(c *varlink.Connection, name_in_ string) (path_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	path_out_, _, err_ = receive()
	return
}

func (m MountContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.MountContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (path_out_ string, flags uint64, err error) {
		var out struct {
			Path string `json:"path"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		path_out_ = out.Path
		return
	}, nil
}

// UnmountContainer umounts a container by its name or full/partial container ID.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.UnmountContainer '{"name": "jolly_shannon", "force": false}'
// {}
// ~~~
type UnmountContainer_methods struct{}

func UnmountContainer() UnmountContainer_methods { return UnmountContainer_methods{} }

func (m UnmountContainer_methods) Call(c *varlink.Connection, name_in_ string, force_in_ bool) (err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, force_in_)
	if err_ != nil {
		return
	}
	_, err_ = receive()
	return
}

func (m UnmountContainer_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, force_in_ bool) (func() (uint64, error), error) {
	var in struct {
		Name  string `json:"name"`
		Force bool   `json:"force"`
	}
	in.Name = name_in_
	in.Force = force_in_
	receive, err := c.Send("io.podman.UnmountContainer", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (flags uint64, err error) {
		flags, err = receive(nil)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		return
	}, nil
}

// ImagesPrune removes all unused images from the local store.  Upon successful pruning,
// the IDs of the removed images are returned.
type ImagesPrune_methods struct{}

func ImagesPrune() ImagesPrune_methods { return ImagesPrune_methods{} }

func (m ImagesPrune_methods) Call(c *varlink.Connection, all_in_ bool) (pruned_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, all_in_)
	if err_ != nil {
		return
	}
	pruned_out_, _, err_ = receive()
	return
}

func (m ImagesPrune_methods) Send(c *varlink.Connection, flags uint64, all_in_ bool) (func() ([]string, uint64, error), error) {
	var in struct {
		All bool `json:"all"`
	}
	in.All = all_in_
	receive, err := c.Send("io.podman.ImagesPrune", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pruned_out_ []string, flags uint64, err error) {
		var out struct {
			Pruned []string `json:"pruned"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pruned_out_ = []string(out.Pruned)
		return
	}, nil
}

// GenerateKube generates a Kubernetes v1 Pod description of a Podman container or pod
// and its containers. The description is in YAML.  See also [ReplayKube](ReplayKube).
type GenerateKube_methods struct{}

func GenerateKube() GenerateKube_methods { return GenerateKube_methods{} }

func (m GenerateKube_methods) Call(c *varlink.Connection, name_in_ string, service_in_ bool) (pod_out_ KubePodService, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, service_in_)
	if err_ != nil {
		return
	}
	pod_out_, _, err_ = receive()
	return
}

func (m GenerateKube_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, service_in_ bool) (func() (KubePodService, uint64, error), error) {
	var in struct {
		Name    string `json:"name"`
		Service bool   `json:"service"`
	}
	in.Name = name_in_
	in.Service = service_in_
	receive, err := c.Send("io.podman.GenerateKube", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pod_out_ KubePodService, flags uint64, err error) {
		var out struct {
			Pod KubePodService `json:"pod"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pod_out_ = out.Pod
		return
	}, nil
}

// ContainerConfig returns a container's config in string form. This call is for
// development of Podman only and generally should not be used.
type ContainerConfig_methods struct{}

func ContainerConfig() ContainerConfig_methods { return ContainerConfig_methods{} }

func (m ContainerConfig_methods) Call(c *varlink.Connection, name_in_ string) (config_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	config_out_, _, err_ = receive()
	return
}

func (m ContainerConfig_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.ContainerConfig", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (config_out_ string, flags uint64, err error) {
		var out struct {
			Config string `json:"config"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		config_out_ = out.Config
		return
	}, nil
}

// ContainerArtifacts returns a container's artifacts in string form.  This call is for
// development of Podman only and generally should not be used.
type ContainerArtifacts_methods struct{}

func ContainerArtifacts() ContainerArtifacts_methods { return ContainerArtifacts_methods{} }

func (m ContainerArtifacts_methods) Call(c *varlink.Connection, name_in_ string, artifactName_in_ string) (config_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, artifactName_in_)
	if err_ != nil {
		return
	}
	config_out_, _, err_ = receive()
	return
}

func (m ContainerArtifacts_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, artifactName_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name         string `json:"name"`
		ArtifactName string `json:"artifactName"`
	}
	in.Name = name_in_
	in.ArtifactName = artifactName_in_
	receive, err := c.Send("io.podman.ContainerArtifacts", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (config_out_ string, flags uint64, err error) {
		var out struct {
			Config string `json:"config"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		config_out_ = out.Config
		return
	}, nil
}

// ContainerInspectData returns a container's inspect data in string form.  This call is for
// development of Podman only and generally should not be used.
type ContainerInspectData_methods struct{}

func ContainerInspectData() ContainerInspectData_methods { return ContainerInspectData_methods{} }

func (m ContainerInspectData_methods) Call(c *varlink.Connection, name_in_ string, size_in_ bool) (config_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, size_in_)
	if err_ != nil {
		return
	}
	config_out_, _, err_ = receive()
	return
}

func (m ContainerInspectData_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, size_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
		Size bool   `json:"size"`
	}
	in.Name = name_in_
	in.Size = size_in_
	receive, err := c.Send("io.podman.ContainerInspectData", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (config_out_ string, flags uint64, err error) {
		var out struct {
			Config string `json:"config"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		config_out_ = out.Config
		return
	}, nil
}

// ContainerStateData returns a container's state config in string form.  This call is for
// development of Podman only and generally should not be used.
type ContainerStateData_methods struct{}

func ContainerStateData() ContainerStateData_methods { return ContainerStateData_methods{} }

func (m ContainerStateData_methods) Call(c *varlink.Connection, name_in_ string) (config_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	config_out_, _, err_ = receive()
	return
}

func (m ContainerStateData_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.ContainerStateData", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (config_out_ string, flags uint64, err error) {
		var out struct {
			Config string `json:"config"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		config_out_ = out.Config
		return
	}, nil
}

// PodStateData returns inspectr level information of a given pod in string form.  This call is for
// development of Podman only and generally should not be used.
type PodStateData_methods struct{}

func PodStateData() PodStateData_methods { return PodStateData_methods{} }

func (m PodStateData_methods) Call(c *varlink.Connection, name_in_ string) (config_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	config_out_, _, err_ = receive()
	return
}

func (m PodStateData_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.PodStateData", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (config_out_ string, flags uint64, err error) {
		var out struct {
			Config string `json:"config"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		config_out_ = out.Config
		return
	}, nil
}

// This call is for the development of Podman only and should not be used.
type CreateFromCC_methods struct{}

func CreateFromCC() CreateFromCC_methods { return CreateFromCC_methods{} }

func (m CreateFromCC_methods) Call(c *varlink.Connection, in_in_ []string) (id_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, in_in_)
	if err_ != nil {
		return
	}
	id_out_, _, err_ = receive()
	return
}

func (m CreateFromCC_methods) Send(c *varlink.Connection, flags uint64, in_in_ []string) (func() (string, uint64, error), error) {
	var in struct {
		In []string `json:"in"`
	}
	in.In = []string(in_in_)
	receive, err := c.Send("io.podman.CreateFromCC", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (id_out_ string, flags uint64, err error) {
		var out struct {
			Id string `json:"id"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		id_out_ = out.Id
		return
	}, nil
}

// Spec returns the oci spec for a container.  This call is for development of Podman only and generally should not be used.
type Spec_methods struct{}

func Spec() Spec_methods { return Spec_methods{} }

func (m Spec_methods) Call(c *varlink.Connection, name_in_ string) (config_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	config_out_, _, err_ = receive()
	return
}

func (m Spec_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.Spec", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (config_out_ string, flags uint64, err error) {
		var out struct {
			Config string `json:"config"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		config_out_ = out.Config
		return
	}, nil
}

// Sendfile allows a remote client to send a file to the host
type SendFile_methods struct{}

func SendFile() SendFile_methods { return SendFile_methods{} }

func (m SendFile_methods) Call(c *varlink.Connection, type_in_ string, length_in_ int64) (file_handle_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, type_in_, length_in_)
	if err_ != nil {
		return
	}
	file_handle_out_, _, err_ = receive()
	return
}

func (m SendFile_methods) Send(c *varlink.Connection, flags uint64, type_in_ string, length_in_ int64) (func() (string, uint64, error), error) {
	var in struct {
		Type   string `json:"type"`
		Length int64  `json:"length"`
	}
	in.Type = type_in_
	in.Length = length_in_
	receive, err := c.Send("io.podman.SendFile", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (file_handle_out_ string, flags uint64, err error) {
		var out struct {
			File_handle string `json:"file_handle"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		file_handle_out_ = out.File_handle
		return
	}, nil
}

// ReceiveFile allows the host to send a remote client a file
type ReceiveFile_methods struct{}

func ReceiveFile() ReceiveFile_methods { return ReceiveFile_methods{} }

func (m ReceiveFile_methods) Call(c *varlink.Connection, path_in_ string, delete_in_ bool) (len_out_ int64, err_ error) {
	receive, err_ := m.Send(c, 0, path_in_, delete_in_)
	if err_ != nil {
		return
	}
	len_out_, _, err_ = receive()
	return
}

func (m ReceiveFile_methods) Send(c *varlink.Connection, flags uint64, path_in_ string, delete_in_ bool) (func() (int64, uint64, error), error) {
	var in struct {
		Path   string `json:"path"`
		Delete bool   `json:"delete"`
	}
	in.Path = path_in_
	in.Delete = delete_in_
	receive, err := c.Send("io.podman.ReceiveFile", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (len_out_ int64, flags uint64, err error) {
		var out struct {
			Len int64 `json:"len"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		len_out_ = out.Len
		return
	}, nil
}

// VolumeCreate creates a volume on a remote host
type VolumeCreate_methods struct{}

func VolumeCreate() VolumeCreate_methods { return VolumeCreate_methods{} }

func (m VolumeCreate_methods) Call(c *varlink.Connection, options_in_ VolumeCreateOpts) (volumeName_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, options_in_)
	if err_ != nil {
		return
	}
	volumeName_out_, _, err_ = receive()
	return
}

func (m VolumeCreate_methods) Send(c *varlink.Connection, flags uint64, options_in_ VolumeCreateOpts) (func() (string, uint64, error), error) {
	var in struct {
		Options VolumeCreateOpts `json:"options"`
	}
	in.Options = options_in_
	receive, err := c.Send("io.podman.VolumeCreate", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (volumeName_out_ string, flags uint64, err error) {
		var out struct {
			VolumeName string `json:"volumeName"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		volumeName_out_ = out.VolumeName
		return
	}, nil
}

// VolumeRemove removes a volume on a remote host
type VolumeRemove_methods struct{}

func VolumeRemove() VolumeRemove_methods { return VolumeRemove_methods{} }

func (m VolumeRemove_methods) Call(c *varlink.Connection, options_in_ VolumeRemoveOpts) (volumeNames_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, options_in_)
	if err_ != nil {
		return
	}
	volumeNames_out_, _, err_ = receive()
	return
}

func (m VolumeRemove_methods) Send(c *varlink.Connection, flags uint64, options_in_ VolumeRemoveOpts) (func() ([]string, uint64, error), error) {
	var in struct {
		Options VolumeRemoveOpts `json:"options"`
	}
	in.Options = options_in_
	receive, err := c.Send("io.podman.VolumeRemove", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (volumeNames_out_ []string, flags uint64, err error) {
		var out struct {
			VolumeNames []string `json:"volumeNames"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		volumeNames_out_ = []string(out.VolumeNames)
		return
	}, nil
}

// GetVolumes gets slice of the volumes on a remote host
type GetVolumes_methods struct{}

func GetVolumes() GetVolumes_methods { return GetVolumes_methods{} }

func (m GetVolumes_methods) Call(c *varlink.Connection, args_in_ []string, all_in_ bool) (volumes_out_ []Volume, err_ error) {
	receive, err_ := m.Send(c, 0, args_in_, all_in_)
	if err_ != nil {
		return
	}
	volumes_out_, _, err_ = receive()
	return
}

func (m GetVolumes_methods) Send(c *varlink.Connection, flags uint64, args_in_ []string, all_in_ bool) (func() ([]Volume, uint64, error), error) {
	var in struct {
		Args []string `json:"args"`
		All  bool     `json:"all"`
	}
	in.Args = []string(args_in_)
	in.All = all_in_
	receive, err := c.Send("io.podman.GetVolumes", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (volumes_out_ []Volume, flags uint64, err error) {
		var out struct {
			Volumes []Volume `json:"volumes"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		volumes_out_ = []Volume(out.Volumes)
		return
	}, nil
}

// VolumesPrune removes unused volumes on the host
type VolumesPrune_methods struct{}

func VolumesPrune() VolumesPrune_methods { return VolumesPrune_methods{} }

func (m VolumesPrune_methods) Call(c *varlink.Connection) (prunedNames_out_ []string, prunedErrors_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	prunedNames_out_, prunedErrors_out_, _, err_ = receive()
	return
}

func (m VolumesPrune_methods) Send(c *varlink.Connection, flags uint64) (func() ([]string, []string, uint64, error), error) {
	receive, err := c.Send("io.podman.VolumesPrune", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (prunedNames_out_ []string, prunedErrors_out_ []string, flags uint64, err error) {
		var out struct {
			PrunedNames  []string `json:"prunedNames"`
			PrunedErrors []string `json:"prunedErrors"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		prunedNames_out_ = []string(out.PrunedNames)
		prunedErrors_out_ = []string(out.PrunedErrors)
		return
	}, nil
}

// ImageSave allows you to save an image from the local image storage to a tarball
type ImageSave_methods struct{}

func ImageSave() ImageSave_methods { return ImageSave_methods{} }

func (m ImageSave_methods) Call(c *varlink.Connection, options_in_ ImageSaveOptions) (reply_out_ MoreResponse, err_ error) {
	receive, err_ := m.Send(c, 0, options_in_)
	if err_ != nil {
		return
	}
	reply_out_, _, err_ = receive()
	return
}

func (m ImageSave_methods) Send(c *varlink.Connection, flags uint64, options_in_ ImageSaveOptions) (func() (MoreResponse, uint64, error), error) {
	var in struct {
		Options ImageSaveOptions `json:"options"`
	}
	in.Options = options_in_
	receive, err := c.Send("io.podman.ImageSave", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (reply_out_ MoreResponse, flags uint64, err error) {
		var out struct {
			Reply MoreResponse `json:"reply"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		reply_out_ = out.Reply
		return
	}, nil
}

// GetPodsByContext allows you to get a list pod ids depending on all, latest, or a list of
// pod names.  The definition of latest pod means the latest by creation date.  In a multi-
// user environment, results might differ from what you expect.
type GetPodsByContext_methods struct{}

func GetPodsByContext() GetPodsByContext_methods { return GetPodsByContext_methods{} }

func (m GetPodsByContext_methods) Call(c *varlink.Connection, all_in_ bool, latest_in_ bool, args_in_ []string) (pods_out_ []string, err_ error) {
	receive, err_ := m.Send(c, 0, all_in_, latest_in_, args_in_)
	if err_ != nil {
		return
	}
	pods_out_, _, err_ = receive()
	return
}

func (m GetPodsByContext_methods) Send(c *varlink.Connection, flags uint64, all_in_ bool, latest_in_ bool, args_in_ []string) (func() ([]string, uint64, error), error) {
	var in struct {
		All    bool     `json:"all"`
		Latest bool     `json:"latest"`
		Args   []string `json:"args"`
	}
	in.All = all_in_
	in.Latest = latest_in_
	in.Args = []string(args_in_)
	receive, err := c.Send("io.podman.GetPodsByContext", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (pods_out_ []string, flags uint64, err error) {
		var out struct {
			Pods []string `json:"pods"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		pods_out_ = []string(out.Pods)
		return
	}, nil
}

// LoadImage allows you to load an image into local storage from a tarball.
type LoadImage_methods struct{}

func LoadImage() LoadImage_methods { return LoadImage_methods{} }

func (m LoadImage_methods) Call(c *varlink.Connection, name_in_ string, inputFile_in_ string, quiet_in_ bool, deleteFile_in_ bool) (reply_out_ MoreResponse, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, inputFile_in_, quiet_in_, deleteFile_in_)
	if err_ != nil {
		return
	}
	reply_out_, _, err_ = receive()
	return
}

func (m LoadImage_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, inputFile_in_ string, quiet_in_ bool, deleteFile_in_ bool) (func() (MoreResponse, uint64, error), error) {
	var in struct {
		Name       string `json:"name"`
		InputFile  string `json:"inputFile"`
		Quiet      bool   `json:"quiet"`
		DeleteFile bool   `json:"deleteFile"`
	}
	in.Name = name_in_
	in.InputFile = inputFile_in_
	in.Quiet = quiet_in_
	in.DeleteFile = deleteFile_in_
	receive, err := c.Send("io.podman.LoadImage", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (reply_out_ MoreResponse, flags uint64, err error) {
		var out struct {
			Reply MoreResponse `json:"reply"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		reply_out_ = out.Reply
		return
	}, nil
}

// GetEvents returns known libpod events filtered by the options provided.
type GetEvents_methods struct{}

func GetEvents() GetEvents_methods { return GetEvents_methods{} }

func (m GetEvents_methods) Call(c *varlink.Connection, filter_in_ []string, since_in_ string, until_in_ string) (events_out_ Event, err_ error) {
	receive, err_ := m.Send(c, 0, filter_in_, since_in_, until_in_)
	if err_ != nil {
		return
	}
	events_out_, _, err_ = receive()
	return
}

func (m GetEvents_methods) Send(c *varlink.Connection, flags uint64, filter_in_ []string, since_in_ string, until_in_ string) (func() (Event, uint64, error), error) {
	var in struct {
		Filter []string `json:"filter"`
		Since  string   `json:"since"`
		Until  string   `json:"until"`
	}
	in.Filter = []string(filter_in_)
	in.Since = since_in_
	in.Until = until_in_
	receive, err := c.Send("io.podman.GetEvents", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (events_out_ Event, flags uint64, err error) {
		var out struct {
			Events Event `json:"events"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		events_out_ = out.Events
		return
	}, nil
}

// Diff returns a diff between libpod objects
type Diff_methods struct{}

func Diff() Diff_methods { return Diff_methods{} }

func (m Diff_methods) Call(c *varlink.Connection, name_in_ string) (diffs_out_ []DiffInfo, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	diffs_out_, _, err_ = receive()
	return
}

func (m Diff_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() ([]DiffInfo, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.Diff", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (diffs_out_ []DiffInfo, flags uint64, err error) {
		var out struct {
			Diffs []DiffInfo `json:"diffs"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		diffs_out_ = []DiffInfo(out.Diffs)
		return
	}, nil
}

// GetLayersMapWithImageInfo is for the development of Podman and should not be used.
type GetLayersMapWithImageInfo_methods struct{}

func GetLayersMapWithImageInfo() GetLayersMapWithImageInfo_methods {
	return GetLayersMapWithImageInfo_methods{}
}

func (m GetLayersMapWithImageInfo_methods) Call(c *varlink.Connection) (layerMap_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0)
	if err_ != nil {
		return
	}
	layerMap_out_, _, err_ = receive()
	return
}

func (m GetLayersMapWithImageInfo_methods) Send(c *varlink.Connection, flags uint64) (func() (string, uint64, error), error) {
	receive, err := c.Send("io.podman.GetLayersMapWithImageInfo", nil, flags)
	if err != nil {
		return nil, err
	}
	return func() (layerMap_out_ string, flags uint64, err error) {
		var out struct {
			LayerMap string `json:"layerMap"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		layerMap_out_ = out.LayerMap
		return
	}, nil
}

// BuildImageHierarchyMap is for the development of Podman and should not be used.
type BuildImageHierarchyMap_methods struct{}

func BuildImageHierarchyMap() BuildImageHierarchyMap_methods { return BuildImageHierarchyMap_methods{} }

func (m BuildImageHierarchyMap_methods) Call(c *varlink.Connection, name_in_ string) (imageInfo_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_)
	if err_ != nil {
		return
	}
	imageInfo_out_, _, err_ = receive()
	return
}

func (m BuildImageHierarchyMap_methods) Send(c *varlink.Connection, flags uint64, name_in_ string) (func() (string, uint64, error), error) {
	var in struct {
		Name string `json:"name"`
	}
	in.Name = name_in_
	receive, err := c.Send("io.podman.BuildImageHierarchyMap", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (imageInfo_out_ string, flags uint64, err error) {
		var out struct {
			ImageInfo string `json:"imageInfo"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		imageInfo_out_ = out.ImageInfo
		return
	}, nil
}

type GenerateSystemd_methods struct{}

func GenerateSystemd() GenerateSystemd_methods { return GenerateSystemd_methods{} }

func (m GenerateSystemd_methods) Call(c *varlink.Connection, name_in_ string, restart_in_ string, timeout_in_ int64, useName_in_ bool) (unit_out_ string, err_ error) {
	receive, err_ := m.Send(c, 0, name_in_, restart_in_, timeout_in_, useName_in_)
	if err_ != nil {
		return
	}
	unit_out_, _, err_ = receive()
	return
}

func (m GenerateSystemd_methods) Send(c *varlink.Connection, flags uint64, name_in_ string, restart_in_ string, timeout_in_ int64, useName_in_ bool) (func() (string, uint64, error), error) {
	var in struct {
		Name    string `json:"name"`
		Restart string `json:"restart"`
		Timeout int64  `json:"timeout"`
		UseName bool   `json:"useName"`
	}
	in.Name = name_in_
	in.Restart = restart_in_
	in.Timeout = timeout_in_
	in.UseName = useName_in_
	receive, err := c.Send("io.podman.GenerateSystemd", in, flags)
	if err != nil {
		return nil, err
	}
	return func() (unit_out_ string, flags uint64, err error) {
		var out struct {
			Unit string `json:"unit"`
		}
		flags, err = receive(&out)
		if err != nil {
			err = Dispatch_Error(err)
			return
		}
		unit_out_ = out.Unit
		return
	}, nil
}

// Generated service interface with all methods

type iopodmanInterface interface {
	GetVersion(c VarlinkCall) error
	GetInfo(c VarlinkCall) error
	ListContainers(c VarlinkCall) error
	Ps(c VarlinkCall, opts_ PsOpts) error
	GetContainersByStatus(c VarlinkCall, status_ []string) error
	Top(c VarlinkCall, nameOrID_ string, descriptors_ []string) error
	GetContainer(c VarlinkCall, id_ string) error
	GetContainersByContext(c VarlinkCall, all_ bool, latest_ bool, args_ []string) error
	CreateContainer(c VarlinkCall, create_ Create) error
	InspectContainer(c VarlinkCall, name_ string) error
	ListContainerProcesses(c VarlinkCall, name_ string, opts_ []string) error
	GetContainerLogs(c VarlinkCall, name_ string) error
	GetContainersLogs(c VarlinkCall, names_ []string, follow_ bool, latest_ bool, since_ string, tail_ int64, timestamps_ bool) error
	ListContainerChanges(c VarlinkCall, name_ string) error
	ExportContainer(c VarlinkCall, name_ string, path_ string) error
	GetContainerStats(c VarlinkCall, name_ string) error
	GetContainerStatsWithHistory(c VarlinkCall, previousStats_ ContainerStats) error
	StartContainer(c VarlinkCall, name_ string) error
	StopContainer(c VarlinkCall, name_ string, timeout_ int64) error
	InitContainer(c VarlinkCall, name_ string) error
	RestartContainer(c VarlinkCall, name_ string, timeout_ int64) error
	KillContainer(c VarlinkCall, name_ string, signal_ int64) error
	PauseContainer(c VarlinkCall, name_ string) error
	UnpauseContainer(c VarlinkCall, name_ string) error
	Attach(c VarlinkCall, name_ string, detachKeys_ string, start_ bool) error
	AttachControl(c VarlinkCall, name_ string) error
	GetAttachSockets(c VarlinkCall, name_ string) error
	WaitContainer(c VarlinkCall, name_ string, interval_ int64) error
	RemoveContainer(c VarlinkCall, name_ string, force_ bool, removeVolumes_ bool) error
	DeleteStoppedContainers(c VarlinkCall) error
	ListImages(c VarlinkCall) error
	GetImage(c VarlinkCall, id_ string) error
	BuildImage(c VarlinkCall, build_ BuildInfo) error
	InspectImage(c VarlinkCall, name_ string) error
	HistoryImage(c VarlinkCall, name_ string) error
	PushImage(c VarlinkCall, name_ string, tag_ string, compress_ bool, format_ string, removeSignatures_ bool, signBy_ string) error
	TagImage(c VarlinkCall, name_ string, tagged_ string) error
	RemoveImage(c VarlinkCall, name_ string, force_ bool) error
	SearchImages(c VarlinkCall, query_ string, limit_ *int64, filter_ ImageSearchFilter) error
	DeleteUnusedImages(c VarlinkCall) error
	Commit(c VarlinkCall, name_ string, image_name_ string, changes_ []string, author_ string, message_ string, pause_ bool, manifestType_ string) error
	ImportImage(c VarlinkCall, source_ string, reference_ string, message_ string, changes_ []string, delete_ bool) error
	ExportImage(c VarlinkCall, name_ string, destination_ string, compress_ bool, tags_ []string) error
	PullImage(c VarlinkCall, name_ string) error
	CreatePod(c VarlinkCall, create_ PodCreate) error
	ListPods(c VarlinkCall) error
	GetPod(c VarlinkCall, name_ string) error
	InspectPod(c VarlinkCall, name_ string) error
	StartPod(c VarlinkCall, name_ string) error
	StopPod(c VarlinkCall, name_ string, timeout_ int64) error
	RestartPod(c VarlinkCall, name_ string) error
	KillPod(c VarlinkCall, name_ string, signal_ int64) error
	PausePod(c VarlinkCall, name_ string) error
	UnpausePod(c VarlinkCall, name_ string) error
	RemovePod(c VarlinkCall, name_ string, force_ bool) error
	TopPod(c VarlinkCall, pod_ string, latest_ bool, descriptors_ []string) error
	GetPodStats(c VarlinkCall, name_ string) error
	GetPodsByStatus(c VarlinkCall, statuses_ []string) error
	ImageExists(c VarlinkCall, name_ string) error
	ContainerExists(c VarlinkCall, name_ string) error
	ContainerCheckpoint(c VarlinkCall, name_ string, keep_ bool, leaveRunning_ bool, tcpEstablished_ bool) error
	ContainerRestore(c VarlinkCall, name_ string, keep_ bool, tcpEstablished_ bool) error
	ContainerRunlabel(c VarlinkCall, runlabel_ Runlabel) error
	ListContainerMounts(c VarlinkCall) error
	MountContainer(c VarlinkCall, name_ string) error
	UnmountContainer(c VarlinkCall, name_ string, force_ bool) error
	ImagesPrune(c VarlinkCall, all_ bool) error
	GenerateKube(c VarlinkCall, name_ string, service_ bool) error
	ContainerConfig(c VarlinkCall, name_ string) error
	ContainerArtifacts(c VarlinkCall, name_ string, artifactName_ string) error
	ContainerInspectData(c VarlinkCall, name_ string, size_ bool) error
	ContainerStateData(c VarlinkCall, name_ string) error
	PodStateData(c VarlinkCall, name_ string) error
	CreateFromCC(c VarlinkCall, in_ []string) error
	Spec(c VarlinkCall, name_ string) error
	SendFile(c VarlinkCall, type_ string, length_ int64) error
	ReceiveFile(c VarlinkCall, path_ string, delete_ bool) error
	VolumeCreate(c VarlinkCall, options_ VolumeCreateOpts) error
	VolumeRemove(c VarlinkCall, options_ VolumeRemoveOpts) error
	GetVolumes(c VarlinkCall, args_ []string, all_ bool) error
	VolumesPrune(c VarlinkCall) error
	ImageSave(c VarlinkCall, options_ ImageSaveOptions) error
	GetPodsByContext(c VarlinkCall, all_ bool, latest_ bool, args_ []string) error
	LoadImage(c VarlinkCall, name_ string, inputFile_ string, quiet_ bool, deleteFile_ bool) error
	GetEvents(c VarlinkCall, filter_ []string, since_ string, until_ string) error
	Diff(c VarlinkCall, name_ string) error
	GetLayersMapWithImageInfo(c VarlinkCall) error
	BuildImageHierarchyMap(c VarlinkCall, name_ string) error
	GenerateSystemd(c VarlinkCall, name_ string, restart_ string, timeout_ int64, useName_ bool) error
}

// Generated service object with all methods

type VarlinkCall struct{ varlink.Call }

// Generated reply methods for all varlink errors

// ImageNotFound means the image could not be found by the provided name or ID in local storage.
func (c *VarlinkCall) ReplyImageNotFound(id_ string, reason_ string) error {
	var out ImageNotFound
	out.Id = id_
	out.Reason = reason_
	return c.ReplyError("io.podman.ImageNotFound", &out)
}

// ContainerNotFound means the container could not be found by the provided name or ID in local storage.
func (c *VarlinkCall) ReplyContainerNotFound(id_ string, reason_ string) error {
	var out ContainerNotFound
	out.Id = id_
	out.Reason = reason_
	return c.ReplyError("io.podman.ContainerNotFound", &out)
}

// NoContainerRunning means none of the containers requested are running in a command that requires a running container.
func (c *VarlinkCall) ReplyNoContainerRunning() error {
	var out NoContainerRunning
	return c.ReplyError("io.podman.NoContainerRunning", &out)
}

// PodNotFound means the pod could not be found by the provided name or ID in local storage.
func (c *VarlinkCall) ReplyPodNotFound(name_ string, reason_ string) error {
	var out PodNotFound
	out.Name = name_
	out.Reason = reason_
	return c.ReplyError("io.podman.PodNotFound", &out)
}

// VolumeNotFound means the volume could not be found by the name or ID in local storage.
func (c *VarlinkCall) ReplyVolumeNotFound(id_ string, reason_ string) error {
	var out VolumeNotFound
	out.Id = id_
	out.Reason = reason_
	return c.ReplyError("io.podman.VolumeNotFound", &out)
}

// PodContainerError means a container associated with a pod failed to perform an operation. It contains
// a container ID of the container that failed.
func (c *VarlinkCall) ReplyPodContainerError(podname_ string, errors_ []PodContainerErrorData) error {
	var out PodContainerError
	out.Podname = podname_
	out.Errors = []PodContainerErrorData(errors_)
	return c.ReplyError("io.podman.PodContainerError", &out)
}

// NoContainersInPod means a pod has no containers on which to perform the operation. It contains
// the pod ID.
func (c *VarlinkCall) ReplyNoContainersInPod(name_ string) error {
	var out NoContainersInPod
	out.Name = name_
	return c.ReplyError("io.podman.NoContainersInPod", &out)
}

// InvalidState indicates that a container or pod was in an improper state for the requested operation
func (c *VarlinkCall) ReplyInvalidState(id_ string, reason_ string) error {
	var out InvalidState
	out.Id = id_
	out.Reason = reason_
	return c.ReplyError("io.podman.InvalidState", &out)
}

// ErrorOccurred is a generic error for an error that occurs during the execution.  The actual error message
// is includes as part of the error's text.
func (c *VarlinkCall) ReplyErrorOccurred(reason_ string) error {
	var out ErrorOccurred
	out.Reason = reason_
	return c.ReplyError("io.podman.ErrorOccurred", &out)
}

// RuntimeErrors generally means a runtime could not be found or gotten.
func (c *VarlinkCall) ReplyRuntimeError(reason_ string) error {
	var out RuntimeError
	out.Reason = reason_
	return c.ReplyError("io.podman.RuntimeError", &out)
}

// The Podman endpoint requires that you use a streaming connection.
func (c *VarlinkCall) ReplyWantsMoreRequired(reason_ string) error {
	var out WantsMoreRequired
	out.Reason = reason_
	return c.ReplyError("io.podman.WantsMoreRequired", &out)
}

// Container is already stopped
func (c *VarlinkCall) ReplyErrCtrStopped(id_ string) error {
	var out ErrCtrStopped
	out.Id = id_
	return c.ReplyError("io.podman.ErrCtrStopped", &out)
}

// Generated reply methods for all varlink methods

func (c *VarlinkCall) ReplyGetVersion(version_ string, go_version_ string, git_commit_ string, built_ string, os_arch_ string, remote_api_version_ int64) error {
	var out struct {
		Version            string `json:"version"`
		Go_version         string `json:"go_version"`
		Git_commit         string `json:"git_commit"`
		Built              string `json:"built"`
		Os_arch            string `json:"os_arch"`
		Remote_api_version int64  `json:"remote_api_version"`
	}
	out.Version = version_
	out.Go_version = go_version_
	out.Git_commit = git_commit_
	out.Built = built_
	out.Os_arch = os_arch_
	out.Remote_api_version = remote_api_version_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetInfo(info_ PodmanInfo) error {
	var out struct {
		Info PodmanInfo `json:"info"`
	}
	out.Info = info_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyListContainers(containers_ []Container) error {
	var out struct {
		Containers []Container `json:"containers"`
	}
	out.Containers = []Container(containers_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyPs(containers_ []PsContainer) error {
	var out struct {
		Containers []PsContainer `json:"containers"`
	}
	out.Containers = []PsContainer(containers_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetContainersByStatus(containerS_ []Container) error {
	var out struct {
		ContainerS []Container `json:"containerS"`
	}
	out.ContainerS = []Container(containerS_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyTop(top_ []string) error {
	var out struct {
		Top []string `json:"top"`
	}
	out.Top = []string(top_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetContainer(container_ Container) error {
	var out struct {
		Container Container `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetContainersByContext(containers_ []string) error {
	var out struct {
		Containers []string `json:"containers"`
	}
	out.Containers = []string(containers_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyCreateContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyInspectContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyListContainerProcesses(container_ []string) error {
	var out struct {
		Container []string `json:"container"`
	}
	out.Container = []string(container_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetContainerLogs(container_ []string) error {
	var out struct {
		Container []string `json:"container"`
	}
	out.Container = []string(container_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetContainersLogs(log_ LogLine) error {
	var out struct {
		Log LogLine `json:"log"`
	}
	out.Log = log_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyListContainerChanges(container_ ContainerChanges) error {
	var out struct {
		Container ContainerChanges `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyExportContainer(tarfile_ string) error {
	var out struct {
		Tarfile string `json:"tarfile"`
	}
	out.Tarfile = tarfile_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetContainerStats(container_ ContainerStats) error {
	var out struct {
		Container ContainerStats `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetContainerStatsWithHistory(container_ ContainerStats) error {
	var out struct {
		Container ContainerStats `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyStartContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyStopContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyInitContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyRestartContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyKillContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyPauseContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyUnpauseContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyAttach() error {
	return c.Reply(nil)
}

func (c *VarlinkCall) ReplyAttachControl() error {
	return c.Reply(nil)
}

func (c *VarlinkCall) ReplyGetAttachSockets(sockets_ Sockets) error {
	var out struct {
		Sockets Sockets `json:"sockets"`
	}
	out.Sockets = sockets_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyWaitContainer(exitcode_ int64) error {
	var out struct {
		Exitcode int64 `json:"exitcode"`
	}
	out.Exitcode = exitcode_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyRemoveContainer(container_ string) error {
	var out struct {
		Container string `json:"container"`
	}
	out.Container = container_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyDeleteStoppedContainers(containers_ []string) error {
	var out struct {
		Containers []string `json:"containers"`
	}
	out.Containers = []string(containers_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyListImages(images_ []Image) error {
	var out struct {
		Images []Image `json:"images"`
	}
	out.Images = []Image(images_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetImage(image_ Image) error {
	var out struct {
		Image Image `json:"image"`
	}
	out.Image = image_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyBuildImage(image_ MoreResponse) error {
	var out struct {
		Image MoreResponse `json:"image"`
	}
	out.Image = image_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyInspectImage(image_ string) error {
	var out struct {
		Image string `json:"image"`
	}
	out.Image = image_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyHistoryImage(history_ []ImageHistory) error {
	var out struct {
		History []ImageHistory `json:"history"`
	}
	out.History = []ImageHistory(history_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyPushImage(reply_ MoreResponse) error {
	var out struct {
		Reply MoreResponse `json:"reply"`
	}
	out.Reply = reply_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyTagImage(image_ string) error {
	var out struct {
		Image string `json:"image"`
	}
	out.Image = image_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyRemoveImage(image_ string) error {
	var out struct {
		Image string `json:"image"`
	}
	out.Image = image_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplySearchImages(results_ []ImageSearchResult) error {
	var out struct {
		Results []ImageSearchResult `json:"results"`
	}
	out.Results = []ImageSearchResult(results_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyDeleteUnusedImages(images_ []string) error {
	var out struct {
		Images []string `json:"images"`
	}
	out.Images = []string(images_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyCommit(reply_ MoreResponse) error {
	var out struct {
		Reply MoreResponse `json:"reply"`
	}
	out.Reply = reply_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyImportImage(image_ string) error {
	var out struct {
		Image string `json:"image"`
	}
	out.Image = image_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyExportImage(image_ string) error {
	var out struct {
		Image string `json:"image"`
	}
	out.Image = image_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyPullImage(reply_ MoreResponse) error {
	var out struct {
		Reply MoreResponse `json:"reply"`
	}
	out.Reply = reply_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyCreatePod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyListPods(pods_ []ListPodData) error {
	var out struct {
		Pods []ListPodData `json:"pods"`
	}
	out.Pods = []ListPodData(pods_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetPod(pod_ ListPodData) error {
	var out struct {
		Pod ListPodData `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyInspectPod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyStartPod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyStopPod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyRestartPod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyKillPod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyPausePod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyUnpausePod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyRemovePod(pod_ string) error {
	var out struct {
		Pod string `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyTopPod(stats_ []string) error {
	var out struct {
		Stats []string `json:"stats"`
	}
	out.Stats = []string(stats_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetPodStats(pod_ string, containers_ []ContainerStats) error {
	var out struct {
		Pod        string           `json:"pod"`
		Containers []ContainerStats `json:"containers"`
	}
	out.Pod = pod_
	out.Containers = []ContainerStats(containers_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetPodsByStatus(pods_ []string) error {
	var out struct {
		Pods []string `json:"pods"`
	}
	out.Pods = []string(pods_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyImageExists(exists_ int64) error {
	var out struct {
		Exists int64 `json:"exists"`
	}
	out.Exists = exists_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerExists(exists_ int64) error {
	var out struct {
		Exists int64 `json:"exists"`
	}
	out.Exists = exists_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerCheckpoint(id_ string) error {
	var out struct {
		Id string `json:"id"`
	}
	out.Id = id_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerRestore(id_ string) error {
	var out struct {
		Id string `json:"id"`
	}
	out.Id = id_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerRunlabel() error {
	return c.Reply(nil)
}

func (c *VarlinkCall) ReplyListContainerMounts(mounts_ map[string]string) error {
	var out struct {
		Mounts map[string]string `json:"mounts"`
	}
	out.Mounts = map[string]string(mounts_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyMountContainer(path_ string) error {
	var out struct {
		Path string `json:"path"`
	}
	out.Path = path_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyUnmountContainer() error {
	return c.Reply(nil)
}

func (c *VarlinkCall) ReplyImagesPrune(pruned_ []string) error {
	var out struct {
		Pruned []string `json:"pruned"`
	}
	out.Pruned = []string(pruned_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGenerateKube(pod_ KubePodService) error {
	var out struct {
		Pod KubePodService `json:"pod"`
	}
	out.Pod = pod_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerConfig(config_ string) error {
	var out struct {
		Config string `json:"config"`
	}
	out.Config = config_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerArtifacts(config_ string) error {
	var out struct {
		Config string `json:"config"`
	}
	out.Config = config_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerInspectData(config_ string) error {
	var out struct {
		Config string `json:"config"`
	}
	out.Config = config_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyContainerStateData(config_ string) error {
	var out struct {
		Config string `json:"config"`
	}
	out.Config = config_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyPodStateData(config_ string) error {
	var out struct {
		Config string `json:"config"`
	}
	out.Config = config_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyCreateFromCC(id_ string) error {
	var out struct {
		Id string `json:"id"`
	}
	out.Id = id_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplySpec(config_ string) error {
	var out struct {
		Config string `json:"config"`
	}
	out.Config = config_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplySendFile(file_handle_ string) error {
	var out struct {
		File_handle string `json:"file_handle"`
	}
	out.File_handle = file_handle_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyReceiveFile(len_ int64) error {
	var out struct {
		Len int64 `json:"len"`
	}
	out.Len = len_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyVolumeCreate(volumeName_ string) error {
	var out struct {
		VolumeName string `json:"volumeName"`
	}
	out.VolumeName = volumeName_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyVolumeRemove(volumeNames_ []string) error {
	var out struct {
		VolumeNames []string `json:"volumeNames"`
	}
	out.VolumeNames = []string(volumeNames_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetVolumes(volumes_ []Volume) error {
	var out struct {
		Volumes []Volume `json:"volumes"`
	}
	out.Volumes = []Volume(volumes_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyVolumesPrune(prunedNames_ []string, prunedErrors_ []string) error {
	var out struct {
		PrunedNames  []string `json:"prunedNames"`
		PrunedErrors []string `json:"prunedErrors"`
	}
	out.PrunedNames = []string(prunedNames_)
	out.PrunedErrors = []string(prunedErrors_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyImageSave(reply_ MoreResponse) error {
	var out struct {
		Reply MoreResponse `json:"reply"`
	}
	out.Reply = reply_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetPodsByContext(pods_ []string) error {
	var out struct {
		Pods []string `json:"pods"`
	}
	out.Pods = []string(pods_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyLoadImage(reply_ MoreResponse) error {
	var out struct {
		Reply MoreResponse `json:"reply"`
	}
	out.Reply = reply_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetEvents(events_ Event) error {
	var out struct {
		Events Event `json:"events"`
	}
	out.Events = events_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyDiff(diffs_ []DiffInfo) error {
	var out struct {
		Diffs []DiffInfo `json:"diffs"`
	}
	out.Diffs = []DiffInfo(diffs_)
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGetLayersMapWithImageInfo(layerMap_ string) error {
	var out struct {
		LayerMap string `json:"layerMap"`
	}
	out.LayerMap = layerMap_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyBuildImageHierarchyMap(imageInfo_ string) error {
	var out struct {
		ImageInfo string `json:"imageInfo"`
	}
	out.ImageInfo = imageInfo_
	return c.Reply(&out)
}

func (c *VarlinkCall) ReplyGenerateSystemd(unit_ string) error {
	var out struct {
		Unit string `json:"unit"`
	}
	out.Unit = unit_
	return c.Reply(&out)
}

// Generated dummy implementations for all varlink methods

// GetVersion returns version and build information of the podman service
func (s *VarlinkInterface) GetVersion(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.GetVersion")
}

// GetInfo returns a [PodmanInfo](#PodmanInfo) struct that describes podman and its host such as storage stats,
// build information of Podman, and system-wide registries.
func (s *VarlinkInterface) GetInfo(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.GetInfo")
}

// ListContainers returns information about all containers.
// See also [GetContainer](#GetContainer).
func (s *VarlinkInterface) ListContainers(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.ListContainers")
}

func (s *VarlinkInterface) Ps(c VarlinkCall, opts_ PsOpts) error {
	return c.ReplyMethodNotImplemented("io.podman.Ps")
}

func (s *VarlinkInterface) GetContainersByStatus(c VarlinkCall, status_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetContainersByStatus")
}

func (s *VarlinkInterface) Top(c VarlinkCall, nameOrID_ string, descriptors_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.Top")
}

// GetContainer returns information about a single container.  If a container
// with the given id doesn't exist, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.  See also [ListContainers](ListContainers) and
// [InspectContainer](#InspectContainer).
func (s *VarlinkInterface) GetContainer(c VarlinkCall, id_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetContainer")
}

// GetContainersByContext allows you to get a list of container ids depending on all, latest, or a list of
// container names.  The definition of latest container means the latest by creation date.  In a multi-
// user environment, results might differ from what you expect.
func (s *VarlinkInterface) GetContainersByContext(c VarlinkCall, all_ bool, latest_ bool, args_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetContainersByContext")
}

// CreateContainer creates a new container from an image.  It uses a [Create](#Create) type for input.
func (s *VarlinkInterface) CreateContainer(c VarlinkCall, create_ Create) error {
	return c.ReplyMethodNotImplemented("io.podman.CreateContainer")
}

// InspectContainer data takes a name or ID of a container returns the inspection
// data in string format.  You can then serialize the string into JSON.  A [ContainerNotFound](#ContainerNotFound)
// error will be returned if the container cannot be found. See also [InspectImage](#InspectImage).
func (s *VarlinkInterface) InspectContainer(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.InspectContainer")
}

// ListContainerProcesses takes a name or ID of a container and returns the processes
// running inside the container as array of strings.  It will accept an array of string
// arguments that represent ps options.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ListContainerProcesses '{"name": "135d71b9495f", "opts": []}'
// {
//   "container": [
//     "  UID   PID  PPID  C STIME TTY          TIME CMD",
//     "    0 21220 21210  0 09:05 pts/0    00:00:00 /bin/sh",
//     "    0 21232 21220  0 09:05 pts/0    00:00:00 top",
//     "    0 21284 21220  0 09:05 pts/0    00:00:00 vi /etc/hosts"
//   ]
// }
// ~~~
func (s *VarlinkInterface) ListContainerProcesses(c VarlinkCall, name_ string, opts_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.ListContainerProcesses")
}

// GetContainerLogs takes a name or ID of a container and returns the logs of that container.
// If the container cannot be found, a [ContainerNotFound](#ContainerNotFound) error will be returned.
// The container logs are returned as an array of strings.  GetContainerLogs will honor the streaming
// capability of varlink if the client invokes it.
func (s *VarlinkInterface) GetContainerLogs(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetContainerLogs")
}

func (s *VarlinkInterface) GetContainersLogs(c VarlinkCall, names_ []string, follow_ bool, latest_ bool, since_ string, tail_ int64, timestamps_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.GetContainersLogs")
}

// ListContainerChanges takes a name or ID of a container and returns changes between the container and
// its base image. It returns a struct of changed, deleted, and added path names.
func (s *VarlinkInterface) ListContainerChanges(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.ListContainerChanges")
}

// ExportContainer creates an image from a container.  It takes the name or ID of a container and a
// path representing the target tarfile.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.
// The return value is the written tarfile.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ExportContainer '{"name": "flamboyant_payne", "path": "/tmp/payne.tar" }'
// {
//   "tarfile": "/tmp/payne.tar"
// }
// ~~~
func (s *VarlinkInterface) ExportContainer(c VarlinkCall, name_ string, path_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.ExportContainer")
}

// GetContainerStats takes the name or ID of a container and returns a single ContainerStats structure which
// contains attributes like memory and cpu usage.  If the container cannot be found, a
// [ContainerNotFound](#ContainerNotFound) error will be returned. If the container is not running, a [NoContainerRunning](#NoContainerRunning)
// error will be returned
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.GetContainerStats '{"name": "c33e4164f384"}'
// {
//   "container": {
//     "block_input": 0,
//     "block_output": 0,
//     "cpu": 2.571123918839990154678e-08,
//     "cpu_nano": 49037378,
//     "id": "c33e4164f384aa9d979072a63319d66b74fd7a128be71fa68ede24f33ec6cfee",
//     "mem_limit": 33080606720,
//     "mem_perc": 2.166828456524753747370e-03,
//     "mem_usage": 716800,
//     "name": "competent_wozniak",
//     "net_input": 768,
//     "net_output": 5910,
//     "pids": 1,
//     "system_nano": 10000000
//   }
// }
// ~~~
func (s *VarlinkInterface) GetContainerStats(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetContainerStats")
}

// GetContainerStatsWithHistory takes a previous set of container statistics and uses libpod functions
// to calculate the containers statistics based on current and previous measurements.
func (s *VarlinkInterface) GetContainerStatsWithHistory(c VarlinkCall, previousStats_ ContainerStats) error {
	return c.ReplyMethodNotImplemented("io.podman.GetContainerStatsWithHistory")
}

// StartContainer starts a created or stopped container. It takes the name or ID of container.  It returns
// the container ID once started.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.  See also [CreateContainer](#CreateContainer).
func (s *VarlinkInterface) StartContainer(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.StartContainer")
}

// StopContainer stops a container given a timeout.  It takes the name or ID of a container as well as a
// timeout value.  The timeout value the time before a forcible stop to the container is applied.  It
// returns the container ID once stopped. If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned instead. See also [KillContainer](KillContainer).
// #### Error
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.StopContainer '{"name": "135d71b9495f", "timeout": 5}'
// {
//   "container": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
// }
// ~~~
func (s *VarlinkInterface) StopContainer(c VarlinkCall, name_ string, timeout_ int64) error {
	return c.ReplyMethodNotImplemented("io.podman.StopContainer")
}

// InitContainer initializes the given container. It accepts a container name or
// ID, and will initialize the container matching that ID if possible, and error
// if not. Containers can only be initialized when they are in the Created or
// Exited states. Initialization prepares a container to be started, but does not
// start the container. It is intended to be used to debug a container's state
// prior to starting it.
func (s *VarlinkInterface) InitContainer(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.InitContainer")
}

// RestartContainer will restart a running container given a container name or ID and timeout value. The timeout
// value is the time before a forcible stop is used to stop the container.  If the container cannot be found by
// name or ID, a [ContainerNotFound](#ContainerNotFound)  error will be returned; otherwise, the ID of the
// container will be returned.
func (s *VarlinkInterface) RestartContainer(c VarlinkCall, name_ string, timeout_ int64) error {
	return c.ReplyMethodNotImplemented("io.podman.RestartContainer")
}

// KillContainer takes the name or ID of a container as well as a signal to be applied to the container.  Once the
// container has been killed, the container's ID is returned.  If the container cannot be found, a
// [ContainerNotFound](#ContainerNotFound) error is returned. See also [StopContainer](StopContainer).
func (s *VarlinkInterface) KillContainer(c VarlinkCall, name_ string, signal_ int64) error {
	return c.ReplyMethodNotImplemented("io.podman.KillContainer")
}

// PauseContainer takes the name or ID of container and pauses it.  If the container cannot be found,
// a [ContainerNotFound](#ContainerNotFound) error will be returned; otherwise the ID of the container is returned.
// See also [UnpauseContainer](#UnpauseContainer).
func (s *VarlinkInterface) PauseContainer(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.PauseContainer")
}

// UnpauseContainer takes the name or ID of container and unpauses a paused container.  If the container cannot be
// found, a [ContainerNotFound](#ContainerNotFound) error will be returned; otherwise the ID of the container is returned.
// See also [PauseContainer](#PauseContainer).
func (s *VarlinkInterface) UnpauseContainer(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.UnpauseContainer")
}

// Attach takes the name or ID of a container and sets up the ability to remotely attach to its console. The start
// bool is whether you wish to start the container in question first.
func (s *VarlinkInterface) Attach(c VarlinkCall, name_ string, detachKeys_ string, start_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.Attach")
}

func (s *VarlinkInterface) AttachControl(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.AttachControl")
}

// GetAttachSockets takes the name or ID of an existing container.  It returns file paths for two sockets needed
// to properly communicate with a container.  The first is the actual I/O socket that the container uses.  The
// second is a "control" socket where things like resizing the TTY events are sent. If the container cannot be
// found, a [ContainerNotFound](#ContainerNotFound) error will be returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/io.podman/io.podman.GetAttachSockets '{"name": "b7624e775431219161"}'
// {
//   "sockets": {
//     "container_id": "b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca",
//     "control_socket": "/var/lib/containers/storage/overlay-containers/b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca/userdata/ctl",
//     "io_socket": "/var/run/libpod/socket/b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca/attach"
//   }
// }
// ~~~
func (s *VarlinkInterface) GetAttachSockets(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetAttachSockets")
}

// WaitContainer takes the name or ID of a container and waits the given interval in milliseconds until the container
// stops.  Upon stopping, the return code of the container is returned. If the container container cannot be found by ID
// or name, a [ContainerNotFound](#ContainerNotFound) error is returned.
func (s *VarlinkInterface) WaitContainer(c VarlinkCall, name_ string, interval_ int64) error {
	return c.ReplyMethodNotImplemented("io.podman.WaitContainer")
}

// RemoveContainer requires the name or ID of container as well a boolean representing whether a running container can be stopped and removed, and a boolean
// indicating whether to remove builtin volumes. Upon successful removal of the
// container, its ID is returned.  If the
// container cannot be found by name or ID, a [ContainerNotFound](#ContainerNotFound) error will be returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.RemoveContainer '{"name": "62f4fd98cb57"}'
// {
//   "container": "62f4fd98cb57f529831e8f90610e54bba74bd6f02920ffb485e15376ed365c20"
// }
// ~~~
func (s *VarlinkInterface) RemoveContainer(c VarlinkCall, name_ string, force_ bool, removeVolumes_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.RemoveContainer")
}

// DeleteStoppedContainers will delete all containers that are not running. It will return a list the deleted
// container IDs.  See also [RemoveContainer](RemoveContainer).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.DeleteStoppedContainers
// {
//   "containers": [
//     "451410b931d00def8aa9b4f8084e4d4a39e5e04ea61f358cf53a5cf95afcdcee",
//     "8b60f754a3e01389494a9581ade97d35c2765b6e2f19acd2d3040c82a32d1bc0",
//     "cf2e99d4d3cad6073df199ed32bbe64b124f3e1aba6d78821aa8460e70d30084",
//     "db901a329587312366e5ecff583d08f0875b4b79294322df67d90fc6eed08fc1"
//   ]
// }
// ~~~
func (s *VarlinkInterface) DeleteStoppedContainers(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.DeleteStoppedContainers")
}

// ListImages returns information about the images that are currently in storage.
// See also [InspectImage](#InspectImage).
func (s *VarlinkInterface) ListImages(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.ListImages")
}

// GetImage returns information about a single image in storage.
// If the image caGetImage returns be found, [ImageNotFound](#ImageNotFound) will be returned.
func (s *VarlinkInterface) GetImage(c VarlinkCall, id_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetImage")
}

// BuildImage takes a [BuildInfo](#BuildInfo) structure and builds an image.  At a minimum, you must provide the
// 'dockerfile' and 'tags' options in the BuildInfo structure. It will return a [MoreResponse](#MoreResponse) structure
// that contains the build logs and resulting image ID.
func (s *VarlinkInterface) BuildImage(c VarlinkCall, build_ BuildInfo) error {
	return c.ReplyMethodNotImplemented("io.podman.BuildImage")
}

// InspectImage takes the name or ID of an image and returns a string representation of data associated with the
// mage.  You must serialize the string into JSON to use it further.  An [ImageNotFound](#ImageNotFound) error will
// be returned if the image cannot be found.
func (s *VarlinkInterface) InspectImage(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.InspectImage")
}

// HistoryImage takes the name or ID of an image and returns information about its history and layers.  The returned
// history is in the form of an array of ImageHistory structures.  If the image cannot be found, an
// [ImageNotFound](#ImageNotFound) error is returned.
func (s *VarlinkInterface) HistoryImage(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.HistoryImage")
}

// PushImage takes two input arguments: the name or ID of an image, the fully-qualified destination name of the image,
// It will return an [ImageNotFound](#ImageNotFound) error if
// the image cannot be found in local storage; otherwise it will return a [MoreResponse](#MoreResponse)
func (s *VarlinkInterface) PushImage(c VarlinkCall, name_ string, tag_ string, compress_ bool, format_ string, removeSignatures_ bool, signBy_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.PushImage")
}

// TagImage takes the name or ID of an image in local storage as well as the desired tag name.  If the image cannot
// be found, an [ImageNotFound](#ImageNotFound) error will be returned; otherwise, the ID of the image is returned on success.
func (s *VarlinkInterface) TagImage(c VarlinkCall, name_ string, tagged_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.TagImage")
}

// RemoveImage takes the name or ID of an image as well as a boolean that determines if containers using that image
// should be deleted.  If the image cannot be found, an [ImageNotFound](#ImageNotFound) error will be returned.  The
// ID of the removed image is returned when complete.  See also [DeleteUnusedImages](DeleteUnusedImages).
// #### Example
// ~~~
// varlink call -m unix:/run/podman/io.podman/io.podman.RemoveImage '{"name": "registry.fedoraproject.org/fedora", "force": true}'
// {
//   "image": "426866d6fa419873f97e5cbd320eeb22778244c1dfffa01c944db3114f55772e"
// }
// ~~~
func (s *VarlinkInterface) RemoveImage(c VarlinkCall, name_ string, force_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.RemoveImage")
}

// SearchImages searches available registries for images that contain the
// contents of "query" in their name. If "limit" is given, limits the amount of
// search results per registry.
func (s *VarlinkInterface) SearchImages(c VarlinkCall, query_ string, limit_ *int64, filter_ ImageSearchFilter) error {
	return c.ReplyMethodNotImplemented("io.podman.SearchImages")
}

// DeleteUnusedImages deletes any images not associated with a container.  The IDs of the deleted images are returned
// in a string array.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.DeleteUnusedImages
// {
//   "images": [
//     "166ea6588079559c724c15223f52927f514f73dd5c5cf2ae2d143e3b2e6e9b52",
//     "da86e6ba6ca197bf6bc5e9d900febd906b133eaa4750e6bed647b0fbe50ed43e",
//     "3ef70f7291f47dfe2b82931a993e16f5a44a0e7a68034c3e0e086d77f5829adc",
//     "59788edf1f3e78cd0ebe6ce1446e9d10788225db3dedcfd1a59f764bad2b2690"
//   ]
// }
// ~~~
func (s *VarlinkInterface) DeleteUnusedImages(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.DeleteUnusedImages")
}

// Commit, creates an image from an existing container. It requires the name or
// ID of the container as well as the resulting image name.  Optionally, you can define an author and message
// to be added to the resulting image.  You can also define changes to the resulting image for the following
// attributes: _CMD, ENTRYPOINT, ENV, EXPOSE, LABEL, ONBUILD, STOPSIGNAL, USER, VOLUME, and WORKDIR_.  To pause the
// container while it is being committed, pass a _true_ bool for the pause argument.  If the container cannot
// be found by the ID or name provided, a (ContainerNotFound)[#ContainerNotFound] error will be returned; otherwise,
// the resulting image's ID will be returned as a string inside a MoreResponse.
func (s *VarlinkInterface) Commit(c VarlinkCall, name_ string, image_name_ string, changes_ []string, author_ string, message_ string, pause_ bool, manifestType_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.Commit")
}

// ImportImage imports an image from a source (like tarball) into local storage.  The image can have additional
// descriptions added to it using the message and changes options. See also [ExportImage](ExportImage).
func (s *VarlinkInterface) ImportImage(c VarlinkCall, source_ string, reference_ string, message_ string, changes_ []string, delete_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.ImportImage")
}

// ExportImage takes the name or ID of an image and exports it to a destination like a tarball.  There is also
// a boolean option to force compression.  It also takes in a string array of tags to be able to save multiple
// tags of the same image to a tarball (each tag should be of the form <image>:<tag>).  Upon completion, the ID
// of the image is returned. If the image cannot be found in local storage, an [ImageNotFound](#ImageNotFound)
// error will be returned. See also [ImportImage](ImportImage).
func (s *VarlinkInterface) ExportImage(c VarlinkCall, name_ string, destination_ string, compress_ bool, tags_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.ExportImage")
}

// PullImage pulls an image from a repository to local storage.  After a successful pull, the image id and logs
// are returned as a [MoreResponse](#MoreResponse).  This connection also will handle a WantsMores request to send
// status as it occurs.
func (s *VarlinkInterface) PullImage(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.PullImage")
}

// CreatePod creates a new empty pod.  It uses a [PodCreate](#PodCreate) type for input.
// On success, the ID of the newly created pod will be returned.
// #### Example
// ~~~
// $ varlink call unix:/run/podman/io.podman/io.podman.CreatePod '{"create": {"name": "test"}}'
// {
//   "pod": "b05dee7bd4ccfee688099fe1588a7a898d6ddd6897de9251d4671c9b0feacb2a"
// }
// # $ varlink call unix:/run/podman/io.podman/io.podman.CreatePod '{"create": {"infra": true, "share": ["ipc", "net", "uts"]}}'
// {
//   "pod": "d7697449a8035f613c1a8891286502aca68fff7d5d49a85279b3bda229af3b28"
// }
// ~~~
func (s *VarlinkInterface) CreatePod(c VarlinkCall, create_ PodCreate) error {
	return c.ReplyMethodNotImplemented("io.podman.CreatePod")
}

// ListPods returns a list of pods in no particular order.  They are
// returned as an array of ListPodData structs.  See also [GetPod](#GetPod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ListPods
// {
//   "pods": [
//     {
//       "cgroup": "machine.slice",
//       "containersinfo": [
//         {
//           "id": "00c130a45de0411f109f1a0cfea2e298df71db20fa939de5cab8b2160a36be45",
//           "name": "1840835294cf-infra",
//           "status": "running"
//         },
//         {
//           "id": "49a5cce72093a5ca47c6de86f10ad7bb36391e2d89cef765f807e460865a0ec6",
//           "name": "upbeat_murdock",
//           "status": "running"
//         }
//       ],
//       "createdat": "2018-12-07 13:10:15.014139258 -0600 CST",
//       "id": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f",
//       "name": "foobar",
//       "numberofcontainers": "2",
//       "status": "Running"
//     },
//     {
//       "cgroup": "machine.slice",
//       "containersinfo": [
//         {
//           "id": "1ca4b7bbba14a75ba00072d4b705c77f3df87db0109afaa44d50cb37c04a477e",
//           "name": "784306f655c6-infra",
//           "status": "running"
//         }
//       ],
//       "createdat": "2018-12-07 13:09:57.105112457 -0600 CST",
//       "id": "784306f655c6200aea321dd430ba685e9b2cc1f7d7528a72f3ff74ffb29485a2",
//       "name": "nostalgic_pike",
//       "numberofcontainers": "1",
//       "status": "Running"
//     }
//   ]
// }
// ~~~
func (s *VarlinkInterface) ListPods(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.ListPods")
}

// GetPod takes a name or ID of a pod and returns single [ListPodData](#ListPodData)
// structure.  A [PodNotFound](#PodNotFound) error will be returned if the pod cannot be found.
// See also [ListPods](ListPods).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.GetPod '{"name": "foobar"}'
// {
//   "pod": {
//     "cgroup": "machine.slice",
//     "containersinfo": [
//       {
//         "id": "00c130a45de0411f109f1a0cfea2e298df71db20fa939de5cab8b2160a36be45",
//         "name": "1840835294cf-infra",
//         "status": "running"
//       },
//       {
//         "id": "49a5cce72093a5ca47c6de86f10ad7bb36391e2d89cef765f807e460865a0ec6",
//         "name": "upbeat_murdock",
//         "status": "running"
//       }
//     ],
//     "createdat": "2018-12-07 13:10:15.014139258 -0600 CST",
//     "id": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f",
//     "name": "foobar",
//     "numberofcontainers": "2",
//     "status": "Running"
//   }
// }
// ~~~
func (s *VarlinkInterface) GetPod(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetPod")
}

// InspectPod takes the name or ID of an image and returns a string representation of data associated with the
// pod.  You must serialize the string into JSON to use it further.  A [PodNotFound](#PodNotFound) error will
// be returned if the pod cannot be found.
func (s *VarlinkInterface) InspectPod(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.InspectPod")
}

// StartPod starts containers in a pod.  It takes the name or ID of pod.  If the pod cannot be found, a [PodNotFound](#PodNotFound)
// error will be returned.  Containers in a pod are started independently. If there is an error starting one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was started with no errors, the pod ID is returned.
// See also [CreatePod](#CreatePod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.StartPod '{"name": "135d71b9495f"}'
// {
//   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6",
// }
// ~~~
func (s *VarlinkInterface) StartPod(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.StartPod")
}

// StopPod stops containers in a pod.  It takes the name or ID of a pod and a timeout.
// If the pod cannot be found, a [PodNotFound](#PodNotFound) error will be returned instead.
// Containers in a pod are stopped independently. If there is an error stopping one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was stopped with no errors, the pod ID is returned.
// See also [KillPod](KillPod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.StopPod '{"name": "135d71b9495f"}'
// {
//   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
// }
// ~~~
func (s *VarlinkInterface) StopPod(c VarlinkCall, name_ string, timeout_ int64) error {
	return c.ReplyMethodNotImplemented("io.podman.StopPod")
}

// RestartPod will restart containers in a pod given a pod name or ID. Containers in
// the pod that are running will be stopped, then all stopped containers will be run.
// If the pod cannot be found by name or ID, a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are restarted independently. If there is an error restarting one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was restarted with no errors, the pod ID is returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.RestartPod '{"name": "135d71b9495f"}'
// {
//   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
// }
// ~~~
func (s *VarlinkInterface) RestartPod(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.RestartPod")
}

// KillPod takes the name or ID of a pod as well as a signal to be applied to the pod.  If the pod cannot be found, a
// [PodNotFound](#PodNotFound) error is returned.
// Containers in a pod are killed independently. If there is an error killing one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was killed with no errors, the pod ID is returned.
// See also [StopPod](StopPod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.KillPod '{"name": "foobar", "signal": 15}'
// {
//   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
// }
// ~~~
func (s *VarlinkInterface) KillPod(c VarlinkCall, name_ string, signal_ int64) error {
	return c.ReplyMethodNotImplemented("io.podman.KillPod")
}

// PausePod takes the name or ID of a pod and pauses the running containers associated with it.  If the pod cannot be found,
// a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are paused independently. If there is an error pausing one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was paused with no errors, the pod ID is returned.
// See also [UnpausePod](#UnpausePod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.PausePod '{"name": "foobar"}'
// {
//   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
// }
// ~~~
func (s *VarlinkInterface) PausePod(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.PausePod")
}

// UnpausePod takes the name or ID of a pod and unpauses the paused containers associated with it.  If the pod cannot be
// found, a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are unpaused independently. If there is an error unpausing one container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was unpaused with no errors, the pod ID is returned.
// See also [PausePod](#PausePod).
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.UnpausePod '{"name": "foobar"}'
// {
//   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
// }
// ~~~
func (s *VarlinkInterface) UnpausePod(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.UnpausePod")
}

// RemovePod takes the name or ID of a pod as well a boolean representing whether a running
// container in the pod can be stopped and removed.  If a pod has containers associated with it, and force is not true,
// an error will occur.
// If the pod cannot be found by name or ID, a [PodNotFound](#PodNotFound) error will be returned.
// Containers in a pod are removed independently. If there is an error removing any container, the ID of those containers
// will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
// If the pod was removed with no errors, the pod ID is returned.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.RemovePod '{"name": "62f4fd98cb57", "force": "true"}'
// {
//   "pod": "62f4fd98cb57f529831e8f90610e54bba74bd6f02920ffb485e15376ed365c20"
// }
// ~~~
func (s *VarlinkInterface) RemovePod(c VarlinkCall, name_ string, force_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.RemovePod")
}

func (s *VarlinkInterface) TopPod(c VarlinkCall, pod_ string, latest_ bool, descriptors_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.TopPod")
}

// GetPodStats takes the name or ID of a pod and returns a pod name and slice of ContainerStats structure which
// contains attributes like memory and cpu usage.  If the pod cannot be found, a [PodNotFound](#PodNotFound)
// error will be returned.  If the pod has no running containers associated with it, a [NoContainerRunning](#NoContainerRunning)
// error will be returned.
// #### Example
// ~~~
// $ varlink call unix:/run/podman/io.podman/io.podman.GetPodStats '{"name": "7f62b508b6f12b11d8fe02e"}'
// {
//   "containers": [
//     {
//       "block_input": 0,
//       "block_output": 0,
//       "cpu": 2.833470544016107524276e-08,
//       "cpu_nano": 54363072,
//       "id": "a64b51f805121fe2c5a3dc5112eb61d6ed139e3d1c99110360d08b58d48e4a93",
//       "mem_limit": 12276146176,
//       "mem_perc": 7.974359265237864966003e-03,
//       "mem_usage": 978944,
//       "name": "quirky_heisenberg",
//       "net_input": 866,
//       "net_output": 7388,
//       "pids": 1,
//       "system_nano": 20000000
//     }
//   ],
//   "pod": "7f62b508b6f12b11d8fe02e0db4de6b9e43a7d7699b33a4fc0d574f6e82b4ebd"
// }
// ~~~
func (s *VarlinkInterface) GetPodStats(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetPodStats")
}

// GetPodsByStatus searches for pods whose status is included in statuses
func (s *VarlinkInterface) GetPodsByStatus(c VarlinkCall, statuses_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetPodsByStatus")
}

// ImageExists talks a full or partial image ID or name and returns an int as to whether
// the image exists in local storage. An int result of 0 means the image does exist in
// local storage; whereas 1 indicates the image does not exists in local storage.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ImageExists '{"name": "imageddoesntexist"}'
// {
//   "exists": 1
// }
// ~~~
func (s *VarlinkInterface) ImageExists(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.ImageExists")
}

// ContainerExists takes a full or partial container ID or name and returns an int as to
// whether the container exists in local storage.  A result of 0 means the container does
// exists; whereas a result of 1 means it could not be found.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.ContainerExists '{"name": "flamboyant_payne"}'{
//   "exists": 0
// }
// ~~~
func (s *VarlinkInterface) ContainerExists(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerExists")
}

// ContainerCheckPoint performs a checkpopint on a container by its name or full/partial container
// ID.  On successful checkpoint, the id of the checkpointed container is returned.
func (s *VarlinkInterface) ContainerCheckpoint(c VarlinkCall, name_ string, keep_ bool, leaveRunning_ bool, tcpEstablished_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerCheckpoint")
}

// ContainerRestore restores a container that has been checkpointed.  The container to be restored can
// be identified by its name or full/partial container ID.  A successful restore will result in the return
// of the container's ID.
func (s *VarlinkInterface) ContainerRestore(c VarlinkCall, name_ string, keep_ bool, tcpEstablished_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerRestore")
}

// ContainerRunlabel runs executes a command as described by a given container image label.
func (s *VarlinkInterface) ContainerRunlabel(c VarlinkCall, runlabel_ Runlabel) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerRunlabel")
}

// ListContainerMounts gathers all the mounted container mount points and returns them as an array
// of strings
// #### Example
// ~~~
// $ varlink call unix:/run/podman/io.podman/io.podman.ListContainerMounts
// {
//   "mounts": {
//     "04e4c255269ed2545e7f8bd1395a75f7949c50c223415c00c1d54bfa20f3b3d9": "/var/lib/containers/storage/overlay/a078925828f57e20467ca31cfca8a849210d21ec7e5757332b72b6924f441c17/merged",
//     "1d58c319f9e881a644a5122ff84419dccf6d138f744469281446ab243ef38924": "/var/lib/containers/storage/overlay/948fcf93f8cb932f0f03fd52e3180a58627d547192ffe3b88e0013b98ddcd0d2/merged"
//   }
// }
// ~~~
func (s *VarlinkInterface) ListContainerMounts(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.ListContainerMounts")
}

// MountContainer mounts a container by name or full/partial ID.  Upon a successful mount, the destination
// mount is returned as a string.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.MountContainer '{"name": "jolly_shannon"}'{
//   "path": "/var/lib/containers/storage/overlay/419eeb04e783ea159149ced67d9fcfc15211084d65e894792a96bedfae0470ca/merged"
// }
// ~~~
func (s *VarlinkInterface) MountContainer(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.MountContainer")
}

// UnmountContainer umounts a container by its name or full/partial container ID.
// #### Example
// ~~~
// $ varlink call -m unix:/run/podman/io.podman/io.podman.UnmountContainer '{"name": "jolly_shannon", "force": false}'
// {}
// ~~~
func (s *VarlinkInterface) UnmountContainer(c VarlinkCall, name_ string, force_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.UnmountContainer")
}

// ImagesPrune removes all unused images from the local store.  Upon successful pruning,
// the IDs of the removed images are returned.
func (s *VarlinkInterface) ImagesPrune(c VarlinkCall, all_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.ImagesPrune")
}

// GenerateKube generates a Kubernetes v1 Pod description of a Podman container or pod
// and its containers. The description is in YAML.  See also [ReplayKube](ReplayKube).
func (s *VarlinkInterface) GenerateKube(c VarlinkCall, name_ string, service_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.GenerateKube")
}

// ContainerConfig returns a container's config in string form. This call is for
// development of Podman only and generally should not be used.
func (s *VarlinkInterface) ContainerConfig(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerConfig")
}

// ContainerArtifacts returns a container's artifacts in string form.  This call is for
// development of Podman only and generally should not be used.
func (s *VarlinkInterface) ContainerArtifacts(c VarlinkCall, name_ string, artifactName_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerArtifacts")
}

// ContainerInspectData returns a container's inspect data in string form.  This call is for
// development of Podman only and generally should not be used.
func (s *VarlinkInterface) ContainerInspectData(c VarlinkCall, name_ string, size_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerInspectData")
}

// ContainerStateData returns a container's state config in string form.  This call is for
// development of Podman only and generally should not be used.
func (s *VarlinkInterface) ContainerStateData(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.ContainerStateData")
}

// PodStateData returns inspectr level information of a given pod in string form.  This call is for
// development of Podman only and generally should not be used.
func (s *VarlinkInterface) PodStateData(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.PodStateData")
}

// This call is for the development of Podman only and should not be used.
func (s *VarlinkInterface) CreateFromCC(c VarlinkCall, in_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.CreateFromCC")
}

// Spec returns the oci spec for a container.  This call is for development of Podman only and generally should not be used.
func (s *VarlinkInterface) Spec(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.Spec")
}

// Sendfile allows a remote client to send a file to the host
func (s *VarlinkInterface) SendFile(c VarlinkCall, type_ string, length_ int64) error {
	return c.ReplyMethodNotImplemented("io.podman.SendFile")
}

// ReceiveFile allows the host to send a remote client a file
func (s *VarlinkInterface) ReceiveFile(c VarlinkCall, path_ string, delete_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.ReceiveFile")
}

// VolumeCreate creates a volume on a remote host
func (s *VarlinkInterface) VolumeCreate(c VarlinkCall, options_ VolumeCreateOpts) error {
	return c.ReplyMethodNotImplemented("io.podman.VolumeCreate")
}

// VolumeRemove removes a volume on a remote host
func (s *VarlinkInterface) VolumeRemove(c VarlinkCall, options_ VolumeRemoveOpts) error {
	return c.ReplyMethodNotImplemented("io.podman.VolumeRemove")
}

// GetVolumes gets slice of the volumes on a remote host
func (s *VarlinkInterface) GetVolumes(c VarlinkCall, args_ []string, all_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.GetVolumes")
}

// VolumesPrune removes unused volumes on the host
func (s *VarlinkInterface) VolumesPrune(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.VolumesPrune")
}

// ImageSave allows you to save an image from the local image storage to a tarball
func (s *VarlinkInterface) ImageSave(c VarlinkCall, options_ ImageSaveOptions) error {
	return c.ReplyMethodNotImplemented("io.podman.ImageSave")
}

// GetPodsByContext allows you to get a list pod ids depending on all, latest, or a list of
// pod names.  The definition of latest pod means the latest by creation date.  In a multi-
// user environment, results might differ from what you expect.
func (s *VarlinkInterface) GetPodsByContext(c VarlinkCall, all_ bool, latest_ bool, args_ []string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetPodsByContext")
}

// LoadImage allows you to load an image into local storage from a tarball.
func (s *VarlinkInterface) LoadImage(c VarlinkCall, name_ string, inputFile_ string, quiet_ bool, deleteFile_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.LoadImage")
}

// GetEvents returns known libpod events filtered by the options provided.
func (s *VarlinkInterface) GetEvents(c VarlinkCall, filter_ []string, since_ string, until_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.GetEvents")
}

// Diff returns a diff between libpod objects
func (s *VarlinkInterface) Diff(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.Diff")
}

// GetLayersMapWithImageInfo is for the development of Podman and should not be used.
func (s *VarlinkInterface) GetLayersMapWithImageInfo(c VarlinkCall) error {
	return c.ReplyMethodNotImplemented("io.podman.GetLayersMapWithImageInfo")
}

// BuildImageHierarchyMap is for the development of Podman and should not be used.
func (s *VarlinkInterface) BuildImageHierarchyMap(c VarlinkCall, name_ string) error {
	return c.ReplyMethodNotImplemented("io.podman.BuildImageHierarchyMap")
}

func (s *VarlinkInterface) GenerateSystemd(c VarlinkCall, name_ string, restart_ string, timeout_ int64, useName_ bool) error {
	return c.ReplyMethodNotImplemented("io.podman.GenerateSystemd")
}

// Generated method call dispatcher

func (s *VarlinkInterface) VarlinkDispatch(call varlink.Call, methodname string) error {
	switch methodname {
	case "GetVersion":
		return s.iopodmanInterface.GetVersion(VarlinkCall{call})

	case "GetInfo":
		return s.iopodmanInterface.GetInfo(VarlinkCall{call})

	case "ListContainers":
		return s.iopodmanInterface.ListContainers(VarlinkCall{call})

	case "Ps":
		var in struct {
			Opts PsOpts `json:"opts"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.Ps(VarlinkCall{call}, in.Opts)

	case "GetContainersByStatus":
		var in struct {
			Status []string `json:"status"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetContainersByStatus(VarlinkCall{call}, []string(in.Status))

	case "Top":
		var in struct {
			NameOrID    string   `json:"nameOrID"`
			Descriptors []string `json:"descriptors"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.Top(VarlinkCall{call}, in.NameOrID, []string(in.Descriptors))

	case "GetContainer":
		var in struct {
			Id string `json:"id"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetContainer(VarlinkCall{call}, in.Id)

	case "GetContainersByContext":
		var in struct {
			All    bool     `json:"all"`
			Latest bool     `json:"latest"`
			Args   []string `json:"args"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetContainersByContext(VarlinkCall{call}, in.All, in.Latest, []string(in.Args))

	case "CreateContainer":
		var in struct {
			Create Create `json:"create"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.CreateContainer(VarlinkCall{call}, in.Create)

	case "InspectContainer":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.InspectContainer(VarlinkCall{call}, in.Name)

	case "ListContainerProcesses":
		var in struct {
			Name string   `json:"name"`
			Opts []string `json:"opts"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ListContainerProcesses(VarlinkCall{call}, in.Name, []string(in.Opts))

	case "GetContainerLogs":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetContainerLogs(VarlinkCall{call}, in.Name)

	case "GetContainersLogs":
		var in struct {
			Names      []string `json:"names"`
			Follow     bool     `json:"follow"`
			Latest     bool     `json:"latest"`
			Since      string   `json:"since"`
			Tail       int64    `json:"tail"`
			Timestamps bool     `json:"timestamps"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetContainersLogs(VarlinkCall{call}, []string(in.Names), in.Follow, in.Latest, in.Since, in.Tail, in.Timestamps)

	case "ListContainerChanges":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ListContainerChanges(VarlinkCall{call}, in.Name)

	case "ExportContainer":
		var in struct {
			Name string `json:"name"`
			Path string `json:"path"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ExportContainer(VarlinkCall{call}, in.Name, in.Path)

	case "GetContainerStats":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetContainerStats(VarlinkCall{call}, in.Name)

	case "GetContainerStatsWithHistory":
		var in struct {
			PreviousStats ContainerStats `json:"previousStats"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetContainerStatsWithHistory(VarlinkCall{call}, in.PreviousStats)

	case "StartContainer":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.StartContainer(VarlinkCall{call}, in.Name)

	case "StopContainer":
		var in struct {
			Name    string `json:"name"`
			Timeout int64  `json:"timeout"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.StopContainer(VarlinkCall{call}, in.Name, in.Timeout)

	case "InitContainer":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.InitContainer(VarlinkCall{call}, in.Name)

	case "RestartContainer":
		var in struct {
			Name    string `json:"name"`
			Timeout int64  `json:"timeout"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.RestartContainer(VarlinkCall{call}, in.Name, in.Timeout)

	case "KillContainer":
		var in struct {
			Name   string `json:"name"`
			Signal int64  `json:"signal"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.KillContainer(VarlinkCall{call}, in.Name, in.Signal)

	case "PauseContainer":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.PauseContainer(VarlinkCall{call}, in.Name)

	case "UnpauseContainer":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.UnpauseContainer(VarlinkCall{call}, in.Name)

	case "Attach":
		var in struct {
			Name       string `json:"name"`
			DetachKeys string `json:"detachKeys"`
			Start      bool   `json:"start"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.Attach(VarlinkCall{call}, in.Name, in.DetachKeys, in.Start)

	case "AttachControl":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.AttachControl(VarlinkCall{call}, in.Name)

	case "GetAttachSockets":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetAttachSockets(VarlinkCall{call}, in.Name)

	case "WaitContainer":
		var in struct {
			Name     string `json:"name"`
			Interval int64  `json:"interval"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.WaitContainer(VarlinkCall{call}, in.Name, in.Interval)

	case "RemoveContainer":
		var in struct {
			Name          string `json:"name"`
			Force         bool   `json:"force"`
			RemoveVolumes bool   `json:"removeVolumes"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.RemoveContainer(VarlinkCall{call}, in.Name, in.Force, in.RemoveVolumes)

	case "DeleteStoppedContainers":
		return s.iopodmanInterface.DeleteStoppedContainers(VarlinkCall{call})

	case "ListImages":
		return s.iopodmanInterface.ListImages(VarlinkCall{call})

	case "GetImage":
		var in struct {
			Id string `json:"id"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetImage(VarlinkCall{call}, in.Id)

	case "BuildImage":
		var in struct {
			Build BuildInfo `json:"build"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.BuildImage(VarlinkCall{call}, in.Build)

	case "InspectImage":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.InspectImage(VarlinkCall{call}, in.Name)

	case "HistoryImage":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.HistoryImage(VarlinkCall{call}, in.Name)

	case "PushImage":
		var in struct {
			Name             string `json:"name"`
			Tag              string `json:"tag"`
			Compress         bool   `json:"compress"`
			Format           string `json:"format"`
			RemoveSignatures bool   `json:"removeSignatures"`
			SignBy           string `json:"signBy"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.PushImage(VarlinkCall{call}, in.Name, in.Tag, in.Compress, in.Format, in.RemoveSignatures, in.SignBy)

	case "TagImage":
		var in struct {
			Name   string `json:"name"`
			Tagged string `json:"tagged"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.TagImage(VarlinkCall{call}, in.Name, in.Tagged)

	case "RemoveImage":
		var in struct {
			Name  string `json:"name"`
			Force bool   `json:"force"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.RemoveImage(VarlinkCall{call}, in.Name, in.Force)

	case "SearchImages":
		var in struct {
			Query  string            `json:"query"`
			Limit  *int64            `json:"limit,omitempty"`
			Filter ImageSearchFilter `json:"filter"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.SearchImages(VarlinkCall{call}, in.Query, in.Limit, in.Filter)

	case "DeleteUnusedImages":
		return s.iopodmanInterface.DeleteUnusedImages(VarlinkCall{call})

	case "Commit":
		var in struct {
			Name         string   `json:"name"`
			Image_name   string   `json:"image_name"`
			Changes      []string `json:"changes"`
			Author       string   `json:"author"`
			Message      string   `json:"message"`
			Pause        bool     `json:"pause"`
			ManifestType string   `json:"manifestType"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.Commit(VarlinkCall{call}, in.Name, in.Image_name, []string(in.Changes), in.Author, in.Message, in.Pause, in.ManifestType)

	case "ImportImage":
		var in struct {
			Source    string   `json:"source"`
			Reference string   `json:"reference"`
			Message   string   `json:"message"`
			Changes   []string `json:"changes"`
			Delete    bool     `json:"delete"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ImportImage(VarlinkCall{call}, in.Source, in.Reference, in.Message, []string(in.Changes), in.Delete)

	case "ExportImage":
		var in struct {
			Name        string   `json:"name"`
			Destination string   `json:"destination"`
			Compress    bool     `json:"compress"`
			Tags        []string `json:"tags"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ExportImage(VarlinkCall{call}, in.Name, in.Destination, in.Compress, []string(in.Tags))

	case "PullImage":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.PullImage(VarlinkCall{call}, in.Name)

	case "CreatePod":
		var in struct {
			Create PodCreate `json:"create"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.CreatePod(VarlinkCall{call}, in.Create)

	case "ListPods":
		return s.iopodmanInterface.ListPods(VarlinkCall{call})

	case "GetPod":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetPod(VarlinkCall{call}, in.Name)

	case "InspectPod":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.InspectPod(VarlinkCall{call}, in.Name)

	case "StartPod":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.StartPod(VarlinkCall{call}, in.Name)

	case "StopPod":
		var in struct {
			Name    string `json:"name"`
			Timeout int64  `json:"timeout"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.StopPod(VarlinkCall{call}, in.Name, in.Timeout)

	case "RestartPod":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.RestartPod(VarlinkCall{call}, in.Name)

	case "KillPod":
		var in struct {
			Name   string `json:"name"`
			Signal int64  `json:"signal"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.KillPod(VarlinkCall{call}, in.Name, in.Signal)

	case "PausePod":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.PausePod(VarlinkCall{call}, in.Name)

	case "UnpausePod":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.UnpausePod(VarlinkCall{call}, in.Name)

	case "RemovePod":
		var in struct {
			Name  string `json:"name"`
			Force bool   `json:"force"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.RemovePod(VarlinkCall{call}, in.Name, in.Force)

	case "TopPod":
		var in struct {
			Pod         string   `json:"pod"`
			Latest      bool     `json:"latest"`
			Descriptors []string `json:"descriptors"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.TopPod(VarlinkCall{call}, in.Pod, in.Latest, []string(in.Descriptors))

	case "GetPodStats":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetPodStats(VarlinkCall{call}, in.Name)

	case "GetPodsByStatus":
		var in struct {
			Statuses []string `json:"statuses"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetPodsByStatus(VarlinkCall{call}, []string(in.Statuses))

	case "ImageExists":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ImageExists(VarlinkCall{call}, in.Name)

	case "ContainerExists":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerExists(VarlinkCall{call}, in.Name)

	case "ContainerCheckpoint":
		var in struct {
			Name           string `json:"name"`
			Keep           bool   `json:"keep"`
			LeaveRunning   bool   `json:"leaveRunning"`
			TcpEstablished bool   `json:"tcpEstablished"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerCheckpoint(VarlinkCall{call}, in.Name, in.Keep, in.LeaveRunning, in.TcpEstablished)

	case "ContainerRestore":
		var in struct {
			Name           string `json:"name"`
			Keep           bool   `json:"keep"`
			TcpEstablished bool   `json:"tcpEstablished"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerRestore(VarlinkCall{call}, in.Name, in.Keep, in.TcpEstablished)

	case "ContainerRunlabel":
		var in struct {
			Runlabel Runlabel `json:"runlabel"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerRunlabel(VarlinkCall{call}, in.Runlabel)

	case "ListContainerMounts":
		return s.iopodmanInterface.ListContainerMounts(VarlinkCall{call})

	case "MountContainer":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.MountContainer(VarlinkCall{call}, in.Name)

	case "UnmountContainer":
		var in struct {
			Name  string `json:"name"`
			Force bool   `json:"force"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.UnmountContainer(VarlinkCall{call}, in.Name, in.Force)

	case "ImagesPrune":
		var in struct {
			All bool `json:"all"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ImagesPrune(VarlinkCall{call}, in.All)

	case "GenerateKube":
		var in struct {
			Name    string `json:"name"`
			Service bool   `json:"service"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GenerateKube(VarlinkCall{call}, in.Name, in.Service)

	case "ContainerConfig":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerConfig(VarlinkCall{call}, in.Name)

	case "ContainerArtifacts":
		var in struct {
			Name         string `json:"name"`
			ArtifactName string `json:"artifactName"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerArtifacts(VarlinkCall{call}, in.Name, in.ArtifactName)

	case "ContainerInspectData":
		var in struct {
			Name string `json:"name"`
			Size bool   `json:"size"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerInspectData(VarlinkCall{call}, in.Name, in.Size)

	case "ContainerStateData":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ContainerStateData(VarlinkCall{call}, in.Name)

	case "PodStateData":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.PodStateData(VarlinkCall{call}, in.Name)

	case "CreateFromCC":
		var in struct {
			In []string `json:"in"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.CreateFromCC(VarlinkCall{call}, []string(in.In))

	case "Spec":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.Spec(VarlinkCall{call}, in.Name)

	case "SendFile":
		var in struct {
			Type   string `json:"type"`
			Length int64  `json:"length"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.SendFile(VarlinkCall{call}, in.Type, in.Length)

	case "ReceiveFile":
		var in struct {
			Path   string `json:"path"`
			Delete bool   `json:"delete"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ReceiveFile(VarlinkCall{call}, in.Path, in.Delete)

	case "VolumeCreate":
		var in struct {
			Options VolumeCreateOpts `json:"options"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.VolumeCreate(VarlinkCall{call}, in.Options)

	case "VolumeRemove":
		var in struct {
			Options VolumeRemoveOpts `json:"options"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.VolumeRemove(VarlinkCall{call}, in.Options)

	case "GetVolumes":
		var in struct {
			Args []string `json:"args"`
			All  bool     `json:"all"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetVolumes(VarlinkCall{call}, []string(in.Args), in.All)

	case "VolumesPrune":
		return s.iopodmanInterface.VolumesPrune(VarlinkCall{call})

	case "ImageSave":
		var in struct {
			Options ImageSaveOptions `json:"options"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.ImageSave(VarlinkCall{call}, in.Options)

	case "GetPodsByContext":
		var in struct {
			All    bool     `json:"all"`
			Latest bool     `json:"latest"`
			Args   []string `json:"args"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetPodsByContext(VarlinkCall{call}, in.All, in.Latest, []string(in.Args))

	case "LoadImage":
		var in struct {
			Name       string `json:"name"`
			InputFile  string `json:"inputFile"`
			Quiet      bool   `json:"quiet"`
			DeleteFile bool   `json:"deleteFile"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.LoadImage(VarlinkCall{call}, in.Name, in.InputFile, in.Quiet, in.DeleteFile)

	case "GetEvents":
		var in struct {
			Filter []string `json:"filter"`
			Since  string   `json:"since"`
			Until  string   `json:"until"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GetEvents(VarlinkCall{call}, []string(in.Filter), in.Since, in.Until)

	case "Diff":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.Diff(VarlinkCall{call}, in.Name)

	case "GetLayersMapWithImageInfo":
		return s.iopodmanInterface.GetLayersMapWithImageInfo(VarlinkCall{call})

	case "BuildImageHierarchyMap":
		var in struct {
			Name string `json:"name"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.BuildImageHierarchyMap(VarlinkCall{call}, in.Name)

	case "GenerateSystemd":
		var in struct {
			Name    string `json:"name"`
			Restart string `json:"restart"`
			Timeout int64  `json:"timeout"`
			UseName bool   `json:"useName"`
		}
		err := call.GetParameters(&in)
		if err != nil {
			return call.ReplyInvalidParameter("parameters")
		}
		return s.iopodmanInterface.GenerateSystemd(VarlinkCall{call}, in.Name, in.Restart, in.Timeout, in.UseName)

	default:
		return call.ReplyMethodNotFound(methodname)
	}
}

// Generated varlink interface name

func (s *VarlinkInterface) VarlinkGetName() string {
	return `io.podman`
}

// Generated varlink interface description

func (s *VarlinkInterface) VarlinkGetDescription() string {
	return `# Podman Service Interface and API description.  The master version of this document can be found
# in the [API.md](https://github.com/containers/libpod/blob/master/API.md) file in the upstream libpod repository.
interface io.podman

type Volume (
  name: string,
  labels: [string]string,
  mountPoint: string,
  driver: string,
  options: [string]string,
  scope: string
)

type NotImplemented (
    comment: string
)

type StringResponse (
    message: string
)

type LogLine (
    device: string,
    parseLogType : string,
    time: string,
    msg: string,
    cid: string
)

# ContainerChanges describes the return struct for ListContainerChanges
type ContainerChanges (
   changed: []string,
   added: []string,
   deleted: []string
)

type ImageSaveOptions (
    name: string,
    format: string,
    output: string,
    outputType: string,
    moreTags: []string,
    quiet: bool,
    compress: bool
)

type VolumeCreateOpts (
   volumeName: string,
   driver: string,
   labels: [string]string,
   options: [string]string
)

type VolumeRemoveOpts (
   volumes: []string,
   all: bool,
   force: bool
)

type Image (
  id: string,
  digest:   string,
  parentId: string,
  repoTags: []string,
  repoDigests: []string,
  created: string, # as RFC3339
  size: int,
  virtualSize: int,
  containers: int,
  labels: [string]string,
  isParent: bool,
  topLayer: string
)

# ImageHistory describes the returned structure from ImageHistory.
type ImageHistory (
    id: string,
    created: string, # as RFC3339
    createdBy: string,
    tags: []string,
    size: int,
    comment: string
)

# Represents a single search result from SearchImages
type ImageSearchResult (
    description: string,
    is_official: bool,
    is_automated: bool,
    registry: string,
    name: string,
    star_count: int
)

type ImageSearchFilter (
    is_official: ?bool,
    is_automated: ?bool,
    star_count: int
)

type KubePodService (
    pod: string,
    service: string
)

type Container (
    id: string,
    image: string,
    imageid: string,
    command: []string,
    createdat: string, # as RFC3339
    runningfor: string,
    status: string,
    ports: []ContainerPortMappings,
    rootfssize: int,
    rwsize: int,
    names: string,
    labels: [string]string,
    mounts: []ContainerMount,
    containerrunning: bool,
    namespaces: ContainerNameSpace
)

# ContainerStats is the return struct for the stats of a container
type ContainerStats (
    id: string,
    name: string,
    cpu: float,
    cpu_nano: int,
    system_nano: int,
    mem_usage: int,
    mem_limit: int,
    mem_perc: float,
    net_input: int,
    net_output: int,
    block_output: int,
    block_input: int,
    pids: int
)

type PsOpts (
    all: bool,
    filters: ?[]string,
    last: ?int,
    latest: ?bool,
    noTrunc: ?bool,
	pod: ?bool,
	quiet: ?bool,
	sort: ?string,
	sync: ?bool
)

type PsContainer (
    id: string,
    image: string,
    command: string,
    created: string,
    ports: string,
    names: string,
    isInfra: bool,
    status: string,
    state: string,
    pidNum: int,
    rootFsSize: int,
    rwSize: int,
    pod: string,
    createdAt: string,
    exitedAt: string,
    startedAt: string,
    labels: [string]string,
    nsPid: string,
    cgroup: string,
    ipc: string,
    mnt: string,
    net: string,
    pidNs: string,
    user: string,
    uts: string,
    mounts: string
)

# ContainerMount describes the struct for mounts in a container
type ContainerMount (
    destination: string,
    type: string,
    source: string,
    options: []string
)

# ContainerPortMappings describes the struct for portmappings in an existing container
type ContainerPortMappings (
    host_port: string,
    host_ip: string,
    protocol: string,
    container_port: string
)

# ContainerNamespace describes the namespace structure for an existing container
type ContainerNameSpace (
    user: string,
    uts: string,
    pidns: string,
    pid: string,
    cgroup: string,
    net: string,
    mnt: string,
    ipc: string
)

# InfoDistribution describes the host's distribution
type InfoDistribution (
    distribution: string,
    version: string
)

# InfoHost describes the host stats portion of PodmanInfo
type InfoHost (
    buildah_version: string,
    distribution: InfoDistribution,
    mem_free: int,
    mem_total: int,
    swap_free: int,
    swap_total: int,
    arch: string,
    cpus: int,
    hostname: string,
    kernel: string,
    os: string,
    uptime: string
)

# InfoGraphStatus describes the detailed status of the storage driver
type InfoGraphStatus (
    backing_filesystem: string,
    native_overlay_diff: string,
    supports_d_type: string
)

# InfoStore describes the host's storage informatoin
type InfoStore (
    containers: int,
    images: int,
    graph_driver_name: string,
    graph_driver_options: string,
    graph_root: string,
    graph_status: InfoGraphStatus,
    run_root: string
)

# InfoPodman provides details on the podman binary
type InfoPodmanBinary (
    compiler: string,
    go_version: string,
    podman_version: string,
    git_commit: string
)

# PodmanInfo describes the Podman host and build
type PodmanInfo (
    host: InfoHost,
    registries: []string,
    insecure_registries: []string,
    store: InfoStore,
    podman: InfoPodmanBinary
)

# Sockets describes sockets location for a container
type Sockets(
    container_id: string,
    io_socket: string,
    control_socket: string
)

# Create is an input structure for creating containers.
type Create (
    args: []string,
    addHost: ?[]string,
    annotation: ?[]string,
    attach: ?[]string,
    blkioWeight: ?string,
    blkioWeightDevice: ?[]string,
    capAdd: ?[]string,
    capDrop: ?[]string,
    cgroupParent: ?string,
    cidFile: ?string,
    conmonPidfile: ?string,
    command: ?[]string,
    cpuPeriod: ?int,
    cpuQuota: ?int,
    cpuRtPeriod: ?int,
    cpuRtRuntime: ?int,
    cpuShares: ?int,
    cpus: ?float,
    cpuSetCpus: ?string,
    cpuSetMems: ?string,
    detach: ?bool,
    detachKeys: ?string,
    device: ?[]string,
    deviceReadBps: ?[]string,
    deviceReadIops: ?[]string,
    deviceWriteBps: ?[]string,
    deviceWriteIops: ?[]string,
    dns: ?[]string,
    dnsOpt: ?[]string,
    dnsSearch: ?[]string,
    dnsServers: ?[]string,
    entrypoint: ?string,
    env:  ?[]string,
    envFile: ?[]string,
    expose: ?[]string,
    gidmap: ?[]string,
    groupadd: ?[]string,
    healthcheckCommand: ?string,
    healthcheckInterval: ?string,
    healthcheckRetries: ?int,
    healthcheckStartPeriod: ?string,
    healthcheckTimeout:?string,
    hostname: ?string,
    imageVolume: ?string,
    init: ?bool,
    initPath: ?string,
    interactive: ?bool,
    ip: ?string,
    ipc: ?string,
    kernelMemory: ?string,
    label: ?[]string,
    labelFile: ?[]string,
    logDriver: ?string,
    logOpt: ?[]string,
    macAddress: ?string,
    memory: ?string,
    memoryReservation: ?string,
    memorySwap: ?string,
    memorySwappiness: ?int,
    name: ?string,
    net: ?string,
    network: ?string,
    noHosts: ?bool,
    oomKillDisable: ?bool,
    oomScoreAdj: ?int,
    pid: ?string,
    pidsLimit: ?int,
    pod: ?string,
    privileged: ?bool,
    publish: ?[]string,
    publishAll: ?bool,
    quiet: ?bool,
    readonly: ?bool,
    readonlytmpfs: ?bool,
    restart: ?string,
    rm: ?bool,
    rootfs: ?bool,
    securityOpt: ?[]string,
    shmSize: ?string,
    stopSignal: ?string,
    stopTimeout: ?int,
    storageOpt: ?[]string,
    subuidname: ?string,
    subgidname: ?string,
    sysctl: ?[]string,
    systemd: ?bool,
    tmpfs: ?[]string,
    tty: ?bool,
    uidmap: ?[]string,
    ulimit: ?[]string,
    user: ?string,
    userns: ?string,
    uts: ?string,
    mount: ?[]string,
    volume: ?[]string,
    volumesFrom: ?[]string,
    workDir: ?string
)

# BuildOptions are are used to describe describe physical attributes of the build
type BuildOptions (
    addHosts: []string,
    cgroupParent: string,
    cpuPeriod: int,
    cpuQuota: int,
    cpuShares: int,
    cpusetCpus: string,
    cpusetMems: string,
    memory: int,
    memorySwap: int,
    shmSize: string,
    ulimit: []string,
    volume: []string
)

# BuildInfo is used to describe user input for building images
type BuildInfo (
    additionalTags: []string,
    annotations: []string,
    buildArgs: [string]string,
    buildOptions: BuildOptions,
    cniConfigDir: string,
    cniPluginDir: string,
    compression: string,
    contextDir: string,
    defaultsMountFilePath: string,
    dockerfiles: []string,
    err: string,
    forceRmIntermediateCtrs: bool,
    iidfile: string,
    label: []string,
    layers: bool,
    nocache: bool,
    out: string,
    output: string,
    outputFormat: string,
    pullPolicy: string,
    quiet: bool,
    remoteIntermediateCtrs: bool,
    reportWriter: string,
    runtimeArgs: []string,
    squash: bool
)

# MoreResponse is a struct for when responses from varlink requires longer output
type MoreResponse (
    logs: []string,
    id: string
)

# ListPodContainerInfo is a returned struct for describing containers
# in a pod.
type ListPodContainerInfo (
    name: string,
    id: string,
    status: string
)

# PodCreate is an input structure for creating pods.
# It emulates options to podman pod create. The infraCommand and
# infraImage options are currently NotSupported.
type PodCreate (
    name: string,
    cgroupParent: string,
    labels: [string]string,
    share: []string,
    infra: bool,
    infraCommand: string,
    infraImage: string,
    publish: []string
)

# ListPodData is the returned struct for an individual pod
type ListPodData (
    id: string,
    name: string,
    createdat: string,
    cgroup: string,
    status: string,
    labels: [string]string,
    numberofcontainers: string,
    containersinfo: []ListPodContainerInfo
)

type PodContainerErrorData (
    containerid: string,
    reason: string
)

# Runlabel describes the required input for container runlabel
type Runlabel(
    image: string,
    authfile: string,
    display: bool,
    name: string,
    pull: bool,
    label: string,
    extraArgs: []string,
    opts: [string]string
)

# Event describes a libpod struct
type Event(
    # TODO: make status and type a enum at some point?
    # id is the container, volume, pod, image ID
    id: string,
    # image is the image name where applicable
    image: string,
    # name is the name of the pod, container, image
    name: string,
    # status describes the event that happened (i.e. create, remove, ...)
    status: string,
    # time the event happened
    time: string,
    # type describes object the event happened with (image, container...)
    type: string
)

type DiffInfo(
    # path that is different
    path: string,
    # Add, Delete, Modify
    changeType: string
)

# GetVersion returns version and build information of the podman service
method GetVersion() -> (
    version: string,
    go_version: string,
    git_commit: string,
    built: string, # as RFC3339
    os_arch: string,
    remote_api_version: int
)

# GetInfo returns a [PodmanInfo](#PodmanInfo) struct that describes podman and its host such as storage stats,
# build information of Podman, and system-wide registries.
method GetInfo() -> (info: PodmanInfo)

# ListContainers returns information about all containers.
# See also [GetContainer](#GetContainer).
method ListContainers() -> (containers: []Container)

method Ps(opts: PsOpts) -> (containers: []PsContainer)

method GetContainersByStatus(status: []string) -> (containerS: []Container)

method Top (nameOrID: string, descriptors: []string) -> (top: []string)

# GetContainer returns information about a single container.  If a container
# with the given id doesn't exist, a [ContainerNotFound](#ContainerNotFound)
# error will be returned.  See also [ListContainers](ListContainers) and
# [InspectContainer](#InspectContainer).
method GetContainer(id: string) -> (container: Container)

# GetContainersByContext allows you to get a list of container ids depending on all, latest, or a list of
# container names.  The definition of latest container means the latest by creation date.  In a multi-
# user environment, results might differ from what you expect.
method GetContainersByContext(all: bool, latest: bool, args: []string) -> (containers: []string)

# CreateContainer creates a new container from an image.  It uses a [Create](#Create) type for input.
method CreateContainer(create: Create) -> (container: string)

# InspectContainer data takes a name or ID of a container returns the inspection
# data in string format.  You can then serialize the string into JSON.  A [ContainerNotFound](#ContainerNotFound)
# error will be returned if the container cannot be found. See also [InspectImage](#InspectImage).
method InspectContainer(name: string) -> (container: string)

# ListContainerProcesses takes a name or ID of a container and returns the processes
# running inside the container as array of strings.  It will accept an array of string
# arguments that represent ps options.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
# error will be returned.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.ListContainerProcesses '{"name": "135d71b9495f", "opts": []}'
# {
#   "container": [
#     "  UID   PID  PPID  C STIME TTY          TIME CMD",
#     "    0 21220 21210  0 09:05 pts/0    00:00:00 /bin/sh",
#     "    0 21232 21220  0 09:05 pts/0    00:00:00 top",
#     "    0 21284 21220  0 09:05 pts/0    00:00:00 vi /etc/hosts"
#   ]
# }
# ~~~
method ListContainerProcesses(name: string, opts: []string) -> (container: []string)

# GetContainerLogs takes a name or ID of a container and returns the logs of that container.
# If the container cannot be found, a [ContainerNotFound](#ContainerNotFound) error will be returned.
# The container logs are returned as an array of strings.  GetContainerLogs will honor the streaming
# capability of varlink if the client invokes it.
method GetContainerLogs(name: string) -> (container: []string)

method GetContainersLogs(names: []string, follow: bool, latest: bool, since: string, tail: int, timestamps: bool) -> (log: LogLine)

# ListContainerChanges takes a name or ID of a container and returns changes between the container and
# its base image. It returns a struct of changed, deleted, and added path names.
method ListContainerChanges(name: string) -> (container: ContainerChanges)

# ExportContainer creates an image from a container.  It takes the name or ID of a container and a
# path representing the target tarfile.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
# error will be returned.
# The return value is the written tarfile.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.ExportContainer '{"name": "flamboyant_payne", "path": "/tmp/payne.tar" }'
# {
#   "tarfile": "/tmp/payne.tar"
# }
# ~~~
method ExportContainer(name: string, path: string) -> (tarfile: string)

# GetContainerStats takes the name or ID of a container and returns a single ContainerStats structure which
# contains attributes like memory and cpu usage.  If the container cannot be found, a
# [ContainerNotFound](#ContainerNotFound) error will be returned. If the container is not running, a [NoContainerRunning](#NoContainerRunning)
# error will be returned
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.GetContainerStats '{"name": "c33e4164f384"}'
# {
#   "container": {
#     "block_input": 0,
#     "block_output": 0,
#     "cpu": 2.571123918839990154678e-08,
#     "cpu_nano": 49037378,
#     "id": "c33e4164f384aa9d979072a63319d66b74fd7a128be71fa68ede24f33ec6cfee",
#     "mem_limit": 33080606720,
#     "mem_perc": 2.166828456524753747370e-03,
#     "mem_usage": 716800,
#     "name": "competent_wozniak",
#     "net_input": 768,
#     "net_output": 5910,
#     "pids": 1,
#     "system_nano": 10000000
#   }
# }
# ~~~
method GetContainerStats(name: string) -> (container: ContainerStats)

# GetContainerStatsWithHistory takes a previous set of container statistics and uses libpod functions
# to calculate the containers statistics based on current and previous measurements.
method GetContainerStatsWithHistory(previousStats: ContainerStats) -> (container: ContainerStats)

# This method has not be implemented yet.
# method ResizeContainerTty() -> (notimplemented: NotImplemented)

# StartContainer starts a created or stopped container. It takes the name or ID of container.  It returns
# the container ID once started.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
# error will be returned.  See also [CreateContainer](#CreateContainer).
method StartContainer(name: string) -> (container: string)

# StopContainer stops a container given a timeout.  It takes the name or ID of a container as well as a
# timeout value.  The timeout value the time before a forcible stop to the container is applied.  It
# returns the container ID once stopped. If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
# error will be returned instead. See also [KillContainer](KillContainer).
# #### Error
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.StopContainer '{"name": "135d71b9495f", "timeout": 5}'
# {
#   "container": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
# }
# ~~~
method StopContainer(name: string, timeout: int) -> (container: string)

# InitContainer initializes the given container. It accepts a container name or
# ID, and will initialize the container matching that ID if possible, and error
# if not. Containers can only be initialized when they are in the Created or
# Exited states. Initialization prepares a container to be started, but does not
# start the container. It is intended to be used to debug a container's state
# prior to starting it.
method InitContainer(name: string) -> (container: string)

# RestartContainer will restart a running container given a container name or ID and timeout value. The timeout
# value is the time before a forcible stop is used to stop the container.  If the container cannot be found by
# name or ID, a [ContainerNotFound](#ContainerNotFound)  error will be returned; otherwise, the ID of the
# container will be returned.
method RestartContainer(name: string, timeout: int) -> (container: string)

# KillContainer takes the name or ID of a container as well as a signal to be applied to the container.  Once the
# container has been killed, the container's ID is returned.  If the container cannot be found, a
# [ContainerNotFound](#ContainerNotFound) error is returned. See also [StopContainer](StopContainer).
method KillContainer(name: string, signal: int) -> (container: string)

# This method has not be implemented yet.
# method UpdateContainer() -> (notimplemented: NotImplemented)

# This method has not be implemented yet.
# method RenameContainer() -> (notimplemented: NotImplemented)

# PauseContainer takes the name or ID of container and pauses it.  If the container cannot be found,
# a [ContainerNotFound](#ContainerNotFound) error will be returned; otherwise the ID of the container is returned.
# See also [UnpauseContainer](#UnpauseContainer).
method PauseContainer(name: string) -> (container: string)

# UnpauseContainer takes the name or ID of container and unpauses a paused container.  If the container cannot be
# found, a [ContainerNotFound](#ContainerNotFound) error will be returned; otherwise the ID of the container is returned.
# See also [PauseContainer](#PauseContainer).
method UnpauseContainer(name: string) -> (container: string)

# Attach takes the name or ID of a container and sets up the ability to remotely attach to its console. The start
# bool is whether you wish to start the container in question first.
method Attach(name: string, detachKeys: string, start: bool) -> ()

method AttachControl(name: string) -> ()

# GetAttachSockets takes the name or ID of an existing container.  It returns file paths for two sockets needed
# to properly communicate with a container.  The first is the actual I/O socket that the container uses.  The
# second is a "control" socket where things like resizing the TTY events are sent. If the container cannot be
# found, a [ContainerNotFound](#ContainerNotFound) error will be returned.
# #### Example
# ~~~
# $ varlink call -m unix:/run/io.podman/io.podman.GetAttachSockets '{"name": "b7624e775431219161"}'
# {
#   "sockets": {
#     "container_id": "b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca",
#     "control_socket": "/var/lib/containers/storage/overlay-containers/b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca/userdata/ctl",
#     "io_socket": "/var/run/libpod/socket/b7624e7754312191613245ce1a46844abee60025818fe3c3f3203435623a1eca/attach"
#   }
# }
# ~~~
method GetAttachSockets(name: string) -> (sockets: Sockets)

# WaitContainer takes the name or ID of a container and waits the given interval in milliseconds until the container
# stops.  Upon stopping, the return code of the container is returned. If the container container cannot be found by ID
# or name, a [ContainerNotFound](#ContainerNotFound) error is returned.
method WaitContainer(name: string, interval: int) -> (exitcode: int)

# RemoveContainer requires the name or ID of container as well a boolean representing whether a running container can be stopped and removed, and a boolean
# indicating whether to remove builtin volumes. Upon successful removal of the
# container, its ID is returned.  If the
# container cannot be found by name or ID, a [ContainerNotFound](#ContainerNotFound) error will be returned.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.RemoveContainer '{"name": "62f4fd98cb57"}'
# {
#   "container": "62f4fd98cb57f529831e8f90610e54bba74bd6f02920ffb485e15376ed365c20"
# }
# ~~~
method RemoveContainer(name: string, force: bool, removeVolumes: bool) -> (container: string)

# DeleteStoppedContainers will delete all containers that are not running. It will return a list the deleted
# container IDs.  See also [RemoveContainer](RemoveContainer).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.DeleteStoppedContainers
# {
#   "containers": [
#     "451410b931d00def8aa9b4f8084e4d4a39e5e04ea61f358cf53a5cf95afcdcee",
#     "8b60f754a3e01389494a9581ade97d35c2765b6e2f19acd2d3040c82a32d1bc0",
#     "cf2e99d4d3cad6073df199ed32bbe64b124f3e1aba6d78821aa8460e70d30084",
#     "db901a329587312366e5ecff583d08f0875b4b79294322df67d90fc6eed08fc1"
#   ]
# }
# ~~~
method DeleteStoppedContainers() -> (containers: []string)

# ListImages returns information about the images that are currently in storage.
# See also [InspectImage](#InspectImage).
method ListImages() -> (images: []Image)

# GetImage returns information about a single image in storage.
# If the image caGetImage returns be found, [ImageNotFound](#ImageNotFound) will be returned.
method GetImage(id: string) -> (image: Image)

# BuildImage takes a [BuildInfo](#BuildInfo) structure and builds an image.  At a minimum, you must provide the
# 'dockerfile' and 'tags' options in the BuildInfo structure. It will return a [MoreResponse](#MoreResponse) structure
# that contains the build logs and resulting image ID.
method BuildImage(build: BuildInfo) -> (image: MoreResponse)

# This function is not implemented yet.
# method CreateImage() -> (notimplemented: NotImplemented)

# InspectImage takes the name or ID of an image and returns a string representation of data associated with the
#image.  You must serialize the string into JSON to use it further.  An [ImageNotFound](#ImageNotFound) error will
# be returned if the image cannot be found.
method InspectImage(name: string) -> (image: string)

# HistoryImage takes the name or ID of an image and returns information about its history and layers.  The returned
# history is in the form of an array of ImageHistory structures.  If the image cannot be found, an
# [ImageNotFound](#ImageNotFound) error is returned.
method HistoryImage(name: string) -> (history: []ImageHistory)

# PushImage takes two input arguments: the name or ID of an image, the fully-qualified destination name of the image,
# It will return an [ImageNotFound](#ImageNotFound) error if
# the image cannot be found in local storage; otherwise it will return a [MoreResponse](#MoreResponse)
method PushImage(name: string, tag: string, compress: bool, format: string, removeSignatures: bool, signBy: string) -> (reply: MoreResponse)

# TagImage takes the name or ID of an image in local storage as well as the desired tag name.  If the image cannot
# be found, an [ImageNotFound](#ImageNotFound) error will be returned; otherwise, the ID of the image is returned on success.
method TagImage(name: string, tagged: string) -> (image: string)

# RemoveImage takes the name or ID of an image as well as a boolean that determines if containers using that image
# should be deleted.  If the image cannot be found, an [ImageNotFound](#ImageNotFound) error will be returned.  The
# ID of the removed image is returned when complete.  See also [DeleteUnusedImages](DeleteUnusedImages).
# #### Example
# ~~~
# varlink call -m unix:/run/podman/io.podman/io.podman.RemoveImage '{"name": "registry.fedoraproject.org/fedora", "force": true}'
# {
#   "image": "426866d6fa419873f97e5cbd320eeb22778244c1dfffa01c944db3114f55772e"
# }
# ~~~
method RemoveImage(name: string, force: bool) -> (image: string)

# SearchImages searches available registries for images that contain the
# contents of "query" in their name. If "limit" is given, limits the amount of
# search results per registry.
method SearchImages(query: string, limit: ?int, filter: ImageSearchFilter) -> (results: []ImageSearchResult)

# DeleteUnusedImages deletes any images not associated with a container.  The IDs of the deleted images are returned
# in a string array.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.DeleteUnusedImages
# {
#   "images": [
#     "166ea6588079559c724c15223f52927f514f73dd5c5cf2ae2d143e3b2e6e9b52",
#     "da86e6ba6ca197bf6bc5e9d900febd906b133eaa4750e6bed647b0fbe50ed43e",
#     "3ef70f7291f47dfe2b82931a993e16f5a44a0e7a68034c3e0e086d77f5829adc",
#     "59788edf1f3e78cd0ebe6ce1446e9d10788225db3dedcfd1a59f764bad2b2690"
#   ]
# }
# ~~~
method DeleteUnusedImages() -> (images: []string)

# Commit, creates an image from an existing container. It requires the name or
# ID of the container as well as the resulting image name.  Optionally, you can define an author and message
# to be added to the resulting image.  You can also define changes to the resulting image for the following
# attributes: _CMD, ENTRYPOINT, ENV, EXPOSE, LABEL, ONBUILD, STOPSIGNAL, USER, VOLUME, and WORKDIR_.  To pause the
# container while it is being committed, pass a _true_ bool for the pause argument.  If the container cannot
# be found by the ID or name provided, a (ContainerNotFound)[#ContainerNotFound] error will be returned; otherwise,
# the resulting image's ID will be returned as a string inside a MoreResponse.
method Commit(name: string, image_name: string, changes: []string, author: string, message: string, pause: bool, manifestType: string) -> (reply: MoreResponse)

# ImportImage imports an image from a source (like tarball) into local storage.  The image can have additional
# descriptions added to it using the message and changes options. See also [ExportImage](ExportImage).
method ImportImage(source: string, reference: string, message: string, changes: []string, delete: bool) -> (image: string)

# ExportImage takes the name or ID of an image and exports it to a destination like a tarball.  There is also
# a boolean option to force compression.  It also takes in a string array of tags to be able to save multiple
# tags of the same image to a tarball (each tag should be of the form <image>:<tag>).  Upon completion, the ID
# of the image is returned. If the image cannot be found in local storage, an [ImageNotFound](#ImageNotFound)
# error will be returned. See also [ImportImage](ImportImage).
method ExportImage(name: string, destination: string, compress: bool, tags: []string) -> (image: string)

# PullImage pulls an image from a repository to local storage.  After a successful pull, the image id and logs
# are returned as a [MoreResponse](#MoreResponse).  This connection also will handle a WantsMores request to send
# status as it occurs.
method PullImage(name: string) -> (reply: MoreResponse)

# CreatePod creates a new empty pod.  It uses a [PodCreate](#PodCreate) type for input.
# On success, the ID of the newly created pod will be returned.
# #### Example
# ~~~
# $ varlink call unix:/run/podman/io.podman/io.podman.CreatePod '{"create": {"name": "test"}}'
# {
#   "pod": "b05dee7bd4ccfee688099fe1588a7a898d6ddd6897de9251d4671c9b0feacb2a"
# }
#
# $ varlink call unix:/run/podman/io.podman/io.podman.CreatePod '{"create": {"infra": true, "share": ["ipc", "net", "uts"]}}'
# {
#   "pod": "d7697449a8035f613c1a8891286502aca68fff7d5d49a85279b3bda229af3b28"
# }
# ~~~
method CreatePod(create: PodCreate) -> (pod: string)

# ListPods returns a list of pods in no particular order.  They are
# returned as an array of ListPodData structs.  See also [GetPod](#GetPod).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.ListPods
# {
#   "pods": [
#     {
#       "cgroup": "machine.slice",
#       "containersinfo": [
#         {
#           "id": "00c130a45de0411f109f1a0cfea2e298df71db20fa939de5cab8b2160a36be45",
#           "name": "1840835294cf-infra",
#           "status": "running"
#         },
#         {
#           "id": "49a5cce72093a5ca47c6de86f10ad7bb36391e2d89cef765f807e460865a0ec6",
#           "name": "upbeat_murdock",
#           "status": "running"
#         }
#       ],
#       "createdat": "2018-12-07 13:10:15.014139258 -0600 CST",
#       "id": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f",
#       "name": "foobar",
#       "numberofcontainers": "2",
#       "status": "Running"
#     },
#     {
#       "cgroup": "machine.slice",
#       "containersinfo": [
#         {
#           "id": "1ca4b7bbba14a75ba00072d4b705c77f3df87db0109afaa44d50cb37c04a477e",
#           "name": "784306f655c6-infra",
#           "status": "running"
#         }
#       ],
#       "createdat": "2018-12-07 13:09:57.105112457 -0600 CST",
#       "id": "784306f655c6200aea321dd430ba685e9b2cc1f7d7528a72f3ff74ffb29485a2",
#       "name": "nostalgic_pike",
#       "numberofcontainers": "1",
#       "status": "Running"
#     }
#   ]
# }
# ~~~
method ListPods() -> (pods: []ListPodData)

# GetPod takes a name or ID of a pod and returns single [ListPodData](#ListPodData)
# structure.  A [PodNotFound](#PodNotFound) error will be returned if the pod cannot be found.
# See also [ListPods](ListPods).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.GetPod '{"name": "foobar"}'
# {
#   "pod": {
#     "cgroup": "machine.slice",
#     "containersinfo": [
#       {
#         "id": "00c130a45de0411f109f1a0cfea2e298df71db20fa939de5cab8b2160a36be45",
#         "name": "1840835294cf-infra",
#         "status": "running"
#       },
#       {
#         "id": "49a5cce72093a5ca47c6de86f10ad7bb36391e2d89cef765f807e460865a0ec6",
#         "name": "upbeat_murdock",
#         "status": "running"
#       }
#     ],
#     "createdat": "2018-12-07 13:10:15.014139258 -0600 CST",
#     "id": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f",
#     "name": "foobar",
#     "numberofcontainers": "2",
#     "status": "Running"
#   }
# }
# ~~~
method GetPod(name: string) -> (pod: ListPodData)

# InspectPod takes the name or ID of an image and returns a string representation of data associated with the
# pod.  You must serialize the string into JSON to use it further.  A [PodNotFound](#PodNotFound) error will
# be returned if the pod cannot be found.
method InspectPod(name: string) -> (pod: string)

# StartPod starts containers in a pod.  It takes the name or ID of pod.  If the pod cannot be found, a [PodNotFound](#PodNotFound)
# error will be returned.  Containers in a pod are started independently. If there is an error starting one container, the ID of those containers
# will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
# If the pod was started with no errors, the pod ID is returned.
# See also [CreatePod](#CreatePod).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.StartPod '{"name": "135d71b9495f"}'
# {
#   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6",
# }
# ~~~
method StartPod(name: string) -> (pod: string)

# StopPod stops containers in a pod.  It takes the name or ID of a pod and a timeout.
# If the pod cannot be found, a [PodNotFound](#PodNotFound) error will be returned instead.
# Containers in a pod are stopped independently. If there is an error stopping one container, the ID of those containers
# will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
# If the pod was stopped with no errors, the pod ID is returned.
# See also [KillPod](KillPod).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.StopPod '{"name": "135d71b9495f"}'
# {
#   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
# }
# ~~~
method StopPod(name: string, timeout: int) -> (pod: string)

# RestartPod will restart containers in a pod given a pod name or ID. Containers in
# the pod that are running will be stopped, then all stopped containers will be run.
# If the pod cannot be found by name or ID, a [PodNotFound](#PodNotFound) error will be returned.
# Containers in a pod are restarted independently. If there is an error restarting one container, the ID of those containers
# will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
# If the pod was restarted with no errors, the pod ID is returned.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.RestartPod '{"name": "135d71b9495f"}'
# {
#   "pod": "135d71b9495f7c3967f536edad57750bfdb569336cd107d8aabab45565ffcfb6"
# }
# ~~~
method RestartPod(name: string) -> (pod: string)

# KillPod takes the name or ID of a pod as well as a signal to be applied to the pod.  If the pod cannot be found, a
# [PodNotFound](#PodNotFound) error is returned.
# Containers in a pod are killed independently. If there is an error killing one container, the ID of those containers
# will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
# If the pod was killed with no errors, the pod ID is returned.
# See also [StopPod](StopPod).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.KillPod '{"name": "foobar", "signal": 15}'
# {
#   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
# }
# ~~~
method KillPod(name: string, signal: int) -> (pod: string)

# PausePod takes the name or ID of a pod and pauses the running containers associated with it.  If the pod cannot be found,
# a [PodNotFound](#PodNotFound) error will be returned.
# Containers in a pod are paused independently. If there is an error pausing one container, the ID of those containers
# will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
# If the pod was paused with no errors, the pod ID is returned.
# See also [UnpausePod](#UnpausePod).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.PausePod '{"name": "foobar"}'
# {
#   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
# }
# ~~~
method PausePod(name: string) -> (pod: string)

# UnpausePod takes the name or ID of a pod and unpauses the paused containers associated with it.  If the pod cannot be
# found, a [PodNotFound](#PodNotFound) error will be returned.
# Containers in a pod are unpaused independently. If there is an error unpausing one container, the ID of those containers
# will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
# If the pod was unpaused with no errors, the pod ID is returned.
# See also [PausePod](#PausePod).
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.UnpausePod '{"name": "foobar"}'
# {
#   "pod": "1840835294cf076a822e4e12ba4152411f131bd869e7f6a4e8b16df9b0ea5c7f"
# }
# ~~~
method UnpausePod(name: string) -> (pod: string)

# RemovePod takes the name or ID of a pod as well a boolean representing whether a running
# container in the pod can be stopped and removed.  If a pod has containers associated with it, and force is not true,
# an error will occur.
# If the pod cannot be found by name or ID, a [PodNotFound](#PodNotFound) error will be returned.
# Containers in a pod are removed independently. If there is an error removing any container, the ID of those containers
# will be returned in a list, along with the ID of the pod in a [PodContainerError](#PodContainerError).
# If the pod was removed with no errors, the pod ID is returned.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.RemovePod '{"name": "62f4fd98cb57", "force": "true"}'
# {
#   "pod": "62f4fd98cb57f529831e8f90610e54bba74bd6f02920ffb485e15376ed365c20"
# }
# ~~~
method RemovePod(name: string, force: bool) -> (pod: string)

# This method has not be implemented yet.
# method WaitPod() -> (notimplemented: NotImplemented)

method TopPod(pod: string, latest: bool, descriptors: []string) -> (stats: []string)

# GetPodStats takes the name or ID of a pod and returns a pod name and slice of ContainerStats structure which
# contains attributes like memory and cpu usage.  If the pod cannot be found, a [PodNotFound](#PodNotFound)
# error will be returned.  If the pod has no running containers associated with it, a [NoContainerRunning](#NoContainerRunning)
# error will be returned.
# #### Example
# ~~~
# $ varlink call unix:/run/podman/io.podman/io.podman.GetPodStats '{"name": "7f62b508b6f12b11d8fe02e"}'
# {
#   "containers": [
#     {
#       "block_input": 0,
#       "block_output": 0,
#       "cpu": 2.833470544016107524276e-08,
#       "cpu_nano": 54363072,
#       "id": "a64b51f805121fe2c5a3dc5112eb61d6ed139e3d1c99110360d08b58d48e4a93",
#       "mem_limit": 12276146176,
#       "mem_perc": 7.974359265237864966003e-03,
#       "mem_usage": 978944,
#       "name": "quirky_heisenberg",
#       "net_input": 866,
#       "net_output": 7388,
#       "pids": 1,
#       "system_nano": 20000000
#     }
#   ],
#   "pod": "7f62b508b6f12b11d8fe02e0db4de6b9e43a7d7699b33a4fc0d574f6e82b4ebd"
# }
# ~~~
method GetPodStats(name: string) -> (pod: string, containers: []ContainerStats)

# GetPodsByStatus searches for pods whose status is included in statuses
method GetPodsByStatus(statuses: []string) -> (pods: []string)

# ImageExists talks a full or partial image ID or name and returns an int as to whether
# the image exists in local storage. An int result of 0 means the image does exist in
# local storage; whereas 1 indicates the image does not exists in local storage.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.ImageExists '{"name": "imageddoesntexist"}'
# {
#   "exists": 1
# }
# ~~~
method ImageExists(name: string) -> (exists: int)

# ContainerExists takes a full or partial container ID or name and returns an int as to
# whether the container exists in local storage.  A result of 0 means the container does
# exists; whereas a result of 1 means it could not be found.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.ContainerExists '{"name": "flamboyant_payne"}'{
#   "exists": 0
# }
# ~~~
method ContainerExists(name: string) -> (exists: int)

# ContainerCheckPoint performs a checkpopint on a container by its name or full/partial container
# ID.  On successful checkpoint, the id of the checkpointed container is returned.
method ContainerCheckpoint(name: string, keep: bool, leaveRunning: bool, tcpEstablished: bool) -> (id: string)

# ContainerRestore restores a container that has been checkpointed.  The container to be restored can
# be identified by its name or full/partial container ID.  A successful restore will result in the return
# of the container's ID.
method ContainerRestore(name: string, keep: bool, tcpEstablished: bool) -> (id: string)

# ContainerRunlabel runs executes a command as described by a given container image label.
method ContainerRunlabel(runlabel: Runlabel) -> ()

# ListContainerMounts gathers all the mounted container mount points and returns them as an array
# of strings
# #### Example
# ~~~
# $ varlink call unix:/run/podman/io.podman/io.podman.ListContainerMounts
# {
#   "mounts": {
#     "04e4c255269ed2545e7f8bd1395a75f7949c50c223415c00c1d54bfa20f3b3d9": "/var/lib/containers/storage/overlay/a078925828f57e20467ca31cfca8a849210d21ec7e5757332b72b6924f441c17/merged",
#     "1d58c319f9e881a644a5122ff84419dccf6d138f744469281446ab243ef38924": "/var/lib/containers/storage/overlay/948fcf93f8cb932f0f03fd52e3180a58627d547192ffe3b88e0013b98ddcd0d2/merged"
#   }
# }
# ~~~
method ListContainerMounts() -> (mounts: [string]string)

# MountContainer mounts a container by name or full/partial ID.  Upon a successful mount, the destination
# mount is returned as a string.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.MountContainer '{"name": "jolly_shannon"}'{
#   "path": "/var/lib/containers/storage/overlay/419eeb04e783ea159149ced67d9fcfc15211084d65e894792a96bedfae0470ca/merged"
# }
# ~~~
method MountContainer(name: string) -> (path: string)

# UnmountContainer umounts a container by its name or full/partial container ID.
# #### Example
# ~~~
# $ varlink call -m unix:/run/podman/io.podman/io.podman.UnmountContainer '{"name": "jolly_shannon", "force": false}'
# {}
# ~~~
method  UnmountContainer(name: string, force: bool) -> ()

# ImagesPrune removes all unused images from the local store.  Upon successful pruning,
# the IDs of the removed images are returned.
method ImagesPrune(all: bool) -> (pruned: []string)

# This function is not implemented yet.
# method ListContainerPorts(name: string) -> (notimplemented: NotImplemented)

# GenerateKube generates a Kubernetes v1 Pod description of a Podman container or pod
# and its containers. The description is in YAML.  See also [ReplayKube](ReplayKube).
method GenerateKube(name: string, service: bool) -> (pod: KubePodService)

# ReplayKube recreates a pod and its containers based on a Kubernetes v1 Pod description (in YAML)
# like that created by GenerateKube. See also [GenerateKube](GenerateKube).
# method ReplayKube() -> (notimplemented: NotImplemented)

# ContainerConfig returns a container's config in string form. This call is for
# development of Podman only and generally should not be used.
method ContainerConfig(name: string) -> (config: string)

# ContainerArtifacts returns a container's artifacts in string form.  This call is for
# development of Podman only and generally should not be used.
method ContainerArtifacts(name: string, artifactName: string) -> (config: string)

# ContainerInspectData returns a container's inspect data in string form.  This call is for
# development of Podman only and generally should not be used.
method ContainerInspectData(name: string, size: bool) -> (config: string)

# ContainerStateData returns a container's state config in string form.  This call is for
# development of Podman only and generally should not be used.
method ContainerStateData(name: string) -> (config: string)

# PodStateData returns inspectr level information of a given pod in string form.  This call is for
# development of Podman only and generally should not be used.
method PodStateData(name: string) -> (config: string)

# This call is for the development of Podman only and should not be used.
method CreateFromCC(in: []string) -> (id: string)

# Spec returns the oci spec for a container.  This call is for development of Podman only and generally should not be used.
method Spec(name: string) -> (config: string)

# Sendfile allows a remote client to send a file to the host
method SendFile(type: string, length: int) -> (file_handle: string)

# ReceiveFile allows the host to send a remote client a file
method ReceiveFile(path: string, delete: bool) -> (len: int)

# VolumeCreate creates a volume on a remote host
method VolumeCreate(options: VolumeCreateOpts) -> (volumeName: string)

# VolumeRemove removes a volume on a remote host
method VolumeRemove(options: VolumeRemoveOpts) -> (volumeNames: []string)

# GetVolumes gets slice of the volumes on a remote host
method GetVolumes(args: []string, all: bool) -> (volumes: []Volume)

# VolumesPrune removes unused volumes on the host
method VolumesPrune() -> (prunedNames: []string, prunedErrors: []string)

# ImageSave allows you to save an image from the local image storage to a tarball
method ImageSave(options: ImageSaveOptions) -> (reply: MoreResponse)

# GetPodsByContext allows you to get a list pod ids depending on all, latest, or a list of
# pod names.  The definition of latest pod means the latest by creation date.  In a multi-
# user environment, results might differ from what you expect.
method GetPodsByContext(all: bool, latest: bool, args: []string) -> (pods: []string)

# LoadImage allows you to load an image into local storage from a tarball.
method LoadImage(name: string, inputFile: string, quiet: bool, deleteFile: bool) -> (reply: MoreResponse)

# GetEvents returns known libpod events filtered by the options provided.
method GetEvents(filter: []string, since: string, until: string) -> (events: Event)

# Diff returns a diff between libpod objects
method Diff(name: string) -> (diffs: []DiffInfo)

# GetLayersMapWithImageInfo is for the development of Podman and should not be used.
method GetLayersMapWithImageInfo() -> (layerMap: string)

# BuildImageHierarchyMap is for the development of Podman and should not be used.
method BuildImageHierarchyMap(name: string) -> (imageInfo: string)

method GenerateSystemd(name: string, restart: string, timeout: int, useName: bool) -> (unit: string)

# ImageNotFound means the image could not be found by the provided name or ID in local storage.
error ImageNotFound (id: string, reason: string)

# ContainerNotFound means the container could not be found by the provided name or ID in local storage.
error ContainerNotFound (id: string, reason: string)

# NoContainerRunning means none of the containers requested are running in a command that requires a running container.
error NoContainerRunning ()

# PodNotFound means the pod could not be found by the provided name or ID in local storage.
error PodNotFound (name: string, reason: string)

# VolumeNotFound means the volume could not be found by the name or ID in local storage.
error VolumeNotFound (id: string, reason: string)

# PodContainerError means a container associated with a pod failed to perform an operation. It contains
# a container ID of the container that failed.
error PodContainerError (podname: string, errors: []PodContainerErrorData)

# NoContainersInPod means a pod has no containers on which to perform the operation. It contains
# the pod ID.
error NoContainersInPod (name: string)

# InvalidState indicates that a container or pod was in an improper state for the requested operation
error InvalidState (id: string, reason: string)

# ErrorOccurred is a generic error for an error that occurs during the execution.  The actual error message
# is includes as part of the error's text.
error ErrorOccurred (reason: string)

# RuntimeErrors generally means a runtime could not be found or gotten.
error RuntimeError (reason: string)

# The Podman endpoint requires that you use a streaming connection.
error WantsMoreRequired (reason: string)

# Container is already stopped
error ErrCtrStopped (id: string)
`
}

// Generated service interface

type VarlinkInterface struct {
	iopodmanInterface
}

func VarlinkNew(m iopodmanInterface) *VarlinkInterface {
	return &VarlinkInterface{m}
}
