// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	// Create mount under /run/user/<uid>/nomad/ to avoid polluting the user's runtime dir
	baseDir := fmt.Sprintf("/var/run/user/%d/nomad", uid)
	if err = os.MkdirAll(baseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create nomad base dir: %w", err)
	}
	mountDir := filepath.Join(baseDir, task.AllocID)

	// Check if mount already exists and is valid (for task/alloc restarts).
	// Reuse existing mount to preserve captured submounts like secrets tmpfs.
	if _, statErr := os.Stat(mountDir); statErr == nil {
		// Mount exists - check if it's still a valid bind mount
		// by verifying expected subdirectory structure exists
		if _, allocErr := os.Stat(filepath.Join(mountDir, "alloc")); allocErr == nil {
			d.logger.Debug("Reusing existing rootless bind mount", "allocID", task.AllocID, "mountDir", mountDir)
			return mountDir, nil
		}
		// Invalid/stale mount, clean up and recreate
		d.logger.Debug("Cleaning up invalid rootless bind mount", "allocID", task.AllocID)
		_ = syscall.Unmount(mountDir, unix.MNT_DETACH)
		_ = os.RemoveAll(mountDir)
	}

	if err = os.Mkdir(mountDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create alloc mount dir: %w", err)
	}

	// Use MS_REC for recursive bind mount - this captures tmpfs submounts
	// for secrets that Nomad creates before StartTask() is called
	err = syscall.Mount(task.AllocDir, mountDir, "", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		os.Remove(mountDir)
		return "", fmt.Errorf("failed to mount user alloc dir: %w", err)
	}

	return mountDir, nil
}

func (d *Driver) removeMount(task *drivers.TaskConfig, driverConfig *TaskConfig) error {
	uid, err := d.socketUid(driverConfig)
	if err != nil {
		return err
	}

	mountDir := fmt.Sprintf("/var/run/user/%d/nomad/%s", uid, task.AllocID)

	// Check if mount dir exists - if not, already cleaned up
	if _, err := os.Stat(mountDir); os.IsNotExist(err) {
		return nil
	}

	// Use MNT_DETACH to handle busy mounts; ignore unmount errors and
	// proceed to removal since the dir might not actually be mounted
	_ = syscall.Unmount(mountDir, unix.MNT_DETACH)

	return os.RemoveAll(mountDir)
}

// removeMountByPath unmounts and removes a bind mount by its exact path.
// Used when we have the path stored directly rather than computing it.
func (d *Driver) removeMountByPath(mountDir string) error {
	if mountDir == "" {
		return nil
	}

	// Check if mount dir exists - if not, already cleaned up
	if _, err := os.Stat(mountDir); os.IsNotExist(err) {
		return nil
	}

	// Use MNT_DETACH to handle busy mounts; ignore unmount errors and
	// proceed to removal since the dir might not actually be mounted
	_ = syscall.Unmount(mountDir, unix.MNT_DETACH)

	return os.RemoveAll(mountDir)
}

func (d *Driver) socketUid(driverConfig *TaskConfig) (int, error) {
	// Normalise empty socket name to "default" to match makePodmanClients behaviour
	socketName := driverConfig.Socket
	if socketName == "" {
		socketName = "default"
	}

	switch {
	case d.config.SocketPath != "":
		return uidFromSocket(d.config.SocketPath)
	case len(d.config.Socket) > 0:
		for _, v := range d.config.Socket {
			if v.Name == socketName {
				return uidFromSocket(v.SocketPath)
			}
		}
	}
	return os.Getuid(), nil
}

// cleanupOrphanedMounts removes bind mounts whose source allocation directories
// no longer exist. This handles cleanup after Nomad GCs allocations.
func (d *Driver) cleanupOrphanedMounts() {
	// Get all UIDs we might have mounts for by checking socket ownership.
	// We can't rely on client.IsRootless() here because this may be called
	// during SetConfig before fingerprinting has set the rootless flag.
	uids := make(map[int]bool)

	// Check all configured sockets for non-root ownership
	for _, sock := range d.config.Socket {
		if uid, err := uidFromSocket(sock.SocketPath); err == nil && uid != 0 {
			uids[uid] = true
		}
	}
	if d.config.SocketPath != "" {
		if uid, err := uidFromSocket(d.config.SocketPath); err == nil && uid != 0 {
			uids[uid] = true
		}
	}
	// Also check default socket location if no explicit config
	if len(d.config.Socket) == 0 && d.config.SocketPath == "" {
		uid := os.Getuid()
		if uid != 0 {
			uids[uid] = true
		}
	}

	for uid := range uids {
		nomadDir := fmt.Sprintf("/var/run/user/%d/nomad", uid)
		entries, err := os.ReadDir(nomadDir)
		if err != nil {
			continue // Directory doesn't exist or can't read
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			allocID := entry.Name()
			mountDir := filepath.Join(nomadDir, allocID)

			// Check if this is an orphaned mount by seeing if the alloc subdir exists
			// If the bind mount is valid, alloc/ will be visible; if orphaned, it won't
			allocSubdir := filepath.Join(mountDir, "alloc")
			if _, err := os.Stat(allocSubdir); os.IsNotExist(err) {
				d.logger.Debug("Cleaning up orphaned rootless bind mount", "allocID", allocID, "mountDir", mountDir)
				_ = syscall.Unmount(mountDir, unix.MNT_DETACH)
				_ = os.RemoveAll(mountDir)
			}
		}
	}
}

// uidFromSocket extracts the owner UID from the podman socket path
func uidFromSocket(path string) (int, error) {
	path = strings.Replace(path, "unix://", "", 1)

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
