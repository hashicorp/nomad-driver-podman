// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestContainerWaitUsesStreamingClient(t *testing.T) {
	var gotMethod string
	var gotPath string
	var gotConditions []string

	client := &API{
		baseUrl: "http://podman",
		httpClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			t.Fatal("regular HTTP client was used for a blocking wait request")
			return nil, nil
		})},
		httpStreamClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotMethod = req.Method
			gotPath = req.URL.Path
			gotConditions = req.URL.Query()["condition"]
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("0")),
				Header:     make(http.Header),
			}, nil
		})},
	}

	err := client.ContainerWait(context.Background(), "container-id", []string{"running", "exited"})
	if err != nil {
		t.Fatalf("ContainerWait returned an error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Fatalf("method = %q; want %q", gotMethod, http.MethodPost)
	}
	if gotPath != "/v1.0.0/libpod/containers/container-id/wait" {
		t.Fatalf("path = %q", gotPath)
	}
	if len(gotConditions) != 2 || gotConditions[0] != "running" || gotConditions[1] != "exited" {
		t.Fatalf("conditions = %#v", gotConditions)
	}
}

func TestContainerWaitPropagatesStreamingError(t *testing.T) {
	wantErr := errors.New("stream disconnected")
	client := &API{
		baseUrl: "http://podman",
		httpStreamClient: &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, wantErr
		})},
	}

	err := client.ContainerWait(context.Background(), "container-id", []string{"exited"})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}
