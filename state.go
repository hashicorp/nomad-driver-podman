package main

import (
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad/plugins/drivers"
)

type taskStore struct {
	taskIdToHandle      map[string]*TaskHandle
	containerIdToHandle map[string]*TaskHandle
	lock                sync.RWMutex
	logger              hclog.Logger
}

func newTaskStore(logger hclog.Logger) *taskStore {
	return &taskStore{
		logger:              logger,
		taskIdToHandle:      map[string]*TaskHandle{},
		containerIdToHandle: map[string]*TaskHandle{},
	}
}

func (ts *taskStore) Set(id string, handle *TaskHandle) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	handle.exitResult = new(drivers.ExitResult)
	handle.diedChannel = make(chan bool)
	ts.taskIdToHandle[id] = handle
	ts.containerIdToHandle[handle.containerID] = handle
}

func (ts *taskStore) Get(id string) (*TaskHandle, bool) {
	ts.lock.RLock()
	defer ts.lock.RUnlock()
	t, ok := ts.taskIdToHandle[id]
	return t, ok
}

func (ts *taskStore) GetByContainerId(containerID string) (*TaskHandle, bool) {
	ts.lock.RLock()
	defer ts.lock.RUnlock()
	t, ok := ts.containerIdToHandle[containerID]
	return t, ok
}

func (ts *taskStore) Delete(id string) {
	ts.lock.Lock()
	defer ts.lock.Unlock()
	t, ok := ts.taskIdToHandle[id]
	if ok {
		delete(ts.containerIdToHandle, t.containerID)
	}
	delete(ts.taskIdToHandle, id)
}

// keep last known containerStats in handle to
// have it available in the stats emitter
func (ts *taskStore) UpdateContainerStats(containerStats api.ContainerStats) bool {
	taskHandle, ok := ts.GetByContainerId(containerStats.ContainerID)
	if ok {
		taskHandle.stateLock.Lock()
		taskHandle.containerStats = containerStats
		taskHandle.stateLock.Unlock()
	}
	return ok
}

// Manage task handle state by consuming libpod events
func (ts *taskStore) HandleLibpodEvent(e interface{}) {
	switch e.(type) {
	case api.ContainerDiedEvent:
		event, _ := e.(api.ContainerDiedEvent)
		taskHandle, ok := ts.GetByContainerId(event.ID)
		if ok {
			taskHandle.onContainerDied(event)
		}
	}

}
