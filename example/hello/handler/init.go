package handler

import (
	"github.com/pubgo/catdog/catdog_log"
	"github.com/pubgo/dix"
	"github.com/pubgo/xlog"
)

var log = catdog_log.GetDevLog().Named("hello")

func init() {
	dix.Go(func(log xlog.XLog) {
		log = log.Named("hello")
	})
}
