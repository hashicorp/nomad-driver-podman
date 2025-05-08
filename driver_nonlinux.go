// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:build !linux
// +build !linux

package main

import "github.com/hashicorp/nomad/plugins/drivers"

// Unimplemented
func (d *Driver) rootlessMount(task *drivers.TaskConfig, driverConfig *TaskConfig) (string, error) {
	return "", nil
}

// Unimplemented
func (d *Driver) removeMount(task *drivers.TaskConfig, driverConfig *TaskConfig) error {
	return nil
}
