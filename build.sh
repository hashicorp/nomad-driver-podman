#!/bin/bash -e

me=$(readlink -f "${BASH_SOURCE[0]}")
project=$(dirname "$me")
cd "$project"

mkdir -p build
mkdir -p build/test

# ensure to build in a isolated GOPATH in order to get predictable dependencies
export GOPATH=$project/build

go env
go version

go install github.com/varlink/go/cmd/varlink-go-interface-generator
go install gotest.tools/gotestsum

go generate github.com/hashicorp/nomad-driver-podman/iopodman
go build
