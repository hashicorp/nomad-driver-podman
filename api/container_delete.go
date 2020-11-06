/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package api

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerDelete deletes a container.
// It takes the name or ID of a container.
func (c *API) ContainerDelete(ctx context.Context, name string, force bool, deleteVolumes bool) error {

	res, err := c.Delete(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s?force=%t&v=%t", name, force, deleteVolumes))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
