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

	"github.com/hashicorp/nomad-driver-podman/registry"
)

// ImagePull pulls a image from a remote location to local storage
func (c *API) ImagePull(ctx context.Context, pullConfig *registry.PullConfig) (string, error) {
	pullConfig.Log(c.logger)

	var (
		headers    = make(map[string]string)
		repository = pullConfig.Repository
		tlsVerify  = pullConfig.TLSVerify
	)

	auth, err := registry.ResolveRegistryAuthentication(repository, pullConfig)
	if err != nil {
		return "", fmt.Errorf("failed to determine authentication for %q", repository)
	}
	auth.SetHeader(headers)
	c.logger.Info("HEADERS", "headers", headers)

	urlPath := fmt.Sprintf("/v1.0.0/libpod/images/pull?reference=%s&tlsVerify=%t", repository, tlsVerify)
	c.logger.Info("URL PATH", "urlPath", urlPath)

	res, err := c.PostWithHeaders(ctx, urlPath, nil, headers)
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	defer ignoreClose(res.Body)

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("cannot pull image, status code: %d: %s", res.StatusCode, body)
	}

	var (
		dec    = json.NewDecoder(res.Body)
		report ImagePullReport
		id     string
	)
	for {
		decodeErr := dec.Decode(&report)
		if errors.Is(decodeErr, io.EOF) {
			break
		} else if decodeErr != nil {
			return "", fmt.Errorf("failed to read image pull response: %w", decodeErr)
		} else if report.Error != "" {
			return "", fmt.Errorf("image pull report indicates error: %s", report.Error)
		} else if report.ID != "" {
			id = report.ID
		}
	}
	return id, nil
}
