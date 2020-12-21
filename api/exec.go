package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	docker "github.com/docker/docker/api/types"
)

type ExecCreateConfig struct {
	docker.ExecConfig
}

type ExecCreateResponse struct {
	docker.IDResponse
}

type ExecStartConfig struct {
	Detach bool `json:"Detach"`
	Tty    bool `json:"Tty"`
}

func (c *API) CreateExec(ctx context.Context, container string, opts ExecCreateConfig) (*ExecCreateResponse, error) {
	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, err
	}

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/exec", container), bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	response := &ExecCreateResponse{}
	if err = json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response, nil
}
