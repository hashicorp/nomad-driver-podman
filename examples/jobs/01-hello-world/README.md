# 01 - Hello World

The smallest useful Podman job: a single [Redis](https://redis.io/) container
with a published port and resource limits. Start here if you have never run a
Nomad + Podman workload before.

## What this demonstrates

- The minimal anatomy of a job: `job` -> `group` -> `task`.
- Selecting an image with the `docker://` transport.
- Publishing a container port through a Nomad dynamic `port`.
- Reserving CPU and memory with a `resources` block.

## Prerequisites

- A running Nomad agent with the `nomad-driver-podman` plugin (see the
  [development setup](../../../README.md#local-development) in the main README).
- `podman` installed on the client node.

## Run

```sh
cd examples/jobs/01-hello-world
nomad job run redis.nomad
```

## Verify

1. Use `nomad job status` to confirm the job is running.

   ```sh
   nomad job status hello-world
   ```

   ```
   Allocations
   ID        Node ID   Task Group  Version  Desired  Status   Created  Modified
   xxxxxxxx  xxxxxxxx  cache       0        run      running  10s ago  3s ago
   ```

2. List all Podman containers matching the `^redis-` filter.

   ```sh
   podman ps --filter name=^redis-
   ```

   ```
   CONTAINER ID  IMAGE                      COMMAND       STATUS         PORTS                                                 NAMES
   a1b2c3d4e5f6  docker.io/library/redis:7  redis-server  Up 10 seconds  127.0.0.1:24089->6379/tcp, 127.0.0.1:24089->6379/udp  redis-xxxxxxxx-...
   ```

3. Read the dynamic host port Nomad assigned into a variable, then use it to
   `PING` the Redis server.

   ```sh
   addr=$(nomad alloc status -json $(nomad job allocs -json hello-world \
     | jq -r '.[0].ID') | jq -r '.Resources.Networks[0].DynamicPorts[0]
     | "\(env.NOMAD_IP // "127.0.0.1"):\(.Value)"')

   redis-cli -u "redis://${addr}" PING
   ```

   ```
   PONG
   ```

   If you don't have `redis-cli` locally, exec into the container instead — it
   returns the same `PONG`:

   ```sh
   podman exec -it $(podman ps -qf name=^redis-) redis-cli PING
   ```

   ```
   PONG
   ```

## Adapt this for your own workload

- Swap `docker://redis:7` for any image your workload needs.
- Change the container port in `network.port "redis" { to = 6379 }` to match the
  port your image listens on.
- Tune `resources.cpu` / `resources.memory` to fit the workload.

## Cleanup

```sh
nomad job stop -purge hello-world
```

## Next

Continue to [02 - Persistent Storage](../02-persistent-storage/) to run a
stateful container whose data survives task restarts.
