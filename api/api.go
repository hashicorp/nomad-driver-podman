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

package api

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

type API struct {
	baseUrl    string
	httpClient *http.Client
	logger     hclog.Logger
}

func NewClient(logger hclog.Logger) *API {
	ac := &API{
		logger: logger,
	}
	ac.SetSocketPath(DefaultSocketPath())
	return ac
}

func (c *API) SetSocketPath(baseUrl string) {
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

// DefaultSocketPath returns the default unix domain socket path for root or non-root users
func DefaultSocketPath() string {
	uid := os.Getuid()
	// are we root?
	if uid == 0 {
		return "unix:/run/podman/podman.sock"
	}
	// not? then let's try the default per-user socket location
	return fmt.Sprintf("unix:/run/user/%d/podman/podman.sock", uid)
}

func (c *API) Do(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	return res, err
}

func (c *API) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *API) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return c.Do(req)
}

func (c *API) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
