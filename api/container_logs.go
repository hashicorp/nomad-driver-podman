// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ContainerLogs gets stdout and stderr logs from a container.
func (c *API) ContainerLogs(ctx context.Context, name string, since time.Time, stdout io.Writer, stderr io.Writer) error {

	c.logger.Debug("Running log stream", "container", name)
	res, err := c.GetStream(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/logs?follow=true&since=%d&stdout=true&stderr=true", name, since.Unix()))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot get logs from container, status code: %d", res.StatusCode)
	}

	defer func() {
		ignoreClose(res.Body)
		c.logger.Debug("Stopped log stream", "container", name)
	}()
	buffer := make([]byte, 1024)
	for {
		fd, l, err := DemuxHeader(res.Body, buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		frame, err := DemuxFrame(res.Body, buffer, l)
		if err != nil {
			return err
		}

		c.logger.Trace("logger", "frame", string(frame), "fd", fd)

		switch fd {
		case 0:
			_, _ = stdout.Write(frame)
		case 1:
			_, _ = stdout.Write(frame)
		case 2:
			_, _ = stderr.Write(frame)
		case 3:
			return fmt.Errorf("Error from log service: %s", string(frame))
		default:
			return fmt.Errorf("Unknown log stream identifier: %d", fd)
		}
	}

}
