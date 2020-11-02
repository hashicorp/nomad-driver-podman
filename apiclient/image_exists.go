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

package apiclient

import (
	"context"
	"fmt"
	"net/http"
)

// ImageExists checks if image exists in local store
func (c *APIClient) ImageExists(ctx context.Context, nameWithTag string) (bool, error) {

	res, err := c.Get(ctx, fmt.Sprintf("/libpod/images/%s/exists", nameWithTag))
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
