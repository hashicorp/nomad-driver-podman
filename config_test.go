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
