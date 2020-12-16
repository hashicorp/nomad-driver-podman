#!/bin/bash -e

me=$(readlink -f "${BASH_SOURCE[0]}")
project=$(dirname "$me")
cd "$project"

mkdir -p build
mkdir -p build/test

# ensure to build in a isolated GOPATH in order to get predictable dependencies
export GOPATH=$project/build 

go build
