package catdog_client

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

		xlog.Debug("client trace")
		fmt.Printf("%#v\n", Default.Client)
		fmt.Println()
	}))
}
