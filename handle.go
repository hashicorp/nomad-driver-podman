// Copyright (c) HashiCorp, Inc.
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

// TaskHandle is the podman specific handle for exactly one container
type TaskHandle struct {
	containerID string
	logger      hclog.Logger
	driver      *Driver

	totalCPUStats  *cpustats.Tracker
	userCPUStats   *cpustats.Tracker
	systemCPUStats *cpustats.Tracker

	collectionInterval time.Duration

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
	h.logger.Debug("Starting statsEmitter", "container", h.containerID)
	h.collectionInterval = interval
	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("Stopping statsEmitter", "container", h.containerID)
			return
		case <-timer.C:
			timer.Reset(interval)
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
			err = h.driver.podman.ContainerLogs(ctx, h.containerID, since, stdout, stderr)
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

func (h *TaskHandle) runContainerMonitor() {

	timer := time.NewTimer(0)
	h.logger.Debug("Monitoring container", "container", h.containerID)

	for {
		select {
		case <-h.driver.ctx.Done():
			return

		case <-timer.C:
			timer.Reset(h.collectionInterval)
		}

		containerStats, statsErr := h.driver.podman.ContainerStats(h.driver.ctx, h.containerID)
		if statsErr != nil {
			gone := false
			if errors.Is(statsErr, api.ContainerNotFound) {
				gone = true
			} else if errors.Is(statsErr, api.ContainerWrongState) {
				gone = true
			}
			if gone {
				h.logger.Debug("Container is not running anymore", "container", h.containerID, "error", statsErr)
				// container was stopped, get exit code and other post mortem infos
				inspectData, err := h.driver.podman.ContainerInspect(h.driver.ctx, h.containerID)
				h.stateLock.Lock()
				h.completedAt = time.Now()
				if err != nil {
					h.exitResult.Err = fmt.Errorf("Driver was unable to get the exit code. %s: %w", h.containerID, err)
					h.logger.Error("Failed to inspect stopped container, can not get exit code", "container", h.containerID, "error", err)
					h.exitResult.Signal = 0
				} else {
					h.exitResult.ExitCode = int(inspectData.State.ExitCode)
					if len(inspectData.State.Error) > 0 {
						h.exitResult.Err = fmt.Errorf(inspectData.State.Error)
						h.logger.Error("Container error", "container", h.containerID, "error", h.exitResult.Err)
					}
					h.completedAt = inspectData.State.FinishedAt
					if inspectData.State.OOMKilled {
						h.exitResult.OOMKilled = true
						h.exitResult.Err = fmt.Errorf("Podman container killed by OOM killer")
						h.logger.Error("Podman container killed by OOM killer", "container", h.containerID)
					}
				}

				h.procState = drivers.TaskStateExited
				h.stateLock.Unlock()
				return
			}
			// continue and wait for next cycle, it should eventually
			// fall into the "TaskStateExited" case
			h.logger.Debug("Could not get container stats, unknown error", "error", fmt.Sprintf("%#v", statsErr))
			continue
		}

		h.stateLock.Lock()
		// keep last known containerStats in handle to
		// have it available in the stats emitter
		h.containerStats = containerStats
		h.stateLock.Unlock()
	}
}
