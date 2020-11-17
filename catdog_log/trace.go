package catdog_log

import (
	"fmt"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xlog/xlog_config"
)

func trace(cfg xlog_config.Config) {
	if !catdog_config.Trace {
		return
	}

	xlog.Debug("log trace")
	fmt.Println(catdog_util.MarshalIndent(cfg))
	fmt.Println()
}
