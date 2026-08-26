// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad/plugins/drivers"
)

type lifecycleStep struct {
	waitErr    error
	inspect    api.InspectContainerData
	inspectErr error
}

type fakeLifecycleClient struct {
	mu           sync.Mutex
	steps        []lifecycleStep
	waitCalls    int
	inspectCalls int
}

type fakeStatsClient struct {
	cancel context.CancelFunc
	err    error
	calls  int
}

func (f *fakeStatsClient) ContainerStats(context.Context, string) (api.Stats, error) {
	f.calls++
	if f.cancel != nil {
		f.cancel()
	}
	return api.Stats{}, f.err
}

func (f *fakeLifecycleClient) ContainerWait(ctx context.Context, _ string, _ []string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.waitCalls >= len(f.steps) {
		return errors.New("unexpected wait call")
	}
	err := f.steps[f.waitCalls].waitErr
	f.waitCalls++
	return err
}

func (f *fakeLifecycleClient) ContainerInspect(context.Context, string) (api.InspectContainerData, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.inspectCalls >= len(f.steps) {
		return api.InspectContainerData{}, errors.New("unexpected inspect call")
	}
	step := f.steps[f.inspectCalls]
	f.inspectCalls++
	return step.inspect, step.inspectErr
}

func (f *fakeLifecycleClient) calls() (int, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.waitCalls, f.inspectCalls
}

func newLifecycleTestHandle() *TaskHandle {
	return &TaskHandle{
		containerID: "container-id",
		logger:      hclog.NewNullLogger(),
		taskConfig:  &drivers.TaskConfig{ID: "task-id", Name: "task"},
		procState:   drivers.TaskStateRunning,
		exitResult:  &drivers.ExitResult{},
	}
}

func TestContainerMonitorRecordsExitedContainer(t *testing.T) {
	finishedAt := time.Date(2026, time.August, 26, 17, 0, 0, 0, time.UTC)
	client := &fakeLifecycleClient{steps: []lifecycleStep{{
		inspect: api.InspectContainerData{State: &api.InspectContainerState{
			Status:     "exited",
			ExitCode:   23,
			FinishedAt: finishedAt,
		}},
	}}}
	handle := newLifecycleTestHandle()

	handle.runContainerMonitorWithClient(context.Background(), client, 0)

	status := handle.taskStatus()
	if status.State != drivers.TaskStateExited {
		t.Fatalf("state = %q; want exited", status.State)
	}
	if status.ExitResult.ExitCode != 23 {
		t.Fatalf("exit code = %d; want 23", status.ExitResult.ExitCode)
	}
	if !status.CompletedAt.Equal(finishedAt) {
		t.Fatalf("completed at = %s; want %s", status.CompletedAt, finishedAt)
	}
}

func TestContainerMonitorReconcilesAfterWaitFailure(t *testing.T) {
	waitErr := errors.New("event stream disconnected")
	client := &fakeLifecycleClient{steps: []lifecycleStep{
		{
			waitErr: waitErr,
			inspect: api.InspectContainerData{State: &api.InspectContainerState{
				Status:  "running",
				Running: true,
			}},
		},
		{
			inspect: api.InspectContainerData{State: &api.InspectContainerState{
				Status:   "exited",
				ExitCode: 7,
			}},
		},
	}}
	handle := newLifecycleTestHandle()

	handle.runContainerMonitorWithClient(context.Background(), client, 0)

	waitCalls, inspectCalls := client.calls()
	if waitCalls != 2 || inspectCalls != 2 {
		t.Fatalf("calls = wait:%d inspect:%d; want wait:2 inspect:2", waitCalls, inspectCalls)
	}
	status := handle.taskStatus()
	if status.State != drivers.TaskStateExited || status.ExitResult.ExitCode != 7 {
		t.Fatalf("status = %#v", status)
	}
}

func TestContainerMonitorFindsExitAfterWaitFailure(t *testing.T) {
	client := &fakeLifecycleClient{steps: []lifecycleStep{{
		waitErr: errors.New("wait request raced with container exit"),
		inspect: api.InspectContainerData{State: &api.InspectContainerState{
			Status:   "exited",
			ExitCode: 42,
		}},
	}}}
	handle := newLifecycleTestHandle()

	handle.runContainerMonitorWithClient(context.Background(), client, 0)

	waitCalls, inspectCalls := client.calls()
	if waitCalls != 1 || inspectCalls != 1 {
		t.Fatalf("calls = wait:%d inspect:%d; want wait:1 inspect:1", waitCalls, inspectCalls)
	}
	status := handle.taskStatus()
	if status.State != drivers.TaskStateExited || status.ExitResult.ExitCode != 42 {
		t.Fatalf("status = %#v", status)
	}
}

func TestContainerExited(t *testing.T) {
	testCases := []struct {
		name  string
		state *api.InspectContainerState
		want  bool
	}{
		{name: "nil state"},
		{name: "running", state: &api.InspectContainerState{Status: "running", Running: true}},
		{name: "paused", state: &api.InspectContainerState{Status: "paused", Paused: true}},
		{name: "created", state: &api.InspectContainerState{Status: "created"}},
		{name: "exited", state: &api.InspectContainerState{Status: "exited"}, want: true},
		{name: "stopped", state: &api.InspectContainerState{Status: "stopped"}, want: true},
		{name: "dead", state: &api.InspectContainerState{Status: "unknown", Dead: true}, want: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := containerExited(tc.state); got != tc.want {
				t.Fatalf("containerExited() = %t; want %t", got, tc.want)
			}
		})
	}
}

func TestStatsMonitorDoesNotUseStatsErrorsAsLifecycleEvents(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &fakeStatsClient{
		cancel: cancel,
		err:    api.ContainerWrongState,
	}
	handle := newLifecycleTestHandle()

	handle.runStatsMonitorWithClient(ctx, client, 0)

	if client.calls != 1 {
		t.Fatalf("stats calls = %d; want 1", client.calls)
	}
	if !handle.isRunning() {
		t.Fatal("stats error changed task lifecycle state")
	}
}

func TestContainerMonitorStopsWhenContextIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	client := &fakeLifecycleClient{steps: []lifecycleStep{{waitErr: context.Canceled}}}
	handle := newLifecycleTestHandle()

	handle.runContainerMonitorWithClient(ctx, client, 0)

	waitCalls, inspectCalls := client.calls()
	if waitCalls != 1 || inspectCalls != 0 {
		t.Fatalf("calls = wait:%d inspect:%d; want wait:1 inspect:0", waitCalls, inspectCalls)
	}
	if !handle.isRunning() {
		t.Fatal("canceled monitor changed task state")
	}
}
