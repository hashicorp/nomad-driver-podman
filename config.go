// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/helper/pluginutils/hclutils"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
)

var (
	// socketBodySpec is the hcl specification for the sockets in the driver object
	socketBodySpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"name":        hclspec.NewAttr("name", "string", false), // If not specified == host_user
		"socket_path": hclspec.NewAttr("socket_path", "string", true),
	})

	// configSpec is the hcl specification returned by the ConfigSchema RPC
	configSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		// image registry authentication options
		"auth": hclspec.NewBlock("auth", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"config": hclspec.NewAttr("config", "string", false),
			"helper": hclspec.NewAttr("helper", "string", false),
		})),

		// volume options
		"volumes": hclspec.NewDefault(hclspec.NewBlock("volumes", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"enabled": hclspec.NewDefault(
				hclspec.NewAttr("enabled", "bool", false),
				hclspec.NewLiteral("true"),
			),
			"selinuxlabel": hclspec.NewAttr("selinuxlabel", "string", false),
		})), hclspec.NewLiteral("{ enabled = true }")),
		// garbage collection options
		// default needed for both if the gc {...} block is not set and
		// if the default fields are missing
		"gc": hclspec.NewDefault(hclspec.NewBlock("gc", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"container": hclspec.NewDefault(
				hclspec.NewAttr("container", "bool", false),
				hclspec.NewLiteral("true"),
			),
		})), hclspec.NewLiteral(`{
			container = true
		}`)),
		// allow TaskRecover to start a still existing, stopped, container during client/driver restart
		"recover_stopped": hclspec.NewDefault(
			hclspec.NewAttr("recover_stopped", "bool", false),
			hclspec.NewLiteral("false"),
		),
		// optional extra_labels to append to all tasks for observability. Globs supported
		"extra_labels": hclspec.NewAttr("extra_labels", "list(string)", false),

		// logging options
		"logging": hclspec.NewDefault(hclspec.NewBlock("logging", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"driver":  hclspec.NewAttr("driver", "string", false),
			"options": hclspec.NewBlockAttrs("options", "string", false),
		})), hclspec.NewLiteral(`{
			driver = "nomad"
			options = {}
		}`)),

		// A list of sockets for the driver to manage
		"socket": hclspec.NewBlockList("socket", socketBodySpec),
		// the path to the podman api socket
		"socket_path": hclspec.NewAttr("socket_path", "string", false),

		// disable_log_collection indicates whether nomad should collect logs of podman
		// task containers.  If true, logs are not forwarded to nomad.
		"disable_log_collection": hclspec.NewAttr("disable_log_collection", "bool", false),
		"client_http_timeout":    hclspec.NewAttr("client_http_timeout", "string", false),
		"dns_servers":            hclspec.NewAttr("dns_servers", "list(string)", false),
	})

	// taskConfigSpec is the hcl specification for the driver config section of
	// a task within a job. It is returned in the TaskConfigSchema RPC.
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"apparmor_profile": hclspec.NewAttr("apparmor_profile", "string", false),
		"args":             hclspec.NewAttr("args", "list(string)", false),
		"auth": hclspec.NewBlock("auth", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"username": hclspec.NewAttr("username", "string", false),
			"password": hclspec.NewAttr("password", "string", false),
			"tls_verify": hclspec.NewDefault(
				hclspec.NewAttr("tls_verify", "bool", false),
				hclspec.NewLiteral("true"),
			),
		})),
		"auth_soft_fail": hclspec.NewAttr("auth_soft_fail", "bool", false),
		"command":        hclspec.NewAttr("command", "string", false),
		"cap_add":        hclspec.NewAttr("cap_add", "list(string)", false),
		"cap_drop":       hclspec.NewAttr("cap_drop", "list(string)", false),
		"selinux_opts":   hclspec.NewAttr("selinux_opts", "list(string)", false),
		"cpu_hard_limit": hclspec.NewAttr("cpu_hard_limit", "bool", false),
		"cpu_cfs_period": hclspec.NewAttr("cpu_cfs_period", "number", false),
		"devices":        hclspec.NewAttr("devices", "list(string)", false),
		"entrypoint":     hclspec.NewAttr("entrypoint", "any", false), // any for compat
		"working_dir":    hclspec.NewAttr("working_dir", "string", false),
		"hostname":       hclspec.NewAttr("hostname", "string", false),
		"image":          hclspec.NewAttr("image", "string", true),
		"image_pull_timeout": hclspec.NewDefault(
			hclspec.NewAttr("image_pull_timeout", "string", false),
			hclspec.NewLiteral(`"5m"`),
		),
		"init":      hclspec.NewAttr("init", "bool", false),
		"init_path": hclspec.NewAttr("init_path", "string", false),
		"labels":    hclspec.NewAttr("labels", "list(map(string))", false),
		"logging": hclspec.NewBlock("logging", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"driver":  hclspec.NewAttr("driver", "string", false),
			"options": hclspec.NewAttr("options", "list(map(string))", false),
		})),
		"memory_reservation": hclspec.NewAttr("memory_reservation", "string", false),
		"memory_swap":        hclspec.NewAttr("memory_swap", "string", false),
		"memory_swappiness":  hclspec.NewAttr("memory_swappiness", "number", false),
		"network_mode":       hclspec.NewAttr("network_mode", "string", false),
		"extra_hosts":        hclspec.NewAttr("extra_hosts", "list(string)", false),
		"pids_limit":         hclspec.NewAttr("pids_limit", "number", false),
		"port_map":           hclspec.NewAttr("port_map", "list(map(number))", false),
		"ports":              hclspec.NewAttr("ports", "list(string)", false),
		"privileged":         hclspec.NewAttr("privileged", "bool", false),
		"socket": hclspec.NewDefault(
			hclspec.NewAttr("socket", "string", false),
			hclspec.NewLiteral(`"default"`),
		),
		"sysctl":          hclspec.NewAttr("sysctl", "list(map(string))", false),
		"tmpfs":           hclspec.NewAttr("tmpfs", "list(string)", false),
		"tty":             hclspec.NewAttr("tty", "bool", false),
		"ulimit":          hclspec.NewAttr("ulimit", "list(map(string))", false),
		"volumes":         hclspec.NewAttr("volumes", "list(string)", false),
		"force_pull":      hclspec.NewAttr("force_pull", "bool", false),
		"readonly_rootfs": hclspec.NewAttr("readonly_rootfs", "bool", false),
		"userns":          hclspec.NewAttr("userns", "string", false),
		"shm_size":        hclspec.NewAttr("shm_size", "string", false),
	})
)

// TaskAuthConfig is the tasks authentication configuration
// (there is also auth_soft_fail on the top level)
type TaskAuthConfig struct {
	Username  string `codec:"username"`
	Password  string `codec:"password"`
	TLSVerify bool   `codec:"tls_verify"`
}

// GCConfig is the driver GarbageCollection configuration
type GCConfig struct {
	Container bool `codec:"container"`
}

// LoggingConfig is the driver logging configuration
// keep in sync with `TaskLoggingConfig`
type LoggingConfig struct {
	Driver  string            `codec:"driver"`
	Options map[string]string `codec:"options"`
}

// LoggingConfig is the tasks logging configuration
// keep in sync with `LoggingConfig`
type TaskLoggingConfig struct {
	Driver  string             `codec:"driver"`
	Options hclutils.MapStrStr `codec:"options"`
}

// Empty returns true if the logging configuration is not set.
func (l *TaskLoggingConfig) Empty() bool {
	return l.Driver == "" && len(l.Options) == 0
}

// VolumeConfig is the drivers volume specific configuration
type VolumeConfig struct {
	Enabled      bool   `codec:"enabled"`
	SelinuxLabel string `codec:"selinuxlabel"`
}

type PluginAuthConfig struct {
	FileConfig string `codec:"config"`
	Helper     string `codec:"helper"`
}

type PluginSocketConfig struct {
	Name       string `codec:"name"`
	SocketPath string `codec:"socket_path"`
}

// PluginConfig is the driver configuration set by the SetConfig RPC call
type PluginConfig struct {
	Auth                 PluginAuthConfig     `codec:"auth"`
	Volumes              VolumeConfig         `codec:"volumes"`
	GC                   GCConfig             `codec:"gc"`
	RecoverStopped       bool                 `codec:"recover_stopped"`
	DisableLogCollection bool                 `codec:"disable_log_collection"`
	Socket               []PluginSocketConfig `codec:"socket"`
	SocketPath           string               `codec:"socket_path"`
	ClientHttpTimeout    string               `codec:"client_http_timeout"`
	ExtraLabels          []string             `codec:"extra_labels"`
	DNSServers           []string             `codec:"dns_servers"`
	Logging              LoggingConfig        `codec:"logging"`
}

// LogWarnings will emit logs about known problematic configurations
func (c *PluginConfig) LogWarnings(logger hclog.Logger) {
	if c.RecoverStopped {
		logger.Error("WARNING - use of recover_stopped may cause Nomad agent to not start on system restarts")
	}
}

// TaskConfig is the driver configuration of a task within a job
type TaskConfig struct {
	ApparmorProfile   string             `codec:"apparmor_profile"`
	Args              []string           `codec:"args"`
	Auth              TaskAuthConfig     `codec:"auth"`
	AuthSoftFail      bool               `codec:"auth_soft_fail"`
	Ports             []string           `codec:"ports"`
	Tmpfs             []string           `codec:"tmpfs"`
	Volumes           []string           `codec:"volumes"`
	CapAdd            []string           `codec:"cap_add"`
	CapDrop           []string           `codec:"cap_drop"`
	SelinuxOpts       []string           `codec:"selinux_opts"`
	Command           string             `codec:"command"`
	Devices           []string           `codec:"devices"`
	Entrypoint        any                `codec:"entrypoint"` // any for compat
	WorkingDir        string             `codec:"working_dir"`
	Hostname          string             `codec:"hostname"`
	Image             string             `codec:"image"`
	ImagePullTimeout  string             `codec:"image_pull_timeout"`
	InitPath          string             `codec:"init_path"`
	Logging           TaskLoggingConfig  `codec:"logging"`
	Labels            hclutils.MapStrStr `codec:"labels"`
	MemoryReservation string             `codec:"memory_reservation"`
	MemorySwap        string             `codec:"memory_swap"`
	NetworkMode       string             `codec:"network_mode"`
	ExtraHosts        []string           `codec:"extra_hosts"`
	CPUCFSPeriod      uint64             `codec:"cpu_cfs_period"`
	MemorySwappiness  int64              `codec:"memory_swappiness"`
	PidsLimit         int64              `codec:"pids_limit"`
	PortMap           hclutils.MapStrInt `codec:"port_map"`
	Socket            string             `codec:"socket"`
	Sysctl            hclutils.MapStrStr `codec:"sysctl"`
	Ulimit            hclutils.MapStrStr `codec:"ulimit"`
	CPUHardLimit      bool               `codec:"cpu_hard_limit"`
	Init              bool               `codec:"init"`
	Tty               bool               `codec:"tty"`
	ForcePull         bool               `codec:"force_pull"`
	Privileged        bool               `codec:"privileged"`
	ReadOnlyRootfs    bool               `codec:"readonly_rootfs"`
	UserNS            string             `codec:"userns"`
	ShmSize           string             `codec:"shm_size"`
}
