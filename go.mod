module github.com/hashicorp/nomad-driver-podman

go 1.14

replace (
	github.com/docker/docker v1.13.1 => github.com/docker/docker v0.7.3-0.20181219122643-d1117e8e1040
	github.com/godbus/dbus => github.com/godbus/dbus v5.0.1+incompatible

	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc2.0.20181210164344-f5b99917df9f
	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20140225173054-eb6ee6f84d0a
)

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/avast/retry-go v2.4.3+incompatible
	github.com/container-storage-interface/spec v1.2.0 // indirect
	github.com/containernetworking/plugins v0.8.5 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/golangci/golangci-lint v1.30.0 // indirect
	github.com/hashicorp/consul v1.7.3 // indirect
	github.com/hashicorp/consul-template v0.24.1
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/go-hclog/hclogvet v0.1.3 // indirect
	github.com/hashicorp/nomad v0.11.3-0.20200630133459-42c2ee4448c8
	github.com/hashicorp/nomad/api v0.0.0-20200630133459-42c2ee4448c8
	github.com/mrunalp/fileutils v0.0.0-20200504145649-7be891c94fd3 // indirect
	github.com/opencontainers/runtime-spec v1.0.2-0.20200307132014-f49fed0d6290
	github.com/opencontainers/selinux v1.5.1 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/varlink/go v0.3.0
	google.golang.org/grpc v1.29.1 // indirect

)
