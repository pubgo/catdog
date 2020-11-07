package catdog_recovery

import (
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/catdog_plugin"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	catdog_plugin.Plugin
}

func init() {
	p := &Plugin{Plugin: catdog_plugin.NewBase("recovery")}
	xerror.Exit(p.clientWrap())
	xerror.Exit(p.handlerWrap())
	xerror.Exit(catdog_plugin.Register(p))
}
