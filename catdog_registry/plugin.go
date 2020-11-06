package catdog_registry

import (
	"github.com/asim/nitro/v3/broker"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/asim/nitro/v3/registry"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	registry.Options
}

func (p *Plugin) Flags() *pflag.FlagSet {
	return nil
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *Plugin) Watch(r reader.Value) error {
	return nil
}

func (p *Plugin) String() string {
	return p.name
}

func init() {
	xerror.Exit(catdog_plugin.Register(&Plugin{
		name:    "registry",
		Options: Default.Options(),
	}))
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
