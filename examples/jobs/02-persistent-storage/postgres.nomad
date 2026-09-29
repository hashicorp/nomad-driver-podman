# Copyright IBM Corp. 2019, 2025
# SPDX-License-Identifier: MPL-2.0

# A stateful workload whose data survives task restarts and rescheduling.
#
# PostgreSQL stores its data under /var/lib/postgresql/data. We bind-mount a
# directory from the host into that path so the database files outlive the
# container lifecycle.
#
# NOTE: a host bind-mount pins this allocation to whichever node holds the data.
# For multi-node durable storage use a CSI volume or a replicated database; this
# example keeps it simple and single-node.
#
# See ./README.md for run and verification steps.

variable "host_data_dir" {
  type        = string
  description = "Host directory used to persist the PostgreSQL data files."
  default     = "/srv/podman-examples/postgres"
}

job "persistent-storage" {
  datacenters = ["dc1"]
  type        = "service"

  group "db" {

    network {
      port "postgres" {
        to = 5432
      }
    }

    task "postgres" {
      driver = "podman"

      env {
        POSTGRES_USER = "demo"
        POSTGRES_DB   = "demo"
        # For a real deployment, source this from a Nomad Variable or Vault
        # instead of hardcoding it. See the "Adapt" section in README.md.
        POSTGRES_PASSWORD = "demo-password"
        # Subdir so Postgres owns an empty dir even if host_data_dir has a
        # lost+found or other entries.
        PGDATA = "/var/lib/postgresql/data/pgdata"
      }

      config {
        image = "docker://postgres:16"
        ports = ["postgres"]

        # Persist the database files on the host. The source directory must
        # already exist (Podman does not create it). The ":U" option tells
        # Podman to chown the mount to the container's user, which is required
        # for rootless Postgres to own its data directory.
        volumes = [
          "${var.host_data_dir}:/var/lib/postgresql/data:U",
        ]
      }

      resources {
        cpu    = 500 # MHz
        memory = 256 # MB
      }
    }
  }
}
