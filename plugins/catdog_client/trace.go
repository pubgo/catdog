package catdog_client

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

		xlog.Debug("client trace")
		fmt.Printf("%#v\n", Default.Client)
		fmt.Println()
	}))
}
