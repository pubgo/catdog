package catdog_watcher

import (
	"github.com/asim/nitro-plugins/config/source/etcd/v3"
	"github.com/asim/nitro/v3/config"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	// 检查是否有watcher
	xerror.Exit(dix.WithBeforeStart(func() {
		//	 获取
		cfg := catdog_config.GetCfg()
		r, err := cfg.Load("watcher")
		if err != nil {
			xlog.Debugf("config [watcher] is error: %v", err)
			return
		}

		cfgMap := r.StringMap(nil)
		if cfgMap["status"] == catdog_config.PluginStop {
			return
		}

		uri := cfgMap[""]

		xerror.Panic(cfg.Init(
			config.WithSource(
				etcd.NewSource(
					// optionally specify etcd address; default to localhost:8500
					etcd.WithAddress("10.0.0.10:8500"),
					// optionally specify prefix; defaults to /micro/config
					etcd.WithPrefix("/my/prefix"),
					// optionally strip the provided prefix from the keys, defaults to false
					etcd.StripPrefix(true),
				),
			)))
	}))
}
