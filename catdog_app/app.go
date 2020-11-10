package catdog_app

import (
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_server"
)

func Start(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	// 启动配置, 初始化组件
	entPlugins := catdog_plugin.List(catdog_plugin.Module(ent.Options().Name))
	for _, pl := range append(catdog_plugin.List(), entPlugins...) {
		xerror.Panic(pl.Init())
		hdlr := pl.Handler()
		if hdlr != nil {
			xerror.Panic(ent.Handler(hdlr, hdlr.Opts...))
		}
	}

	xerror.Panic(dix.BeforeStart())
	xerror.Panic(catdog_server.Default.Start())
	xerror.Panic(dix.AfterStart())

	return
}

func Stop() (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(dix.BeforeStop())
	xerror.Panic(catdog_server.Default.Stop())
	xerror.Panic(dix.AfterStop())

	return nil
}
