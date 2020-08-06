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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SystemInfo returns information on the system and libpod configuration
func (c *APIClient) SystemInfo(ctx context.Context) (Info, error) {

	var infoData Info

	// the libpod/info endpoint seems to have some trouble
	// using "compat" endpoint and minimal struct
	// until podman returns proper data.
	res, err := c.Get(ctx, "/info")
	if err != nil {
		return infoData, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return infoData, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return infoData, err
	}
	err = json.Unmarshal(body, &infoData)
	if err != nil {
		return infoData, err
	}

	return infoData, nil
}

// for now, we only use the podman version, so it's enough to
// declare this tiny struct
type Info struct {
	Version string `json:"ServerVersion"`
}
