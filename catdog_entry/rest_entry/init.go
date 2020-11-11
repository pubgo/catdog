package catdog_entry

import (
	"github.com/pubgo/catdog/plugins/catdog_log"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

const Name = "rest"

var log xlog.XLog

func init() {
	xerror.Exit(catdog_log.Watch(func(logs xlog.XLog) {
		log = logs.Named(Name)
	}))
}
