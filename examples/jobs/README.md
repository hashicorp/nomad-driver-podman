# Podman Driver Examples

A curated, progressive set of Nomad + Podman examples. They start from a
single-container "hello world" and build up to a production-shaped deployment.
Each example is self-contained: a runnable job spec plus a README explaining
what it demonstrates, exact run/verify steps, and the expected output.

## Learning path

Work through them in order — each one introduces one new operational concept on
top of the previous.

| # | Example | What you learn | Key options |
| --- | --- | --- | --- |
| 1 | [Hello World](01-hello-world/) | The minimal runnable job | `image`, `ports`, `resources` |
| 2 | [Configuration & Templates](02-config-templates/) | Inject config/files into a stock image | `template`, `env`, `volumes` |
| 3 | [Persistent Storage](03-persistent-storage/) | Stateful data that survives restarts | bind-mount `volumes`, job `variable` |
| 4 | [Service Discovery & Health](04-service-health/) | Register a service, gate on health | `service` (Nomad provider), `check` |
| 5 | [Sidecar / Shared Network](05-sidecar-network/) | Multiple containers, one namespace | `network_mode = "task:..."`, `lifecycle` |
| 6 | [Rootless & Hardened](06-rootless-hardened/) | Unprivileged, locked-down containers | `socket`, `user`, `cap_drop`, `readonly_rootfs`, `security_opt` |
| 7 | [Production Scaling & Rollouts](07-production-scaling/) | Scale + safely upgrade a service | `count`, `update` (canary), `restart`, `reschedule` |

## Prerequisites (all examples)

- A running Nomad agent with the `nomad-driver-podman` plugin and `podman`
  installed on the client node. The repository ships a ready-to-use dev setup —
  see [Local Development](../../README.md#local-development) in the main README,
  or use the sample agent configs in [../nomad/](../nomad/).
- The CLI verification snippets use [`jq`](https://stedolan.github.io/jq/) and
  `curl`.

To bring up a local single-node agent quickly:

```sh
# from the repository root
make dev
cp ./build/nomad-driver-podman examples/plugins/
nomad agent -config=examples/nomad/server.hcl  > server.log 2>&1 &
sudo nomad agent -config=examples/nomad/client.hcl > client.log 2>&1 &
```

Then run any example, e.g.:

```sh
nomad job run examples/jobs/01-hello-world/redis.nomad
```

## Additional scenario examples

These older, feature-focused jobs live alongside the learning path and remain
valid references:

- Network namespace sharing variants: `nats_simple_pod.nomad`, `nats_pod.nomad`,
  `nats_group.nomad` (see the main README's
  [Network Configuration](../../README.md#network-configuration) section).
- Port mapping: `redis_ports.nomad`, `redis_deprecated.nomad` (legacy `port_map`).
- Rootless with a named socket and volumes: `rootless_kanboard.hcl`.
