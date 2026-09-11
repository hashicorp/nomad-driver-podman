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

func TestRootlessTaskDirRewritePath(t *testing.T) {
	testCases := []struct {
		name     string
		mountDir string
		allocDir string
		path     string
		expected string
	}{
		{
			name:     "empty mount directory",
			allocDir: filepath.Join("a", "alloc"),
			path:     filepath.Join("a", "alloc", "task", "file"),
			expected: filepath.Join("a", "alloc", "task", "file"),
		},
		{
			name:     "empty allocation directory",
			mountDir: filepath.Join("run", "mount"),
			path:     filepath.Join("a", "alloc", "task", "file"),
			expected: filepath.Join("a", "alloc", "task", "file"),
		},
		{
			name:     "allocation directory itself",
			mountDir: filepath.Join("run", "mount"),
			allocDir: filepath.Join("a", "alloc"),
			path:     filepath.Join("a", "alloc"),
			expected: filepath.Join("run", "mount"),
		},
		{
			name:     "path under allocation directory",
			mountDir: filepath.Join("run", "mount"),
			allocDir: filepath.Join("a", "alloc"),
			path:     filepath.Join("a", "alloc", "task", "file"),
			expected: filepath.Join("run", "mount", "task", "file"),
		},
		{
			name:     "path outside allocation directory",
			mountDir: filepath.Join("run", "mount"),
			allocDir: filepath.Join("a", "alloc"),
			path:     filepath.Join("var", "log", "task.log"),
			expected: filepath.Join("var", "log", "task.log"),
		},
		{
			name:     "sibling prefix is currently rewritten",
			mountDir: "/run/mount",
			allocDir: "/a/alloc",
			path:     "/a/alloc-other/x",
			expected: "/run/mount-other/x",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rootlessDir := &rootlessTaskDir{
				mountDir: tc.mountDir,
				allocDir: tc.allocDir,
			}

			// NOTE: rewritePath uses strings.HasPrefix, so sibling names that
			// share the allocDir prefix are also rewritten.
			must.Eq(t, tc.expected, rootlessDir.rewritePath(tc.path))
		})
	}
}
