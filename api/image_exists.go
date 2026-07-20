// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type imageListEntry struct {
	Id       string   `json:"Id"`
	RepoTags []string `json:"RepoTags"`
}

// ImageExists reports whether an image with the given reference is present in
// local storage, returning its image ID when found.
//
// Unlike ImageInspectID, which resolves an image by name via
// GET /libpod/images/{name}/json, this lists all local images and matches the
// reference against each image's RepoTags. The name-based inspect can fail to
// resolve locally-built images depending on how their names were recorded in
// storage, returning a spurious 404. This affects both "localhost/"-prefixed
// references and bare shortnames stored under the implicit "localhost/"
// registry. Listing avoids that name-resolution path and mirrors
// `podman images` semantics.
//
// The unfiltered list endpoint is used deliberately: the server-side
// {"reference":[...]} filter returns HTTP 500 on some Podman versions (e.g.
// 5.7.x), so matching is performed client-side for version robustness.
func (c *API) ImageExists(ctx context.Context, image string) (string, bool, error) {
	res, err := c.Get(ctx, "/v1.0.0/libpod/images/json")
	if err != nil {
		return "", false, err
	}
	defer ignoreClose(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", false, err
	}
	if res.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("cannot list images, status code: %d: %s", res.StatusCode, body)
	}

	var images []imageListEntry
	if err := json.Unmarshal(body, &images); err != nil {
		return "", false, err
	}

	if id, found := imageIDForReference(image, images); found {
		return id, true, nil
	}
	return "", false, nil
}

// imageIDForReference returns the image ID of the first listed image whose
// RepoTags contain a match for reference.
func imageIDForReference(reference string, images []imageListEntry) (string, bool) {
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if referenceMatchesTag(reference, tag) {
				return img.Id, true
			}
		}
	}
	return "", false
}

// referenceMatchesTag reports whether an image reference matches a stored
// RepoTag. It matches exactly, and also tolerates a reference that omits the
// implicit "localhost/" registry prefix that Podman records for locally-built
// images (e.g. reference "busybox:local" matching RepoTag
// "localhost/busybox:local").
func referenceMatchesTag(reference, tag string) bool {
	if reference == tag {
		return true
	}
	return "localhost/"+reference == tag
}
