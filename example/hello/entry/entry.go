package entry

import (
	"github.com/pubgo/catdog"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_log"
	"github.com/pubgo/catdog/example/hello/handler"
	helloworld "github.com/pubgo/catdog/example/hello/proto"
	"github.com/pubgo/catdog/version"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

var log = catdog_log.GetDevLog().Named("hello")

func init() {
	xerror.Exit(dix.Dix(func(plugin xlog.XLog) {
		log = plugin.Named("hello")
	}))
}

func GetEntry() catdog.Entry {
	ent := catdog.NewEntry()
	xerror.Exit(ent.Name("hello", "hello 服务"))
	xerror.Exit(ent.Version(version.Version))
	ent.Init(catdog_abc.BeforeStart(func() error {
		log.Info("init Hello")
		return nil
	}))

	xerror.Exit(ent.Handler(helloworld.RegisterHelloworldHandler, handler.NewHelloworld()))
	xerror.Exit(ent.Handler(helloworld.RegisterTestApiHandler, handler.NewTestAPIHandler()))
	xerror.Exit(ent.Handler(helloworld.RegisterTestApiV2Handler, handler.NewTestAPIHandler()))
	return ent
}
