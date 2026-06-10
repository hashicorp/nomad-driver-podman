// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// NetworkExists checks if a Podman network with the given name exists.
// Returns true if the network exists (HTTP 204), false if not found (HTTP 404).
// Note: The /networks/{name}/exists endpoint was introduced in Podman API v4.
// Unlike other endpoints that work with /v1.0.0, this one requires /v4.0.0+.
func (c *API) NetworkExists(ctx context.Context, name string) (bool, error) {

	res, err := c.Get(ctx, fmt.Sprintf("/v4.0.0/libpod/networks/%s/exists", url.PathEscape(name)))
	if err != nil {
		return false, err
	}

	defer ignoreClose(res.Body)

	switch res.StatusCode {
	case http.StatusNoContent:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status code checking network %q: %d", name, res.StatusCode)
	}
}
