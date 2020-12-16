module github.com/hashicorp/nomad-driver-podman

go 1.14

replace (
	github.com/docker/docker v1.13.1 => github.com/docker/docker v0.7.3-0.20181219122643-d1117e8e1040
	github.com/godbus/dbus => github.com/godbus/dbus v5.0.1+incompatible

	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc2.0.20181210164344-f5b99917df9f
	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20140225173054-eb6ee6f84d0a
)

require (
	contrib.go.opencensus.io/exporter/ocagent v0.4.12 // indirect
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/checkpoint-restore/go-criu v0.0.0-20190109184317-bdb7599cd87b // indirect
	github.com/container-storage-interface/spec v1.2.0 // indirect
	github.com/containernetworking/plugins v0.8.5 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.0 // indirect
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/nomad v1.0.0
	github.com/hashicorp/nomad/api v0.0.0-20201208134522-a480eed0815c
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20200728170252-4d89ac9fbff6
	github.com/stretchr/testify v1.6.1
	golang.org/x/build v0.0.0-20190111050920-041ab4dc3f9d // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e // indirect
	golang.org/x/tools v0.0.0-20200724022722-7017fd6b1305 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	istio.io/gogo-genproto v0.0.0-20190124151557-6d926a6e6feb // indirect

)
