# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# Two cooperating containers in one shared network namespace.
#
# Pattern: an application binds a *private* port on localhost, and a sidecar
# (here an nginx reverse proxy) joins the same network namespace to reach it
# over localhost. Only the proxy's port is published to the network; the app
# port is never exposed directly.
#
# This reuses Nomad lifecycle hooks plus the driver's `network_mode = "task:..."`
# to share one namespace, without using native Podman pods. See the main
# README's "Network Configuration" section for background.
#
# See ./README.md for run and verification steps.

job "sidecar-network" {
  datacenters = ["dc1"]
  type        = "service"

  group "web" {

    # Ports must be declared on the task that *owns* the namespace (the main
    # workload, "app", started first). Only the proxy's port is published.
    network {
      port "http" {
        to = 8080
      }
    }

    # Main workload: started first, owns the network namespace.
    # It listens on localhost:5678 only - never exposed to the network.
    task "app" {
      driver = "podman"

      config {
        image = "docker://hashicorp/http-echo:1.0"
        ports = ["http"]
        args = [
          "-listen=127.0.0.1:5678",
          "-text=response from the private app",
        ]
      }

      resources {
        cpu    = 100 # MHz
        memory = 64  # MB
      }
    }

    # Sidecar: starts after the app and joins its network namespace, so it can
    # reach the app on localhost:5678 and re-expose it on the public port 8080.
    task "proxy" {
      driver = "podman"

      lifecycle {
        hook    = "poststart"
        sidecar = "true"
      }

      template {
        data        = file("templates/proxy.conf.tpl")
        destination = "local/proxy.conf"
        change_mode = "restart"
      }

      config {
        image = "docker://nginx:1.27"

        # Join the app's network namespace.
        network_mode = "task:app"

        volumes = [
          "local/proxy.conf:/etc/nginx/conf.d/default.conf:ro",
        ]
      }

      resources {
        cpu    = 100 # MHz
        memory = 64  # MB
      }
    }
  }
}
