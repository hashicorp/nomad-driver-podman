# 03 - Persistent Storage

Run a stateful [PostgreSQL](https://www.postgresql.org/) database whose data
survives task restarts, by bind-mounting a host directory into the container.

## What this demonstrates

- Persisting container data with a host bind-mount in the `volumes` config.
- Configuring a stateful image through `env` (`POSTGRES_*`).
- Using a job `variable` to parameterize the host data path.
- The trade-off of host bind-mounts: the allocation is pinned to the node that
  holds the data.

## Prerequisites

- A running Nomad agent with the `nomad-driver-podman` plugin.
- The driver's `volumes` feature enabled (default).
- A client node where Podman may create/write the host data directory
  (`/srv/podman-examples/postgres` by default).

## Run

The bind-mount source directory must exist before the job starts (Podman does
not create it). Create it first, then run the job.

**Rootful Podman** (default path):

```sh
sudo mkdir -p /srv/podman-examples/postgres
nomad job run postgres.nomad
```

**Rootless Podman** — point the job at a directory your user owns:

```sh
mkdir -p "$HOME/pg-data"
nomad job run -var "host_data_dir=$HOME/pg-data" postgres.nomad
```

The job mounts the directory with the `:U` option so Podman chowns it to the
Postgres user inside the container.

## Verify

```sh
nomad job status persistent-storage

# Create a table and a row.
cid=$(podman ps -qf name=postgres)
podman exec -i "$cid" psql -U demo -d demo -c \
  "CREATE TABLE IF NOT EXISTS t (msg text); INSERT INTO t VALUES ('it persists');"

# Restart the task, then confirm the row is still there.
nomad alloc restart $(nomad job allocs -json persistent-storage | jq -r '.[0].ID')
sleep 5
cid=$(podman ps -qf name=postgres)
podman exec -i "$cid" psql -U demo -d demo -c "SELECT msg FROM t;"
```

## Expected output

After the restart the previously inserted row is still returned:

```
     msg
-------------
 it persists
(1 row)
```

You can also confirm the data files exist on the host:

```sh
sudo ls /srv/podman-examples/postgres/pgdata
# base  global  pg_wal  postgresql.conf  ...
```

## Adapt this for your own workload

- **Secrets**: replace the inline `POSTGRES_PASSWORD` with a `template` that
  reads a [Nomad Variable](https://developer.hashicorp.com/nomad/docs/concepts/variables)
  or a Vault secret, e.g.
  ```hcl
  template {
    data        = "POSTGRES_PASSWORD={{ with nomadVar \"nomad/jobs/persistent-storage\" }}{{ .password }}{{ end }}"
    destination = "secrets/db.env"
    env         = true
  }
  ```
- **Durable multi-node storage**: swap the host bind-mount for a
  [CSI volume](https://developer.hashicorp.com/nomad/docs/other-specifications/volume)
  so the workload can reschedule across nodes.

## Cleanup

```sh
nomad job stop -purge persistent-storage
# Remove persisted data (irreversible):
sudo rm -rf /srv/podman-examples/postgres
```

## Next

Continue to [04 - Service Discovery & Health Checks](../04-service-health/) to
register a service and gate traffic on health checks.
