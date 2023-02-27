// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"github.com/hashicorp/go-hclog"
)

func newApi() *API {
	return NewClient(hclog.Default(), DefaultClientConfig())
}
