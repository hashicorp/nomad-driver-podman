// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/hashicorp/nomad/plugins/drivers"
	"golang.org/x/sys/unix"
)

func (d *Driver) rootlessMount(task *drivers.TaskConfig, driverConfig *TaskConfig) (string, error) {
	uid, err := d.socketUid(driverConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get socket uid: %w", err)
	}

	// This root uid check should have already been done
	if uid == 0 {
		return "", fmt.Errorf("cannot rootless alloc mount when using root socket")
	}

	mountDir := fmt.Sprintf("/var/run/user/%d/%s", uid, task.AllocID)

	// The user dir is only accessible by root and the user
	err = os.Mkdir(mountDir, 755)
	if err != nil {
		return "", fmt.Errorf("failed to create user alloc dir: %w", err)
	}

	err = syscall.Mount(task.AllocDir, mountDir, "", syscall.MS_BIND, "")
	if err != nil {
		return "", fmt.Errorf("failed to mount user alloc dir: %w", err)
	}

	return mountDir, nil
}

func (d *Driver) removeMount(task *drivers.TaskConfig, driverConfig *TaskConfig) error {
	uid, err := d.socketUid(driverConfig)
	if err != nil {
		return err
	}

	mountDir := fmt.Sprintf("/var/run/user/%d/%s", uid, task.AllocID)

	err = syscall.Unmount(mountDir, unix.MNT_DETACH)
	if err != nil {
		return err
	}

	return os.RemoveAll(mountDir)
}

func (d *Driver) socketUid(driverConfig *TaskConfig) (int, error) {
	switch {
	case d.config.SocketPath != "":
		return uidFromSocket(d.config.SocketPath)
	case len(d.config.Socket) > 0:
		for _, v := range d.config.Socket {
			if v.Name == driverConfig.Socket {
				return uidFromSocket(v.SocketPath)
			}
		}
	}
	return os.Getuid(), nil
}

// this only works on linux
func uidFromSocket(path string) (int, error) {
	path = strings.Replace(path, "unix:/", "/var", 1)

	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, fmt.Errorf("failed to stat socket")
	}

	return int(stat.Uid), nil
}
