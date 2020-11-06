/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.
*/

package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ImagePull pulls a image from a remote location to local storage
func (c *API) ImagePull(ctx context.Context, nameWithTag string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/images/pull?reference=%s", nameWithTag), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}

	return nil
}
