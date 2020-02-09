module github.com/pascomnet/nomad-driver-podman

go 1.12

replace (
	// fix Sirupsen/logrus vs sirupsen/logrus problems
	github.com/docker/docker v1.13.1 => github.com/docker/docker v0.7.3-0.20181219122643-d1117e8e1040
	// following https://github.com/hashicorp/nomad-skeleton-driver-plugin/blob/master/go.mod:  don't use ugorji/go, use the hashicorp fork
	github.com/hashicorp/go-msgpack => github.com/hashicorp/go-msgpack v0.0.0-20191101193846-96ddbed8d05b
	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc2.0.20181210164344-f5b99917df9f
	// following https://github.com/hashicorp/nomad-skeleton-driver-plugin/blob/master/go.mod:  fix the version of hashicorp/go-msgpack to 96ddbed8d05b
	github.com/ugorji/go => github.com/hashicorp/go-msgpack v0.0.0-20190927083313-23165f7bc3c2
)

require (
	github.com/LK4D4/joincontext v0.0.0-20171026170139-1724345da6d5 // indirect
	github.com/TylerBrock/colorjson v0.0.0-20180527164720-95ec53f28296 // indirect
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/fsouza/go-dockerclient v1.5.0 // indirect
	github.com/hashicorp/consul v1.4.5
	github.com/hashicorp/consul-template v0.20.0 // indirect
	github.com/hashicorp/go-hclog v0.10.0
	github.com/hashicorp/hcl2 v0.0.0-20191002203319-fb75b3253c80 // indirect
	github.com/hashicorp/nomad v0.10.1
	github.com/hashicorp/nomad/api v0.0.0-20191203164002-b31573ae7206 // indirect
	github.com/hashicorp/raft v1.1.2 // indirect
	github.com/hashicorp/serf v0.8.5 // indirect
	github.com/hashicorp/vault/api v1.0.5-0.20190730042357-746c0b111519 // indirectt
	github.com/mitchellh/hashstructure v1.0.0 // indirect
	github.com/opencontainers/runc v0.1.1
	github.com/opencontainers/runtime-spec v1.0.1 // indirect
	github.com/shirou/gopsutil v2.18.12+incompatible // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/ugorji/go v1.1.7 // indirect
	github.com/varlink/go v0.3.0
	gotest.tools/gotestsum v0.4.0 // indirect
)
