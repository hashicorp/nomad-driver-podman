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

// ImagePull pulls a image from a remote location to local storage
func (c *API) ImagePull(ctx context.Context, nameWithTag string, auth ImageAuthConfig) (string, error) {
	var id string
	headers := map[string]string{}

	// handle authentication
	if auth.Username != "" && auth.Password != "" {
		authHeader, err := NewAuthHeader(auth)
		if err != nil {
			return "", err
		}
		headers["X-Registry-Auth"] = authHeader
	}

	res, err := c.PostWithHeaders(ctx, fmt.Sprintf("/v1.0.0/libpod/images/pull?reference=%s", nameWithTag), nil, headers)
	if err != nil {
		return "", err
	}

	defer ignoreClose(res.Body)
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("cannot pull image, status code: %d: %s", res.StatusCode, body)
	}

	dec := json.NewDecoder(res.Body)
	var report ImagePullReport
	for {
		if err = dec.Decode(&report); err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("Error reading response: %w", err)
		}

		if report.Error != "" {
			return "", errors.New(report.Error)
		}

		if report.ID != "" {
			id = report.ID
		}
	}
	return id, nil
}
