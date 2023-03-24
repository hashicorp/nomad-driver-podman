// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func (c *API) ExecResize(ctx context.Context, execId string, height int, width int) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/exec/%s/resize?h=%d&w=%d", execId, height, width), nil)
	if err != nil {
		return err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(res.Body)
		return fmt.Errorf("cannot resize exec session, status code: %d: %s", res.StatusCode, body)
	}

	return err
}
