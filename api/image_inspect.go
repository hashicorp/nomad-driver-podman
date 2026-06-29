// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var ImageNotFound = errors.New("No such Image")

type inspectIDImageResponse struct {
	Id string `json:"Id"`
}

// ImageInspectID image and returns the image unique identifier
func (c *API) ImageInspectID(ctx context.Context, image string) (string, error) {
	var inspectData inspectIDImageResponse

	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/images/%s/json", image))
	if err != nil {
		return "", err
	}

	defer ignoreClose(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode == http.StatusNotFound {
		return "", ImageNotFound
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot inspect image, status code: %d: %s", res.StatusCode, body)
	}
	err = json.Unmarshal(body, &inspectData)
	if err != nil {
		return "", err
	}

	return inspectData.Id, nil
}

type inspectPlatformImageResponse struct {
	Id           string `json:"Id"`
	Os           string `json:"Os"`
	Architecture string `json:"Architecture"`
	Variant      string `json:"Variant"`
}

// ImageInspectIDForPlatform inspects a local image by name and returns its image
// ID, verifying that the image's recorded operating system, architecture and
// variant match the requested platform. Empty platform fields are not checked.
//
// Unlike ImageInspectID, which resolves an image by name regardless of its
// platform, this returns an error when the stored image does not match the
// requested os/arch/variant. This lets callers confirm that a platform-specific
// pull produced the expected image rather than just any image with that name.
func (c *API) ImageInspectIDForPlatform(ctx context.Context, image, os, arch, variant string) (string, error) {
	var inspectData inspectPlatformImageResponse

	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/images/%s/json", image))
	if err != nil {
		return "", err
	}

	defer ignoreClose(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode == http.StatusNotFound {
		return "", ImageNotFound
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot inspect image, status code: %d: %s", res.StatusCode, body)
	}
	if err := json.Unmarshal(body, &inspectData); err != nil {
		return "", err
	}

	if os != "" && !strings.EqualFold(os, inspectData.Os) {
		return "", fmt.Errorf("image %s os %q does not match requested %q", image, inspectData.Os, os)
	}
	if arch != "" && !strings.EqualFold(arch, inspectData.Architecture) {
		return "", fmt.Errorf("image %s architecture %q does not match requested %q", image, inspectData.Architecture, arch)
	}
	if variant != "" && !strings.EqualFold(variant, inspectData.Variant) {
		return "", fmt.Errorf("image %s variant %q does not match requested %q", image, inspectData.Variant, variant)
	}

	return inspectData.Id, nil
}
