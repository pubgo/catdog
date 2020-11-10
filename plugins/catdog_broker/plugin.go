package catdog_broker

import (
	"fmt"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/catdog/plugins/catdog_client"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/spf13/pflag"
)

func init() {
	opts := Default.Options()
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "broker",
		OnFlags: func(flags *pflag.FlagSet) {
			_ = opts
		},
		OnInit: func() {
			xerror.Exit(catdog_abc.WithBeforeStart(func() {
				xerror.Exit(catdog_client.Default.Init(client.Broker(Default.Broker)))
				xerror.Exit(catdog_server.Default.Init(server.Broker(Default.Broker)))
			}))

			xerror.Exit(catdog_abc.WithAfterStart(func() {
				if !catdog_config.Trace {
					return
				}

				xlog.Debug("broker trace")
				fmt.Printf("%v\n", Default)
			}))
		},
	}))
}
