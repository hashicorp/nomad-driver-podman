# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

job "kanboard" {
  datacenters = ["dc1"]

  group "kanboard" {
    network {
      mode = "host"
      port "kanboard" {
        static = 8000
        to     = 80
      }
    }

    task "deploy_kanboard" {
      driver = "podman"
      env {
        ENABLE_URL_REWRITE = "false"
      }
      config {
        image      = "docker.io/kanboard/kanboard:v1.2.39"
        ports      = ["kanboard"]
        privileged = false
        socket     = "app1" # Will use the rootless socket defined in client.hcl called "app1"
        logging = {
          driver = "journald"
        }
        volumes = [
          "/srv/kanboard/data:/var/www/app/data:noexec"
        ]
      }
    }
  }
}
