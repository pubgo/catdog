package catdog_broker_plugin

import (
	"github.com/micro/go-micro/v3/broker"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_broker"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	opts broker.Options
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

func (p *Plugin) Init(cat catdog_abc.CatDog) error {
	return nil
}

func NewPlugin() *Plugin {
	return &Plugin{
		name: "broker",
		opts: catdog_broker.Default.Options(),
	}
}
