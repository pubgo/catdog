package catdog_config

import (
	"fmt"
	"github.com/pubgo/dix"
	"strings"

	"github.com/asim/nitro/v3/config"
	"github.com/asim/nitro/v3/config/memory"
	mEnv "github.com/asim/nitro/v3/config/source/env"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"

	"github.com/pubgo/catdog/catdog_env"
)

func init() {
	// 从环境变量中获取系统默认值
	// 获取系统默认的前缀, 环境变量前缀等
	catdog_env.Get(&Domain, "catdog", "catdog_domain", "catdog_prefix", "env_prefix")
	if Domain = strings.TrimSpace(strings.ToLower(Domain)); Domain == "" {
		Domain = "catdog"
		xlog.Warnf("[domain] prefix should be set, default: %s", Domain)
	}

	// 设置系统环境变量前缀
	catdog_env.Prefix(Domain)

	// 使用前缀获取系统环境变量
	catdog_env.Get(&Home, "cfg_dir", "config_dir", "home_dir", "home", "project_dir", "dir")
	catdog_env.Get(&Project, "project_name", "service_name", "server_name", "project", "name")

	// 加载env source
	cfg = &Config{Config: xerror.ExitErr(memory.NewConfig()).(config.Config)}
	xerror.Exit(cfg.Init(config.WithSource(mEnv.NewSource(mEnv.WithStrippedPrefix(Domain)))))
}

func init() {
	xerror.Exit(dix.WithBeforeStart(func() {
		if Trace {
			fmt.Println("config", string(LoadBytes()))
			fmt.Println("deps", dix.Graph())
		}
	}))
}
