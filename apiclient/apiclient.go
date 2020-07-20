package apiclient

import (
	"context"
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
		Timeout: 5 * time.Second,
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

func (c *APIClient) Post(ctx context.Context, path string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.baseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	return c.Do(req)
}
