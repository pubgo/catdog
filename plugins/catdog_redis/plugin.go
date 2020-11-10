package catdog_redis

import (
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/catdog_plugin"
)

func init() {
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "redis",
	}))
}
