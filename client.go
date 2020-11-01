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
	"github.com/varlink/go/varlink"

	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/user"
	"strings"
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
	for _, l := range procFilesystems {
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
