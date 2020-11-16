package catdog_app

import (
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
)

func Start(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(ent.Init())

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

	for _, fn := range catdog_abc.GetBeforeStart() {
		fn()
	}
	xerror.Panic(ent.Start())
	for _, fn := range catdog_abc.GetAfterStart() {
		fn()
	}

	return
}

func Stop(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	for _, fn := range catdog_abc.GetBeforeStop() {
		fn()
	}
	xerror.Panic(ent.Stop())
	for _, fn := range catdog_abc.GetAfterStop() {
		fn()
	}

	return nil
}
