package catdog_server_plugin

import (
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_server"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	server.Options
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
	flags := pflag.NewFlagSet(p.name, pflag.PanicOnError)
	flags.StringVar(&p.Address, "server_addr", p.Address, "server address")
	flags.StringVar(&p.Name, "server_name", p.Name, "server name")
	return flags
}

func (p *Plugin) catDogWatcher(cat catdog_app.CatDog) (err error) {
	defer xerror.RespErr(&err)
	return xerror.Wrap(dix.Dix(p))
}

func New() *Plugin {
	p := &Plugin{
		Options: catdog_server.Default.Options(),
	}
	xerror.Exit(catdog_app.Watch(p.catDogWatcher))

	return p
}

// Server sets the catdog_server_plugin for handling requests
func Server(s server.Server) Option {
	return func(o *Options) {
		o.Server = s
	}
}
