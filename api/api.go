// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

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
	apiVersion       string
	baseUrl          string
	defaultPodman    bool
	appArmor         bool
	cgroupV2         bool
	cgroupMgr        string
	rootless         bool
	httpClient       *http.Client
	httpStreamClient *http.Client
	logger           hclog.Logger
	defaultTimeout   time.Duration
}

type ClientConfig struct {
	SocketPath    string
	HttpTimeout   time.Duration
	DefaultPodman bool
}

func DefaultClientConfig() ClientConfig {
	cfg := ClientConfig{
		HttpTimeout: 60 * time.Second,
	}
	uid := os.Getuid()
	if uid == 0 {
		cfg.SocketPath = "unix:///run/podman/podman.sock"
	} else {
		cfg.SocketPath = fmt.Sprintf("unix:///run/user/%d/podman/podman.sock", uid)
	}
	cfg.DefaultPodman = true
	return cfg
}

func NewClient(logger hclog.Logger, config ClientConfig) *API {
	ac := &API{
		logger:        logger,
		defaultPodman: config.DefaultPodman,
	}

	ac.logger.Debug("http baseurl", "url", config.SocketPath)
	ac.defaultTimeout = config.HttpTimeout
	ac.httpClient = ac.CreateHttpClient(config.HttpTimeout, config.SocketPath, false)
	ac.httpStreamClient = ac.CreateHttpClient(config.HttpTimeout, config.SocketPath, true)
	if strings.HasPrefix(config.SocketPath, "unix:") {
		ac.baseUrl = "http://u"
	} else {
		ac.baseUrl = config.SocketPath
	}

	return ac
}

func (c *API) SetAPIVersion(v string) {
	c.apiVersion = v
}

func (c *API) GetAPIVersion() string {
	return c.apiVersion
}

func (c *API) IsDefaultClient() bool {
	return c.defaultPodman
}

func (c *API) SetClientAsDefault(d bool) {
	c.defaultPodman = d
}

func (c *API) SetCgroupV2(isV2 bool) {
	c.cgroupV2 = isV2
}

func (c *API) IsCgroupV2() bool {
	return c.cgroupV2
}

func (c *API) SetCgroupMgr(mgr string) {
	c.cgroupMgr = mgr
}

func (c *API) GetCgroupMgr() string {
	return c.cgroupMgr
}

func (c *API) SetRootless(isRootless bool) {
	c.rootless = isRootless
}

func (c *API) IsRootless() bool {
	return c.rootless
}

func (c *API) SetAppArmor(appArmorEnabled bool) {
	c.appArmor = appArmorEnabled
}

func (c *API) IsAppArmorEnabled() bool {
	return c.appArmor
}

func (c *API) Do(req *http.Request, streaming bool) (*http.Response, error) {
	if !streaming {
		return c.httpClient.Do(req)
	} else {
		return c.httpStreamClient.Do(req)
	}
}

func (c *API) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, false)
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
	// If a context was passed with a Deadline longer than our HTTP client timeout,
	// the deadline will ignored. So use the streaming HTTP client (with no timeout)
	// so that the deadline is respected.
	deadline, _ := ctx.Deadline()
	if time.Until(deadline) > c.defaultTimeout {
		return c.Do(req, true)
	} else {
		return c.Do(req, false)
	}
}

func (c *API) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, false)
}

func ignoreClose(c io.Closer) {
	_ = c.Close()
}

func (c *API) CreateHttpClient(timeout time.Duration, baseUrl string, streaming bool) *http.Client {
	var httpClient *http.Client
	if !streaming {
		httpClient = &http.Client{Timeout: timeout}
	} else { // Streaming doesn't have a timeout
		httpClient = &http.Client{}
	}

	if strings.HasPrefix(baseUrl, "unix:") {
		path := strings.TrimPrefix(baseUrl, "unix:")
		httpClient.Transport = &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", path)
			},
		}
	}

	return httpClient
}
