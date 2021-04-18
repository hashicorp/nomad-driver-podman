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
			hclspec.NewLiteral("true"),
		),
		// the path to the podman api socket
		"socket_path": hclspec.NewAttr("socket_path", "string", false),
	})

	// taskConfigSpec is the hcl specification for the driver config section of
	// a task within a job. It is returned in the TaskConfigSchema RPC
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"args": hclspec.NewAttr("args", "list(string)", false),
		"auth": hclspec.NewBlock("auth", false, hclspec.NewObject(map[string]*hclspec.Spec{
			"username": hclspec.NewAttr("username", "string", false),
			"password": hclspec.NewAttr("password", "string", false),
		})),
		"command":            hclspec.NewAttr("command", "string", false),
		"cap_add":            hclspec.NewAttr("cap_add", "list(string)", false),
		"cap_drop":           hclspec.NewAttr("cap_drop", "list(string)", false),
		"entrypoint":         hclspec.NewAttr("entrypoint", "string", false),
		"working_dir":        hclspec.NewAttr("working_dir", "string", false),
		"hostname":           hclspec.NewAttr("hostname", "string", false),
		"image":              hclspec.NewAttr("image", "string", true),
		"init":               hclspec.NewAttr("init", "bool", false),
		"init_path":          hclspec.NewAttr("init_path", "string", false),
		"log_driver":         hclspec.NewAttr("log_driver", "string", false),
		"memory_reservation": hclspec.NewAttr("memory_reservation", "string", false),
		"memory_swap":        hclspec.NewAttr("memory_swap", "string", false),
		"memory_swappiness":  hclspec.NewAttr("memory_swappiness", "number", false),
		"network_mode":       hclspec.NewAttr("network_mode", "string", false),
		"port_map":           hclspec.NewAttr("port_map", "list(map(number))", false),
		"ports":              hclspec.NewAttr("ports", "list(string)", false),
		"sysctl":             hclspec.NewAttr("sysctl", "list(map(string))", false),
		"tmpfs":              hclspec.NewAttr("tmpfs", "list(string)", false),
		"tty":                hclspec.NewAttr("tty", "bool", false),
		"volumes":            hclspec.NewAttr("volumes", "list(string)", false),
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

// VolumeConfig is the drivers volume specific configuration
type VolumeConfig struct {
	Enabled      bool   `codec:"enabled"`
	SelinuxLabel string `codec:"selinuxlabel"`
}

// PluginConfig is the driver configuration set by the SetConfig RPC call
type PluginConfig struct {
	Volumes        VolumeConfig `codec:"volumes"`
	GC             GCConfig     `codec:"gc"`
	RecoverStopped bool         `codec:"recover_stopped"`
	SocketPath     string       `codec:"socket_path"`
}

// TaskConfig is the driver configuration of a task within a job
type TaskConfig struct {
	Args              []string           `codec:"args"`
	Auth              AuthConfig         `codec:"auth"`
	Ports             []string           `codec:"ports"`
	Tmpfs             []string           `codec:"tmpfs"`
	Volumes           []string           `codec:"volumes"`
	CapAdd            []string           `codec:"cap_add"`
	CapDrop           []string           `codec:"cap_drop"`
	Command           string             `codec:"command"`
	Entrypoint        string             `codec:"entrypoint"`
	WorkingDir        string             `codec:"working_dir"`
	Hostname          string             `codec:"hostname"`
	Image             string             `codec:"image"`
	InitPath          string             `codec:"init_path"`
	LogDriver         string             `codec:"log_driver"`
	MemoryReservation string             `codec:"memory_reservation"`
	MemorySwap        string             `codec:"memory_swap"`
	NetworkMode       string             `codec:"network_mode"`
	MemorySwappiness  int64              `codec:"memory_swappiness"`
	PortMap           hclutils.MapStrInt `codec:"port_map"`
	Sysctl            hclutils.MapStrStr `codec:"sysctl"`
	Init              bool               `codec:"init"`
	Tty               bool               `codec:"tty"`
}
