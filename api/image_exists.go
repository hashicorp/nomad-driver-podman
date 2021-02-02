package api

import (
	"context"
	"fmt"
	"net/http"
)

// ImageExists checks if image exists in local store
func (c *API) ImageExists(ctx context.Context, nameWithTag string) (bool, error) {

	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/images/%s/exists", nameWithTag))
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return false, nil
	}
	if res.StatusCode == http.StatusNoContent {
		return true, nil
	}
	if res.StatusCode == http.StatusInternalServerError {
		return true, nil
	}
	return false, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
