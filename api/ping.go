// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"net/http"
)

// Ping podman service
// Return lipod api version
func (c *API) Ping(ctx context.Context) (string, error) {

	res, err := c.Get(ctx, "/libpod/_ping")
	if err != nil {
		return "", err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot ping Podman api, status code: %d", res.StatusCode)
	}
	version := res.Header.Get("Libpod-API-Version")
	if version == "" {
		return "", fmt.Errorf("Unable to get libpod api version from response header")
	}
	return version, nil
}
