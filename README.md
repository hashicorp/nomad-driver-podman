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
* Container log is forwarded to [Nomad logger](https://www.nomadproject.io/docs/commands/alloc/logs.html)
* Utilize podmans --init feature
* Set username or UID used for the specified command within the container (podman --user option).
* Fine tune memory usage: standard [Nomad memory resource](https://www.nomadproject.io/docs/job-specification/resources.html#memory) plus additional driver specific swap, swappiness and reservation parameters, OOM handling
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

- [Nomad](https://www.nomadproject.io/downloads.html) 0.12.9+
- Linux host with `podman` installed
- For rootless containers you need a system supporting cgroup V2 and a few other things, follow [this tutorial](https://github.com/containers/libpod/blob/master/docs/tutorials/rootless_tutorial.md)

You need a 2.x podman binary and a system socket activation unit,
see https://www.redhat.com/sysadmin/podmans-new-rest-api

Nomad agent, nomad-driver-podman and podman will reside on the same host, so you
do not have to worry about the ssh aspects of the podman api.

Ensure that Nomad can find the plugin, see [plugin_dir](https://www.nomadproject.io/docs/configuration/index.html#plugin_dir)

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

* **cap_add** - (Optional)  A list of Linux capabilities as strings to pass to --cap-add.

```
config {
  cap_add = [
    "SYS_TIME"
  ]
}
```

* **cap_drop** - (Optional)  A list of Linux capabilities as strings to pass to --cap-drop.

```
config {
  cap_add = [
    "MKNOD"
  ]
}
```

* **dns** - (Optional)  A list of dns servers. Replaces the default from podman binary and containers.conf.

```
config {
  dns = [
    "1.1.1.1"
  ]
}
```

## Example job

```
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

### Rootless on ubuntu

edit `/etc/default/grub` to enable cgroups v2
```
GRUB_CMDLINE_LINUX_DEFAULT="quiet cgroup_enable=memory swapaccount=1 systemd.unified_cgroup_hierarchy=1"
```

`sudo update-grub`

ensure that podman socket is running
```
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

```
crun -V
crun version 0.13.227-d38b
commit: d38b8c28fc50a14978a27fa6afc69a55bfdd2c11
spec: 1.0.0
+SYSTEMD +SELINUX +APPARMOR +CAP +SECCOMP +EBPF +YAJL
```

`nomad job run example.nomad`
```
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

```
$ podman ps
CONTAINER ID  IMAGE                           COMMAND       CREATED        STATUS            PORTS                                                 NAMES
2423ae3efa21  docker.io/library/redis:latest  redis-server  7 seconds ago  Up 6 seconds ago  127.0.0.1:21510->6379/tcp, 127.0.0.1:21510->6379/udp  redis-b640480f-4b93-65fd-7bba-c15722886395
```

### Local Development

#### Requirements

  - Vagrant >= 2.2
  - VirtualBox >= v6.0

#### Vagrant Environment Setup

```
# create the vm
vagrant up

# ssh into the vm
vagrant ssh

# Build the task driver plugin
sudo -E ./build.sh

# Copy the build nomad-driver-plugin executable to examples/plugins/
cp nomad-driver-podman examples/plugins/

# Start Nomad
nomad agent -config=examples/nomad/server.hcl 2>&1 > server.log &

# Run the client as sudo
sudo nomad agent -config=examples/nomad/client.hcl 2>&1 > client.log &

# Run a job
nomad job run examples/redis_ports.nomad

# Verify
nomad job status redis

sudo podman ps
```
