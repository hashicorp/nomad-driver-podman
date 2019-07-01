package main

import (
	"context"
	"os"
	"strconv"
	"sync"
	"syscall"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad/client/stats"
	"github.com/hashicorp/nomad/plugins/drivers"
)

const (
	// containerMonitorIntv is the interval at which the driver checks if the
	// container is still alive
	containerMonitorIntv = 2 * time.Second
)

type TaskHandle struct {
	//container *lxc.Container
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
	for {
		select {
		case <-ctx.Done():
			return

		case <-timer.C:
			timer.Reset(interval)
		}

		//fixme implement stats
		t := time.Now()

		taskResUsage := drivers.TaskResourceUsage{
			ResourceUsage: &drivers.ResourceUsage{},
			Timestamp:     t.UTC().UnixNano(),
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
	// FIXME: implement
	return nil
}

// waitTillStopped blocks and returns true when container stops;
// returns false with an error message if the container processes cannot be identified.
//
// fixme comment
// Use this in preference to c.Wait() - lxc Wait() function holds a write lock on the container
// blocking any other operation on container, including looking up container stats
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
