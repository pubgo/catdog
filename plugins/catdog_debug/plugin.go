package catdog_debug

import (
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xerror"
)

func init() {
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "catdog_debug",
		OnHandler: func() *catdog_handler.Handler {
			return catdog_handler.New(NewHandler())
		},
	}))
}
