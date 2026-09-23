// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad/helper/uuid"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/shoenig/test/must"
)

func TestPodmanDriver_CreateDestroyNetwork(t *testing.T) {
	const hostname = "network-test"

	socketPath := "unix:///run/podman/podman.sock"
	if os.Getuid() != 0 {
		socketPath = fmt.Sprintf("unix:///run/user/%d/podman/podman.sock", os.Getuid())
	}
	harness := podmanDriverHarness(t, map[string]interface{}{
		"Socket": []PluginSocketConfig{{Name: hostname, SocketPath: socketPath}},
	})
	driver := getPodmanDriver(t, harness)
	client, err := driver.getPodmanClient(hostname)
	must.NoError(t, err)
	if _, err = client.Ping(context.Background()); err != nil {
		t.Skipf("podman is not available for network integration test: %v", err)
	}

	allocID := uuid.Generate()
	spec, created, err := driver.CreateNetwork(allocID, &drivers.NetworkCreateRequest{Hostname: hostname})
	must.NoError(t, err)
	must.True(t, created)
	destroyed := false
	defer func() {
		if !destroyed {
			_ = driver.DestroyNetwork(allocID, spec)
		}
	}()

	must.Eq(t, drivers.NetIsolationModeGroup, spec.Mode)
	var pid int
	_, scanErr := fmt.Sscanf(spec.Path, "/proc/%d/ns/net", &pid)
	must.NoError(t, scanErr)
	must.True(t, pid > 0)
	must.Eq(t, fmt.Sprintf("/proc/%d/ns/net", pid), spec.Path)
	must.Eq(t, hostname, spec.Labels["podman_pause_hostname"])

	must.NoError(t, driver.DestroyNetwork(allocID, spec))
	destroyed = true
	_, err = client.ContainerInspect(context.Background(), fmt.Sprintf("pause-%s", allocID))
	must.ErrorIs(t, err, api.ContainerNotFound)
}
