# 02 - Configuration & Templates

Run an off-the-shelf [nginx](https://nginx.org/) image but supply its
configuration at deploy time instead of building a custom image. This is the
single most common real-world pattern: take a stock image and feed it your
config.

## What this demonstrates

- Rendering files with Nomad `template` blocks into the task's `local/`
  directory.
- Interpolating environment variables (`SITE_NAME`) and Nomad runtime metadata
  (`NOMAD_ALLOC_ID`, `node.unique.name`) into a rendered file.
- Bind-mounting rendered files into the container with the `volumes` config,
  using relative paths that resolve against the task directory.
- `change_mode = "restart"` so the task restarts when its config changes.

## Files

| File | Purpose |
| --- | --- |
| `nginx.nomad` | The job spec. |
| `templates/nginx.conf.tpl` | nginx server config, mounted at `/etc/nginx/conf.d/default.conf`. |
| `templates/index.html.tpl` | Landing page that interpolates env + runtime metadata. |

## Prerequisites

- A running Nomad agent with the `nomad-driver-podman` plugin.
- The driver's `volumes` feature must be enabled (it is by default):
  ```hcl
  plugin "nomad-driver-podman" {
    config {
      volumes { enabled = true }
    }
  }
  ```

## Run

The job uses HCL's `file("templates/...")` function, which resolves paths
relative to your **current working directory** (not the job file). Run from
this directory so the templates are found:

```sh
cd examples/jobs/02-config-templates
nomad job run nginx.nomad
```

> If you run it from elsewhere you'll see `Unsuitable value: value must be
> known` — that just means `file()` couldn't locate the template; `cd` into
> this directory first.

## Verify

```sh
nomad job status config-templates

# Fetch the rendered page through the published port.
addr=$(nomad alloc status -json $(nomad job allocs -json config-templates \
  | jq -r '.[0].ID') | jq -r '.Resources.Networks[0].DynamicPorts[0]
  | "127.0.0.1:\(.Value)"')

curl -s "http://${addr}/"
```

## Expected output

The `curl` returns the rendered HTML, with the allocation ID and node name
filled in:

```html
<!doctype html>
<html>
  <head>
    <title>Nomad + Podman</title>
  </head>
  <body>
    <h1>Hello from Nomad + Podman!</h1>
    <p>Served by allocation <code>a1b2c3d4-...</code></p>
    <p>Running on node <code>podmanclient</code></p>
  </body>
</html>
```

## Adapt this for your own workload

- Replace the two templates with your application's real config files.
- Add more `template` blocks for additional files (TLS certs, app configs).
- Use `change_mode = "signal"` with `change_signal = "SIGHUP"` if your app
  reloads config without a full restart.

## Cleanup

```sh
nomad job stop -purge config-templates
```

## Next

Continue to [03 - Persistent Storage](../03-persistent-storage/) to run a
stateful workload whose data survives restarts.
