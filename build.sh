#!/bin/bash -e

# FIXME
project=$(pwd)

mkdir -p build
GOPATH=$project/build go install github.com/varlink/go/cmd/varlink
GOPATH=$project/build go install github.com/varlink/go/cmd/varlink-go-interface-generator

go generate github.com/pascomnet/nomad-driver-podman/iopodman
go build
cp -v nomad-driver-podman $project/../dist/