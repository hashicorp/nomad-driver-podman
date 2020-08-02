package apiclient

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerDelete deletes a container.
// It takes the name or ID of a container.
func (c *APIClient) ContainerDelete(ctx context.Context, name string, force bool, deleteVolumes bool) error {

	res, err := c.Delete(ctx, fmt.Sprintf("/containers/%s?force=%t&v=%t", name, force, deleteVolumes))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
