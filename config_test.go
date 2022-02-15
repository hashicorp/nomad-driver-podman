package main

import (
	"testing"

	"github.com/hashicorp/nomad/helper/pluginutils/hclutils"
	"github.com/stretchr/testify/require"
)

func TestConfig_Ports(t *testing.T) {
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
	require.EqualValues(t, expectedPorts, tc.Ports)
}

func TestConfig_Logging(t *testing.T) {
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
	require.EqualValues(t, expectedDriver, tc.Logging.Driver)
	require.EqualValues(t, expectedTag, tc.Logging.Options["tag"])
}

func TestConfig_Labels(t *testing.T) {
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
	require.EqualValues(t, "job", tc.Labels["nomad"])
}

func TestConfig_ForcePull(t *testing.T) {
	parser := hclutils.NewConfigParser(taskConfigSpec)

	validHCL := `
  config {
		image = "docker://redis"
		force_pull = true
  }
`

	var tc *TaskConfig
	parser.ParseHCL(t, validHCL, &tc)
	require.EqualValues(t, true, tc.ForcePull)
}

func TestConfig_CPUHardLimit(t *testing.T) {
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
	require.EqualValues(t, true, tc.CPUHardLimit)
	require.EqualValues(t, 200000, tc.CPUCFSPeriod)
}
