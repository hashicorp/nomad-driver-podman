package api

import (
	"github.com/hashicorp/go-hclog"
)

func newApi() *API {
	return NewClient(hclog.Default(), DefaultClientConfig())
}
