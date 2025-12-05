// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"encoding/base64"
	"strings"
)

// noBackend may be used when it is known no AuthBackend exists
// that can be used.
func noBackend(string) (*RegistryAuthConfig, error) {
	return nil, nil
}

// ImageSpecs holds each of the components of a full container image identifier,
// such that we can ask questions about it. There are three parts,
// - domain
// - path
// - tag
type ImageSpecs struct {
	Domain string
	Path   []string
	Tag    string
}

func (is *ImageSpecs) Index() string {
	s := []string{is.Domain}
	s = append(s, is.Path...)
	return strings.Join(s, "/")
}

// Match returns true if the repository belongs to the given registry index.
//
// Given repository example.com/foo/bar, each of these index values would
// match, e.g.
// - example.com
// - examle.com/foo
// - example.com/foo/bar
//
// Whereas other index values would not match, e.g.
// - other.com
// - example.com/baz
// - example.com/foo/bar/baz
func (is *ImageSpecs) Match(index string) bool {
	repository := is.Index()
	if len(index) > len(repository) {
		return false
	}
	if strings.HasPrefix(repository, index) {
		return true
	}
	return false
}

// parse the repository string into a useful object we can interact with to
// ask questions about that repository
func parse(repository string) *ImageSpecs {
	repository = strings.TrimPrefix(repository, "https://")
	repository = strings.TrimPrefix(repository, "http://")
	tagIdx := strings.LastIndex(repository, ":")
	var tag string
	if tagIdx != -1 {
		tag = repository[tagIdx+1:]
		repository = repository[0:tagIdx]
	}
	if tag == "" {
		tag = "latest"
	}
	parts := strings.Split(repository, "/")
	return &ImageSpecs{
		Domain: parts[0],
		Path:   parts[1:],
		Tag:    tag,
	}
}

func decode(b64 string) (string, string) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", ""
	}
	parts := strings.SplitN(string(data), ":", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}
