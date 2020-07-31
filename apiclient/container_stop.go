package apiclient

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerStop stops a container given a timeout.  It takes the name or ID of a container as well as a
// timeout value.  The timeout value the time before a forcible stop to the container is applied.
// If the container cannot be found, a [ContainerNotFound](#ContainerNotFound)
// error will be returned instead.
func (c *APIClient) ContainerStop(ctx context.Context, name string, timeout int) error {

	res, err := c.Post(ctx, fmt.Sprintf("/containers/%s/stop?timeout=%d", name, timeout))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
