package catdog_watcher

import (
	"fmt"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Panic(dix_run.WithAfterStart(func(ctx *dix_run.AfterStartCtx) {
		if !catdog_config.Trace {
			return
		}

		xlog.Debugf("watcher trace")
		fmt.Printf("%#v\n", catdog_config.GetCfg().Config)
		fmt.Println()
	}))
}
