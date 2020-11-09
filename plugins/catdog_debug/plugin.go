package catdog_debug

import (
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xerror"
)

type Plugin struct {
	catdog_plugin.Plugin
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return catdog_handler.New(NewHandler())
}

func init() {
	p := &Plugin{Plugin: catdog_plugin.NewBase("catdog_debug")}
	xerror.Exit(catdog_plugin.Register(p))
}
