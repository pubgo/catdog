package catdog_abc

import (
	"crypto/tls"
	"time"

	"github.com/pubgo/xerror"

	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/client"
	grpcClient "github.com/micro/go-micro/v3/client/grpc"
	"github.com/micro/go-micro/v3/model"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/server"
	grpcServer "github.com/micro/go-micro/v3/server/grpc"
)

func OptionOf(opts ...Option) []Option {
	return opts
}

// HandleSignal toggles automatic installation of the signal handler that
// traps TERM, INT, and QUIT.  Users of this feature to disable the signal
// handler, should control liveness of the catdog_service through the context.
func HandleSignal(b bool) Option {
	return func(o *Options) {
		o.IsSignal = b
	}
}

// Address sets the address of the internal_catdog_server
func Address(addr string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Address(addr)))
	}
}

// name of the catdog_service
func Name(n string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Name(n)))
	}
}

func IsSignal(b bool) Option {
	return func(o *Options) {
		o.IsSignal = b
	}
}

// Version of the catdog_service
func Version(v string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Version(v)))
	}
}

// Metadata associated with the catdog_service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Metadata(md)))
	}
}

// RegisterTTL specifies the TTL to use when registering the catdog_service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.RegisterTTL(t)))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.RegisterInterval(t)))
	}
}

// WrapClient is a convenience method for wrapping a Client with
// some middleware component. A list of wrappers can be provided.
// Wrappers are applied in reverse order so the last is executed first.
func WrapClient(w ...client.Wrapper) Option {
	return func(o *Options) {
		// apply in reverse
		for i := len(w); i > 0; i-- {
			o.Client = w[i-1](o.Client)
		}
	}
}

// WrapCall is a convenience method for wrapping a Client CallFunc
func WrapCall(w ...client.CallWrapper) Option {
	return func(o *Options) {
		xerror.Exit(o.Client.Init(client.WrapCall(w...)))
	}
}

// WrapHandler adds a handler Wrapper to a list of options passed into the internal_catdog_server
func WrapHandler(w ...server.HandlerWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapHandler(wrap))
		}

		// initCatDog once
		xerror.Exit(o.Server.Init(wrappers...))
	}
}

// WrapSubscriber adds a subscriber Wrapper to a list of options passed into the internal_catdog_server
func WrapSubscriber(w ...server.SubscriberWrapper) Option {
	return func(o *Options) {
		var wrappers []server.Option

		for _, wrap := range w {
			wrappers = append(wrappers, server.WrapSubscriber(wrap))
		}

		// initCatDog once
		xerror.Exit(o.Server.Init(wrappers...))
	}
}

// Before and Afters

// BeforeStart Run funcs before catdog_service starts
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop Run funcs before catdog_service stops
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart Run funcs after catdog_service starts
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop Run funcs after catdog_service stops
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}

// WithTLS sets the TLS config for the catdog_service
func WithTLS(t *tls.Config) Option {
	return func(o *Options) {
		xerror.Exit(o.Client.Init(grpcClient.AuthTLS(t)))
		xerror.Exit(o.Server.Init(grpcServer.AuthTLS(t)))
	}
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		xerror.Exit(o.Client.Init(client.Broker(b)))
		xerror.Exit(o.Server.Init(server.Broker(b)))
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

// Server sets the catdog_server_plugin for handling requests
func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}

// Model sets the Model for data access
func Model(m model.Model) Option {
	return func(o *Options) {
		o.Model = m
	}
}

// Registry sets the Registry for the catdog_service
// and the underlying components
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Registry(r)))
		xerror.Exit(o.Broker.Init(broker.Registry(r)))
		xerror.Exit(o.Client.Init(client.Registry(r)))
	}
}
