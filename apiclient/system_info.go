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
	res, err := c.Get(ctx, fmt.Sprintf("/info"))
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
