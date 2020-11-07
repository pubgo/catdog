package catdog_mongo

import (
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xlog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	log  xlog.XLog
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

func NewPlugin() *Plugin {
	return &Plugin{
		name: "mongo",
	}
}
