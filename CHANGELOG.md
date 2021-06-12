## Unreleased

* config: Image registry authentication [[GH-71](https://github.com/hashicorp/nomad-driver-podman/issues/71)]
* config: Added tty option
* config: Support for sysctl configuration [[GH-82](https://github.com/hashicorp/nomad-driver-podman/issues/82)]
* config: Fixed a bug where we always pulled an image if image name has a transport prefix [[GH-88](https://github.com/hashicorp/nomad-driver-podman/pull/88)]
* config: Added labels option
* config: Add force_pull option

BUG FIXES:
* [[GH-93](https://github.com/hashicorp/nomad-driver-podman/issues/93)] use slirp4netns as default network mode if running rootless
* [[GH-92](https://github.com/hashicorp/nomad-driver-podman/issues/92)] parse rootless info correctly from podman 3.0.x struct

## 0.2.0

FEATURES:

* core: Support for Podman V2 HTTP API [[GH-51](https://github.com/hashicorp/nomad-driver-podman/issues/51)]
* config: Support for group allocated ports [[GH-74](https://github.com/hashicorp/nomad-driver-podman/issues/74)]
* config:  Ability to configure dns server list [[GH-54](https://github.com/hashicorp/nomad-driver-podman/issues/54)]
* runtime:  Add support for SignalTask [[GH-64](https://github.com/hashicorp/nomad-driver-podman/issues/64)]

BUG FIXES:

* [[GH-67](https://github.com/hashicorp/nomad-driver-podman/issues/67)] run container from oci-archive image


__BACKWARDS INCOMPATIBILITIES:__

* core: The driver no longer supports varlink communication with Podman
* config: `port_map` is deprecated in favor or group network ports and labels

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
