package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mitchellh/go-linereader"
)

// PodmanEvent is the common header for all events
type PodmanEvent struct {
	Type   string
	Action string
}

// ContainerEvent is a generic PodmanEvent for a single container
// example json:
// {"Type":"container","Action":"create","Actor":{"ID":"cc0d7849692360df2cba94eafb2715b9deec0cbd96ec41c3329dd8636cd070ce","Attributes":{"containerExitCode":"0","image":"docker.io/library/redis:latest","name":"redis-6f2b07a8-73e9-7098-83e1-55939851d46d"}},"scope":"local","time":1609413164,"timeNano":1609413164982188073}
type ContainerEvent struct {
	// create/init/start/stop/died
	Action   string              `json:"Action"`
	Scope    string              `json:"scope"`
	TimeNano uint64              `json:"timeNano"`
	Time     uint32              `json:"time"`
	Actor    ContainerEventActor `json:"Actor"`
}

type ContainerEventActor struct {
	ID         string                   `json:"ID"`
	Attributes ContainerEventAttributes `json:"Attributes"`
}

type ContainerEventAttributes struct {
	Image             string `json:"image"`
	Name              string `json:"name"`
	ContainerExitCode string `json:"containerExitCode"`
}

// ContainerStartEvent is emitted when a container completely started
type ContainerStartEvent struct {
	ID   string
	Name string
}

// ContainerDiedEvent is emitted when a container exited
type ContainerDiedEvent struct {
	ID       string
	Name     string
	ExitCode int
}

// LibpodEventStream streams podman events
func (c *API) LibpodEventStream(ctx context.Context) (chan interface{}, error) {

	res, err := c.GetStream(ctx, "/v1.0.0/libpod/events?stream=true")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	eventsChannel := make(chan interface{}, 5)
	lr := linereader.New(res.Body)

	go func() {
		c.logger.Debug("Running libpod event stream")
		defer func() {
			res.Body.Close()
			close(eventsChannel)
			c.logger.Debug("Stopped libpod event stream")
		}()
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("Stopping libpod event stream")
				return
			case line, ok := <-lr.Ch:
				if !ok {
					c.logger.Debug("Event reader channel was closed")
					return
				}
				var podmanEvent PodmanEvent
				err := json.Unmarshal([]byte(line), &podmanEvent)
				if err != nil {
					c.logger.Error("Unable to parse libpod event", "error", err)
					// no need to stop the stream, maybe we can parse the next event
					continue
				}
				c.logger.Trace("libpod event", "event", line)
				if podmanEvent.Type == "container" {
					var containerEvent ContainerEvent
					err := json.Unmarshal([]byte(line), &containerEvent)
					if err != nil {
						c.logger.Error("Unable to parse ContainerEvent", "error", err)
						// no need to stop the stream, maybe we can parse the next event
						continue
					}
					switch containerEvent.Action {
					case "start":
						eventsChannel <- ContainerStartEvent{
							ID:   containerEvent.Actor.ID,
							Name: containerEvent.Actor.Attributes.Name,
						}
						continue
					case "died":
						i, err := strconv.Atoi(containerEvent.Actor.Attributes.ContainerExitCode)
						if err != nil {
							c.logger.Error("Unable to parse ContainerEvent exitCode", "error", err)
							// no need to stop the stream, maybe we can parse the next event
							continue
						}
						eventsChannel <- ContainerDiedEvent{
							ID:       containerEvent.Actor.ID,
							Name:     containerEvent.Actor.Attributes.Name,
							ExitCode: i,
						}
						continue
					}
					// no action specific parser? emit what we've got
					eventsChannel <- containerEvent
					continue
				}

				// emit a generic event if we do not have a parser for it
				eventsChannel <- podmanEvent
			}
		}
	}()

	return eventsChannel, nil
}
