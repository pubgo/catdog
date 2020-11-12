package catdog_version

import (
	"fmt"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !catdog_config.Trace {
			return
		}

		for name, v := range List() {
			xlog.Debug(name)
			fmt.Println(catdog_util.MarshalIndent(v))
		}
		fmt.Println()
	}))
}
