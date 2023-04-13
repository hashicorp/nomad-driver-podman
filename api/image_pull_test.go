// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"testing"

	"github.com/shoenig/test/must"
)

func TestApi_Image_Pull(t *testing.T) {
	api := newApi()
	ctx := context.Background()

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
		id, err := api.ImagePull(ctx, testCase.Image, ImageAuthConfig{})
		if testCase.Exists {
			must.NoError(t, err)
			must.NotEq(t, "", id)
		} else {
			must.Error(t, err)
		}
	}
}
