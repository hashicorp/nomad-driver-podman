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
	usesAuth := auth.Username != "" && auth.Password != ""
	if usesAuth {
		authHeader, err := NewAuthHeader(auth)
		if err != nil {
			return "", err
		}
		headers["X-Registry-Auth"] = authHeader
	}

	c.logger.Trace("image pull details", "tls_verify", auth.TLSVerify, "reference", nameWithTag, "uses_auth", usesAuth)

	urlPath := fmt.Sprintf("/v1.0.0/libpod/images/pull?reference=%s&tlsVerify=%t", nameWithTag, auth.TLSVerify)
	res, err := c.PostWithHeaders(ctx, urlPath, nil, headers)
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
		decErr := dec.Decode(&report)
		if errors.Is(decErr, io.EOF) {
			break
		} else if decErr != nil {
			return "", fmt.Errorf("Error reading response: %w", decErr)
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
