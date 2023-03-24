// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"encoding/binary"
	"fmt"
	"io"
)

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

	_, err = io.ReadFull(r, buffer[0:length])
	if err != nil {
		return
	}
	return buffer[0:length], nil
}
