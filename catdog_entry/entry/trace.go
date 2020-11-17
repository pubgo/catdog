package entry

import (
	"fmt"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func (t *BaseEntry) trace() {
	xerror.Panic(dix_run.WithAfterStart(func(ctx *dix_run.AfterStartCtx) {
		if !catdog_config.Trace || !t.opts.Initialized {
			return
		}

		xlog.Debug("BaseEntry rest trace")
		for _, stacks := range t.app.Stack() {
			for _, stack := range stacks {
				if stack.Path == "/" {
					continue
				}

				log.Debugf("%s %s", stack.Method, stack.Path)
			}
		}
		fmt.Println()

		xlog.Debugf("http server config trace")
		fmt.Println(catdog_util.MarshalIndent(t.app.Config()))
		fmt.Println()
	}))
}
