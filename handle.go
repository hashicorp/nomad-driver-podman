// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad/client/lib/cpustats"
	"github.com/hashicorp/nomad/plugins/drivers"
)

var (
	measuredCPUStats = []string{"System Mode", "User Mode", "Percent"}
	measuredMemStats = []string{"Usage", "Max Usage"}
)

const (
	containerMonitorRetryInterval = time.Second
	containerStatsPollInterval    = time.Second
)

type containerLifecycleClient interface {
	ContainerInspect(context.Context, string) (api.InspectContainerData, error)
	ContainerWait(context.Context, string, []string) error
}

type containerStatsClient interface {
	ContainerStats(context.Context, string) (api.Stats, error)
}

// TaskHandle is the podman specific handle for exactly one container
type TaskHandle struct {
	containerID  string
	logger       hclog.Logger
	driver       *Driver
	podmanClient *api.API

	totalCPUStats  *cpustats.Tracker
	userCPUStats   *cpustats.Tracker
	systemCPUStats *cpustats.Tracker

	// stateLock syncs access to all fields below
	stateLock sync.RWMutex

	taskConfig  *drivers.TaskConfig
	procState   drivers.TaskState
	startedAt   time.Time
	logPointer  time.Time
	completedAt time.Time
	exitResult  *drivers.ExitResult

	containerStats        api.Stats
	removeContainerOnExit bool
	logStreamer           bool
}

func (h *TaskHandle) taskStatus() *drivers.TaskStatus {
	h.stateLock.RLock()
	defer h.stateLock.RUnlock()

	return &drivers.TaskStatus{
		ID:               h.taskConfig.ID,
		Name:             h.taskConfig.Name,
		State:            h.procState,
		StartedAt:        h.startedAt,
		CompletedAt:      h.completedAt,
		ExitResult:       h.exitResult,
		DriverAttributes: map[string]string{
			// we do not need custom attributes yet
		},
	}
}

func (h *TaskHandle) isRunning() bool {
	h.stateLock.RLock()
	defer h.stateLock.RUnlock()
	return h.procState == drivers.TaskStateRunning
}

func (h *TaskHandle) runExitWatcher(ctx context.Context, exitChannel chan *drivers.ExitResult) {
	timer := time.NewTimer(0)
	h.logger.Debug("Starting exitWatcher", "container", h.containerID)

	defer func() {
		h.logger.Debug("Stopping exitWatcher", "container", h.containerID)
		// be sure to get the whole result
		h.stateLock.Lock()
		result := h.exitResult
		h.stateLock.Unlock()
		exitChannel <- result
		close(exitChannel)
	}()

	for {
		if !h.isRunning() {
			return
		}
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			timer.Reset(time.Second)
		}
	}
}

func (h *TaskHandle) runStatsEmitter(ctx context.Context, statsChannel chan *drivers.TaskResourceUsage, interval time.Duration) {
	timer := time.NewTimer(0)
	defer timer.Stop()
	defer close(statsChannel)
	h.logger.Debug("Starting statsEmitter", "container", h.containerID)
	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("Stopping statsEmitter", "container", h.containerID)
			return
		case <-timer.C:
			timer.Reset(interval)
		}
		if !h.isRunning() {
			h.logger.Debug("Stopping statsEmitter for exited container", "container", h.containerID)
			return
		}

		h.stateLock.Lock()
		t := time.Now()

		// FIXME implement cpu stats correctly
		totalPercent := h.totalCPUStats.Percent(float64(h.containerStats.CPUStats.CPUUsage.TotalUsage))
		cs := &drivers.CpuStats{
			SystemMode: h.systemCPUStats.Percent(float64(h.containerStats.CPUStats.CPUUsage.UsageInKernelmode)),
			UserMode:   h.userCPUStats.Percent(float64(h.containerStats.CPUStats.CPUUsage.UsageInUsermode)),
			Percent:    totalPercent,
			TotalTicks: h.systemCPUStats.TicksConsumed(totalPercent),
			Measured:   measuredCPUStats,
		}

		ms := &drivers.MemoryStats{
			MaxUsage: h.containerStats.MemoryStats.MaxUsage,
			Usage:    h.containerStats.MemoryStats.Usage,
			RSS:      h.containerStats.MemoryStats.Usage,
			Measured: measuredMemStats,
		}
		h.stateLock.Unlock()

		// update uasge
		usage := drivers.TaskResourceUsage{
			ResourceUsage: &drivers.ResourceUsage{
				CpuStats:    cs,
				MemoryStats: ms,
			},
			Timestamp: t.UTC().UnixNano(),
		}
		// send stats to nomad
		statsChannel <- &usage
	}
}
func (h *TaskHandle) runLogStreamer(ctx context.Context) {
	stdout, err := os.OpenFile(h.taskConfig.StdoutPath, os.O_WRONLY|syscall.O_NONBLOCK, 0600)
	if err != nil {
		h.logger.Warn("Unable to open stdout fifo", "error", err)
		return
	}
	defer stdout.Close()
	stderr, err := os.OpenFile(h.taskConfig.StderrPath, os.O_WRONLY|syscall.O_NONBLOCK, 0600)
	if err != nil {
		h.logger.Warn("Unable to open stderr fifo", "error", err)
		return
	}
	defer stderr.Close()

	init := true
	since := h.logPointer
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !init {
				// throttle logger reconciliation
				time.Sleep(2 * time.Second)
			}
			err = h.podmanClient.ContainerLogs(ctx, h.containerID, since, stdout, stderr)
			if err != nil {
				h.logger.Warn("Log stream was interrupted", "error", err)
				init = false
				since = time.Now()
				// increment logPointer
				h.stateLock.Lock()
				h.logPointer = since
				h.stateLock.Unlock()
			} else {
				h.logger.Trace("runLogStreamer loop exit")
				return
			}
		}
	}

}

func (h *TaskHandle) runStatsMonitor() {
	h.runStatsMonitorWithClient(h.driver.ctx, h.podmanClient, containerStatsPollInterval)
}

func (h *TaskHandle) runStatsMonitorWithClient(ctx context.Context, client containerStatsClient, interval time.Duration) {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	h.logger.Debug("Monitoring container stats", "container", h.containerID)

	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			timer.Reset(interval)
		}

		if !h.isRunning() {
			return
		}

		containerStats, statsErr := client.ContainerStats(ctx, h.containerID)
		if statsErr != nil {
			h.logger.Debug("Could not get container stats", "container", h.containerID, "error", statsErr)
			if ctx.Err() != nil {
				return
			}
			continue
		}

		h.stateLock.Lock()
		// keep last known containerStats in handle to
		// have it available in the stats emitter
		h.containerStats = containerStats
		h.stateLock.Unlock()
	}
}

func (h *TaskHandle) runContainerMonitor() {
	h.runContainerMonitorWithClient(h.driver.ctx, h.podmanClient, containerMonitorRetryInterval)
}

func (h *TaskHandle) runContainerMonitorWithClient(ctx context.Context, client containerLifecycleClient, retryInterval time.Duration) {
	h.logger.Debug("Monitoring container lifecycle", "container", h.containerID)

	for h.isRunning() {
		waitErr := client.ContainerWait(ctx, h.containerID, []string{"exited"})
		if ctx.Err() != nil {
			return
		}

		inspectData, inspectErr := client.ContainerInspect(ctx, h.containerID)
		switch {
		case inspectErr == nil && containerExited(inspectData.State):
			h.markContainerExited(inspectData, nil)
			return
		case errors.Is(inspectErr, api.ContainerNotFound):
			h.markContainerExited(api.InspectContainerData{}, inspectErr)
			return
		case inspectErr == nil:
			h.logger.Warn("Container wait returned while container is not exited", "container", h.containerID, "state", inspectData.State, "error", waitErr)
		default:
			h.logger.Warn("Could not reconcile container after wait returned", "container", h.containerID, "wait_error", waitErr, "inspect_error", inspectErr)
		}

		timer := time.NewTimer(retryInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func containerExited(state *api.InspectContainerState) bool {
	if state == nil || state.Running || state.Paused {
		return false
	}

	switch state.Status {
	case "exited", "stopped":
		return true
	default:
		return state.Dead
	}
}

func (h *TaskHandle) markContainerExited(inspectData api.InspectContainerData, inspectErr error) {
	h.stateLock.Lock()
	defer h.stateLock.Unlock()

	if h.procState == drivers.TaskStateExited {
		return
	}

	h.completedAt = time.Now()
	if inspectErr != nil {
		h.exitResult.Err = fmt.Errorf("Driver was unable to get the exit code. %s: %w", h.containerID, inspectErr)
		h.logger.Error("Failed to inspect stopped container, can not get exit code", "container", h.containerID, "error", inspectErr)
		h.exitResult.Signal = 0
	} else if inspectData.State != nil {
		h.exitResult.ExitCode = int(inspectData.State.ExitCode)
		if inspectData.State.Error != "" {
			h.exitResult.Err = errors.New(inspectData.State.Error)
			h.logger.Error("Container error", "container", h.containerID, "error", h.exitResult.Err)
		}
		if !inspectData.State.FinishedAt.IsZero() {
			h.completedAt = inspectData.State.FinishedAt
		}
		if inspectData.State.OOMKilled {
			h.exitResult.OOMKilled = true
			h.exitResult.Err = errors.New("Podman container killed by OOM killer")
			h.logger.Error("Podman container killed by OOM killer", "container", h.containerID)
		}
	} else {
		h.exitResult.Err = fmt.Errorf("Driver was unable to get the exit code. %s: inspect response has no state", h.containerID)
		h.logger.Error("Stopped container inspect response has no state", "container", h.containerID)
	}

	h.procState = drivers.TaskStateExited
}
