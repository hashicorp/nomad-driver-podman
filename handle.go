package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad/client/stats"
	"github.com/hashicorp/nomad/plugins/drivers"
)

var (
	measuredCPUStats = []string{"System Mode", "User Mode", "Percent"}
	measuredMemStats = []string{"Usage"}
)

// TaskHandle is the podman specific handle for exactly one container
type TaskHandle struct {
	containerID string
	logger      hclog.Logger
	driver      *Driver

	totalCPUStats  *stats.CpuStats
	userCPUStats   *stats.CpuStats
	systemCPUStats *stats.CpuStats
	diedChannel    chan bool

	// stateLock syncs access to all fields below
	stateLock sync.RWMutex

	taskConfig  *drivers.TaskConfig
	procState   drivers.TaskState
	startedAt   time.Time
	completedAt time.Time
	exitResult  *drivers.ExitResult

	removeContainerOnExit bool

	containerStats api.ContainerStats
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

	if !h.isRunning() {
		h.logger.Debug("No need to run exitWatcher on a stopped container")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-h.diedChannel:
			return
		}
	}
}

func (h *TaskHandle) runStatsEmitter(ctx context.Context, statsChannel chan *drivers.TaskResourceUsage, interval time.Duration) {
	timer := time.NewTimer(0)
	h.logger.Debug("Starting statsEmitter", "container", h.containerID)
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

		totalPercent := h.containerStats.CPU
		cs := &drivers.CpuStats{
			SystemMode: h.systemCPUStats.Percent(float64(h.containerStats.CPUSystemNano)),
			UserMode:   h.userCPUStats.Percent(float64(h.containerStats.CPUNano)),
			Percent:    totalPercent,
			TotalTicks: h.systemCPUStats.TicksConsumed(totalPercent),
			Measured:   measuredCPUStats,
		}

		ms := &drivers.MemoryStats{
			Usage:    h.containerStats.MemUsage,
			RSS:      h.containerStats.MemUsage,
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

func (h *TaskHandle) onContainerDied(event api.ContainerDiedEvent) {
	h.logger.Debug("Container is not running anymore", "event", event)
	// container was stopped, get exit code and other post mortem infos
	inspectData, err := h.driver.podman.ContainerInspect(h.driver.ctx, h.containerID)
	h.stateLock.Lock()
	h.completedAt = time.Now()
	if err != nil {
		h.exitResult.Err = fmt.Errorf("Driver was unable to get the exit code. %s: %v", h.containerID, err)
		h.logger.Error("Failed to inspect stopped container, can not get exit code", "container", h.containerID, "err", err)
		h.exitResult.Signal = 0
	} else {
		h.exitResult.ExitCode = int(inspectData.State.ExitCode)
		if len(inspectData.State.Error) > 0 {
			h.exitResult.Err = fmt.Errorf(inspectData.State.Error)
			h.logger.Error("Container error", "container", h.containerID, "err", h.exitResult.Err)
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
	// unblock exitWatcher go routine
	close(h.diedChannel)
}
