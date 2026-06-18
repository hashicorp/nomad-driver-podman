# 04 - Rootless & Hardened

Run a container as an unprivileged user against a **rootless** Podman socket,
with a defense-in-depth security profile. This is the recommended posture for
untrusted or internet-facing workloads.

## What this demonstrates

- Targeting a named rootless socket with the `socket` config option.
- Running the container process as a non-root `user`.
- Dropping all Linux capabilities (`cap_drop = ["ALL"]`).
- An immutable root filesystem (`readonly_rootfs = true`) with writable `tmpfs`
  mounts for the few paths the app needs.
- `security_opt = ["no-new-privileges"]` to block privilege escalation.

## Prerequisites

- A **rootless** Podman socket running for an unprivileged user, and the driver
  configured with a matching `socket` block. The provided
  [client.hcl](../../nomad/client.hcl) defines a socket named `app1`:
  ```hcl
  plugin "nomad-driver-podman" {
    config {
      socket {
        name        = "app1"
        socket_path = "unix://run/user/1337/podman/podman.sock"
      }
    }
  }
  ```
  The `1337` is a placeholder. Find your own rootless UID and use it for the
  socket path:
  ```sh
  id -u                                   # e.g. 1000 — your rootless UID
  systemctl --user status podman.socket   # make sure the socket is active
  ```
  Replace `1337` in the `socket_path` with that UID (e.g.
  `unix://run/user/1000/podman/podman.sock`) and reload the agent.
- A cgroup v2 host (required for rootless resource control). See
  [Rootless on ubuntu](../../../README.md#rootless-on-ubuntu) in the main README.

## Run

```sh
cd examples/jobs/04-rootless-hardened
nomad job run web.nomad
```

## Verify

```sh
nomad job status rootless-hardened

# Page is served.
addr=$(nomad alloc status -json $(nomad job allocs -json rootless-hardened \
  | jq -r '.[0].ID') | jq -r '.Resources.Networks[0].DynamicPorts[0]
  | "127.0.0.1:\(.Value)"')
curl -sI "http://${addr}/" | head -n1

# Run these as the rootless user that owns the `app1` socket so podman can see
# the container. The container runs rootless (owned by that user, not root).
podman ps --filter name=^web-

# Confirm the hardening took effect.
cid=$(podman ps -qf name=^web-)
podman inspect "$cid" --format \
  'ReadonlyRootfs={{.HostConfig.ReadonlyRootfs}} CapDrop={{.HostConfig.CapDrop}} SecurityOpt={{.HostConfig.SecurityOpt}}'
```

## Expected output

- `curl -I` returns `HTTP/1.1 200 OK`.
- `podman inspect` shows the read-only rootfs and dropped capabilities. Podman
  expands `cap_drop = ["ALL"]` into the concrete set of dropped capabilities, so
  the output looks like:
  ```
  ReadonlyRootfs=true CapDrop=[CAP_CHOWN CAP_DAC_OVERRIDE CAP_FOWNER CAP_FSETID CAP_KILL CAP_NET_BIND_SERVICE CAP_SETFCAP CAP_SETGID CAP_SETPCAP CAP_SETUID CAP_SYS_CHROOT] SecurityOpt=[no-new-privileges]
  ```
  The key signals are `ReadonlyRootfs=true`, a non-empty `CapDrop`, and
  `SecurityOpt=[no-new-privileges]`.
- Because the root filesystem is read-only, a write outside the tmpfs paths is
  rejected:
  ```sh
  podman exec "$cid" sh -c 'echo x > /etc/test' 2>&1
  # sh: 1: cannot create /etc/test: Read-only file system
  ```

## Adapt this for your own workload

- If your image insists on running as root, prefer an unprivileged variant
  (like `nginxinc/nginx-unprivileged`) or set a non-root `user`.
- Add back only the specific capabilities your app needs with `cap_add` instead
  of running with the full default set.
- Extend `tmpfs` to cover any additional writable paths your app requires.

## Cleanup

```sh
nomad job stop -purge rootless-hardened
```

## Next

That's the end of the learning path. Head back to the
[examples index](../README.md) to revisit any example or explore the additional
scenario jobs.
