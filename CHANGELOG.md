## Unreleased

FEATURES:

* config: Ability to configure dns server list

## 0.1.0

FEATURES:

* config: Add ability to configure container network_mode [[GH-33](https://github.com/hashicorp/nomad-driver-podman/issues/33)]
* network: (Consul Connect) Ability to accept a bridge network namespace from Nomad. [[GH-38](https://github.com/hashicorp/nomad-driver-podman/issues/38)]
* runtime: Ability to run podman rootless [[GH-42](https://github.com/hashicorp/nomad-driver-podman/issues/42)]
* config: Ability to specify varlink socket path [[GH-42](https://github.com/hashicorp/nomad-driver-podman/issues/42)]
* runtime: Conditionally set memory swappiness only if cgroupv1 is running [[GH-42](https://github.com/hashicorp/nomad-driver-podman/issues/42)]
* config: Ability to configure linux capabilities (cap_add/cap_drop) [[GH-44](https://github.com/hashicorp/nomad-driver-podman/issues/44)]

## 0.0.2 (June 11, 2020)

FEATURES:

* #8 podman --init support
* #14 oom killer handling, logging
* #10 support for --user option
* #15 configurable swap and memory reservation
* Add recover_stopped driver option

IMPROVEMENTS:

* varlink retries in case of socket issues

BUG FIXES:

* fixed problem with container naming conflict on startup/recovery
