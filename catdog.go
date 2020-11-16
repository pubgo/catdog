package catdog

import (
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_entry/rest_entry"
	"github.com/pubgo/catdog/catdog_entry/rpc_entry"
	"github.com/pubgo/xerror"
)

func Run(entries ...Entry) (err error) {
	defer xerror.RespErr(&err)
	xerror.Panic(catdog_app.Run(entries...))
	return nil
}

func Init() (err error) {
	defer xerror.RespErr(&err)

	// 初始化配置文件
	xerror.Panic(catdog_config.Init())
	return nil
}

type Entry = catdog_entry.Entry

func NewRpcEntry(name string) catdog_entry.Entry {
	return rpc_entry.New(name)
}

func NewRestEntry(name string) catdog_entry.Entry {
	return rest_entry.New(name)
}
