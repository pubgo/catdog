package catdog_client

import (
	"crypto/tls"
	"fmt"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xlog"

	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_server"
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
		OnInit: func() {
			xerror.Exit(catdog_abc.WithBeforeStart(func() {
				var t *tls.Config
				// WithTLS sets the TLS config for the catdog_service
				xerror.Exit(Default.Init(grpcC.AuthTLS(t)))
				xerror.Exit(catdog_server.Default.Init(grpcS.AuthTLS(t)))
			}))

			xerror.Exit(catdog_abc.WithAfterStart(func() {
				if catdog_config.Trace {
					xlog.Debug("client trace")
					fmt.Printf("%v\n", Default)
				}
			}))
		},
	}))
}
