// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"path/filepath"
	"strings"
)

type rootlessTaskDir struct {
	mountDir string
	allocDir string // original alloc dir path for rewriting
	taskName string
}

func (r *rootlessTaskDir) sharedAllocDir() string {
	return filepath.Join(r.mountDir, "alloc")
}

// rewritePath rewrites a path under the original allocDir to use the bind-mounted path.
// Returns the original path unchanged if it's not under allocDir or if mountDir is empty.
func (r *rootlessTaskDir) rewritePath(path string) string {
	if r.mountDir == "" || r.allocDir == "" {
		return path
	}
	if strings.HasPrefix(path, r.allocDir) {
		return r.mountDir + path[len(r.allocDir):]
	}
	return path
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
