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
	"sync"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/iopodman"
	"github.com/hashicorp/nomad/client/stats"
	"github.com/hashicorp/nomad/plugins/drivers"
)

var (
	measuredCPUStats = []string{"System Mode", "User Mode", "Percent"}
	measuredMemStats = []string{"RSS"}
)

// TaskHandle is the podman specific handle for exactly one container
type TaskHandle struct {
	containerID string
	logger      hclog.Logger
	driver      *Driver

	totalCPUStats  *stats.CpuStats
	userCPUStats   *stats.CpuStats
	systemCPUStats *stats.CpuStats

	// stateLock syncs access to all fields below
	stateLock sync.RWMutex

	taskConfig  *drivers.TaskConfig
	procState   drivers.TaskState
	startedAt   time.Time
	completedAt time.Time
	exitResult  *drivers.ExitResult

	removeContainerOnExit bool

	containerStats iopodman.ContainerStats
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

		//FIXME implement cpu stats correctly
		//available := shelpers.TotalTicksAvailable()
		//cpus := shelpers.CPUNumCores()

		totalPercent := h.totalCPUStats.Percent(h.containerStats.Cpu * 10e16)
		cs := &drivers.CpuStats{
			SystemMode: h.systemCPUStats.Percent(float64(h.containerStats.System_nano)),
			UserMode:   h.userCPUStats.Percent(float64(h.containerStats.Cpu_nano)),
			Percent:    totalPercent,
			TotalTicks: h.systemCPUStats.TicksConsumed(totalPercent),
			Measured:   measuredCPUStats,
		}

		//h.driver.logger.Info("stats", "cpu", containerStats.Cpu, "system", containerStats.System_nano, "user", containerStats.Cpu_nano, "percent", totalPercent, "ticks", cs.TotalTicks, "cpus", cpus, "available", available)

		ms := &drivers.MemoryStats{
			RSS:      uint64(h.containerStats.Mem_usage),
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

func (h *TaskHandle) runContainerMonitor() {

	timer := time.NewTimer(0)
	interval := time.Second * 1
	h.logger.Debug("Monitoring container", "container", h.containerID)

	cleanup := func() {
		h.logger.Debug("Container monitor exits", "container", h.containerID)
	}
	defer cleanup()

	for {
		select {
		case <-h.driver.ctx.Done():
			return

		case <-timer.C:
			timer.Reset(interval)
		}

		containerStats, err := h.driver.podmanClient.GetContainerStats(h.containerID)
		h.logger.Debug("Container stats", "container", h.containerID, "stats", fmt.Sprintf("%#v", containerStats))
		if err != nil {
			if _, ok := err.(*iopodman.NoContainerRunning); ok {
				h.logger.Debug("Container is not running anymore", "container", h.containerID)
				// container was stopped, get exit code and other post mortem infos
				inspectData, err := h.driver.podmanClient.InspectContainer(h.containerID)
				h.stateLock.Lock()
				if err != nil {
					h.exitResult.Err = fmt.Errorf("Driver was unable to get the exit code. %s: %v", h.containerID, err)
					h.logger.Error("Failed to inspect stopped container, can not get exit code", "container", h.containerID, "err", err)
					h.exitResult.Signal = 0
					h.completedAt = time.Now()
				} else {
					h.exitResult.ExitCode = int(inspectData.State.ExitCode)
					if len(inspectData.State.Error) > 0 {
						h.exitResult.Err = fmt.Errorf(inspectData.State.Error)
						h.logger.Error("Container error", "container", h.containerID, "err", fmt.Sprintf("%s", h.exitResult.Err.Error()))
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
			h.logger.Debug("Could not get container stats, unknown error", "err", fmt.Sprintf("%#v", err))
			continue
		}

		h.stateLock.Lock()
		// keep last known containerStats in handle to
		// have it available in the stats emitter
		h.containerStats = *containerStats
		h.stateLock.Unlock()
	}
}
