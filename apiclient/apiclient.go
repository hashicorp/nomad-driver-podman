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
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

type APIClient struct {
	baseUrl    string
	httpClient *http.Client
}

func NewClient(baseUrl string) *APIClient {
	httpClient := http.Client{
		Timeout: 60 * time.Second,
	}
	if strings.HasPrefix(baseUrl, "unix:") {
		path := strings.TrimPrefix(baseUrl, "unix:")
		baseUrl = "http://localhost"
		httpClient.Transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		}
	}
	return &APIClient{
		baseUrl:    baseUrl,
		httpClient: &httpClient,
	}
}

func (c *APIClient) Do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	return res, err
}

func (c *APIClient) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.Do(req)
}

func (c *APIClient) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

func (c *APIClient) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.Do(req)
}
