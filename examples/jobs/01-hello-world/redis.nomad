# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# The simplest possible Podman workload: a single Redis container.
#
# It demonstrates the three things almost every job needs:
#   - an image to run
#   - a dynamically allocated, published port
#   - CPU and memory resources
#
# See ./README.md for run and verification steps.

job "hello-world" {
  datacenters = ["dc1"]
  type        = "service"

  group "cache" {

    network {
      # Nomad picks a free host port and maps it to the container's 6379.
      port "redis" {
        to = 6379
      }
    }

    task "redis" {
      driver = "podman"

      config {
        image = "docker://redis:7"
        ports = ["redis"]
      }

      resources {
        cpu    = 200 # MHz
        memory = 128 # MB
      }
    }
  }
}
