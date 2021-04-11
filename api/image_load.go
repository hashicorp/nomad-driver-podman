package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// ImageLoad uploads a tar archive as an image
func (c *API) ImageLoad(ctx context.Context, path string) (string, error) {
	archive, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer archive.Close()
	names := []string{}

	res, err := c.Post(ctx, "/v1.0.0/libpod/images/load", archive)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}
	err = json.Unmarshal(body, &names)
	if err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("unknown error: image load successful but reply empty")
	}
	return names[0], nil
}
