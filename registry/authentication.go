package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-hclog"
)

// PullConfig represents the necessary information needed for an image pull query.
type PullConfig struct {
	Repository     string
	RegistryConfig *RegistryAuthConfig
	TLSVerify      bool

	// TODO: config file
	// TODO: creds helper
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

// Log the components of PullConfig in a safe way.
func (pc *PullConfig) Log(logger hclog.Logger) {
	var (
		repository = pc.Repository
		tls        bool
		username   = "<unset>"
		email      = "<unset>"
		server     = "<unset>"
		password   bool
		token      bool
	)

	if r := pc.RegistryConfig; r != nil {
		tls = pc.TLSVerify
		username = r.Username
		email = r.Email
		server = r.ServerAddress
		password = r.Password != ""
		token = r.IdentityToken != ""
	}

	// todo: trace
	logger.Info("pull config",
		"repository", repository,
		"username", username,
		"password", password,
		"email", email,
		"server", server,
		"tls", tls,
		"token", token,
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
	return firstValidAuth(repository, []AuthBackend{
		authFromTaskConfig(pullConfig.BasicConfig()),
	})
}

// An AuthBackend is a function that resolves registry credentials for a given
// repository. If no auth exitsts for the given repository, (nil, nil) is
// returned and should be skipped.
type AuthBackend func(string) (*RegistryAuthConfig, error)

func firstValidAuth(repository string, authBackends []AuthBackend) (*RegistryAuthConfig, error) {
	for _, backend := range authBackends {
		auth, err := backend(repository)
		if auth != nil || err != nil {
			return auth, err
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

func authFromFileConfig(filename string) (*RegistryAuthConfig, error) {
	return nil, fmt.Errorf("not implemented yet")
}

func authFromCredsHelper(helper string) (*RegistryAuthConfig, error) {
	return nil, fmt.Errorf("not implemented yet")
}
