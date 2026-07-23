# 03 - Sidecar with a Shared Network Namespace

Run two cooperating containers in a single network namespace: an application
that binds a **private** port on localhost, and a reverse-proxy sidecar that
joins the same namespace to reach it and re-expose it on the public port.

This is the same building block the repository's `nats_*` examples use, distilled
into a minimal, self-contained reverse-proxy pattern. See the main README's
[Network Configuration](../../../README.md#network-configuration) section for
the broader background.

## What this demonstrates

- Sharing one network namespace between tasks with
  `network_mode = "task:<name>"`.
- Ordering with lifecycle hooks: the main workload starts first, the
  `poststart` sidecar joins afterward.
- Keeping an internal port private (bound to `127.0.0.1`) while exposing only
  the proxy.
- Declaring all published ports on the namespace-owning task.

## Files

| File | Purpose |
| --- | --- |
| `app.nomad` | The job spec (app + proxy sidecar). |
| `templates/proxy.conf.tpl` | nginx reverse-proxy config. |

## Prerequisites

- A running Nomad agent with the `nomad-driver-podman` plugin.

## Run

The job uses HCL's `file("templates/...")` function, which resolves paths
relative to your **current working directory** (not the job file). Run from
this directory so the template is found:

```sh
cd examples/jobs/03-sidecar-network
nomad job run app.nomad
```

> If you run it from elsewhere you'll see `Unsuitable value: value must be
> known` — that just means `file()` couldn't locate the template; `cd` into
> this directory first.

## Verify

1. Confirm both tasks (app + proxy sidecar) are running.

   ```sh
   nomad job status sidecar-network
   ```

   ```
   Allocations
   ID        Node ID   Task Group  Version  Desired  Status   Created  Modified
   xxxxxxxx  xxxxxxxx  web         0        run      running  25s ago  10s ago
   ```

2. Read the dynamic port into a variable and send traffic to the published port.
   It reaches the proxy, which forwards it to the private app over localhost.

   ```sh
   addr=$(nomad alloc status -json $(nomad job allocs -json sidecar-network \
     | jq -r '.[0].ID') | jq -r '.Resources.Networks[0].DynamicPorts[0]
     | "127.0.0.1:\(.Value)"')

   curl "http://${addr}/"
   ```

   ```
   response from the private app
   ```

3. Confirm the app's own port `5678` is not publicly reachable — it stays
   private to the shared namespace.

   ```sh
   curl --max-time 2 "http://${addr%:*}:5678/" || echo "app port not exposed (expected)"
   ```

   ```
   curl: (7) Failed to connect to 127.0.0.1 port 5678 after 0 ms: Connection refused
   app port not exposed (expected)
   ```

## Adapt this for your own workload

- Replace the proxy with a metrics exporter that scrapes the app's private
  admin port (the original `nats_*` examples do exactly this).
- Add TLS termination in the proxy while the app stays plain HTTP on localhost.
- Add more sidecars (log shipper, auth proxy); each joins with
  `network_mode = "task:app"`.

## Cleanup

```sh
nomad job stop -purge sidecar-network
```

## Next

Continue to [04 - Service Discovery & Health Checks](../04-service-health/) to
register a workload in Nomad's service catalog and gate it on a health check.
