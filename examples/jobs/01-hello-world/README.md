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
nomad job run redis.nomad
```

## Verify

```sh
# 1. The allocation should be "running".
nomad job status hello-world

# 2. Podman should show the container.
podman ps --filter name=redis

# 3. Talk to Redis through the published port. Read the host:port that Nomad
#    assigned, then PING it.
addr=$(nomad alloc status -json $(nomad job allocs -json hello-world \
  | jq -r '.[0].ID') | jq -r '.Resources.Networks[0].DynamicPorts[0]
  | "\(env.NOMAD_IP // "127.0.0.1"):\(.Value)"')

redis-cli -u "redis://${addr}" PING
```

If you don't have `redis-cli` locally, exec into the container instead:

```sh
podman exec -it $(podman ps -qf name=redis) redis-cli PING
```

## Expected output

`nomad job status hello-world` shows one allocation in the `running` state:

```
Allocations
ID        Node ID   Task Group  Version  Desired  Status   Created  Modified
xxxxxxxx  xxxxxxxx  cache       0        run      running  10s ago  3s ago
```

`podman ps` lists the container with a `redis-server` command and a published
port, and the `PING` returns:

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

Continue to [02 - Configuration & Templates](../02-config-templates/) to learn
how to inject configuration files and environment variables into a container.
