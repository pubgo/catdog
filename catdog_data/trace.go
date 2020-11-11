package catdog_data

import (
	"fmt"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xerror"
)

func init() {
	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !catdog_config.Trace {
			return
		}

		for k := range List() {
			fmt.Println(k)
		}
	}))
}
