package handler

import (
	"github.com/pubgo/catdog/catdog_log"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

var log xlog.XLog

func init() {
	xerror.Exit(catdog_log.Watch("hello.handler", &log))
}
