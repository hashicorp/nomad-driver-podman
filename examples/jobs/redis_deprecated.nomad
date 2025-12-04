# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# The following example redis job uses the old deprecated port map

job "redis" {
  datacenters = ["dc1"]
  type        = "service"

  group "cache" {
    task "redis" {
      driver = "podman"

      env {
        foo = "bar"
      }
      config {
        port_map {
          redis = 6379
        }

        image = "docker://redis"
      }
      resources {
        network {
          port "redis" {}
        }
      }
    }
  }
}

