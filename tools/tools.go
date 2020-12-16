//+build tools

// Package tools anonymously imports packages of tools used to build nomad-driver-podman.
// See the GNUMakefile for 'go get` commands.
package tools

import (
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/hashicorp/go-hclog/hclogvet"
	_ "gotest.tools/gotestsum"
)
