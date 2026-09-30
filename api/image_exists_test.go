// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestApi_imageIDForReference(t *testing.T) {
	images := []imageListEntry{
		{Id: "aaa111", RepoTags: []string{"localhost/testapp:local"}},
		{Id: "bbb222", RepoTags: []string{"docker.io/library/busybox:latest", "docker.io/library/busybox:1.36"}},
		{Id: "ccc333", RepoTags: nil},
	}

	testCases := []struct {
		Name      string
		Reference string
		WantID    string
		WantFound bool
	}{
		{
			Name:      "exact localhost match",
			Reference: "localhost/testapp:local",
			WantID:    "aaa111",
			WantFound: true,
		},
		{
			Name:      "localhost prefix tolerated",
			Reference: "testapp:local",
			WantID:    "aaa111",
			WantFound: true,
		},
		{
			Name:      "bare shortname does not match docker.io registry tag",
			Reference: "busybox:latest",
			WantID:    "",
			WantFound: false,
		},
		{
			Name:      "localhost prefix does not match different tag",
			Reference: "localhost/testapp:other",
			WantID:    "",
			WantFound: false,
		},
		{
			Name:      "fully qualified registry match",
			Reference: "docker.io/library/busybox:latest",
			WantID:    "bbb222",
			WantFound: true,
		},
		{
			Name:      "second tag on same image",
			Reference: "docker.io/library/busybox:1.36",
			WantID:    "bbb222",
			WantFound: true,
		},
		{
			Name:      "no match returns not found",
			Reference: "localhost/other:local",
			WantID:    "",
			WantFound: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			id, found := imageIDForReference(testCase.Reference, images)
			must.Eq(t, testCase.WantFound, found)
			must.Eq(t, testCase.WantID, id)
		})
	}
}
