// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"github.com/hashicorp/go-hclog"
	_ "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/plugins"
)

// trigger build

func main() {
	// Serve the plugin
	plugins.Serve(factory)
}

// factory returns a new instance of the LXC driver plugin
func factory(log hclog.Logger) interface{} {
	return NewPodmanDriver(log)
}
