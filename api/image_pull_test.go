// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/nomad-driver-podman/registry"
	"github.com/shoenig/test/must"
)

func TestApi_Image_Pull(t *testing.T) {
	api := newApi()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	testCases := []struct {
		Image  string
		Exists bool
	}{
		{Image: "docker.io/library/busybox", Exists: true},
		{Image: "docker.io/library/busybox:unstable", Exists: true},
		{Image: "docker.io/library/busybox:notag", Exists: false},
		{Image: "docker.io/non-existent/image:tag", Exists: false},
	}

	for _, testCase := range testCases {
		id, err := api.ImagePull(ctx, &registry.PullConfig{
			Image: testCase.Image,
		})
		if testCase.Exists {
			must.NoError(t, err)
			must.NotEq(t, "", id)
		} else {
			must.Error(t, err)
		}
	}
}
