package catdog_data

import (
	"fmt"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/internal/catdog_action"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Exit(catdog_action.WithAfterStart(func() {
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
