package catdog_server

import (
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
)

// WrapHandler
// adds a handler Wrapper to a list of options passed into the internal_catdog_server
func WrapHandler(w ...server.HandlerWrapper) error {
	return xerror.Wrap(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// initCatDog once
		xerror.Exit(Default.Server.Init(wrappers...))
	}))
}

// WrapSubscriber
// adds a subscriber Wrapper to a list of options passed into the internal_catdog_server
func WrapSubscriber(w ...server.SubscriberWrapper) error {
	return xerror.Wrap(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		xerror.Exit(Default.Server.Init(wrappers...))
	}))
}
