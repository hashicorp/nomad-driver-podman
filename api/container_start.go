package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
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
		return fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}

	return err
}
