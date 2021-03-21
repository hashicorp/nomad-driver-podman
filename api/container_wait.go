package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ContainerWait waits on a container to met a given condition
func (c *API) ContainerWait(ctx context.Context, name string, conditions []string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/wait?condition=%s", name, strings.Join(conditions, "&condition=")), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
