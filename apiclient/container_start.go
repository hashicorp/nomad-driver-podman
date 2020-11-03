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
	"time"
)

// ContainerStart starts a container via id or name
func (c *APIClient) ContainerStart(ctx context.Context, name string) error {

	res, err := c.Post(ctx, fmt.Sprintf("/%s/libpod/containers/%s/start", PODMAN_API_VERSION, name), nil)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}

	// wait max 10 seconds for running state
	timeout, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	err = c.ContainerWait(timeout, name, "running")
	return err
}
