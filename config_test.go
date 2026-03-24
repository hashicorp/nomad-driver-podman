// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"testing"

	"github.com/hashicorp/nomad-driver-podman/ci"
	"github.com/hashicorp/nomad/helper/pluginutils/hclutils"
	"github.com/shoenig/test/must"
)

func TestConfig_Ports(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	expectedPorts := []string{"redis"}
	validHCL := `
  config {
	image = "docker://redis"
	ports = ["redis"]
  }
`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.SliceContainsAll(t, expectedPorts, tc.Ports)
}

func TestConfig_Logging(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	expectedDriver := "journald"
	expectedTag := "redis"
	validHCL := `
  config {
	  image = "docker://redis"
	  logging = {
			driver = "journald"
			options = [
			  {
				  "tag" = "redis"
			  }
			]
	  }
  }
`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.Eq(t, expectedDriver, tc.Logging.Driver)
	must.Eq(t, expectedTag, tc.Logging.Options["tag"])
}

func TestConfig_Labels(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	validHCL := `
  config {
	  image = "docker://redis"
		labels = {
		  "nomad" = "job"
		 }
  }
`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.Eq(t, "job", tc.Labels["nomad"])
}

func TestConfig_ForcePull(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	validHCL := `
  config {
		image = "docker://redis"
		force_pull = true
  }
`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.Eq(t, true, tc.ForcePull)
}

func TestConfig_CPUHardLimit(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	validHCL := `
  config {
		image = "docker://redis"
		cpu_hard_limit = true
		cpu_cfs_period = 200000
  }
`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.True(t, tc.CPUHardLimit)
	must.Eq(t, 200000, tc.CPUCFSPeriod)
}

func TestConfig_ImagePullTimeout(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	validHCL := `
  config {
		image = "docker://redis"
		image_pull_timeout = "10m"
  }
`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.Eq(t, "10m", tc.ImagePullTimeout)
}

func TestConfig_ExtraHosts(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	validHCL := `
		config {
		image = "docker://redis"
		extra_hosts = ["myhost:127.0.0.2", "example.com:10.0.0.1"]
	}
	`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.Eq(t, []string{"myhost:127.0.0.2", "example.com:10.0.0.1"}, tc.ExtraHosts)
}

func TestConfig_PodmanSocketDefaultIfNotGiven(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	validHCL := `
	config {
		image = "docker://redis"
	}
	`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.Eq(t, "default", tc.Socket)
}

func TestConfig_PodmanOOMScoreAdj(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(taskConfigSpec)
	validHCL := `
	config {
		image = "docker://redis"
		oom_score_adj = -1000
	}
	`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	must.Eq(t, "default", tc.Socket)
}

func TestPluginConfig_Parsing(t *testing.T) {
	ci.Parallel(t)

	parser := hclutils.NewConfigParser(configSpec)
	validHCL := `config {

  socket {
    name        = "default"
    socket_path = "unix:///run/user/1000/podman/podman.sock"
  }

  auth {
    config = "/etc/podman-auth.json"
	helper = "ecr-login"
  }

  gc {
    container = false
  }

  recover_stopped        = true
  extra_labels           = ["nomadproject.io:foo"]
  disable_log_collection = true
  client_http_timeout    = "2m"
  dns_servers            = ["9.9.9.9"]

  logging {
    driver = "journald"
    options {
      tag="{{.ImageName}}"
    }
  }

  networking {
    default_rootless_mode = "slirp4netns"
  }

  volumes {
    selinuxlabel = "z"
  }
}
	`

	var pc *PluginConfig
	parser.ParseHCL(t, validHCL, &pc)
	must.Eq(t, pc, &PluginConfig{
		Auth: PluginAuthConfig{
			FileConfig: "/etc/podman-auth.json",
			Helper:     "ecr-login",
		},
		Volumes: VolumeConfig{
			Enabled:      true,
			SelinuxLabel: "z",
		},
		GC:                   GCConfig{},
		RecoverStopped:       true,
		DisableLogCollection: true,
		Socket: []PluginSocketConfig{{
			Name:       "default",
			SocketPath: "unix:///run/user/1000/podman/podman.sock",
		}},
		SocketPath:        "",
		ClientHttpTimeout: "2m",
		ExtraLabels:       []string{"nomadproject.io:foo"},
		DNSServers:        []string{"9.9.9.9"},
		Logging: LoggingConfig{
			Driver: "journald",
			Options: map[string]string{
				"tag": "{{.ImageName}}",
			},
		},
		Networking: NetworkingConfig{DefaultRootlessMode: "slirp4netns"},
	})
}
