package catdog_app

import (
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_client"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
)

func Start(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)
	catdog_server.Default.Server = ent.Server()
	catdog_client.Default.Client = ent.Client()

	// 启动配置, 初始化组件
	entPlugins := catdog_plugin.List(catdog_plugin.Module(ent.Options().Name))
	for _, pl := range append(catdog_plugin.List(), entPlugins...) {
		key := pl.String()
		r, err := catdog_config.Load(key)
		xerror.PanicF(err, "plugin [%s] load error", key)
		xerror.PanicF(pl.Init(r), "plugin [%s] init error", key)

		hdlr := pl.Handler()
		if hdlr != nil {
			xerror.Panic(ent.Handler(hdlr, hdlr.Opts...))
		}
	}

	xerror.Panic(dix.BeforeStart())
	xerror.Panic(ent.Start())
	xerror.Panic(dix.AfterStart())

	return
}

func Stop(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(dix.BeforeStop())
	xerror.Panic(ent.Stop())
	xerror.Panic(dix.AfterStop())

	return nil
}
