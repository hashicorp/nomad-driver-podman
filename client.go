/*
Copyright 2019 Thomas Weber

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/avast/retry-go"
	"github.com/hashicorp/go-hclog"
	"github.com/pascomnet/nomad-driver-podman/iopodman"
	"github.com/varlink/go/varlink"

	"os"
	"bufio"
	"os/user"
	"context"
	"encoding/json"
	"strings"
	"fmt"
	"net"
	"time"
)

const (
	// number of varlink op retries
	varlinkRetries = 4
)

// PodmanClient encapsulates varlink operations
type PodmanClient struct {

	// ctx is the context for the driver. It is passed to other subsystems to
	// coordinate shutdown
	ctx context.Context

	// logger will log to the Nomad agent
	logger hclog.Logger

	varlinkSocketPath string
}

// withVarlink calls a podman varlink function and retries N times in case of network failures
func (c *PodmanClient) withVarlink(cb func(*varlink.Connection) error) error {
	err := retry.Do(
		// invoke the callback in a fresh varlink connection
		func() error {
			connection, err := c.getConnection()
			if err != nil {
				c.logger.Debug("Failed to open varlink connection", "err", err)
				return err
			}
			defer connection.Close()
			return cb(connection)
		},
		// ... and retry up to N times
		retry.Attempts(varlinkRetries),
		// ... but only if it failed with net.OpError
		retry.RetryIf(func(err error) bool {
			if _, ok := err.(*net.OpError); ok {
				c.logger.Debug("Failed to invoke varlink method, will retry", "err", err)
				return true
			}
			return false
		}),
		// wait 1 sec between retries
		retry.Delay(time.Second),
		// and return last error directly (not wrapped)
		retry.LastErrorOnly(true),
	)
	return err
}

// GetContainerStats takes the name or ID of a container and returns a single ContainerStats structure which
// contains attributes like memory and cpu usage.  If the container cannot be found, a
// [ContainerNotFound](#ContainerNotFound) error will be returned. If the container is not running, a [NoContainerRunning](#NoContainerRunning)
// error will be returned
func (c *PodmanClient) GetContainerStats(containerID string) (*iopodman.ContainerStats, error) {
	var containerStats *iopodman.ContainerStats
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		result, err := iopodman.GetContainerStats().Call(c.ctx, varlinkConnection, containerID)
		containerStats = &result
		return err
	})
	return containerStats, err
}

// StopContainer stops a container given a timeout.  It takes the name or ID of a container as well as a
// timeout value.  The timeout value the time before a forcible stop to the container is applied.
// If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned instead.
func (c *PodmanClient) StopContainer(containerID string, timeoutSeconds int64) error {
	c.logger.Debug("Stopping container", "container", containerID, "timeout", timeoutSeconds)
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		_, err := iopodman.StopContainer().Call(c.ctx, varlinkConnection, containerID, timeoutSeconds)
		return err
	})
	return err
}

// ForceRemoveContainer requires the name or ID of a container and will stop and remove a container and it's volumes.
// iI container cannot be found by name or ID, a [ContainerNotFound](#ContainerNotFound) error will be returned.
func (c *PodmanClient) ForceRemoveContainer(containerID string) error {
	c.logger.Debug("Removing container", "container", containerID)
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		_, err := iopodman.RemoveContainer().Call(c.ctx, varlinkConnection, containerID, true, true)
		return err
	})
	return err
}

// GetInfo returns a [PodmanInfo](#PodmanInfo) struct that describes podman and its host such as storage stats,
// build information of Podman, and system-wide registries
func (c *PodmanClient) GetInfo() (*iopodman.PodmanInfo, error) {
	var ret *iopodman.PodmanInfo
	c.logger.Debug("Get podman info")
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		result, err := iopodman.GetInfo().Call(c.ctx, varlinkConnection)
		ret = &result
		return err
	})
	return ret, err
}

// PsID returns a PsContainer struct that describes the process state of exactly
// one container.
func (c *PodmanClient) PsID(containerID string) (*iopodman.PsContainer, error) {
	filters := []string{"id=" + containerID}
	psInfo, err := c.Ps(filters)
	if err == nil {
		if len(psInfo) == 1 {
			return &psInfo[0], nil
		} else {
			return nil, fmt.Errorf("No such container: %s", containerID)
		}
	}
	return nil, err
}

// PsByName returns a PsContainer struct that describes the process state of exactly
// one container.
func (c *PodmanClient) PsByName(containerName string) (*iopodman.PsContainer, error) {
	filters := []string{"name=" + containerName}
	psInfo, err := c.Ps(filters)
	if err == nil {
		if len(psInfo) == 1 {
			return &psInfo[0], nil
		} else {
			return nil, fmt.Errorf("No such container: %s", containerName)
		}
	}
	return nil, err
}

// Ps finds process info for one or more containers by applying a filter
func (c *PodmanClient) Ps(filters []string) ([]iopodman.PsContainer, error) {
	var ret []iopodman.PsContainer
	c.logger.Debug("Get container list", "filters", filters)
	psOpts := iopodman.PsOpts{
		Filters: &filters,
		All:     true,
	}
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		result, err := iopodman.Ps().Call(c.ctx, varlinkConnection, psOpts)
		if err == nil {
			ret = result
		}
		return err
	})
	return ret, err
}

// CreateContainer creates a new container from an image.  It uses a [Create](#Create) type for input.
func (c *PodmanClient) CreateContainer(createOpts iopodman.Create) (string, error) {
	ret := ""
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		result, err := iopodman.CreateContainer().Call(c.ctx, varlinkConnection, createOpts)
		if err == nil {
			ret = result
			c.logger.Debug("Created container", "container", ret)
		}
		return err
	})
	return ret, err
}

// StartContainer starts a created or stopped container. It takes the name or ID of container.  It returns
// the container ID once started.  If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned.
func (c *PodmanClient) StartContainer(containerID string) error {
	c.logger.Debug("Starting container", "container", containerID)
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		_, err := iopodman.StartContainer().Call(c.ctx, varlinkConnection, containerID)
		return err
	})
	return err
}

// InspectContainer data takes a name or ID of a container returns the inspection
// data as iopodman.InspectContainerData.
func (c *PodmanClient) InspectContainer(containerID string) (iopodman.InspectContainerData, error) {
	var ret iopodman.InspectContainerData
	c.logger.Debug("Inspect container", "container", containerID)
	err := c.withVarlink(func(varlinkConnection *varlink.Connection) error {
		inspectJSON, err := iopodman.InspectContainer().Call(c.ctx, varlinkConnection, containerID)
		if err == nil {
			err = json.Unmarshal([]byte(inspectJSON), &ret)
			if err != nil {
				c.logger.Error("failed to unmarshal inspect container", "err", err)
				return err
			}

		}
		return err
	})
	return ret, err
}


func guessSocketPath(user *user.User, procFilesystems []string) string {
	rootVarlinkPath := "unix://run/podman/io.podman"
	if user.Uid == "0" {
		return rootVarlinkPath
	}

	cgroupv2 := isCGroupV2(procFilesystems)

	if cgroupv2 {
		return fmt.Sprintf("unix://run/user/%s/podman/io.podman", user.Uid)
	}

	return rootVarlinkPath
}

func isCGroupV2(procFilesystems []string) bool {
	cgroupv2 := false
	for _,l := range procFilesystems {
		if strings.HasSuffix(l, "cgroup2") {
			cgroupv2 = true
		}
	}

	return cgroupv2
}

func getProcFilesystems() ([]string, error) {
	file, err := os.Open("/proc/filesystems")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// getConnection opens a new varlink connection
func (c *PodmanClient) getConnection() (*varlink.Connection, error) {
	varlinkConnection, err := varlink.NewConnection(c.ctx, c.varlinkSocketPath)
	return varlinkConnection, err
}
