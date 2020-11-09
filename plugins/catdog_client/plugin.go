package catdog_client

import (
	"crypto/tls"
	"fmt"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/spf13/pflag"

	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
)

func init() {
	opts := Default.Options()
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "client",
		OnFlags: func(flags *pflag.FlagSet) {
			_ = opts
		},
		OnInit: func() {
			xerror.Exit(dix.WithBeforeStart(func() {
				var t *tls.Config
				// WithTLS sets the TLS config for the catdog_service
				xerror.Exit(Default.Init(grpcC.AuthTLS(t)))
				xerror.Exit(catdog_server.Default.Init(grpcS.AuthTLS(t)))
			}))

			xerror.Exit(dix.WithAfterStart(func() {
				if catdog_config.Trace {
					fmt.Printf("%v\n", Default)
				}
			}))
		},
	}))
}
