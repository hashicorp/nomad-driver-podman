package api

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerKill sends a signal to a container
func (c *API) ContainerKill(ctx context.Context, name string, signal string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/kill?signal=%s", name, signal), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
