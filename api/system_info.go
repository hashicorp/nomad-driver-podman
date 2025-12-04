// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SystemInfo returns information on the system and libpod configuration
func (c *API) SystemInfo(ctx context.Context) (Info, error) {

	var infoData Info

	res, err := c.Get(ctx, "/v1.0.0/libpod/info")
	if err != nil {
		return infoData, err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return infoData, fmt.Errorf("cannot fetch Podman system info, status code: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return infoData, err
	}
	err = json.Unmarshal(body, &infoData)
	if err != nil {
		return infoData, err
	}

	return infoData, nil
}
