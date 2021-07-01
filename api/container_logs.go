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

	res, err := c.GetStream(ctx, fmt.Sprintf("/v1.0.0/libpod/containers/%s/logs?follow=true&since=%d&stdout=true&stderr=true", name, since.Unix()))
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	c.logger.Debug("Running log stream", "container", name)
	defer func() {
		res.Body.Close()
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
			stdout.Write(frame)
			stdout.Write([]byte("\n"))
		case 1:
			stdout.Write(frame)
			stdout.Write([]byte("\n"))
		case 2:
			stderr.Write(frame)
			stderr.Write([]byte("\n"))
		case 3:
			return fmt.Errorf("Error from log service: %s", string(frame))
		default:
			return fmt.Errorf("Unknown log stream identifier: %d", fd)
		}
	}

}
