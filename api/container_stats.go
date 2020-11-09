package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
