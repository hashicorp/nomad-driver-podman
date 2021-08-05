package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mitchellh/go-linereader"
)

var ContainerNotFound = errors.New("No such Container")
var ContainerWrongState = errors.New("Container has wrong state")

// ContainerStats data takes a name or ID of a container returns stats data
func (c *API) ContainerStats(ctx context.Context, name string) (Stats, error) {

	var stats Stats
	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/stats?stream=false", name))
	if err != nil {
		return stats, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return stats, ContainerNotFound
	}

	if res.StatusCode == http.StatusConflict {
		return stats, ContainerWrongState
	}
	if res.StatusCode != http.StatusOK {
		return stats, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return stats, err
	}
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return stats, err
	}

	return stats, nil
}

// ContainerStatsStream streams stats for all containers
func (c *API) ContainerStatsStream(ctx context.Context) (chan ContainerStats, error) {

	res, err := c.GetStream(ctx, "/v1.0.0/libpod/containers/stats?stream=true")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	statsChannel := make(chan ContainerStats, 5)
	lr := linereader.New(res.Body)

	go func() {
		c.logger.Debug("Running stats stream")
		defer func() {
			res.Body.Close()
			close(statsChannel)
			c.logger.Debug("Stopped stats stream")
		}()
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("Stopping stats stream")
				return
			case line, ok := <-lr.Ch:
				if !ok {
					c.logger.Debug("Stats reader channel was closed")
					return
				}
				var statsReport ContainerStatsReport
				if jerr := json.Unmarshal([]byte(line), &statsReport); jerr != nil {
					c.logger.Error("Unable to unmarshal statsreport", "err", jerr)
					return
				}
				if statsReport.Error != nil {
					c.logger.Error("Stats stream is broken", "error", statsReport.Error)
					return
				}
				for _, stat := range statsReport.Stats {
					statsChannel <- stat
				}
			}
		}
	}()

	return statsChannel, nil
}
