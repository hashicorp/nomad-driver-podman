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
	"github.com/hashicorp/go-hclog"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type APIClient struct {
	baseUrl    string
	httpClient *http.Client
	logger     hclog.Logger
}

func NewClient(logger hclog.Logger) *APIClient {
	ac := &APIClient{
		logger: logger,
	}
	ac.SetSocketPath(GuessSocketPath())
	return ac
}

func (c *APIClient) SetSocketPath(baseUrl string) {
	c.logger.Debug("http baseurl", "url", baseUrl)
	c.httpClient = &http.Client{
		Timeout: 60 * time.Second,
	}
	if strings.HasPrefix(baseUrl, "unix:") {
		c.baseUrl = "http://u"
		path := strings.TrimPrefix(baseUrl, "unix:")
		c.httpClient.Transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		}
	} else {
		c.baseUrl = baseUrl
	}
}

// GuessSocketPath returns the default unix domain socket path for root or non-root users
func GuessSocketPath() string {
	uid := os.Getuid()
	// are we root?
	if uid == 0 {
		return "unix:/run/podman/podman.sock"
	}
	// not? then let's try the default per-user socket location
	return fmt.Sprintf("unix:/run/user/%d/podman/podman.sock", uid)
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
