package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *API) ExecResize(ctx context.Context, execId string, height int, width int) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/exec/%s/resize?h=%d&w=%d", execId, height, width), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("cannot resize exec session, status code: %d: %s", res.StatusCode, body)
	}

	return err
}
