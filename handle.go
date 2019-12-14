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
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/client/stats"
	"github.com/hashicorp/nomad/plugins/drivers"
	"github.com/pascomnet/nomad-driver-podman/iopodman"
)

const (
	// containerMonitorIntv is the interval at which the driver checks if the
	// container is still alive
	containerMonitorIntv = 2 * time.Second
)

var (
	MeasuredCpuStats = []string{"System Mode", "User Mode", "Percent"}
	MeasuredMemStats = []string{"RSS"}
)

type TaskHandle struct {
	initPid     int
	containerID string
	logger      hclog.Logger
	driver      *Driver

	totalCpuStats  *stats.CpuStats
	userCpuStats   *stats.CpuStats
	systemCpuStats *stats.CpuStats

	// stateLock syncs access to all fields below
	stateLock sync.RWMutex

	taskConfig  *drivers.TaskConfig
	procState   drivers.TaskState
	startedAt   time.Time
	completedAt time.Time
	exitResult  *drivers.ExitResult

	removeContainerOnExit bool
}

func (h *TaskHandle) TaskStatus() *drivers.TaskStatus {
	h.stateLock.RLock()
	defer h.stateLock.RUnlock()

	return &drivers.TaskStatus{
		ID:          h.taskConfig.ID,
		Name:        h.taskConfig.Name,
		State:       h.procState,
		StartedAt:   h.startedAt,
		CompletedAt: h.completedAt,
		ExitResult:  h.exitResult,
		DriverAttributes: map[string]string{
			"pid": strconv.Itoa(h.initPid),
		},
	}
}

func (h *TaskHandle) IsRunning() bool {
	h.stateLock.RLock()
	defer h.stateLock.RUnlock()
	return h.procState == drivers.TaskStateRunning
}

func (h *TaskHandle) run() {
	h.stateLock.Lock()
	if h.exitResult == nil {
		h.exitResult = &drivers.ExitResult{}
	}
	h.stateLock.Unlock()
	h.logger.Debug("Monitoring process", "container", h.containerID, "pid", h.initPid)

	if ok, err := waitTillStopped(h.initPid); !ok {
		h.logger.Error("failed to find container process", "error", err)
		return
	}

	h.logger.Debug("Process stopped", "container", h.containerID, "pid", h.initPid)

	h.stateLock.Lock()
	defer h.stateLock.Unlock()

	h.procState = drivers.TaskStateExited
	h.exitResult.ExitCode = 0
	h.exitResult.Signal = 0
	h.completedAt = time.Now()

	// TODO: detect if the task OOMed
}

func (h *TaskHandle) stats(ctx context.Context, interval time.Duration) (<-chan *drivers.TaskResourceUsage, error) {
	ch := make(chan *drivers.TaskResourceUsage)
	go h.handleStats(ctx, ch, interval)
	return ch, nil
}

func (h *TaskHandle) handleStats(ctx context.Context, ch chan *drivers.TaskResourceUsage, interval time.Duration) {
	defer close(ch)
	timer := time.NewTimer(0)

	// stats are polled relatively often so it should be better
	// to open a long living varlink connection instead
	// of re-opening it on each poll cycle in the for-loop.
	varlinkConnection, err := h.driver.getConnection()
	if err != nil {
		h.logger.Error("failed to get varlink connection for stats", "err", err)
		return
	}
	defer varlinkConnection.Close()

	failures := 0

	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			timer.Reset(interval)
		}

		containerStats, err := iopodman.GetContainerStats().Call(h.driver.ctx, varlinkConnection, h.containerID)
		if err != nil {
			h.logger.Debug("Could not get container stats, trying to reconnect", "err", err)
			// maybe varlink connection was lost, we should check and try to reconnect
			varlinkConnection.Close()
			varlinkConnection, err = h.driver.getConnection()
			if err != nil {
				h.logger.Error("failed to reconnect varlink for stats", "err", err)
				failures ++
			}

			if failures<5 {
				// try again
				continue
			} else {
				// stop, makes no sense
				h.logger.Error("too many reconnect errors")
				return
			}

		}

		// reset failure count if we successfully got stats
		failures = 0
		t := time.Now()

		//FIXME implement cpu stats correctly
		//available := shelpers.TotalTicksAvailable()
		//cpus := shelpers.CPUNumCores()

		totalPercent := h.totalCpuStats.Percent(containerStats.Cpu * 10e16)
		cs := &drivers.CpuStats{
			SystemMode: h.systemCpuStats.Percent(float64(containerStats.System_nano)),
			UserMode:   h.userCpuStats.Percent(float64(containerStats.Cpu_nano)),
			Percent:    totalPercent,
			TotalTicks: h.systemCpuStats.TicksConsumed(totalPercent),
			Measured:   MeasuredCpuStats,
		}

		//h.driver.logger.Info("stats", "cpu", containerStats.Cpu, "system", containerStats.System_nano, "user", containerStats.Cpu_nano, "percent", totalPercent, "ticks", cs.TotalTicks, "cpus", cpus, "available", available)

		ms := &drivers.MemoryStats{
			RSS:      uint64(containerStats.Mem_usage),
			Measured: MeasuredMemStats,
		}

		taskResUsage := drivers.TaskResourceUsage{
			ResourceUsage: &drivers.ResourceUsage{
				CpuStats:    cs,
				MemoryStats: ms,
			},
			Timestamp: t.UTC().UnixNano(),
		}
		select {
		case <-ctx.Done():
			return
		case ch <- &taskResUsage:
		}
	}
}

// shutdown shuts down the container, with `timeout` grace period
// before killing the container with SIGKILL.
func (h *TaskHandle) shutdown(timeout time.Duration) error {
	varlinkConnection, err := h.driver.getConnection()
	if err != nil {
		return fmt.Errorf("executor Shutdown failed, could not get podman connection: %v", err)
	}
	defer varlinkConnection.Close()

	h.driver.logger.Debug("Stopping podman container", "container", h.containerID)
	// TODO: we should respect the "signal" parameter here
	if _, err := iopodman.StopContainer().Call(h.driver.ctx, varlinkConnection, h.containerID, int64(timeout)); err != nil {
		h.driver.logger.Warn("Could not stop container gracefully, killing it now", "containerID", h.containerID, "err", err)
		if _, err := iopodman.KillContainer().Call(h.driver.ctx, varlinkConnection, h.containerID, 9); err != nil {
			h.driver.logger.Error("Could not kill container", "containerID", h.containerID, "err", err)
			return fmt.Errorf("Could not kill container: %v", err)
		}
	}
	return nil
}

// waitTillStopped blocks and returns true when container stops;
// returns false with an error message if the container processes cannot be identified.
//
func waitTillStopped(pid int) (bool, error) {
	ps, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}

	for {
		if err := ps.Signal(syscall.Signal(0)); err != nil {
			return true, nil
		}

		time.Sleep(containerMonitorIntv)
	}
}
