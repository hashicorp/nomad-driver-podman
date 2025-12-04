// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"os"
	"testing"

	"github.com/shoenig/test/must"
)

func Test_authFromCredsHelper(t *testing.T) {
	cases := []struct {
		name       string
		repository string
		exp        *RegistryAuthConfig
	}{
		{
			name:       "basic",
			repository: "docker.io/library/bash:5",
			exp: &RegistryAuthConfig{
				Username: "user1",
				Password: "pass1",
			},
		},
		{
			name:       "example.com",
			repository: "example.com/some/silly/thing:v1",
			exp: &RegistryAuthConfig{
				Username: "user2",
				Password: "pass2",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("PATH", os.ExpandEnv("${PWD}/tests:${PATH}"))
			be := authFromCredsHelper("fake.sh")
			rac, err := be(tc.repository)
			must.NoError(t, err)
			must.Eq(t, tc.exp, rac)
		})
	}
}
