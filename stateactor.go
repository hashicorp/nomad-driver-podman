package main

import (
	"context"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/nomad-driver-podman/api"
	"github.com/hashicorp/nomad/plugins/drivers"
)

type GetTaskHandleMsg struct {
	TaskID string
	Result chan *TaskHandle
}

type TaskStartedMsg struct {
	TaskID     string
	TaskHandle *TaskHandle
	Done       chan bool
}

type TaskDeletedMsg struct {
	TaskID string
}

type StartStatesEmitterMsg struct {
	TaskID              string
	Interval            time.Duration
	TaskResourceChannel chan *drivers.TaskResourceUsage
}

type WaitTaskMsg struct {
	TaskID            string
	Ctx               context.Context
	ExitResultChannel chan *drivers.ExitResult
}

// the state actor serializes requests to the TaskHandle by using a message channel.
// it also holds taskid->handle and containerid->handle maps
// and routes incoming podman stats and podman events to TaskHandles
func runStateActor(ctx context.Context, actorMessageChannel <-chan interface{}, logger hclog.Logger) error {

	taskIdToHandle := map[string]*TaskHandle{}
	containerIdToHandle := map[string]*TaskHandle{}

	go func() {
		logger.Debug("Starting state actor")

		for m := range actorMessageChannel {
			switch message := m.(type) {

			case GetTaskHandleMsg:
				taskHandle, ok := taskIdToHandle[message.TaskID]
				if ok {
					message.Result <- taskHandle
				} else {
					close(message.Result)
				}

			case TaskStartedMsg:
				logger.Trace("on started", "message", message)
				taskIdToHandle[message.TaskID] = message.TaskHandle
				containerIdToHandle[message.TaskHandle.containerID] = message.TaskHandle
				// unblock caller
				message.Done <- true

			case WaitTaskMsg:
				logger.Trace("on wait task", "message", message)
				taskHandle, ok := taskIdToHandle[message.TaskID]
				if ok {
					go taskHandle.runExitWatcher(message.Ctx, message.ExitResultChannel)
				}

			case StartStatesEmitterMsg:
				logger.Trace("on start emitter", "message", message)
				taskHandle, ok := taskIdToHandle[message.TaskID]
				if ok {
					go taskHandle.runStatsEmitter(ctx, message.TaskResourceChannel, message.Interval)
					taskHandle.statsEmitterRunning = true
				}

			case api.ContainerStats:
				taskHandle, ok := containerIdToHandle[message.ContainerID]
				// avoid to block the actor if the target stats emitter is not running
				if ok && taskHandle.statsEmitterRunning {
					taskHandle.containerStatsChannel <- message
				}

			case api.ContainerDiedEvent:
				logger.Trace("on died", "message", message)
				taskHandle, ok := containerIdToHandle[message.ID]
				if ok {
					go taskHandle.onContainerDied()
				}

			case TaskDeletedMsg:
				logger.Trace("on deleted", "message", message)
				taskHandle, ok := taskIdToHandle[message.TaskID]
				if ok {
					delete(taskIdToHandle, message.TaskID)
					delete(containerIdToHandle, taskHandle.containerID)
				}

			}
		}
		logger.Debug("State actor stopped")
	}()
	return nil
}
