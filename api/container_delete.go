// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerDelete deletes a container.
// It takes the name or ID of a container.
func (c *API) ContainerDelete(ctx context.Context, name string, force bool, deleteVolumes bool) error {

	res, err := c.Delete(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s?force=%t&volumes=%t", name, force, deleteVolumes))
	if err != nil {
		return err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("cannot delete container, status code: %d", res.StatusCode)
}
