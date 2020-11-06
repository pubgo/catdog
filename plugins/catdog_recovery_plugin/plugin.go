package catdog_recovery_plugin

import (
	"github.com/pubgo/catdog/catdog_app"
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

func (p *Plugin) catdogWatcher(cat catdog_app.CatDog) (rErr error) {
	defer xerror.RespErr(&rErr)

	cat.Init(p.clientWrap(), p.handlerWrap())

	return xerror.Wrap(dix.Dix(p))
}

func (p *Plugin) Flags() *pflag.FlagSet {
	return nil
}

func New() *Plugin {
	p := &Plugin{
		name: "recovery",
	}

	catdog_app.Watch(p.catdogWatcher)

	return p

}
