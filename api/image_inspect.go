package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ImageNotFound = errors.New("No such Image")

type inspectIDImageResponse struct {
	Id string `json:"Id"`
}

// Inspects image and returns the image unique identifier
func (c *API) ImageInspectID(ctx context.Context, image string) (string, error) {
	var inspectData inspectIDImageResponse

	res, err := c.Get(ctx, fmt.Sprintf("/v1.0.0/libpod/images/%s/json", image))
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode == http.StatusNotFound {
		return "", ImageNotFound
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}
	err = json.Unmarshal(body, &inspectData)
	if err != nil {
		return "", err
	}

	return inspectData.Id, nil
}
