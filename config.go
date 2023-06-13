// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/hashicorp/nomad/helper/pluginutils/hclutils"
	"github.com/hashicorp/nomad/plugins/shared/hclspec"
)

var (
	// configSpec is the hcl specification returned by the ConfigSchema RPC
	configSpec = hclspec.NewObject(map[string]*hclspec.Spec{
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
		// the path to the podman api socket
		"socket_path": hclspec.NewAttr("socket_path", "string", false),
		// disable_log_collection indicates whether nomad should collect logs of podman
		// task containers.  If true, logs are not forwarded to nomad.
		"disable_log_collection": hclspec.NewAttr("disable_log_collection", "bool", false),
		"client_http_timeout":    hclspec.NewAttr("client_http_timeout", "string", false),
	})

	// taskConfigSpec is the hcl specification for the driver config section of
	// a task within a job. It is returned in the TaskConfigSchema RPC.
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"apparmor_profile": hclspec.NewAttr("apparmor_profile", "string", false),
		"args":             hclspec.NewAttr("args", "list(string)", false),
		"auth": hclspec.NewBlock("auth", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"username": hclspec.NewAttr("username", "string", false),
			"password": hclspec.NewAttr("password", "string", false),
		})),
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
		"sysctl":             hclspec.NewAttr("sysctl", "list(map(string))", false),
		"tmpfs":              hclspec.NewAttr("tmpfs", "list(string)", false),
		"tty":                hclspec.NewAttr("tty", "bool", false),
		"ulimit":             hclspec.NewAttr("ulimit", "list(map(string))", false),
		"volumes":            hclspec.NewAttr("volumes", "list(string)", false),
		"force_pull":         hclspec.NewAttr("force_pull", "bool", false),
		"readonly_rootfs":    hclspec.NewAttr("readonly_rootfs", "bool", false),
		"userns":             hclspec.NewAttr("userns", "string", false),
	})
)

// AuthConfig is the tasks authentication configuration
type AuthConfig struct {
	Username string `codec:"username"`
	Password string `codec:"password"`
}

// GCConfig is the driver GarbageCollection configuration
type GCConfig struct {
	Container bool `codec:"container"`
}

// LoggingConfig is the tasks logging configuration
type LoggingConfig struct {
	Driver  string             `codec:"driver"`
	Options hclutils.MapStrStr `codec:"options"`
}

// VolumeConfig is the drivers volume specific configuration
type VolumeConfig struct {
	Enabled      bool   `codec:"enabled"`
	SelinuxLabel string `codec:"selinuxlabel"`
}

// PluginConfig is the driver configuration set by the SetConfig RPC call
type PluginConfig struct {
	Volumes              VolumeConfig `codec:"volumes"`
	GC                   GCConfig     `codec:"gc"`
	RecoverStopped       bool         `codec:"recover_stopped"`
	DisableLogCollection bool         `codec:"disable_log_collection"`
	SocketPath           string       `codec:"socket_path"`
	ClientHttpTimeout    string       `codec:"client_http_timeout"`
	ExtraLabels          []string     `codec:"extra_labels"`
}

// TaskConfig is the driver configuration of a task within a job
type TaskConfig struct {
	ApparmorProfile   string             `codec:"apparmor_profile"`
	Args              []string           `codec:"args"`
	Auth              AuthConfig         `codec:"auth"`
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
	Logging           LoggingConfig      `codec:"logging"`
	Labels            hclutils.MapStrStr `codec:"labels"`
	MemoryReservation string             `codec:"memory_reservation"`
	MemorySwap        string             `codec:"memory_swap"`
	NetworkMode       string             `codec:"network_mode"`
	ExtraHosts        []string           `codec:"extra_hosts"`
	CPUCFSPeriod      uint64             `codec:"cpu_cfs_period"`
	MemorySwappiness  int64              `codec:"memory_swappiness"`
	PidsLimit         int64              `codec:"pids_limit"`
	PortMap           hclutils.MapStrInt `codec:"port_map"`
	Sysctl            hclutils.MapStrStr `codec:"sysctl"`
	Ulimit            hclutils.MapStrStr `codec:"ulimit"`
	CPUHardLimit      bool               `codec:"cpu_hard_limit"`
	Init              bool               `codec:"init"`
	Tty               bool               `codec:"tty"`
	ForcePull         bool               `codec:"force_pull"`
	Privileged        bool               `codec:"privileged"`
	ReadOnlyRootfs    bool               `codec:"readonly_rootfs"`
	UserNS            string             `codec:"userns"`
}
