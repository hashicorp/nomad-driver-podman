// Copyright IBM Corp. 2019, 2026
// SPDX-License-Identifier: MPL-2.0

// Package hclcompat provides HCL compatibility types that were previously
// imported from github.com/hashicorp/nomad/helper/pluginutils/hclutils.
// These types are replicated here to eliminate the transitive dependency
// on Docker packages that came through that import path.
package hclcompat

import (
	"github.com/hashicorp/go-msgpack/v2/codec"
)

// MapStrInt is a wrapper for map[string]int that handles
// deserialization from different HCL2 JSON representations
// that were supported since Nomad 0.8.
type MapStrInt map[string]int

func (s *MapStrInt) CodecEncodeSelf(enc *codec.Encoder) {
	v := []map[string]int{*s}
	enc.MustEncode(v)
}

func (s *MapStrInt) CodecDecodeSelf(dec *codec.Decoder) {
	ms := []map[string]int{}
	dec.MustDecode(&ms)

	r := map[string]int{}
	for _, m := range ms {
		for k, v := range m {
			r[k] = v
		}
	}
	*s = r
}

// MapStrStr is a wrapper for map[string]string that handles
// deserialization from different HCL2 JSON representations
// that were supported since Nomad 0.8.
type MapStrStr map[string]string

func (s *MapStrStr) CodecEncodeSelf(enc *codec.Encoder) {
	v := []map[string]string{*s}
	enc.MustEncode(v)
}

func (s *MapStrStr) CodecDecodeSelf(dec *codec.Decoder) {
	ms := []map[string]string{}
	dec.MustDecode(&ms)

	r := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			r[k] = v
		}
	}
	*s = r
}
