package catdog_plugin

import (
	"fmt"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !catdog_config.Trace {
			return
		}

		xlog.Debug("plugin trace")
		fmt.Println(String())
		fmt.Println()
	}))
}
