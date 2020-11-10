package catdog_rabbitmq

import (
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xerror"
)

func init() {
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "rabbit",
	}))
}
