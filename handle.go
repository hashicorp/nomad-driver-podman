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

	diedChannel chan bool

	// receive container stats from global podman stats streamer
	containerStatsChannel chan api.ContainerStats

	// stateLock syncs access to all fields below
	stateLock sync.RWMutex

	taskConfig  *drivers.TaskConfig
	procState   drivers.TaskState
	startedAt   time.Time
	completedAt time.Time
	exitResult  *drivers.ExitResult

	removeContainerOnExit bool
	statsEmitterRunning   bool
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

func (h *TaskHandle) runStatsEmitter(ctx context.Context, taskResourceChannel chan *drivers.TaskResourceUsage, interval time.Duration) {
	timer := time.NewTimer(0)

	containerStats := api.ContainerStats{}
	userCPUStats := stats.NewCpuStats()
	systemCPUStats := stats.NewCpuStats()

	h.logger.Debug("Starting statsEmitter", "container", h.containerID)
	for {
		select {

		case <-ctx.Done():
			h.logger.Debug("Stopping statsEmitter", "container", h.containerID)
			return

		case s := <-h.containerStatsChannel:
			// keep latest known container stats in this go routine
			// and convert/emit it to nomad based on interval
			containerStats = s
			continue

		case <-timer.C:
			timer.Reset(interval)

			totalPercent := containerStats.CPU

			cs := &drivers.CpuStats{
				SystemMode: systemCPUStats.Percent(float64(containerStats.CPUSystemNano)),
				UserMode:   userCPUStats.Percent(float64(containerStats.CPUNano)),
				Percent:    totalPercent,
				TotalTicks: systemCPUStats.TicksConsumed(totalPercent),
				Measured:   measuredCPUStats,
			}

			ms := &drivers.MemoryStats{
				Usage:    containerStats.MemUsage,
				RSS:      containerStats.MemUsage,
				Measured: measuredMemStats,
			}

			// send stats to nomad
			taskResourceChannel <- &drivers.TaskResourceUsage{
				ResourceUsage: &drivers.ResourceUsage{
					CpuStats:    cs,
					MemoryStats: ms,
				},
				Timestamp: time.Now().UTC().UnixNano(),
			}
		}
	}
}

func (h *TaskHandle) onContainerDied() {
	h.logger.Debug("Container is not running anymore")
	// container was stopped, get exit code and other post mortem infos
	inspectData, err := h.driver.podman.ContainerInspect(h.driver.ctx, h.containerID)
	h.stateLock.Lock()
	h.completedAt = time.Now()
	if err != nil {
		h.exitResult.Err = fmt.Errorf("Driver was unable to get the exit code. %s: %v", h.containerID, err)
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
	// unblock exitWatcher go routine
	close(h.diedChannel)
}
