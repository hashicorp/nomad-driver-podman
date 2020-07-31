package apiclient

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerStart starts a container via id or name
func (c *APIClient) ContainerStart(ctx context.Context, name string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/containers/%s/start", name))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
