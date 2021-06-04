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
