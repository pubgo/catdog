package catdog_server

import (
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

/*
// WithTLS sets the TLS config for the catdog_service
func WithTLS(t *tls.Config) Option {
	return func() {
		xerror.Exit(Default.Server.Init(grpcS.AuthTLS(t)))
	}
}

//sets the address of the internal_catdog_server
func Address(addr string) Option {
	return func(o *Options) {
		xerror.Exit(Default.Server.Init(server.Address(addr)))
	}
}

// name of the catdog_service
func Description(n string) Option {
	return func(o *Options) {
		xerror.Exit(Default.Server.Init(server.Description(n)))
	}
}

// Version of the catdog_service
func Version(v string) Option {
	return func(o *Options) {
		xerror.Exit(Default.Server.Init(server.Version(v)))
	}
}

// Metadata associated with the catdog_service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		xerror.Exit(Default.Server.Init(server.Metadata(md)))
	}
}

// RegisterTTL specifies the TTL to use when registering the catdog_service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		xerror.Exit(Default.Server.Init(server.RegisterTTL(t)))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		xerror.Exit(Default.Server.Init(server.RegisterInterval(t)))
	}
}
*/
