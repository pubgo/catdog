package catdog_broker

import (
	"github.com/asim/nitro/v3/broker"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_client"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	client.Options
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *Plugin) String() string {
	return p.name
}

func (p *Plugin) Flags() *pflag.FlagSet {
	return nil
}

func (p *Plugin) Init(cat catdog_app.CatDog) error {

	//Default.initCatDog()

	return xerror.Wrap(dix.Dix(p))
}

func New() *Plugin {
	return &Plugin{
		name:    "client",
		Options: catdog_client.Default.Options(),
	}
}

func Broker(b broker.Broker) Option {
	return func(o *Options) {
		xerror.Exit(o.Client.Init(client.Broker(b)))
		xerror.Exit(o.Server.Init(server.Broker(b)))
	}
}
