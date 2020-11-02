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
	"io/ioutil"
	"net/http"
)

// ImagePull pulls a image from a remote location to local storage
func (c *APIClient) ImagePull(ctx context.Context, nameWithTag string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/%s/libpod/images/pull?reference=%s", PODMAN_API_VERSION, nameWithTag), nil)
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
