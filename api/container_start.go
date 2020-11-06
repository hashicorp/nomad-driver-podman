/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ContainerStart starts a container via id or name
func (c *API) ContainerStart(ctx context.Context, name string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/start", name), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}

	// wait max 10 seconds for running state
	timeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = c.ContainerWait(timeout, name, "running")
	return err
}
