// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	defer ignoreClose(res.Body)

	if res.StatusCode == http.StatusNotFound {
		return stats, ContainerNotFound
	}

	if res.StatusCode == http.StatusConflict {
		return stats, ContainerWrongState
	}
	if res.StatusCode != http.StatusOK {
		return stats, fmt.Errorf("cannot get stats of container, status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return stats, err
	}

	// Since podman 4.1.1, an empty 200 response is returned for stopped containers.
	if len(body) == 0 {
		return stats, ContainerNotFound
	}

	// Since podman 4.6.0, a 200 response with `container is stopped` is returned for stopped containers.
	var errResponse Error
	if _ := json.Unmarshal(body, &errResponse); errResponse.Cause == "container is stopped" {
		return stats, ContainerNotFound
	}

	err = json.Unmarshal(body, &stats)
	if err != nil {
		return stats, err
	}

	return stats, nil
}
