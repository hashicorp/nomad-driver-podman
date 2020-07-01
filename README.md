Nomad podman Driver
==================


[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/hashicorp/nomad-driver-podman/blob/master/LICENSE)
![](https://github.com/hashicorp/nomad-driver-podman/workflows/build/badge.svg)

*THIS IS A WORK IN PROGRESS PLUGIN*. Do not run it in production!
Contributions are welcome, of course.

Many thanks to [@towe75](https://github.com/towe75) and [Pascom](https://www.pascom.net/) for contributing
this plugin to Nomad!

## Features

* Use the jobs driver config to define the image for your container
* Start/stop containers with default or customer entrypoint and arguments
* [Nomad runtime environment](https://www.nomadproject.io/docs/runtime/environment.html) is populated
* Use nomad alloc data in the container.
* Bind mount custom volumes into the container
* Publish ports
* Monitor the memory consumption
* Monitor CPU usage
* Task config cpu value is used to populate podman CpuShares
* Container log is forwarded to [Nomad logger](https://www.nomadproject.io/docs/commands/alloc/logs.html)
* Utilize podmans --init feature
* Set username or UID used for the specified command within the container (podman --user option).
* Fine tune memory usage: standard [nomad memory resource](https://www.nomadproject.io/docs/job-specification/resources.html#memory) plus additional driver specific swap, swappiness and reservation parameters, OOM handling
* Supports rootless containers with cgroup V2


## Building The Driver from source

This project has a go.mod definition. So you can clone it to whatever directory you want.
It is not necessary to setup a go path at all.
Ensure that you use go 1.13 or newer.

```sh
$ git clone git@github.com:hashicorp/nomad-driver-podman
cd nomad-driver-podman
./build.sh
```

## Runtime dependencies

- [Nomad](https://www.nomadproject.io/downloads.html) 0.9+
- Linux host with `podman` installed
- For rootless containers you need a system supporting cgroup V2 and a few other things, follow [this tutorial](https://github.com/containers/libpod/blob/master/docs/tutorials/rootless_tutorial.md)

You need a varlink enabled podman binary and a system socket activation unit,
see https://podman.io/blogs/2019/01/16/podman-varlink.html.

nomad agent, nomad-driver-podman and podman will reside on the same host, so you
do not have to worry about the ssh aspects of podman varlink.

Ensure that nomad can find the plugin, see [plugin_dir](https://www.nomadproject.io/docs/configuration/index.html#plugin_dir)

## Driver Configuration

* volumes stanza:

  * enabled - Defaults to true. Allows tasks to bind host paths (volumes) inside their container.
  * selinuxlabel - Allows the operator to set a SELinux label to the allocation and task local bind-mounts to containers. If used with _volumes.enabled_ set to false, the labels will still be applied to the standard binds in the container.

```
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

```
plugin "nomad-driver-podman" {
  config {
    gc {
      container = false
    }
  }
}
```

* recover_stopped (bool) Defaults to true. Allows the driver to start and resuse a previously stopped container after
  a Nomad client restart.
  Consider a simple single node system and a complete reboot. All previously managed containers
  will be reused instead of disposed and recreated.

```
plugin "nomad-driver-podman" {
  config {
    recover_stopped = false
  }
}
```

* socket_path (string) Defaults to `"unix://run/podman/io.podman"` when running as root or a cgroup V1 system, and `"unix://run/user/<USER_ID>/podman/io.podman"` for rootless cgroup V2 systems


```
plugin "nomad-driver-podman" {
  config {
    socket_path = "unix://run/podman/io.podman"
  }
}
```
## Task Configuration

* **image** - The image to run,

```
config {
  image = "docker://redis"
}
```

* **entrypoint** - (Optional) The entrypoint for the container. Defaults to the entrypoint set in the image.

```
config {
  entrypoint = "/entrypoint.sh"
}
```

* **command** - (Optional) The command to run when starting the container.

```
config {
  command = "some-command"
}
```

* **args** - (Optional) A list of arguments to the optional command. If no *command* is specified, the arguments are passed directly to the container.

```
config {
  args = [
    "arg1",
    "arg2",
  ]
}
```

* **working_dir** - (Optional) The working directory for the container. Defaults to the default set in the image.

```
config {
  working_dir = "/data"
}
```

* **volumes** - (Optional) A list of host_path:container_path strings to bind host paths to container paths.

```
config {
  volumes = [
    "/some/host/data:/container/data"
  ]
}
```

* **tmpfs** - (Optional) A list of /container_path strings for tmpfs mount points. See podman run --tmpfs options for details.

```
config {
  tmpfs = [
    "/var"
  ]
}
```

* **hostname** -  (Optional) The hostname to assign to the container. When launching more than one of a task (using count) with this option set, every container the task starts will have the same hostname.

* **Forwarding and Exposing Ports** - (Optional) See [Docker Driver Configuration](https://www.nomadproject.io/docs/drivers/docker.html#forwarding-and-exposing-ports) for details.

* **init** - Run an init inside the container that forwards signals and reaps processes.

```
config {
  init = true
}
```

* **init_path** - Path to the container-init binary.

```
config {
  init = true
  init_path = /usr/libexec/podman/catatonit
}
```

* **user** - Run the command as a specific user/uid within the container. See [Task configuration](https://www.nomadproject.io/docs/job-specification/task.html#user)

```
user = nobody

config {
}

```

* **memory_reservation** - Memory soft limit (nit = b (bytes), k (kilobytes), m (megabytes), or g (gigabytes))

After setting memory reservation, when the system detects memory contention or low memory, containers are forced to restrict their consumption to their reservation. So you should always set the value below --memory, otherwise the hard limit will take precedence. By default, memory reservation will be the same as memory limit.

```
config {
  memory_reservation = "100m"
}
```

* **memory_swap** - A limit value equal to memory plus swap. The swap LIMIT should always be larger than the [memory value](https://www.nomadproject.io/docs/job-specification/resources.html#memory).

Unit can be b (bytes), k (kilobytes), m (megabytes), or g (gigabytes). If you don't specify a unit, b is used. Set LIMIT to -1 to enable unlimited swap.

```
config {
  memory_swap = "180m"
}
```

* **memory_swappiness** - Tune a container's memory swappiness behavior. Accepts an integer between 0 and 100.

```
config {
  memory_swappiness = 60
}
```

* **network_mode** - Set the [network mode](http://docs.podman.io/en/latest/markdown/podman-run.1.html#options) for the container.

- `bridge`: (default for rootful) create a network stack on the default bridge
- `none`: no networking
- `container:id`: reuse another container's network stack
- `host`: use the Podman host network stack. Note: the host mode gives the
  container full access to local system services such as D-bus and is therefore
  considered insecure
- `slirp4netns`: use `slirp4netns` to create a user network stack. This is the
  default for rootless containers. Podman currently does not support it for root
  containers [issue](https://github.com/containers/libpod/issues/6097).

```
config {
  network_mode = "bridge"
}
```

## Example job

```
job "redis" {
  datacenters = ["dc1"]
  type        = "service"

  group "redis" {
    task "redis" {
      driver = "podman"

        config {
          image = "docker://redis"
          port_map {
              redis = 6379
          }
        }

      resources {
        cpu    = 500
        memory = 256
        network {
          mbits = 20
          port "redis" {}
        }
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
