// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ExecInspect returns low-level information about an exec instance.
func (c *API) ExecInspect(ctx context.Context, sessionId string) (InspectExecSession, error) {

	var inspectData InspectExecSession

	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/exec/%s/json", sessionId))
	if err != nil {
		return inspectData, err
	}

	defer ignoreClose(res.Body)

	if res.StatusCode != http.StatusOK {
		return inspectData, fmt.Errorf("cannot inspect exec session, status code: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return inspectData, err
	}
	err = json.Unmarshal(body, &inspectData)
	if err != nil {
		return inspectData, err
	}

	return inspectData, nil
}
