// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build linux
// +build linux

package main

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/hashicorp/nomad-driver-podman/ci"
	"github.com/shoenig/test/must"
)

func TestUidFromSocket(t *testing.T) {
	ci.Parallel(t)

	path := filepath.Join(t.TempDir(), "podman.sock")
	must.NoError(t, os.WriteFile(path, nil, 0600))

	info, err := os.Stat(path)
	must.NoError(t, err)
	expectedUID := int(info.Sys().(*syscall.Stat_t).Uid)

	uid, err := uidFromSocket("unix://" + path)
	must.NoError(t, err)
	must.Eq(t, expectedUID, uid)
}

func TestSocketUidEmptySocketUsesDefault(t *testing.T) {
	ci.Parallel(t)

	path := filepath.Join(t.TempDir(), "podman.sock")
	must.NoError(t, os.WriteFile(path, nil, 0600))

	info, err := os.Stat(path)
	must.NoError(t, err)
	expectedUID := int(info.Sys().(*syscall.Stat_t).Uid)

	driver := &Driver{
		config: &PluginConfig{
			Socket: []PluginSocketConfig{{
				Name:       "default",
				SocketPath: "unix://" + path,
			}},
		},
	}
	uid, err := driver.socketUid(&TaskConfig{})
	must.NoError(t, err)
	must.Eq(t, expectedUID, uid)
}
