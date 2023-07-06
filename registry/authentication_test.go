// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestPullConfig_NoAuth(t *testing.T) {
	cases := []struct {
		name string
		pc   *PullConfig
		exp  bool
	}{
		{
			name: "task config",
			pc: &PullConfig{
				RegistryConfig: &RegistryAuthConfig{
					Username: "user",
					Password: "pass",
				},
			},
			exp: true,
		},
		{
			name: "creds helper",
			pc: &PullConfig{
				CredentialsHelper: "helper.sh",
			},
			exp: true,
		},
		{
			name: "creds file",
			pc: &PullConfig{
				CredentialsFile: "auth.json",
			},
			exp: true,
		},
		{
			name: "none",
			pc:   &PullConfig{},
			exp:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.pc.AuthAvailable()
			must.Eq(t, tc.exp, result)
		})
	}
}
