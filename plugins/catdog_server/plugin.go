package catdog_server

import (
	"crypto/tls"
	"github.com/pubgo/dix/dix_run"

	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xerror"
	"github.com/spf13/pflag"
)

func init() {
	opts := Default.Options()

	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "server",
		OnFlags: func(flags *pflag.FlagSet) {
			flags.StringVar(&opts.Address, "server_addr", opts.Address, "server address")
			flags.StringVar(&opts.Name, "server_name", opts.Name, "server name")
		},
		OnInit: func(r reader.Value) {
			xerror.Exit(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
				xerror.Exit(Default.Server.Init(server.Name(opts.Name)))
				xerror.Exit(Default.Server.Init(server.Address(opts.Address)))

				var t *tls.Config
				// WithTLS sets the TLS config for the catdog_service
				xerror.Exit(Default.Init(grpcS.AuthTLS(t)))
			}))
		},
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
