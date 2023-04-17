#!/usr/bin/env bash

# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

set -euo pipefail

echo "====== Install tools from apt"
apt-get update
apt-get install -y ca-certificates podman curl build-essential

echo "====== Install catatonit"
curl -v -L -o /usr/local/bin/catatonit https://github.com/openSUSE/catatonit/releases/download/v0.1.7/catatonit.x86_64
chmod +x /usr/local/bin/catatonit

echo "====== Podman info"
podman version
podman info

echo "====== Setup archives"
podman pull alpine:3
podman save --format docker-archive --output /tmp/docker-archive alpine:3
podman save --format oci-archive --output /tmp/oci-archive alpine:3
podman image rm alpine:3

echo "===== Configure registries"
cat <<EOF > /etc/containers/registries.conf
unqualified-search-registries = ["docker.io", "quay.io"]
[[registry]]
location = "localhost:5000"
insecure = true
EOF
