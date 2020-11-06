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

// ContainerWait waits on a container to met a given condition
func (c *API) ContainerWait(ctx context.Context, name string, condition string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/wait?condition=%s", name, condition), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
