# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# Capstone: a production-shaped deployment that combines the earlier patterns.
#
#   - multiple instances (count) for availability
#   - safe rollouts with canary + rolling `update`
#   - automatic recovery via `restart` (in-place) and `reschedule` (to new nodes)
#   - native service discovery + HTTP health check
#   - resource reservations
#
# Roll out a new version safely by changing the rendered VERSION (via the
# `app_version` variable) and running `nomad job run` again; Nomad will deploy
# one canary first and only promote when healthy.
#
# See ./README.md for run, rollout and verification steps.

variable "app_version" {
  type        = string
  description = "Value rendered into the response body; bump it to trigger a rollout."
  default     = "v1"
}

job "production-scaling" {
  datacenters = ["dc1"]
  type        = "service"

  # Safe, observable rollouts.
  update {
    max_parallel     = 1     # update one allocation at a time
    canary           = 1     # deploy 1 canary before touching the rest
    min_healthy_time = "10s" # must stay healthy this long to count as healthy
    healthy_deadline = "2m"  # fail the deployment if not healthy in time
    auto_revert      = true  # roll back automatically on a failed deployment
    auto_promote     = false # require explicit promotion of the canary
  }

  group "api" {
    count = 3

    # Restart in place a few times before giving up on a node...
    restart {
      attempts = 3
      interval = "5m"
      delay    = "15s"
      mode     = "fail"
    }

    # ...then reschedule onto another eligible node with backoff.
    reschedule {
      delay          = "10s"
      delay_function = "exponential"
      max_delay      = "2m"
      unlimited      = true
    }

    network {
      port "http" {
        to = 5678
      }
    }

    service {
      name     = "scaled-api"
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
          "-text=scaled-api ${var.app_version}",
        ]
      }

      resources {
        cpu    = 100 # MHz
        memory = 64  # MB
      }
    }
  }
}
