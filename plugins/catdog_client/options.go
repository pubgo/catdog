package catdog_client

import (
	"github.com/asim/nitro/v3/client"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
)

// WrapClient
// is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) error {
	return xerror.Wrap(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
		for i := len(w); i > 0; i-- {
			Default.Client = w[i-1](Default.Client)
		}
	}))
}

// WrapCall
// is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) error {
	return xerror.Wrap(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
		xerror.Exit(Default.Client.Init(client.WrapCall(w...)))
	}))
}
