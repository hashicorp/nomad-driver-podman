// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"testing"

	"github.com/shoenig/test/must"
)

func Test_authFromConfigFile(t *testing.T) {
	ab := authFromFileConfig("tests/auth.json")
	must.NotNil(t, ab)

	cases := []struct {
		name    string
		image   string
		expUser string
		expPass string
	}{
		{
			name:    "complete",
			image:   "one.example.com/library/bash:5",
			expUser: "user1",
			expPass: "pass1",
		},
		{
			name:    "partial path",
			image:   "two.example.com/library/bash:5",
			expUser: "user2",
			expPass: "pass2",
		},
		{
			name:    "domain only",
			image:   "three.example.com/library/bash:5",
			expUser: "user3",
			expPass: "pass3",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rac, err := ab(tc.image)
			must.NoError(t, err)
			must.NotNil(t, rac, must.Sprintf("RAC should not be nil"))
			must.Eq(t, tc.expUser, rac.Username)
			must.Eq(t, tc.expPass, rac.Password)
		})
	}
}
