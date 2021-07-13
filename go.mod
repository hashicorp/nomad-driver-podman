module github.com/hashicorp/nomad-driver-podman

go 1.14

replace (
	github.com/docker/docker v1.13.1 => github.com/docker/docker v0.7.3-0.20181219122643-d1117e8e1040
	github.com/godbus/dbus => github.com/godbus/dbus v5.0.1+incompatible

	github.com/opencontainers/runc v0.1.1 => github.com/opencontainers/runc v1.0.0-rc2.0.20181210164344-f5b99917df9f
	launchpad.net/gocheck => github.com/go-check/check v0.0.0-20140225173054-eb6ee6f84d0a
)

require (
	github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e
	github.com/containers/image/v5 v5.11.1
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/nomad v1.0.4-0.20210415141937-ee9bb3cc4f40
	github.com/hashicorp/nomad/api v0.0.0-20210415141937-ee9bb3cc4f40
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mitchellh/go-linereader v0.0.0-20190213213312-1b945b3263eb
	github.com/mitchellh/go-ps v1.0.0 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/opencontainers/runtime-spec v1.0.3-0.20200929063507-e6143ca7d51d
	github.com/stretchr/testify v1.7.0

)
