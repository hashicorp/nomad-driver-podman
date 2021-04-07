package api

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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
		err := api.ImagePull(ctx, testCase.Image)
		if testCase.Exists {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}
	}
}
