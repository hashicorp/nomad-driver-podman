// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"encoding/base64"
	"encoding/json"

	"github.com/hashicorp/go-hclog"
)

// PullConfig represents the necessary information needed for an image pull query.
type PullConfig struct {
	Image             string
	RegistryConfig    *RegistryAuthConfig
	TLSVerify         bool
	CredentialsFile   string
	CredentialsHelper string
	AuthSoftFail      bool
}

// BasicConfig returns the Basic-Auth level of configuration.
//
// For the podman driver, this is what can be supplied in the task config block.
func (pc *PullConfig) BasicConfig() *TaskAuthConfig {
	if pc == nil || pc.RegistryConfig == nil {
		return nil
	}
	return &TaskAuthConfig{
		Username:      pc.RegistryConfig.Username,
		Password:      pc.RegistryConfig.Password,
		Email:         pc.RegistryConfig.Email,
		ServerAddress: pc.RegistryConfig.ServerAddress,
	}
}

// AuthAvailable returns true if any of the supported authentication methods are
// set, indicating we should try to find credentials for the given image by
// checking each of the options for looking up credentials.
// - RegistryConfig (from task configuration)
// - CredentialsFile (from plugin configuration)
// - CredentialsHelper (from plugin configuration)
func (pc *PullConfig) AuthAvailable() bool {
	switch {
	case !pc.RegistryConfig.Empty():
		return true
	case pc.CredentialsFile != "":
		return true
	case pc.CredentialsHelper != "":
		return true
	default:
		return false
	}
}

// Log the components of PullConfig in a safe way.
func (pc *PullConfig) Log(logger hclog.Logger) {
	var (
		repository = pc.Image
		tls        bool
		username   string
		email      string
		server     string
		password   bool
		token      bool
		helper     = pc.CredentialsHelper
		file       = pc.CredentialsFile
	)

	if r := pc.RegistryConfig; r != nil {
		tls = pc.TLSVerify
		username = r.Username
		email = r.Email
		server = r.ServerAddress
		password = r.Password != ""
		token = r.IdentityToken != ""
	}

	logger.Info("pull config",
		"repository", repository,
		"username", username,
		"password", password,
		"email", email,
		"server", server,
		"tls", tls,
		"token", token,
		"helper", helper,
		"file", file,
	)
}

// RegistryAuthConfig represents the actualized authentication for accessing a
// container registry.
//
// Note that IdentityToken is mutually exclusive with the remaining fields.
type RegistryAuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
	IdentityToken string `json:"identitytoken,omitempty"`
}

// SetHeader will apply the X-Registry-Auth header (if necessary) to headers.
func (r *RegistryAuthConfig) SetHeader(headers map[string]string) {
	if !r.Empty() {
		b, err := json.Marshal(r)
		if err != nil {
			return
		}
		header := base64.StdEncoding.EncodeToString(b)
		headers["X-Registry-Auth"] = header
	}
}

// Empty returns true if all of username, password, and identity token are unset.
func (r *RegistryAuthConfig) Empty() bool {
	return r == nil || (r.Username == "" && r.Password == "" && r.IdentityToken == "")
}

// TaskAuthConfig represents the "auth" section of the config block for
// the podman task driver.
type TaskAuthConfig struct {
	Username      string
	Password      string
	Email         string
	ServerAddress string
}

func (c *TaskAuthConfig) IsEmpty() bool {
	return c == nil || (c.Username == "" && c.Password == "")
}

// ResolveRegistryAuthentication will find a compatible AuthBackend for the given repository.
// In order, try using
// - auth block from task
// - auth from auth.json file specified by plugin config
// - auth from a credentials helper specified by plugin config
func ResolveRegistryAuthentication(repository string, pullConfig *PullConfig) (*RegistryAuthConfig, error) {
	return firstValidAuth(repository, pullConfig.AuthSoftFail, []AuthBackend{
		authFromTaskConfig(pullConfig.BasicConfig()),
		authFromFileConfig(pullConfig.CredentialsFile),
		authFromCredsHelper(pullConfig.CredentialsHelper),
	})
}

// An AuthBackend is a function that resolves registry credentials for a given
// repository. If no auth exitsts for the given repository, (nil, nil) is
// returned and should be skipped.
type AuthBackend func(string) (*RegistryAuthConfig, error)

// firstValidAuth returns the first RegistryAuthConfig associated with repository.
//
// If softFail is set, ignore error return values from an authBackend and pretend
// like they simply do not recognize repository, and continue searching.
func firstValidAuth(repository string, softFail bool, authBackends []AuthBackend) (*RegistryAuthConfig, error) {
	for _, backend := range authBackends {
		auth, err := backend(repository)
		switch {
		case auth != nil:
			return auth, nil
		case err != nil && softFail:
			continue
		case err != nil:
			return nil, err
		}
	}
	return nil, nil
}

func authFromTaskConfig(taskAuthConfig *TaskAuthConfig) AuthBackend {
	return func(string) (*RegistryAuthConfig, error) {
		if taskAuthConfig.IsEmpty() {
			return nil, nil
		}
		return &RegistryAuthConfig{
			Username:      taskAuthConfig.Username,
			Password:      taskAuthConfig.Password,
			Email:         taskAuthConfig.Email,
			ServerAddress: taskAuthConfig.ServerAddress,
		}, nil
	}
}
