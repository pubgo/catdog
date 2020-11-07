package hello

import (
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"

	"github.com/pubgo/catdog"
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/example/hello/handler"
	helloworld "github.com/pubgo/catdog/example/hello/proto"
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
	catdog_config.Domain = "dev"
	catdog_config.Project = "hello"

	ent := catdog.NewEntry()
	xerror.Exit(ent.Name(catdog_config.Project, "hello 服务"))
	xerror.Exit(ent.Version(version.Version))
	ent.Init(catdog_app.BeforeStart(func() error {
		log.Info("init Hello")
		return nil
	}))

	xerror.Exit(ent.Handler(helloworld.RegisterHelloworldHandler, handler.NewHelloworld()))
	xerror.Exit(ent.Handler(helloworld.RegisterTestApiHandler, handler.NewTestAPIHandler()))
	xerror.Exit(ent.Handler(helloworld.RegisterTestApiV2Handler, handler.NewTestAPIHandler()))
	return ent
}
