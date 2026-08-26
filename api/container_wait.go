// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ContainerWait waits on a container to meet one of the given conditions.
// It uses the streaming client because a container may run longer than the
// regular API client timeout.
func (c *API) ContainerWait(ctx context.Context, name string, conditions []string) error {
	path := fmt.Sprintf("/v1.0.0/libpod/containers/%s/wait?condition=%s", name, strings.Join(conditions, "&condition="))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseUrl+path, nil)
	if err != nil {
		return err
	}
	res, err := c.Do(req, true)
	if err != nil {
		return err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("cannot wait for container, status code: %d", res.StatusCode)
}
