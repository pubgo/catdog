package catdog_app

import (
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_client"
	"github.com/pubgo/catdog/plugins/catdog_server"
)

func Start(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)
	catdog_server.Default.Server = ent.Server()
	catdog_client.Default.Client = ent.Client()
	catdog_config.Project = ent.Options().Name

	// 启动配置, 初始化组件, 初始化插件
	plugins := catdog_plugin.List(catdog_plugin.Module(ent.Options().Name))
	for _, pg := range append(catdog_plugin.List(), plugins...) {
		key := pg.String()
		r, err := catdog_config.Load(key)
		xerror.PanicF(err, "plugin [%s] load error", key)
		xerror.PanicF(pg.Init(r), "plugin [%s] init error", key)

		hdlr := pg.Handler()
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
