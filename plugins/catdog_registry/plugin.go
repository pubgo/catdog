package catdog_registry

import (
	"github.com/asim/nitro/v3/broker"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_broker"
	"github.com/pubgo/catdog/plugins/catdog_client"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
)

func init() {
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "registry",
		OnInit: func(r reader.Value) {
			xerror.Exit(dix_run.WithBeforeStart(func(ctx *dix_run.BeforeStartCtx) {
				xerror.Exit(catdog_server.Default.Init(server.Registry(Default)))
				xerror.Exit(catdog_broker.Default.Init(broker.Registry(Default)))
				xerror.Exit(catdog_client.Default.Init(client.Registry(Default)))
			}))
		},
	}))
}
