/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

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
	return false, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
