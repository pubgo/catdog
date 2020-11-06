package catdog_app

import (
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_server"
)

func Start(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	// 启动配置, 初始化组件
	nameModule := catdog_plugin.Module(ent.Options().Name)
	plugins := append(catdog_plugin.List(), catdog_plugin.List(nameModule)...)
	for _, pl := range plugins {
		hdlr := pl.Handler()
		if hdlr == nil {
			continue
		}

		xerror.Panic(ent.Handler(hdlr, hdlr.Opts...))
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
