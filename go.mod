module github.com/pascomnet/nomad-driver-podman

go 1.14

replace (
	// fix Sirupsen/logrus vs sirupsen/logrus problems
	github.com/docker/docker v1.13.1 => github.com/docker/docker v0.7.3-0.20181219122643-d1117e8e1040
	// following https://github.com/hashicorp/nomad-skeleton-driver-plugin/blob/master/go.mod:  don't use ugorji/go, use the hashicorp fork
	github.com/hashicorp/go-msgpack => github.com/hashicorp/go-msgpack v0.0.0-20191101193846-96ddbed8d05b
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc2.0.20181210164344-f5b99917df9f
	// following https://github.com/hashicorp/nomad-skeleton-driver-plugin/blob/master/go.mod:  fix the version of hashicorp/go-msgpack to 96ddbed8d05b
	github.com/ugorji/go => github.com/hashicorp/go-msgpack v0.0.0-20190927123313-23165f7bc3c2
)

require (
	github.com/LK4D4/joincontext v0.0.0-20171026170139-1724345da6d5 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/TylerBrock/colorjson v0.0.0-20180527164720-95ec53f28296 // indirect
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/containerd/console v1.0.0 // indirect
	github.com/containernetworking/plugins v0.8.5 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/cyphar/filepath-securejoin v0.2.2 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/fsouza/go-dockerclient v1.5.0 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/hashicorp/consul v1.4.5
	github.com/hashicorp/consul-template v0.20.0 // indirect
	github.com/hashicorp/go-hclog v0.10.0
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80 // indirect
	github.com/hashicorp/nomad v0.10.1
	github.com/hashicorp/nomad/api v0.0.0-20191203164002-b31573ae7206 // indirect
	github.com/hashicorp/raft v1.1.2 // indirect
	github.com/hashicorp/serf v0.8.5 // indirect
	github.com/hashicorp/vault/api v1.0.5-0.20190730042357-746c0b111519 // indirect; indirectt
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/moby/moby v1.13.1 // indirect
	github.com/mrunalp/fileutils v0.0.0-20200504145649-7be891c94fd3 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.0.1 // indirect
	github.com/opencontainers/selinux v1.5.1 // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible // indirect
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/gocapability v0.0.0-20180916011248-d98352740cb2 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	github.com/varlink/go v0.3.0
	github.com/vishvananda/netlink v1.1.0 // indirect
	gotest.tools/gotestsum v0.4.0 // indirect
)
