package catdog_plugin

import (
	"fmt"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Exit(dix_run.WithAfterStart(func(ctx *dix_run.AfterStartCtx) {
		if !catdog_config.Trace {
			return
		}

		xlog.Debug("plugin trace")
		fmt.Println(String())
		fmt.Println()
	}))
}
