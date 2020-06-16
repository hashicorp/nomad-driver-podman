## 0.0.3 UNRELEASED

FEATURES:

* config: Add ability to configure container network_mode [[GH-33](https://github.com/hashicorp/nomad-driver-podman/issues/33)]

FEATURES:

* #33 configurable network_mode

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
