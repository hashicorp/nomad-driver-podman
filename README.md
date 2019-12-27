Nomad podman Driver
==================


[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/pascomnet/nomad-driver-podman/blob/master/LICENSE)
![](https://github.com/pascomnet/nomad-driver-podman/workflows/build/badge.svg)

*THIS IS A PROOF OF CONCEPT PLUGIN*. Do not run it in production!
Contributions are welcome, of course.

## Features

* use the jobs driver config to define the image for your container
* start/stop containers with default or customer entrypoint and arguments
* [Nomad runtime environment](https://www.nomadproject.io/docs/runtime/environment.html) is populated
* use nomad alloc data in the container.
* bind mount custom volumes into the container
* publish ports
* monitor the memory consumption
* monitor CPU usage
* task config cpu value is used to populate podman CpuShares
* Container log is forwarded to [Nomad logger](https://www.nomadproject.io/docs/commands/alloc/logs.html) 
* utilize podmans --init feature
* set username or UID used for the specified command within the container (podman --user option).
* fine tune memory usage: standard [nomad memory resource](https://www.nomadproject.io/docs/job-specification/resources.html#memory) plus additional driver specific swap, swappiness and reservation parameters, OOM handling


## Building The Driver from source

This project has a go.mod definition. So you can clone it to whatever directory you want.
It is not necessary to setup a go path at all.

```sh
$ git clone git@github.com:pascomnet/nomad-driver-podman
cd nomad-driver-podman
./build.sh
```

## Runtime dependencies

- [Nomad](https://www.nomadproject.io/downloads.html) 0.9+
- Linux host with `podman` installed

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


## Task Configuration

* **image** - The image to run, 

```
config {
  image = "docker://redis"
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

* **volumes** - (Optional) A list of host_path:container_path strings to bind host paths to container paths. 

```
config {
  volumes = [
    "/some/host/data:/container/data"
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
