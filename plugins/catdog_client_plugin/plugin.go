package catdog_client_plugin

import (
	"github.com/asim/nitro/v3/client"
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_client"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
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
		name: "client",
		Options: catdog_client.Default.Options(),
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Client = c
	}
}

