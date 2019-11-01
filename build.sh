#!/bin/bash -e

me=$(readlink -f "${BASH_SOURCE[0]}")
project=$(dirname "$me")
cd "$project"

mkdir -p build

# ensure to build in a isolated GOPATH in order to get predictable dependencies
export GOPATH=$project/build 

go install github.com/varlink/go/cmd/varlink-go-interface-generator

go generate github.com/pascomnet/nomad-driver-podman/iopodman
go build