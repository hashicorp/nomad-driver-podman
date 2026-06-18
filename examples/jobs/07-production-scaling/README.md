# 07 - Production Scaling & Rollouts

The capstone example. It combines the earlier building blocks into a
production-shaped deployment: multiple instances, safe canary + rolling
upgrades, automatic recovery, and health-gated service discovery.

## What this demonstrates

- Running multiple instances with `count`.
- Safe rollouts with an `update` block: one canary first, rolling one at a time,
  `auto_revert` on failure, manual `auto_promote = false`.
- Automatic recovery with `restart` (in place) and `reschedule` (to other
  nodes).
- Native `service` discovery with an HTTP health `check`.
- Parameterizing the deployed version with a job `variable`.

## Prerequisites

- A running Nomad agent with the `nomad-driver-podman` plugin (Nomad 1.3+ for
  native services).

## Run

```sh
cd examples/jobs/07-production-scaling
nomad job run api.nomad
```

## Verify scaling and health

```sh
# Three allocations, all healthy.
nomad job status production-scaling

# Three service instances registered.
nomad service info scaled-api

# Round-robin across instances.
addr=$(nomad service info -json scaled-api | jq -r '.[0].Address + ":" + (.[0].Port|tostring)')
curl -s "http://${addr}/"
```

Expected: `nomad job status` shows `Healthy` with `3` running allocations, and
the `curl` returns `scaled-api v1`.

## Perform a safe rollout

Bump the version to trigger a canary deployment. Use `-detach` so the CLI
returns immediately instead of blocking on the canary, which waits for manual
promotion:

```sh
nomad job run -detach -var 'app_version=v2' api.nomad

# Watch the deployment: one canary is placed alongside the v1 allocations.
nomad job status production-scaling
nomad deployment status $(nomad job deployments -json production-scaling | jq -r '.[0].ID')
```

The deployment pauses with the canary running and waits for promotion. Inspect
the canary, then promote (or fail) it:

```sh
# Promote the canary -> Nomad rolls the change to the rest, one at a time.
nomad deployment promote $(nomad job deployments -json production-scaling | jq -r '.[0].ID')
```

## Expected output

- Before promotion: the canary serves `scaled-api v2` while the stable
  allocations still serve `scaled-api v1`.
- After promotion: all instances serve `scaled-api v2`.
- If the new version never becomes healthy within `healthy_deadline`,
  `auto_revert = true` rolls the job back to `v1` automatically.

## Test self-healing

```sh
# Kill one allocation's task; `restart`/`reschedule` brings capacity back.
nomad alloc stop $(nomad job allocs -json production-scaling | jq -r '.[0].ID')
nomad job status production-scaling   # Nomad replaces it automatically
```

## Adapt this for your own workload

- Increase `count` and `update.max_parallel` for larger fleets.
- Set `auto_promote = true` if you trust health checks to gate promotion.
- Combine with [02 - Config & Templates](../02-config-templates/) for app
  config and [06 - Rootless & Hardened](../06-rootless-hardened/) for the
  security profile.

## Cleanup

```sh
nomad job stop -purge production-scaling
```
