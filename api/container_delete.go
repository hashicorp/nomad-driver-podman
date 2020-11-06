/*
Copyright 2019 Thomas Weber

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"context"
	"fmt"
	"net/http"
)

// ContainerDelete deletes a container.
// It takes the name or ID of a container.
func (c *API) ContainerDelete(ctx context.Context, name string, force bool, deleteVolumes bool) error {

	res, err := c.Delete(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s?force=%t&v=%t", name, force, deleteVolumes))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
}
