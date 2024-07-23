// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// ContainerCreate creates a new container
func (c *API) ContainerCreate(ctx context.Context, create SpecGenerator) (ContainerCreateResponse, error) {

	response := ContainerCreateResponse{}

	jsonString, err := json.Marshal(create)
	if err != nil {
		return response, err
	} else {
		f, _ := os.OpenFile("/tmp/debug.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()
		_, _ = f.Write([]byte(jsonString));
		_, _ = f.WriteString("\n---\n");
	}


	res, err := c.Post(ctx, "/v1.0.0/libpod/containers/create", bytes.NewBuffer(jsonString))
	if err != nil {
		return response, err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		return response, fmt.Errorf("cannot create container, status code: %d: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, err
}

type ContainerCreateRequest struct {
	// Name is the name the container will be given.
	// If no name is provided, one will be randomly generated.
	// Optional.
	Name string `json:"name,omitempty"`

	// Command is the container's command.
	// If not given and Image is specified, this will be populated by the
	// image's configuration.
	// Optional.
	Command []string `json:"command,omitempty"`

	// Entrypoint is the container's entrypoint.
	// If not given and Image is specified, this will be populated by the
	// image's configuration.
	// Optional.
	Entrypoint []string `json:"entrypoint,omitempty"`

	// WorkDir is the container's working directory.
	// If unset, the default, /, will be used.
	// Optional.
	WorkDir string `json:"work_dir,omitempty"`
	// Env is a set of environment variables that will be set in the
	// container.
	// Optional.
	Env map[string]string `json:"env,omitempty"`
}

type ContainerCreateResponse struct {
	Id       string
	Warnings []string
}
