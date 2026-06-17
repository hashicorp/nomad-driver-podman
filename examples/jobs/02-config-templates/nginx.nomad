# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# Render configuration into a container at deploy time.
#
# This job runs nginx, but instead of baking a custom image it injects:
#   - an nginx.conf rendered from a Nomad `template`
#   - an index.html rendered from a `template` that interpolates env + runtime
#     metadata (the allocation ID and node name)
#
# The rendered files live in the task's `local/` directory. We bind-mount them
# to the paths nginx expects using relative `volumes` sources (the driver
# resolves them against the task directory on the host).
#
# See ./README.md for run and verification steps.

job "config-templates" {
  datacenters = ["dc1"]
  type        = "service"

  group "web" {

    network {
      port "http" {
        to = 8080
      }
    }

    task "nginx" {
      driver = "podman"

      env {
        # Demonstrates passing plain environment variables to the container and
        # reusing them inside a template (see index.html.tpl).
        SITE_NAME = "Nomad + Podman"
      }

      # nginx main config, rendered from a template in this directory.
      template {
        data        = file("templates/nginx.conf.tpl")
        destination = "local/nginx.conf"
        change_mode = "restart"
      }

      # A tiny landing page that interpolates env vars and Nomad runtime metadata.
      template {
        data        = file("templates/index.html.tpl")
        destination = "local/html/index.html"
        change_mode = "restart"
      }

      config {
        image = "docker://nginx:1.27"
        ports = ["http"]

        # Mount the rendered files into the locations nginx expects. Relative
        # volume sources are resolved against the task directory on the host,
        # whose `local/` subdir is where the templates above are written.
        volumes = [
          "local/nginx.conf:/etc/nginx/conf.d/default.conf:ro",
          "local/html:/usr/share/nginx/html:ro",
        ]
      }

      resources {
        cpu    = 200 # MHz
        memory = 128 # MB
      }
    }
  }
}
