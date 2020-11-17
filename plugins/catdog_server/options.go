package catdog_server

import (
	"context"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/xerror"
)

// Adds a handler Wrapper to a list of options passed into the server
func WrapHandler(w server.HandlerWrapper) error {
	return xerror.Wrap(Default.Server.Init(func(o *server.Options) {
		o.HdlrWrappers = append(o.HdlrWrappers, w)
	}))
}

// Adds a subscriber Wrapper to a list of options passed into the server
func WrapSubscriber(w server.SubscriberWrapper) error {
	return xerror.Wrap(Default.Server.Init(func(o *server.Options) {
		o.SubWrappers = append(o.SubWrappers, w)
	}))
}

func init() {
	WrapHandler(func(handlerFunc server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {

		}
	})
}
