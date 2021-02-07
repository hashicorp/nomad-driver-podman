package api

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mitchellh/go-linereader"
)

// ContainerLogs gets stdout and stderr logs from a container.
func (c *API) ContainerLogs(ctx context.Context, name string, stdout bool, stderr bool, target io.Writer) error {

	res, err := c.GetStream(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/logs?follow=true&tail=5&stdout=%t&stderr=%t", name, stdout, stderr))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	lr := linereader.New(res.Body)

	go func() {
		c.logger.Debug("Running log stream", "container", name, "stdout", stdout, "stderr", stderr)
		defer func() {
			res.Body.Close()
			c.logger.Debug("Stopped log stream", "container", name, "stdout", stdout, "stderr", stderr)
		}()
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("Stopping log stream", "container", name, "stdout", stdout, "stderr", stderr)
				return
			case line, ok := <-lr.Ch:
				if !ok {
					c.logger.Debug("Log stream was closed", "container", name, "stdout", stdout, "stderr", stderr)
					return
				}
				target.Write([]byte(line + "\n"))
			}
		}
	}()

	return nil
}
