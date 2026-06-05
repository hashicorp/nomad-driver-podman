// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

func TestApi_Network_Exists(t *testing.T) {
	api := newApi()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testCases := []struct {
		Name        string
		NetworkName string
		Exists      bool
	}{
		{Name: "default podman network exists", NetworkName: "podman", Exists: true},
		{Name: "non-existent network returns false", NetworkName: "this-network-does-not-exist-xyz", Exists: false},
	}

	for _, testCase := range testCases {
		exists, err := api.NetworkExists(ctx, testCase.NetworkName)
		must.NoError(t, err)
		must.Eq(t, testCase.Exists, exists)
	}
}
