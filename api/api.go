// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

type API struct {
	baseUrl          string
	httpClient       *http.Client
	httpStreamClient *http.Client
	logger           hclog.Logger
}

type ClientConfig struct {
	SocketPath    string
	HttpTimeout   time.Duration
	HostUser      string
	DefaultPodman bool
}

func DefaultClientConfig() ClientConfig {
	cfg := ClientConfig{
		HttpTimeout: 60 * time.Second,
	}
	cfg.SocketPath = "unix:///run/podman/podman.sock"
	cfg.HostUser = "root"
	cfg.DefaultPodman = true
	return cfg
}

func NewClient(logger hclog.Logger, config ClientConfig) *API {
	ac := &API{
		logger: logger,
	}

	baseUrl := config.SocketPath
	ac.logger.Debug("http baseurl", "url", baseUrl)
	ac.httpClient = &http.Client{
		Timeout: config.HttpTimeout,
	}
	// we do not want a timeout for streaming requests.
	ac.httpStreamClient = &http.Client{}
	if strings.HasPrefix(baseUrl, "unix:") {
		ac.baseUrl = "http://u"
		path := strings.TrimPrefix(baseUrl, "unix:")
		ac.httpClient.Transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		}
		ac.httpStreamClient.Transport = ac.httpClient.Transport
	} else {
		ac.baseUrl = baseUrl
	}

	return ac
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

func (c *API) GetStream(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	return c.httpStreamClient.Do(req)
}

func (c *API) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	return c.PostWithHeaders(ctx, path, body, map[string]string{})
}

func (c *API) PostWithHeaders(ctx context.Context, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return c.Do(req)
}

func (c *API) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func ignoreClose(c io.Closer) {
	_ = c.Close()
}
