// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ContainerInspect data takes a name or ID of a container returns the inspection data
func (c *API) ContainerInspect(ctx context.Context, name string) (InspectContainerData, error) {

	var inspectData InspectContainerData

	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/json", name))
	if err != nil {
		return inspectData, err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode == http.StatusNotFound {
		return inspectData, ContainerNotFound
	}

	if res.StatusCode != http.StatusOK {
		return inspectData, fmt.Errorf("cannot inspect container, status code: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return inspectData, err
	}
	err = json.Unmarshal(body, &inspectData)
	if err != nil {
		return inspectData, err
	}

	return inspectData, nil
}
