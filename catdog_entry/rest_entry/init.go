package rest_entry

import (
	"github.com/pubgo/catdog/catdog_log"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

const Name = "rest_entry"

var log xlog.XLog

func init() {
	xerror.Exit(catdog_log.Watch(func(logs xlog.XLog) {
		log = logs.Named(Name)
	}))
}
