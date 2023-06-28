// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"testing"

	"github.com/shoenig/test/must"
)

func Test_parse(t *testing.T) {
	cases := []struct {
		name       string
		repository string
		exp        *ImageSpecs
	}{
		{
			name:       "normal format",
			repository: "docker.io/library/bash:5",
			exp: &ImageSpecs{
				Domain: "docker.io",
				Path:   []string{"library", "bash"},
				Tag:    "5",
			},
		},
		{
			name:       "no tag",
			repository: "docker.io/library/bash",
			exp: &ImageSpecs{
				Domain: "docker.io",
				Path:   []string{"library", "bash"},
				Tag:    "latest",
			},
		},
		{
			name:       "http prefix",
			repository: "http://docker.io/library/bash:5",
			exp: &ImageSpecs{
				Domain: "docker.io",
				Path:   []string{"library", "bash"},
				Tag:    "5",
			},
		},
		{
			name:       "https prefix",
			repository: "https://docker.io/library/bash:5",
			exp: &ImageSpecs{
				Domain: "docker.io",
				Path:   []string{"library", "bash"},
				Tag:    "5",
			},
		},
		{
			name:       "single path element",
			repository: "example.com/app:version",
			exp: &ImageSpecs{
				Domain: "example.com",
				Path:   []string{"app"},
				Tag:    "version",
			},
		},
		{
			name:       "triple path element",
			repository: "example.com/one/two/three:v1",
			exp: &ImageSpecs{
				Domain: "example.com",
				Path:   []string{"one", "two", "three"},
				Tag:    "v1",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := parse(tc.repository)
			must.Eq(t, tc.exp, repo)
		})
	}
}

func TestRepository_Match(t *testing.T) {
	cases := []struct {
		name       string
		repository string
		index      string
		exp        bool
	}{
		{
			name:       "exact",
			repository: "docker.io/library/bash",
			index:      "docker.io/library/bash",
			exp:        true,
		},
		{
			name:       "domain",
			repository: "docker.io/library/bash",
			index:      "docker.io",
			exp:        true,
		},
		{
			name:       "sub path",
			repository: "docker.io/library/bash",
			index:      "docker.io/library",
			exp:        true,
		},
		{
			name:       "wrong domain",
			repository: "docker.io/library/bash",
			index:      "other.com",
			exp:        false,
		},
		{
			name:       "wrong path",
			repository: "docker.io/library/bash",
			index:      "docker.io/other/bash",
			exp:        false,
		},
		{
			name:       "extended path",
			repository: "docker.io/library/bash",
			index:      "docker.io/library/bash/other",
			exp:        false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := parse(tc.repository)
			result := repo.Match(tc.index)
			must.Eq(t, tc.exp, result, must.Sprintf(
				"expect %t, repository: %s, index: %s",
				tc.exp,
				tc.repository,
				tc.index,
			))
		})
	}
}

func Test_decode(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		expUser string
		expPass string
	}{
		{
			name:    "ok",
			input:   "dXNlcjE6cGFzczE=",
			expUser: "user1",
			expPass: "pass1",
		},
		{
			name:  "no split",
			input: "dXNlcjE=",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user, pass := decode(tc.input)
			must.Eq(t, tc.expUser, user)
			must.Eq(t, tc.expPass, pass)
		})
	}
}
