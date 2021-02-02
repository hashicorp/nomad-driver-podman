package api

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/hashicorp/nomad/plugins/drivers"
)

// ExecStartRequest prepares to stream a exec session
type ExecStartRequest struct {

	// streams
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	// terminal size channel
	ResizeCh <-chan drivers.TerminalSize

	// Tty indicates whether pseudo-terminal is to be allocated
	Tty bool

	// AttachOutput is whether to attach to STDOUT
	// If false, stdout will not be attached
	AttachOutput bool
	// AttachError is whether to attach to STDERR
	// If false, stdout will not be attached
	AttachError bool
	// AttachInput is whether to attach to STDIN
	// If false, stdout will not be attached
	AttachInput bool
}

// DemuxHeader reads header for stream from server multiplexed stdin/stdout/stderr/2nd error channel
func DemuxHeader(r io.Reader, buffer []byte) (fd, sz int, err error) {
	n, err := io.ReadFull(r, buffer[0:8])
	if err != nil {
		return
	}
	if n < 8 {
		err = io.ErrUnexpectedEOF
		return
	}

	fd = int(buffer[0])
	if fd < 0 || fd > 3 {
		err = fmt.Errorf(`channel "%d" found, 0-3 supported`, fd)
		return
	}

	sz = int(binary.BigEndian.Uint32(buffer[4:8]))
	return
}

// DemuxFrame reads contents for frame from server multiplexed stdin/stdout/stderr/2nd error channel
func DemuxFrame(r io.Reader, buffer []byte, length int) (frame []byte, err error) {
	if len(buffer) < length {
		buffer = append(buffer, make([]byte, length-len(buffer)+1)...)
	}

	n, err := io.ReadFull(r, buffer[0:length])
	if err != nil {
		return nil, nil
	}
	if n < length {
		err = io.ErrUnexpectedEOF
		return
	}

	return buffer[0:length], nil
}

// This is intended to be run as a goroutine, handling resizing for a container
// or exec session.
func (c *API) attachHandleResize(ctx context.Context, resizeChannel <-chan drivers.TerminalSize, sessionId string) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Warn("Resize handler is done")
			return
		case size := <-resizeChannel:
			c.logger.Trace("Resize terminal", "sessionId", sessionId, "height", size.Height, "width", size.Width)
			rerr := c.ExecResize(ctx, sessionId, size.Height, size.Width)
			if rerr != nil {
				c.logger.Warn("failed to resize TTY", "err", rerr)
			}
		}
	}
}

// ExecStartAndAttach starts and attaches to a given exec session.
func (c *API) ExecStart(ctx context.Context, sessionID string, options ExecStartRequest) error {
	client := new(http.Client)
	*client = *c.httpClient
	client.Timeout = 0

	var socket net.Conn
	socketSet := false
	dialContext := client.Transport.(*http.Transport).DialContext
	t := &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			c, err := dialContext(ctx, network, address)
			if err != nil {
				return nil, err
			}
			if !socketSet {
				socket = c
				socketSet = true
			}
			return c, err
		},
		IdleConnTimeout: time.Duration(0),
	}
	client.Transport = t

	// Detach is always false.
	// podman reference doc states that "true" is not supported
	execStartReq := struct {
		Detach bool `json:"Detach"`
	}{
		Detach: false,
	}
	jsonBytes, err := json.Marshal(execStartReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/exec/%s/start", c.baseUrl, sessionID), bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return err
	}

	if options.Tty {
		go c.attachHandleResize(ctx, options.ResizeCh, sessionID)
	}

	if options.AttachInput {
		go func() {
			_, err := io.Copy(socket, options.Stdin)
			if err != nil {
				c.logger.Error("Failed to send stdin to exec session", "err", err)
			}
		}()
	}

	defer func() {
		c.logger.Debug("Finish exec session attach")
	}()

	buffer := make([]byte, 1024)
	if options.Tty {
		// If not multiplex'ed, read from server and write to stdout
		_, err := io.Copy(options.Stdout, socket)
		if err != nil {
			return err
		}
	} else {
		for {
			// Read multiplexed channels and write to appropriate stream
			fd, l, err := DemuxHeader(socket, buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
			frame, err := DemuxFrame(socket, buffer, l)
			if err != nil {
				return err
			}

			switch {
			case fd == 0:
				// Write STDIN to STDOUT (echoing characters
				// typed by another attach session)
				if options.AttachInput {
					if _, err := options.Stdout.Write(frame[0:l]); err != nil {
						return err
					}
				}
			case fd == 1:
				if options.AttachOutput {
					if _, err := options.Stdout.Write(frame[0:l]); err != nil {
						return err
					}
				}
			case fd == 2:
				if options.AttachError {
					if _, err := options.Stderr.Write(frame[0:l]); err != nil {
						return err
					}
				}
			case fd == 3:
				return fmt.Errorf("error from service from stream: %s", frame)
			default:
				return fmt.Errorf("unrecognized channel '%d' in header, 0-3 supported", fd)
			}
		}
	}
	return nil
}
