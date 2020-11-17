package catdog_client

import (
	"crypto/tls"

	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
	"github.com/spf13/pflag"
)

func init() {
	opts := Default.Options()
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "client",
		OnFlags: func(flags *pflag.FlagSet) {
			_ = opts
		},
		OnInit: func(r reader.Value) {
			xerror.Exit(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
				var t *tls.Config
				// WithTLS sets the TLS config for the catdog_service
				xerror.Exit(Default.Init(grpcC.AuthTLS(t)))
				xerror.Exit(catdog_server.Default.Init(grpcS.AuthTLS(t)))
			}))
		},
	}))
}
