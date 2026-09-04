# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# Register a workload in Nomad's native service catalog and gate it on a
# health check.
#
# This job runs a small HTTP echo server, registers it as a Nomad service
# (provider = "nomad", so no Consul is required), and attaches an HTTP health
# check. Nomad only considers the allocation healthy once the check passes.
#
# See ./README.md for run and verification steps.

job "service-health" {
  datacenters = ["dc1"]
  type        = "service"

  group "api" {

    network {
      port "http" {
        to = 5678
      }
    }

    # Nomad-native service registration (no Consul dependency).
    service {
      name     = "echo-api"
      port     = "http"
      provider = "nomad"

      check {
        name     = "http-alive"
        type     = "http"
        path     = "/health"
        interval = "10s"
        timeout  = "2s"
      }
    }

    task "echo" {
      driver = "podman"

      config {
        image = "docker://hashicorp/http-echo:1.0"
        ports = ["http"]
        args = [
          "-listen=:5678",
          "-text=hello from a healthy service",
        ]
      }

      resources {
        cpu    = 100 # MHz
        memory = 64  # MB
      }
    }
  }
}
