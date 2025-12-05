// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	// helperTimeout is the maximum amount of time to allow a credential helper
	// to execute before force killing it
	helperTimeout = 10 * time.Second
)

func helpCmd(ctx context.Context, name, index string) (*exec.Cmd, error) {
	executable := fmt.Sprintf("docker-credential-%s", name)

	path, err := exec.LookPath(executable)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, path)
	cmd.Args = []string{"get"}
	cmd.Stdin = strings.NewReader(index)
	return cmd, nil
}

type Credential struct {
	Username string `json:"Username"`
	Secret   string `json:"Secret"`
}

func (c *Credential) Empty() bool {
	return c == nil || (c.Username == "" && c.Secret == "")
}

func authFromCredsHelper(name string) AuthBackend {
	if name == "" {
		return noBackend
	}
	return func(repository string) (*RegistryAuthConfig, error) {
		repo := parse(repository)
		index := repo.Index()

		ctx, cancel := context.WithTimeout(context.Background(), helperTimeout)
		defer cancel()

		cmd, err := helpCmd(ctx, name, index)
		if err != nil {
			return nil, err
		}

		b, err := cmd.CombinedOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to run credential-helper %q: %w", name, err)
		}

		var credential Credential
		if err := json.Unmarshal(b, &credential); err != nil {
			return nil, fmt.Errorf("failed to read credential from credential-helper %q: %w", name, err)
		}

		if credential.Empty() {
			return nil, nil
		}

		return &RegistryAuthConfig{
			Username: credential.Username,
			Password: credential.Secret,
		}, nil
	}
}
