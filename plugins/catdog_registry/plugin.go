package catdog_registry

import (
	"github.com/asim/nitro/v3/broker"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/catdog/plugins/catdog_broker"
	"github.com/pubgo/catdog/plugins/catdog_client"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/xerror"
)

func init() {
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "registry",
		OnInit: func() {
			xerror.Exit(catdog_abc.WithBeforeStart(func() {
				xerror.Exit(catdog_server.Default.Init(server.Registry(Default)))
				xerror.Exit(catdog_broker.Default.Init(broker.Registry(Default)))
				xerror.Exit(catdog_client.Default.Init(client.Registry(Default)))
			}))
		},
	}))
}
