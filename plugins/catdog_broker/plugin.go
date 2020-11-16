package catdog_broker

import (
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_client"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
	"github.com/spf13/pflag"
)

func init() {
	opts := Default.Options()
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "broker",
		OnFlags: func(flags *pflag.FlagSet) {
			_ = opts
		},
		OnInit: func(r reader.Value) {
			xerror.Exit(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
				xerror.Exit(catdog_client.Default.Init(client.Broker(Default.Broker)))
				xerror.Exit(catdog_server.Default.Init(server.Broker(Default.Broker)))
			}))
		},
	}))
}

//xerror.Exit(r.s.Init(server.Broker(catdog_broker.Default.Broker)))
