package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ExecInspect returns low-level information about an exec instance.
func (c *API) ExecInspect(ctx context.Context, sessionId string) (InspectExecSession, error) {

	var inspectData InspectExecSession

	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/exec/%s/json", sessionId))
	if err != nil {
		return inspectData, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return inspectData, fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return inspectData, err
	}
	err = json.Unmarshal(body, &inspectData)
	if err != nil {
		return inspectData, err
	}

	return inspectData, nil
}
