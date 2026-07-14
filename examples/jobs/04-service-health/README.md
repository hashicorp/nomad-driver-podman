# 04 - Service Discovery & Health Checks

Register a workload in Nomad's built-in service catalog and attach a health
check, so other workloads can discover it and Nomad can tell whether it is
actually serving traffic. Uses the **native Nomad provider**, so no Consul is
required.

## What this demonstrates

- A `service` block with `provider = "nomad"` for service discovery.
- An HTTP `check` that gates allocation health.
- Discovering the service from the CLI and from a template.

## Prerequisites

- A running Nomad agent with the `nomad-driver-podman` plugin.
- Nomad 1.3+ (native service discovery).

## Run

```sh
cd examples/jobs/04-service-health
nomad job run echo.nomad
```

## Verify

The allocation becomes healthy once the check passes:

```sh
nomad job status service-health
```

```
Allocations
ID        Node ID   Task Group  Version  Desired  Status   Created  Modified
xxxxxxxx  xxxxxxxx  api         0        run      running  15s ago  2s ago
```

Confirm the service is registered in Nomad's catalog:

```sh
nomad service info echo-api
```

```
Job ID          Address          Tags  Node ID   Alloc ID
service-health  127.0.0.1:24658  []    a1b2c3d4  e5f6a7b8
```

Hit the application through the registered address:

```sh
addr=$(nomad service info -json echo-api | jq -r '.[0].Address + ":" + (.[0].Port|tostring)')
curl "http://${addr}/"
```

```
hello from a healthy service
```

Check the `/health` endpoint used by the `check` block (built into the
`http-echo` image):

```sh
curl "http://${addr}/health"
```

```
{"status":"ok"}
```

To see the check in action, watch how an unhealthy task is held back: change
`-listen=:5678` to a different port (so the check can't reach it) and re-run —
the deployment will not mark the allocation healthy.

## Adapt this for your own workload

- Point `check.path` at your app's real readiness/liveness endpoint.
- Add `tags` to the `service` block for routing (e.g. with a load balancer).
- Discover the service from another job using a template:
  ```hcl
  template {
    data = <<EOH
  {{ range nomadService "echo-api" }}
  UPSTREAM={{ .Address }}:{{ .Port }}
  {{ end }}
  EOH
    destination = "local/upstream.env"
    env         = true
  }
  ```

## Cleanup

```sh
nomad job stop -purge service-health
```

## Next

Continue to [05 - Rootless & Hardened](../05-rootless-hardened/) to run
containers as an unprivileged user with a tightened security profile.
