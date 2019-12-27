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
	})

	// taskConfigSpec is the hcl specification for the driver config section of
	// a task within a job. It is returned in the TaskConfigSchema RPC
	taskConfigSpec = hclspec.NewObject(map[string]*hclspec.Spec{
		"image":              hclspec.NewAttr("image", "string", true),
		"command":            hclspec.NewAttr("command", "string", false),
		"hostname":           hclspec.NewAttr("hostname", "string", false),
		"init":               hclspec.NewAttr("init", "bool", false),
		"init_path":          hclspec.NewAttr("init_path", "string", false),
		"args":               hclspec.NewAttr("args", "list(string)", false),
		"volumes":            hclspec.NewAttr("volumes", "list(string)", false),
		"port_map":           hclspec.NewAttr("port_map", "list(map(number))", false),
		"memory_reservation": hclspec.NewAttr("memory_reservation", "string", false),
		"memory_swap":        hclspec.NewAttr("memory_swap", "string", false),
		"memory_swappiness":  hclspec.NewAttr("memory_swappiness", "number", false),
	})
)

// GCConfig is the driver GarbageCollection configuration
type GCConfig struct {
	Container bool `codec:"container"`
}

// VolumeConfig
type VolumeConfig struct {
	Enabled      bool   `codec:"enabled"`
	SelinuxLabel string `codec:"selinuxlabel"`
}

// PluginConfig is the driver configuration set by the SetConfig RPC call
type PluginConfig struct {
	Volumes VolumeConfig `codec:"volumes"`
	GC      GCConfig     `codec:"gc"`
}

// TaskConfig is the driver configuration of a task within a job
type TaskConfig struct {
	Image             string             `codec:"image"`
	Command           string             `codec:"command"`
	Args              []string           `codec:"args"`
	Volumes           []string           `codec:"volumes"`
	Hostname          string             `codec:"hostname"`
	Init              bool               `codec:"init"`
	InitPath          string             `codec:"init_path"`
	PortMap           hclutils.MapStrInt `codec:"port_map"`
	MemoryReservation string             `codec:"memory_reservation"`
	MemorySwap        string             `codec:"memory_swap"`
	MemorySwappiness  int64              `codec:"memory_swappiness"`
}
