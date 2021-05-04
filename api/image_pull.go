package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// ImagePull pulls a image from a remote location to local storage
func (c *API) ImagePull(ctx context.Context, nameWithTag string) (string, error) {
	var id string

	res, err := c.Post(ctx, fmt.Sprintf("/v1.0.0/libpod/images/pull?reference=%s", nameWithTag), nil)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf("unknown error, status code: %d: %s", res.StatusCode, body)
	}

	dec := json.NewDecoder(res.Body)
	var report ImagePullReport
	for {
		if err = dec.Decode(&report); err == io.EOF {
			break
		} else if err != nil {
			return "", fmt.Errorf("Error reading response: %w", err)
		}

		if report.Error != "" {
			return "", errors.New(report.Error)
		}

		if report.ID != "" {
			id = report.ID
		}
	}
	return id, nil
}
