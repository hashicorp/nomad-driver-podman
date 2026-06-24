// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type imageLoadResponse struct {
	Names []string `json:"Names"`
}

// ImageLoad uploads a tar archive on the local filesystem as an image.
func (c *API) ImageLoad(ctx context.Context, path string) (string, error) {
	archive, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer ignoreClose(archive)
	return c.ImageLoadReader(ctx, archive)
}

// ImageLoadReader uploads a tar archive read from archive as an image. The
// caller is responsible for closing archive. This allows loading an archive
// from any source (e.g. a local file or an HTTP response body) without first
// staging it on disk.
func (c *API) ImageLoadReader(ctx context.Context, archive io.Reader) (string, error) {
	response := imageLoadResponse{}

	res, err := c.Post(ctx, "/v1.0.0/libpod/images/load", archive)
	if err != nil {
		return "", err
	}

	defer ignoreClose(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot load image, status code: %d: %s", res.StatusCode, body)
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	if len(response.Names) == 0 {
		return "", fmt.Errorf("unknown error: image load successful but reply empty")
	}
	return response.Names[0], nil
}
