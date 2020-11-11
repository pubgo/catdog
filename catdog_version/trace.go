package catdog_version

import (
	"fmt"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xerror"
)

func init() {
	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !catdog_config.Trace {
			return
		}

		for name, v := range List() {
			fmt.Println(name, catdog_util.MarshalIndent(v))
		}
	}))
}
