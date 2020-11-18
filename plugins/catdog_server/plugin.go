package catdog_server

import (
	"crypto/tls"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/dix/dix_run"

	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
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
		OnInit: func(ent catdog_entry.Entry) {
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
