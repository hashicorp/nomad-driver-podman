// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"encoding/json"
	"os"
)

// [Glossary]
// registry: The service that serves container images.
// domain: The DNS domain name of a registry.
// repository: The domain + complete path of an image (no tag).
// index: The domain[+[sub]path] of an image (no tag) (i.e. repository prefix).
// image: A complete domain + path + tag
// credential: A base64 encoded username + password associated with an index.

// CredentialsFile is the struct that represents the contents of the "auth.json" file
// that stores or references credentials that enable authentication into specific
// registries or even repositories.
//
// reference: https://www.mankier.com/5/containers-auth.json
// alternate: https://man.archlinux.org/man/containers-auth.json.5
type CredentialsFile struct {
	Auths map[string]EncAuth `json:"auths"`

	// CredHelpers is currently not supported by the podman task driver
	// CredHelpers map[string]string  `json:"credHelpers"`
}

// EncAuth represents a single registry (or specific repository) and the associated
// base64 encoded auth token.
type EncAuth struct {
	Auth string `json:"auth"`
}

func authFromFileConfig(filename string) AuthBackend {
	return func(repository string) (*RegistryAuthConfig, error) {
		repo := parse(repository)

		cFile, err := loadCredentialsFile(filename)
		if err != nil {
			return nil, err
		}

		rac := cFile.lookup(repo)
		return rac, nil
	}
}

func loadCredentialsFile(path string) (*CredentialsFile, error) {
	if path == "" {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var cFile CredentialsFile

	dec := json.NewDecoder(f)
	if err := dec.Decode(&cFile); err != nil {
		return nil, err
	}
	return &cFile, nil
}

func (cf *CredentialsFile) lookup(repository *ImageSpecs) *RegistryAuthConfig {
	if cf == nil {
		return nil
	}

	// first look for any static auth that applies
	for index, credential := range cf.Auths {
		if repository.Match(index) {
			user, pass := decode(credential.Auth)
			return &RegistryAuthConfig{
				Username: user,
				Password: pass,
			}
		}
	}

	// TODO: add support for specifying credentials helpers that can be used
	// via credsHelpers in this credentials file.

	return nil
}
