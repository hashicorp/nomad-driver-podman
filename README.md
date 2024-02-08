Nomad podman Driver
==================

![](https://github.com/hashicorp/nomad-driver-podman/workflows/build/badge.svg)

Many thanks to [@towe75](https://github.com/towe75) and [Pascom](https://www.pascom.net/) for contributing
this plugin to Nomad!

## Features

* Use the jobs driver config to define the image for your container
* Start/stop containers with default or customer entrypoint and arguments
* [Nomad runtime environment](https://www.nomadproject.io/docs/runtime/environment.html) is populated
* Use Nomad alloc data in the container.
* Bind mount custom volumes into the container
* Publish ports
* Monitor the memory consumption
* Monitor CPU usage
* Task config cpu value is used to populate podman CpuShares
* Task config cores value is used to populate podman Cpuset
* Container log is forwarded to [Nomad logger](https://www.nomadproject.io/docs/commands/alloc/logs.html)
* Utilize podmans --init feature
* Set username or UID used for the specified command within the container (podman --user option).
* Fine tune memory usage: standard [Nomad memory resource](https://www.nomadproject.io/docs/job-specification/resources.html#memory) plus additional driver specific swap, swappiness and reservation parameters, OOM handling
* Supports rootless containers with cgroup V2
* Set DNS servers, searchlist and options via [Nomad dns parameters](https://www.nomadproject.io/docs/job-specification/network#dns-parameters)
* Support for nomad shared network namespaces and consul connect
* Quite flexible [network configuration](#network-configuration), allows to simply build pod-like structures within a nomad group

## Redis Example job

Here is a simple redis "hello world" Example:

```hcl
job "redis" {
  datacenters = ["dc1"]
  type        = "service"

  group "redis" {
    network {
      port "redis" { to = 6379 }
    }

    task "redis" {
      driver = "podman"

        config {
          image = "docker://redis"
          ports = ["redis"]
        }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
```

```sh
nomad run redis.nomad

==> Monitoring evaluation "9fc25b88"
    Evaluation triggered by job "redis"
    Allocation "60fdc69b" created: node "f6bccd6d", group "redis"
    Evaluation status changed: "pending" -> "complete"
==> Evaluation "9fc25b88" finished with status "complete"

podman ps

CONTAINER ID  IMAGE                           COMMAND               CREATED         STATUS             PORTS  NAMES
6d2d700cbce6  docker.io/library/redis:latest  docker-entrypoint...  16 seconds ago  Up 16 seconds ago         redis-60fdc69b-65cb-8ece-8554-df49321b3462
```

## Building The Driver from source

This project has a `go.mod` definition. So you can clone it to whatever directory you want.
It is not necessary to setup a go path at all.
Ensure that you use go 1.17 or newer.

```shell-session
git clone git@github.com:hashicorp/nomad-driver-podman
cd nomad-driver-podman
make dev
```

The compiled binary will be located at `./build/nomad-driver-podman`.

## Runtime dependencies

* [Nomad](https://www.nomadproject.io/downloads.html) 0.12.9+
* Linux host with `podman` installed
* For rootless containers you need a system supporting cgroup V2 and a few other things, follow [this tutorial](https://github.com/containers/libpod/blob/master/docs/tutorials/rootless_tutorial.md)

You need a 3.0.x podman binary and a system socket activation unit,
see <https://www.redhat.com/sysadmin/podmans-new-rest-api>

Nomad agent, nomad-driver-podman and podman will reside on the same host, so you
do not have to worry about the ssh aspects of the podman api.

Ensure that Nomad can find the plugin, see [plugin_dir](https://www.nomadproject.io/docs/configuration/index.html#plugin_dir)

## Driver Configuration

* volumes stanza:

  * enabled - Defaults to true. Allows tasks to bind host paths (volumes) inside their container.
  * selinuxlabel - Allows the operator to set a SELinux label to the allocation and task local bind-mounts to containers. If used with _volumes.enabled_ set to false, the labels will still be applied to the standard binds in the container.

```hcl
plugin "nomad-driver-podman" {
  config {
    volumes {
      enabled      = true
      selinuxlabel = "z"
    }
  }
}
```

* gc stanza:

  * container - Defaults to true. This option can be used to disable Nomad from removing a container when the task exits.

```hcl
plugin "nomad-driver-podman" {
  config {
    gc {
      container = false
    }
  }
}
```

* recover_stopped (bool) Defaults to false. Allows the driver to start and reuse a previously stopped container after
  a Nomad client restart.
  Consider a simple single node system and a complete reboot. All previously managed containers
  will be reused instead of disposed and recreated.

  WARNING - use of recover_stopped may cause Nomad agent to not start on system restarts. This setting has been left in place for compatibility.

```hcl
plugin "nomad-driver-podman" {
  config {
    recover_stopped = true
  }
}
```

* socket_path (string) Defaults to `"unix:///run/podman/podman.sock"` when running as root or a cgroup V1 system, and `"unix:///run/user/<USER_ID>/podman/podman.sock"` for rootless cgroup V2 systems

```hcl
plugin "nomad-driver-podman" {
  config {
    socket_path = "unix:///run/podman/podman.sock"
  }
}
```

* disable_log_collection (string) Defaults to `false`. Setting this to `true` will disable Nomad logs collection of Podman tasks. If you don't rely on nomad log capabilities and exclusively use host based log aggregation, you may consider this option to disable nomad log collection overhead. Beware to you also loose automatic log rotation.

```hcl
plugin "nomad-driver-podman" {
  config {
    disable_log_collection = false
  }
}
```

* extra_labels ([]string) Defaults to `[]`. Setting this will automatically append Nomad-related labels to Podman tasks. Supports glob matching such as `task*`. Possible values are:

```
job_name
job_id
task_group_name
task_name
namespace
node_name
node_id
```

```hcl
plugin "nomad-driver-podman" {
  config {
    extra_labels = ["job_name", "job_id", "task_group_name", "task_name", "namespace", "node_name", "node_id"]
  }
}
```

* logging stanza:

  * type - Defaults to `"nomad"`. See the task configuration for details.
  * options - Defaults to `{}`. See the task configuration for details.

* client_http_timeout (string) Defaults to `60s` default timeout used by http.Client requests

```hcl
plugin "nomad-driver-podman" {
  config {
    client_http_timeout = "60s"
  }
```

## Task Configuration

* **image** - The image to run. Accepted transports are `docker` (default if missing), `oci-archive` and `docker-archive`. Images reference as [short-names](https://github.com/containers/image/blob/master/docs/containers-registries.conf.5.md#short-name-aliasing) will be treated according to user-configured preferences.

```hcl
config {
  image = "docker://redis"
}
```

* **auth** - (Optional) Authenticate to the image registry using a static credential. `tls_verify` can be disabled for insecure registries.

```hcl
config {
  image = "your.registry.tld/some/image"
  auth {
    username   = "someuser"
    password   = "sup3rs3creT"
    tls_verify = true
  }
}
```

* **entrypoint** - (Optional) A string list overriding the image's entrypoint. Defaults to the entrypoint set in the image.

```hcl
config {
  entrypoint = [
    "/bin/bash",
    "-c"
  ]
}
```

* **command** - (Optional) The command to run when starting the container.

```hcl
config {
  command = "some-command"
}
```

* **args** - (Optional) A list of arguments to the optional command. If no _command_ is specified, the arguments are passed directly to the container.

```hcl
config {
  args = [
    "arg1",
    "arg2",
  ]
}
```

* **working_dir** - (Optional) The working directory for the container. Defaults to the default set in the image.

```hcl
config {
  working_dir = "/data"
}
```

* **volumes** - (Optional) A list of host_path:container_path:options strings to bind host paths to container paths. Named volumes are not supported.

```hcl
config {
  volumes = [
    "/some/host/data:/container/data:ro,noexec"
  ]
}
```

* **tmpfs** - (Optional) A list of /container_path strings for tmpfs mount points. See podman run --tmpfs options for details.

```hcl
config {
  tmpfs = [
    "/var"
  ]
}
```

* **devices** - (Optional) A list of `host-device[:container-device][:permissions]` definitions.
Each entry adds a host device to the container. Optional permissions can be used to specify device permissions, it is combination of r for read, w for write, and m for mknod(2). See podman documentation for more details.

```hcl
config {
  devices = [
    "/dev/net/tun"
  ]
}
```

* **hostname** -  (Optional) The hostname to assign to the container. When launching more than one of a task (using count) with this option set, every container the task starts will have the same hostname.

* **Forwarding and Exposing Ports** - (Optional) See [Docker Driver Configuration](https://www.nomadproject.io/docs/drivers/docker.html#forwarding-and-exposing-ports) for details.

* **init** - Run an init inside the container that forwards signals and reaps processes.

```hcl
config {
  init = true
}
```

* **init_path** - Path to the container-init binary.

```hcl
config {
  init = true
  init_path = /usr/libexec/podman/catatonit
}
```

* **user** - Run the command as a specific user/uid within the container. See [Task configuration](https://www.nomadproject.io/docs/job-specification/task.html#user)

```hcl
user = nobody

config {
}

```

* **logging** - Configure logging. See also plugin option **disable_log_collection**

`driver = "nomad"` (default) Podman redirects its combined stdout/stderr logstream directly to a Nomad fifo.
Benefits of this mode are: zero overhead, don't have to worry about log rotation at system or Podman level. Downside: you cannot easily ship the logstream to a log aggregator plus stdout/stderr is multiplexed into a single stream..

```hcl
config {
  logging = {
    driver = "nomad"
  }
}
```

`driver = "journald"` The container log is forwarded from Podman to the journald on your host. Next, it's pulled by the Podman API back from the journal into the Nomad fifo (controllable by **disable_log_collection**)
Benefits: all containers can log into the host journal, you can ship a structured stream incl. metadata to your log aggregator. No log rotation at Podman level. You can add additional tags to the journal.
Drawbacks: a bit more overhead, depends on Journal (will not work on WSL2). You should configure some rotation policy for your Journal.
Ensure you're running Podman 3.1.0 or higher because of bugs in older versions.

```hcl
config {
  logging = {
    driver = "journald"
    options = {
      "tag" = "redis"
    }
  }
}
```

* **memory_reservation** - Memory soft limit (nit = b (bytes), k (kilobytes), m (megabytes), or g (gigabytes))

After setting memory reservation, when the system detects memory contention or low memory, containers are forced to restrict their consumption to their reservation. So you should always set the value below --memory, otherwise the hard limit will take precedence. By default, memory reservation will be the same as memory limit.

```hcl
config {
  memory_reservation = "100m"
}
```

* **memory_swap** - A limit value equal to memory plus swap. The swap LIMIT should always be larger than the [memory value](https://www.nomadproject.io/docs/job-specification/resources.html#memory).

Unit can be b (bytes), k (kilobytes), m (megabytes), or g (gigabytes). If you don't specify a unit, b is used. Set LIMIT to -1 to enable unlimited swap.

```hcl
config {
  memory_swap = "180m"
}
```

* **memory_swappiness** - Tune a container's memory swappiness behavior. Accepts an integer between 0 and 100.

```hcl
config {
  memory_swappiness = 60
}
```

* **network_mode** - Set the [network mode](http://docs.podman.io/en/latest/markdown/podman-run.1.html#options) for the container.

By default the task uses the network stack defined in the task group, see [network Stanza](https://www.nomadproject.io/docs/job-specification/network). If the groups network behavior is also undefined, it will fallback to `bridge` in rootful mode or `slirp4netns` for rootless containers.

* `bridge`: create a network stack on the default podman bridge.
* `none`: no networking
* `host`: use the Podman host network stack. Note: the host mode gives the
  container full access to local system services such as D-bus and is therefore
  considered insecure
* `slirp4netns`: use `slirp4netns` to create a user network stack. This is the
  default for rootless containers. Podman currently does not support it for root
  containers [issue](https://github.com/containers/libpod/issues/6097).
* `container:id`: reuse another podman containers network stack
* `task:name-of-other-task`: join the network of another task in the same allocation.

```hcl
config {
  network_mode = "bridge"
}
```

* **cap_add** - (Optional)  A list of Linux capabilities as strings to pass to --cap-add.

```hcl
config {
  cap_add = [
    "SYS_TIME"
  ]
}
```

* **cap_drop** - (Optional)  A list of Linux capabilities as strings to pass to --cap-drop.

```hcl
config {
  cap_add = [
    "MKNOD"
  ]
}
```

* **selinux_opts** - (Optional)  A list of process labels the container will use.

```
config {
  selinux_opts = [
    "type:my_container.process"
  ]
}
```

* **sysctl** - (Optional)  A key-value map of sysctl configurations to set to the containers on start.

```hcl
config {
  sysctl = {
    "net.core.somaxconn" = "16384"
  }
}
```

* **privileged** - (Optional)  true or false (default). A privileged container turns off the security features that isolate the container from the host. Dropped Capabilities, limited devices, read-only mount points, Apparmor/SELinux separation, and Seccomp filters are all disabled.

* **tty** - (Optional)  true or false (default). Allocate a pseudo-TTY for the container.

* **labels** - (Optional)  Set labels on the container.

```hcl
config {
  labels = {
    "nomad" = "job"
  }
}
```

* **apparmor_profile** - (Optional) Name of a apparmor profile to be used instead of the default profile. The special value `unconfined` disables apparmor for this container:

```
config {
  apparmor_profile = "your-profile"
}
```

* **force_pull** - (Optional)  true or false (default). Always pull the latest image on container start.

```hcl
config {
  force_pull = true
}
```

* **readonly_rootfs** - (Optional)  true or false (default). Mount the rootfs as read-only.

```hcl
config {
  readonly_rootfs = true
}
```

* **ulimit** - (Optional) A key-value map of ulimit configurations to set to the containers to start.

```hcl
config {
  ulimit {
    nproc = "4242"
    nofile = "2048:4096"
  }
```

* **userns** - (Optional) Set the [user namespace mode](https://docs.podman.io/en/latest/markdown/podman-run.1.html#userns-mode) for the container.

```hcl
config {
  userns = "keep-id:uid=200,gid=210"
}
```

* **pids_limit** - (Optional) An integer value that specifies the pid limit for the container.

```hcl
config {
  pids_limit = 64
}
```

* **image_pull_timeout** - (Optional) time duration for your pull timeout (default to 5m).

```
config {
  image_pull_timeout = "5m"
}
```

## Network Configuration

[nomad lifecycle hooks](https://www.nomadproject.io/docs/job-specification/lifecycle) combined with the drivers `network_mode` allows very flexible network namespace definitions. This feature does not build upon the native podman pod structure but simply reuses the networking namespace of one container for other tasks in the same group.

A typical example is a network server and a metric exporter or log shipping sidecar. The metric exporter needs access to i.E. a private monitoring Port which should not be exposed the the network and thus is usually bound to localhost.

The repository includes three different examples jobs for such a setup. All of them will start a [nats](https://nats.io/) server and a [prometheus-nats-exporter](https://github.com/nats-io/prometheus-nats-exporter) using different approaches.

You can use `curl` to proof that the job is working correctly and that you can get prometheus metrics:

`curl http://your-machine:7777/metrics`

### 2 Task setup, server defines the network

See `examples/jobs/nats_simple_pod.nomad`

Here, the _server_ task is started as main workload and the _exporter_ runs as a poststart sidecar.
Because of that, Nomad guarantees that the server is started first and thus the exporter can
easily join the servers network namespace via `network_mode = "task:server"`.

Note, that the _server_ configuration file binds the _http_port_ to localhost.

Be aware that ports must be defined in the parent network namespace, here _server_.

### 3 Task setup, a pause container defines the network

See `examples/jobs/nats_pod.nomad`

A slightly different setup is demonstrated in this job. It reassembles more closely the idea of a _pod_ by starting a
pause task, named _pod_ via a prestart/sidecar [hook](https://www.nomadproject.io/docs/job-specification/lifecycle).

Next, the main workload, _server_ is started and joins the network namespace by using the `network_mode = "task:pod"` stanza.
Finally, Nomad starts the poststart/sidecar _exporter_ which also joins the network.

Note that all ports must be defined on the _pod_ level.

### 2 Task setup, shared Nomad network namespace

See `examples/jobs/nats_group.nomad`

This example is very different. Both _server_ and _exporter_ join a network namespace which is created and managed
by Nomad itself. See [nomad network stanza](https://www.nomadproject.io/docs/job-specification/network) to get started with this generic approach.

## Rootless on ubuntu

edit `/etc/default/grub` to enable cgroups v2

```sh
GRUB_CMDLINE_LINUX_DEFAULT="quiet cgroup_enable=memory swapaccount=1 systemd.unified_cgroup_hierarchy=1"
```

`sudo update-grub`

ensure that podman socket is running

```console
$ systemctl --user status podman.socket
* podman.socket - Podman API Socket
     Loaded: loaded (/usr/lib/systemd/user/podman.socket; disabled; vendor preset: disabled)
     Active: active (listening) since Sat 2020-10-31 19:21:29 CET; 22h ago
   Triggers: * podman.service
       Docs: man:podman-system-service(1)
     Listen: /run/user/1000/podman/podman.sock (Stream)
     CGroup: /user.slice/user-1000.slice/user@1000.service/podman.socket
```

ensure that you have a recent version of [crun](https://github.com/containers/crun/)

```console
$ crun -V
crun version 0.13.227-d38b
commit: d38b8c28fc50a14978a27fa6afc69a55bfdd2c11
spec: 1.0.0
+SYSTEMD +SELINUX +APPARMOR +CAP +SECCOMP +EBPF +YAJL
```

`nomad job run example.nomad`

```hcl
job "example" {
  datacenters = ["dc1"]
  type        = "service"

  group "cache" {
    count = 1
    restart {
      attempts = 2
      interval = "30m"
      delay    = "15s"
      mode     = "fail"
    }
    network {
      port "redis" { to = 6379 }
    }
    task "redis" {
      driver = "podman"

      config {
        image = "redis"
        ports = ["redis"]
      }

      resources {
        cpu    = 500 # 500 MHz
        memory = 256 # 256MB
      }
    }
  }
}
```

verify `podman ps`

```console
$ podman ps
CONTAINER ID  IMAGE                           COMMAND       CREATED        STATUS            PORTS                                                 NAMES
2423ae3efa21  docker.io/library/redis:latest  redis-server  7 seconds ago  Up 6 seconds ago  127.0.0.1:21510->6379/tcp, 127.0.0.1:21510->6379/udp  redis-b640480f-4b93-65fd-7bba-c15722886395
```

## Local Development

### Requirements

* Vagrant >= 2.2
* VirtualBox >= v6.0

### Vagrant Environment Setup

```sh
# create the vm
vagrant up

# ssh into the vm
vagrant ssh
````

Running a Nomad dev agent with the Podman plugin:

```
# Build the task driver plugin
make dev

# Copy the build nomad-driver-plugin executable to examples/plugins/
cp ./build/nomad-driver-podman examples/plugins/

# Start Nomad
nomad agent -config=examples/nomad/server.hcl 2>&1 > server.log &

# Run the client as sudo
sudo nomad agent -config=examples/nomad/client.hcl 2>&1 > client.log &

# Run a job
nomad job run examples/jobs/redis_ports.nomad

# Verify
nomad job status redis

sudo podman ps
```

Running the tests:

```
# Start the Podman server
systemctl --user start podman.socket

# Run the tests
CI=1 ./build/bin/gotestsum --junitfile ./build/test/result.xml -- -timeout=15m . ./api
```
