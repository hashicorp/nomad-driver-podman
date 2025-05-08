// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import "path/filepath"

type rootlessTaskDir struct {
	mountDir string
	taskName string
}

func (r *rootlessTaskDir) sharedAllocDir() string {
	return filepath.Join(r.mountDir, "alloc")
}

func (r *rootlessTaskDir) localDir() string {
	return filepath.Join(r.mountDir, r.taskName, "local")
}

func (r *rootlessTaskDir) dir() string {
	return filepath.Join(r.mountDir, r.taskName)
}

func (r *rootlessTaskDir) secretsDir() string {
	return filepath.Join(r.mountDir, r.taskName, "secrets")
}
