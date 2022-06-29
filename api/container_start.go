package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// ContainerStart starts a container via id or name
func (c *API) ContainerStart(ctx context.Context, name string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/start", name), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("cannot start container, status code: %d: %s", res.StatusCode, body)
	}

	// wait max 10 seconds for running state
	// TODO: make timeout configurable
	timeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = c.ContainerWait(timeout, name, []string{"running", "exited"})
	return err
}
