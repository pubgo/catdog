package catdog_config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/asim/nitro/v3/config"
	"github.com/asim/nitro/v3/config/memory"
	"github.com/asim/nitro/v3/config/source"
	mEnv "github.com/asim/nitro/v3/config/source/env"
	mFile "github.com/asim/nitro/v3/config/source/file"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"

	"github.com/pubgo/catdog/catdog_env"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/catdog/internal/plugins/config/encoder/yaml"
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
	catdog_env.Get(&Home, "home", "dir")
	catdog_env.Get(&Project, "project", "name")
	catdog_env.Get(&Mode, "mode", "run")

	CfgPath = filepath.Join(Home, "config", "config.yaml")
	if !catdog_util.PathExist(Home) {
		xerror.Exit(xerror.Fmt("home path %s not exists", Home))
	}
	if !catdog_util.PathExist(CfgPath) {
		xerror.Exit(xerror.Fmt("config path %s not exists", CfgPath))
	}

	cfg = &Config{Config: xerror.ExitErr(memory.NewConfig()).(config.Config)}
	xerror.Exit(cfg.Init(
		// 加载env source
		config.WithSource(mEnv.NewSource(mEnv.WithStrippedPrefix(Domain))),
		// 加载file source
		config.WithSource(mFile.NewSource(
			mFile.WithPath(filepath.Join(Home, "config", "config.yaml")),
			source.WithEncoder(yaml.NewEncoder()),
		)),
	))

	// debug and trace
	xerror.Exit(dix.WithAfterStart(func() {
		if !Trace {
			return
		}

		fmt.Println("config", string(LoadBytes()))
		fmt.Println("deps", dix.Graph())
	}))

	// 运行环境检查
	xerror.Exit(dix.WithBeforeStart(func() {
		var m = RunMode
		switch Mode {
		case m.Dev, m.Stag, m.Prod, m.Test, m.Release:
		default:
			xerror.Exit(xerror.Fmt("running mode does not match, mode: %s", Mode))
		}
	}))
}
