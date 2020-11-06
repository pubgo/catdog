package catdog_app

import (
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"

	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"

	"github.com/pubgo/catdog/catdog_client"
	"github.com/pubgo/catdog/catdog_server"
)

func init() {
	// 在start前执行wrap或者handler
	xerror.Exit(dix.WithBeforeStart(func() {
		xerror.Exit(dix.Dix(enableClientWrap{}))
		xerror.Exit(dix.Dix(enableCallWrap{}))
		xerror.Exit(dix.Dix(enableServerHandler{}))
		xerror.Exit(dix.Dix(enableSubscriber{}))
	}))
}

type enableClientWrap struct{ dix.Model }

// Address
// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) error {
	return xerror.Wrap(dix.Dix(func(*enableClientWrap) {
		for i := len(w); i > 0; i-- {
			catdog_client.Default.Client = w[i-1](catdog_client.Default.Client)
		}
	}))
}

type enableCallWrap struct{ dix.Model }

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) error {
	return xerror.Wrap(dix.Dix(func(*enableClientWrap) {
		xerror.Exit(catdog_client.Default.Client.Init(client.WrapCall(w...)))
	}))
}

type enableServerHandler struct{ dix.Model }

// WrapHandler adds a handler Wrapper to a list of options passed into the internal_catdog_server
func WrapHandler(w ...server.HandlerWrapper) error {
	return xerror.Wrap(dix.Dix(func(*enableServerHandler) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// initCatDog once
		xerror.Exit(catdog_server.Default.Server.Init(wrappers...))
	}))
}

type enableSubscriber struct{ dix.Model }

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the internal_catdog_server
func WrapSubscriber(w ...server.SubscriberWrapper) error {
	return xerror.Wrap(dix.Dix(func(*enableSubscriber) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		xerror.Exit(catdog_server.Default.Server.Init(wrappers...))
	}))
}
