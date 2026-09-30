# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# Run a container rootlessly with a hardened security profile.
#
# This job targets a rootless Podman socket (the "app1" socket defined in
# examples/nomad/client.hcl) and applies defense-in-depth settings:
#   - non-root user inside the container
#   - all Linux capabilities dropped
#   - read-only root filesystem (with a writable tmpfs for runtime scratch)
#   - no-new-privileges
#
# See ./README.md for run, prerequisites and verification steps.

job "rootless-hardened" {
  datacenters = ["dc1"]
  type        = "service"

  group "web" {

    network {
      port "http" {
        to = 8080
      }
    }

    task "web" {
      driver = "podman"

      # Run the container process as an unprivileged user.
      user = "101:101" # nginx unprivileged uid:gid in nginxinc/nginx-unprivileged

      config {
        # Use the rootless socket named "app1" from the client config.
        socket = "app1"

        image = "docker://nginxinc/nginx-unprivileged:1.27"
        ports = ["http"]

        # Hardening:
        cap_drop        = ["ALL"] # drop every Linux capability
        readonly_rootfs = true    # immutable root filesystem
        security_opt    = ["no-new-privileges"]

        # nginx needs a few writable paths; provide them as tmpfs since the root
        # filesystem is read-only.
        tmpfs = [
          "/tmp",
          "/var/cache/nginx",
          "/var/run",
        ]
      }

      resources {
        cpu    = 100 # MHz
        memory = 64  # MB
      }
    }
  }
}
