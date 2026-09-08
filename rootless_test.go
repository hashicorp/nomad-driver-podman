// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"path/filepath"
	"testing"

	"github.com/shoenig/test/must"
)

func TestRootlessTaskDirPaths(t *testing.T) {
	mountDir := filepath.Join("run", "user", "12080", "nomad", "alloc-id")
	taskName := "example-task"
	rootlessDir := &rootlessTaskDir{
		mountDir: mountDir,
		taskName: taskName,
	}

	testCases := []struct {
		name     string
		path     func() string
		expected string
	}{
		{
			name:     "shared alloc directory",
			path:     rootlessDir.sharedAllocDir,
			expected: filepath.Join(mountDir, "alloc"),
		},
		{
			name:     "task local directory",
			path:     rootlessDir.localDir,
			expected: filepath.Join(mountDir, taskName, "local"),
		},
		{
			name:     "task directory",
			path:     rootlessDir.dir,
			expected: filepath.Join(mountDir, taskName),
		},
		{
			name:     "task secrets directory",
			path:     rootlessDir.secretsDir,
			expected: filepath.Join(mountDir, taskName, "secrets"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			must.Eq(t, tc.expected, tc.path())
		})
	}
}
