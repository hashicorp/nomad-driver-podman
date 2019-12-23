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

	exitChannel    chan *drivers.ExitResult
	containerStats iopodman.ContainerStats
}

func (h *TaskHandle) TaskStatus() *drivers.TaskStatus {
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

func (h *TaskHandle) IsRunning() bool {
	h.stateLock.RLock()
	defer h.stateLock.RUnlock()
	return h.procState == drivers.TaskStateRunning
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

		totalPercent := h.totalCpuStats.Percent(h.containerStats.Cpu * 10e16)
		cs := &drivers.CpuStats{
			SystemMode: h.systemCpuStats.Percent(float64(h.containerStats.System_nano)),
			UserMode:   h.userCpuStats.Percent(float64(h.containerStats.Cpu_nano)),
			Percent:    totalPercent,
			TotalTicks: h.systemCpuStats.TicksConsumed(totalPercent),
			Measured:   MeasuredCpuStats,
		}

		//h.driver.logger.Info("stats", "cpu", containerStats.Cpu, "system", containerStats.System_nano, "user", containerStats.Cpu_nano, "percent", totalPercent, "ticks", cs.TotalTicks, "cpus", cpus, "available", available)

		ms := &drivers.MemoryStats{
			RSS:      uint64(h.containerStats.Mem_usage),
			Measured: MeasuredMemStats,
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

func (h *TaskHandle) MonitorContainer() {

	timer := time.NewTimer(0)
	interval := time.Second * 1
	h.logger.Debug("Monitoring container", "container", h.containerID)

	cleanup := func() {
		if h.exitChannel != nil {
			close(h.exitChannel)
		}
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
		if err != nil {
			if _, ok := err.(*iopodman.NoContainerRunning); ok {
				h.logger.Debug("Container is not running anymore", "container", h.containerID)
				// container was stopped, get exit code and other post mortem infos
				inspectData, err := h.driver.podmanClient.InspectContainer(h.containerID)
				h.stateLock.Lock()
				h.procState = drivers.TaskStateExited
				if err != nil {
					h.logger.Error("Failt to inspect stopped container, can not get exit code", "err", err)
					h.exitResult.Signal = 0
					h.completedAt = time.Now()
				} else {
					h.exitResult.ExitCode = int(inspectData.State.ExitCode)
					h.exitResult.OOMKilled = inspectData.State.OOMKilled
					h.completedAt = inspectData.State.FinishedAt
				}
				if h.exitChannel != nil {
					h.exitChannel <- h.exitResult
				} else {
					// this is bad....
					h.logger.Error("Can not send exit signal to nomad, channel is empty", "container", h.containerID)
				}
				h.stateLock.Unlock()
				return
			} else {
				h.logger.Debug("Could not get container stats, unknown error", "err", fmt.Sprintf("%#v", err))
			}
			// continue and wait for next cycle, it should then
			// fall into the "TaskStateExited" case
			continue
		}

		h.stateLock.Lock()
		// keep last known containerStats in handle to
		// have it available in the stats emitter
		h.containerStats = *containerStats
		h.stateLock.Unlock()
	}
}
