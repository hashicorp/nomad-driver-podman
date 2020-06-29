module github.com/nomad-driver-podman

go 1.14

replace (
	// fix Sirupsen/logrus vs sirupsen/logrus problems
	github.com/docker/docker v1.13.1 => github.com/docker/docker v0.7.3-0.20181219122643-d1117e8e1040
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc2.0.20181210164344-f5b99917df9f

	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20140225173054-eb6ee6f84d0a

)

exclude (
	github.com/NVIDIA/gpu-monitoring-tools v0.0.0-20200622050622-c34507425bdb
	github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml v0.0.0-20200622050622-c34507425bdb
	github.com/hashicorp/nomad/devices/gpu/nvidia v0.11.2
	github.com/hashicorp/nomad/devices/gpu/nvidia/nvml v0.0.0
)

require (
	github.com/LK4D4/joincontext v0.0.0-20171026170139-1724345da6d5 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/TylerBrock/colorjson v0.0.0-20180527164720-95ec53f28296 // indirect
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/burntsushi/toml v0.3.1 // indirect
	github.com/container-storage-interface/spec v1.2.0 // indirect
	github.com/containerd/console v1.0.0 // indirect
	github.com/containerd/go-cni v1.0.0 // indirect
	github.com/containernetworking/plugins v0.8.5 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/fsouza/go-dockerclient v1.5.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-openapi/errors v0.19.6 // indirect
	github.com/go-openapi/runtime v0.19.19 // indirect
	github.com/go-openapi/strfmt v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.9 // indirect
	github.com/go-openapi/validate v0.19.10 // indirect
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/hashicorp/consul v1.7.3 // indirect
	github.com/hashicorp/consul/sdk v0.4.0 // indirect
	github.com/hashicorp/cronexpr v1.1.0 // indirect
	github.com/hashicorp/go-envparse v0.0.0-20200406174449-d9cfd743a15e // indirect
	github.com/hashicorp/go-getter v1.4.1 // indirect
	github.com/hashicorp/go-hclog v0.12.1
	github.com/hashicorp/go-msgpack v1.1.5 // indirect
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80 // indirect
	github.com/hashicorp/nomad v0.11.2
	github.com/hashicorp/vault v1.4.2 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/moby/moby v1.13.1 // indirect
	github.com/mrunalp/fileutils v0.0.0-20200504145649-7be891c94fd3 // indirect
	github.com/opencontainers/runtime-spec v1.0.1 // indirect
	github.com/opencontainers/selinux v1.5.1 // indirect
	github.com/seccomp/libseccomp-golang v0.9.1 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/varlink/go v0.3.0
	github.com/vishvananda/netlink v1.1.0 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gotest.tools/gotestsum v0.4.0 // indirect

)
