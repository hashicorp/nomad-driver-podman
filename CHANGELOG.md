## UNRELEASED

* api: Support pasta network_mode and use it as default for rootless containers on Podman 5.0+ [[GH-476](https://github.com/hashicorp/nomad-driver-podman/pull/476)]
* build: Update Nomad version to 1.11.0 [[GH-478](https://github.com/hashicorp/nomad-driver-podman/pull/478)]
* build: Updated to Go 1.25.5 [[GH-477](https://github.com/hashicorp/nomad-driver-podman/pull/477)]

## 0.6.3 (June 18, 2025)

IMPROVEMENTS:

* api: Added support for ipv6_address, ipv4_address, static_ips, and static_mac options to networking configuration [[GH-359](https://github.com/hashicorp/nomad-driver-podman/pull/359)]
* api: Address a backwards incompatible change in Podman 4.0 for secondary IP addresses [[GH-443](https://github.com/hashicorp/nomad-driver-podman/pull/443)]
* build: Update Nomad version to 1.10.0 [[GH-431](https://github.com/hashicorp/nomad-driver-podman/pull/431)]
* build: Updated to Go 1.24.2 [[GH-430](https://github.com/hashicorp/nomad-driver-podman/pull/430)]
* config: Adds support for oom_score_adj [[GH-425](https://github.com/hashicorp/nomad-driver-podman)]

BUG FIXES:

* fingerprint: Check for nil pointers in SystemInfo client response. [[GH-412](https://github.com/hashicorp/nomad-driver-podman/issues/412)]
* recovery: Fixes driver failing to recover task [[GH-433](https://github.com/hashicorp/nomad-driver-podman/pull/433)]
* runtime: Fixes parsing options when creating a tmpfs mount [[GH-439](https://github.com/hashicorp/nomad-driver-podman/pull/439)]
* logging: Fixes bug where driver did not honor default logging configuration [[GH-411](https://github.com/hashicorp/nomad-driver-podman/pull/411)]
* volumes: Fixes volume deletion when deleting container [[GH-424](https://github.com/hashicorp/nomad-driver-podman/pull/424)]

## 0.6.2 (December 16, 2024)

IMPROVEMENTS:

* config: Allow setting `security_opt` option. [[GH-382](https://github.com/hashicorp/nomad-driver-podman/pull/382)]
* config: Add `socket` stanza to allow multiple Podman sockets to be used. [[GH-371](https://github.com/hashicorp/nomad-driver-podman/pull/371)]

## 0.6.1 (July 15, 2024)

BUG FIXES:

* build: Removed CGO dependency accidentally introduced in 0.6.0

## 0.6.0 (July 12, 2024)

IMPROVEMENTS:

 * api: Address a backwards incompatible change in Podman 5.0 [[GH-332](https://github.com/hashicorp/nomad-driver-podman/issues/332)]
 * config: Add `logging` options to the plugin configuration [[GH-285](https://github.com/hashicorp/nomad-driver-podman/pull/285)]
 * build: Updated Nomad 1.8.0 [[GH-347](https://github.com/hashicorp/nomad-driver-podman/pull/347)]
 * config: Add granular control of SELinux labels for host mounts [[GH-321]](https://github.com/hashicorp/nomad-driver-podman/pull/321)

## 0.5.2 (February 5, 2024)

SECURITY:

 * deps: Updated runc to 1.1.12 to address CVE-2024-21626313 [[GH-313](https://github.com/hashicorp/nomad-driver-podman/pull/313)]

IMPROVEMENTS:

 * build: Updated to Go 1.21.5 [[GH-303](https://github.com/hashicorp/nomad-driver-podman/pull/303)]

## 0.5.1 (August 14, 2023)

* api: Address a backwards incompatible change in Podman 4.6.0 preventing jobs from restarting [[GH-278](https://github.com/hashicorp/nomad-driver-podman/pull/278)]

## 0.5.0 (July 19, 2023)

IMPROVEMENTS:

* config: Add support for auth helper or config file in plugin configuration [[GH-265](https://github.com/hashicorp/nomad-driver-podman/pull/265)]
* config: Add support for setting tlsVerify in task configuration [[GH-262](https://github.com/hashicorp/nomad-driver-podman/pull/262)]
* config: Add support for setting extra_hosts in task configuration [[GH-255](https://github.com/hashicorp/nomad-driver-podman/pull/255)]

BUG FIXES:

* config: Set recover_stopped to false by default since Client may hang if enabled [[GH-260](https://github.com/hashicorp/nomad-driver-podman/pull/260)]
* runtime: Correctly configure cpuset on system using cgroups v2 [[GH-252](https://github.com/hashicorp/nomad-driver-podman/pull/252)]
* runtime: Fixed a bug where driver would panic on systems without a nobody user in /etc/passwd [[GH-266](https://github.com/hashicorp/nomad-driver-podman/pull/266)]

## 0.4.2 (March 24, 2023)

IMPROVEMENTS:

* config: Add `extra_labels` option [[GH-215](https://github.com/hashicorp/nomad-driver-podman/pull/215)]
* config: Allow setting `pids_limit` option. [[GH-203](https://github.com/hashicorp/nomad-driver-podman/pull/203)]
* config: Allow setting `userns` option. [[GH-212](https://github.com/hashicorp/nomad-driver-podman/pull/212)]
* config: Allow setting `entrypoint` as a list of strings. [[GH-209](https://github.com/hashicorp/nomad-driver-podman/pull/209)]
* runtime: Set mount propagation from TaskConfig [[GH-204](https://github.com/hashicorp/nomad-driver-podman/pull/204)]

BUG FIXES:

* driver: Fixed a bug that caused `image_pull_timeout` to the capped by the value of `client_http_timeout` [[GH-218](https://github.com/hashicorp/nomad-driver-podman/pull/218)]

## 0.4.1 (November 15, 2022)

FEATURES:

* config: Set custom apparmor profile or disable apparmor. [[GH-188](https://github.com/hashicorp/nomad-driver-podman/pull/188)]

IMPROVEMENTS:

* config: Add `selinux_opts` option [[GH-139](https://github.com/hashicorp/nomad-driver-podman/pull/139)]
* perf: Use ping api instead of system info for fingerprinting [[GH-186](https://github.com/hashicorp/nomad-driver-podman/pull/186)]
* runtime: Prevent concurrent image pulls of same imageRef [[GH-159](https://github.com/hashicorp/nomad-driver-podman/pull/159)]

BUG FIXES:

* runtime: Don't apply SELinux labels to volumes of privileged containers [[GH-196](https://github.com/hashicorp/nomad-driver-podman/pull/196)]
* runtime: Fixed a bug caused by a Podman API change that prevented the task driver to detect stopped containers [[GH-183](https://github.com/hashicorp/nomad-driver-podman/pull/183)]

## 0.4.0 (July 14, 2022)

FEATURES:

* config: Map host devices into container. [[GH-41](https://github.com/hashicorp/nomad-driver-podman/pull/41)]
* config: Stream logs via API, support journald log driver. [[GH-99](https://github.com/hashicorp/nomad-driver-podman/pull/99)]
* config: Privileged containers. [[GH-137](https://github.com/hashicorp/nomad-driver-podman/pull/137)]
* config: Add `cpu_hard_limit` and `cpu_cfs_period` options [[GH-149](https://github.com/hashicorp/nomad-driver-podman/pull/149)]
* config: Allow mounting rootfs as read-only. [[GH-133](https://github.com/hashicorp/nomad-driver-podman/pull/133)]
* config: Allow setting `ulimit` configuration. [[GH-166](https://github.com/hashicorp/nomad-driver-podman/pull/166)]
* config: Allow setting `image_pull_timeout` and `client_http_timeout ` [[GH-131](https://github.com/hashicorp/nomad-driver-podman/pull/131)]
* runtime: Add support for host and CSI volumes and using podman tasks as CSI plugins [[GH-169](https://github.com/hashicorp/nomad-driver-podman/pull/169)][[GH-152](https://github.com/hashicorp/nomad-driver-podman/pull/152)]

IMPROVEMENTS:

* log: Improve log messages on errors. [[GH-177](https://github.com/hashicorp/nomad-driver-podman/pull/177)]

BUG FIXES:

* log: Use error key context to log errors rather than Go err style. [[GH-126](https://github.com/hashicorp/nomad-driver-podman/pull/126)]
* telemetry: respect telemetry.collection_interval to reduce cpu churn when running many containers [[GH-130](https://github.com/hashicorp/nomad-driver-podman/pull/130)]

## 0.3.0

* config: Image registry authentication [[GH-71](https://github.com/hashicorp/nomad-driver-podman/issues/71)]
* config: Added tty option
* config: Support for sysctl configuration [[GH-82](https://github.com/hashicorp/nomad-driver-podman/issues/82)]
* config: Fixed a bug where we always pulled an image if image name has a transport prefix [[GH-88](https://github.com/hashicorp/nomad-driver-podman/pull/88)]
* config: Added labels option
* config: Add force_pull option
* config: Added logging options

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
