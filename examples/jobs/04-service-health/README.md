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

```sh
# The allocation should become "healthy" once the check passes.
nomad job status service-health

# The service is registered in Nomad's catalog.
nomad service info echo-api

# Hit the application through the registered address.
addr=$(nomad service info -json echo-api | jq -r '.[0].Address + ":" + (.[0].Port|tostring)')
curl -s "http://${addr}/"

# The health-check endpoint used by the check block.
curl -s "http://${addr}/health"
```

## Expected output

`nomad service info echo-api` lists one healthy instance:

```
Job ID          Address          Tags  Node ID   Alloc ID
service-health  127.0.0.1:24658  []    a1b2c3d4  e5f6a7b8
```

`nomad job status service-health` shows the deployment as healthy. The root
path returns the configured text:

```
hello from a healthy service
```

The `/health` endpoint (built into the `http-echo` image and used by the
`check` block) returns:

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

Continue to [05 - Sidecar with a Shared Network](../05-sidecar-network/) to run
multiple cooperating containers in one network namespace.
