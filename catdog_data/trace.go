package catdog_data

import (
	"fmt"
	"github.com/pubgo/dix/dix_run"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Exit(dix_run.WithAfterStart(func(ctx *dix_run.AfterStartCtx) {
		if !catdog_config.Trace {
			return
		}

		xlog.Debug("data trace")
		for k, v := range List() {
			fmt.Printf("%s: %#v\n", k, v)
		}
		fmt.Println()
	}))
}
