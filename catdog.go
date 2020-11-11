package catdog

import (
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_entry/rpc_entry"
	"github.com/pubgo/catdog/plugins/catdog_pidfile"
	"github.com/pubgo/xerror"
)

func Run(entries ...Entry) (err error) {
	defer xerror.RespErr(&err)
	xerror.Panic(catdog_app.Run(entries...))
	return nil
}

func Init() (err error) {
	defer xerror.RespErr(&err)
	catdog_pidfile.Init()

	// 初始化配置文件
	xerror.Panic(catdog_config.Init())
	return nil
}

type Entry = catdog_entry.Entry

func NewEntry() catdog_entry.Entry {
	return rpc_entry.New()
}
