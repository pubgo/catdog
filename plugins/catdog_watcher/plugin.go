package catdog_watcher

import (
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/catdog_plugin"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	catdog_plugin.Plugin
}

func init() {
	p := &Plugin{Plugin: catdog_plugin.NewBase("watcher")}
	xerror.Exit(catdog_plugin.Register(p))
}
