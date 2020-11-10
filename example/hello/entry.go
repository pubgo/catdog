package hello

import (
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"

	"github.com/pubgo/catdog"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/example/hello/handler"
	"github.com/pubgo/catdog/plugins/catdog_log"
	"github.com/pubgo/catdog/version"
)

var log xlog.XLog

func init() {
	xerror.Exit(catdog_log.Watch(func(logs xlog.XLog) {
		log = logs.Named("hello")
	}))
}

func GetEntry() catdog.Entry {
	//catdog_config.Domain = "dev"
	catdog_config.Project = "hello"

	ent := catdog.NewEntry()
	xerror.Exit(ent.Name(catdog_config.Project, "hello 服务"))
	xerror.Exit(ent.Version(version.Version))

	xerror.Exit(ent.Handler(handler.NewHelloworld()))
	xerror.Exit(ent.Handler(handler.NewTestAPIHandler()))
	//xerror.Exit(ent.Handler(handler.NewTestAPIHandler1()))
	return ent
}
