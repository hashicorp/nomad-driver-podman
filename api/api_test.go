// Copyright IBM Corp. 2019, 2025
// SPDX-License-Identifier: MPL-2.0

package api

import (
	"github.com/hashicorp/go-hclog"
)

func newApi() *API {
	return NewClient(hclog.Default(), DefaultClientConfig())
}
