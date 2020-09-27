package catdog_config

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v3/config"
	mEnv "github.com/micro/go-micro/v3/config/source/env"
	"github.com/pubgo/xerror"
)

// 默认的全局配置
var (
	Domain  = "catdog"
	Project = "catdog"
	Debug   = true
	Trace   = false
	Mode    = "dev"
	CfgDir  = xerror.PanicStr(filepath.Abs(filepath.Dir("")))
	CfgPath = filepath.Join(CfgDir, "config", "config.yaml")
	cfg     config.Config
)

func init() {
	// 从环境变量中获取系统默认值
	getSysEnv(&Domain, "catdog", "domain", "domain_name", "env_prefix", "catdog_domain")
	Domain = strings.ToLower(strings.TrimSpace(Domain))
	if Domain == "" {
		xerror.Exit(errors.New("domain should not be empty"))
	}

	getSysEnv(&CfgDir, "cfg_dir", "config_dir", "home_dir", "catdog_home", "project_dir")
	getSysEnv(&CfgPath, "cfg_path", "config_path", "catdog_config", "catdog_cfg")
	getSysEnv(&Project, "project_name", "service_name", "server_name")

	// 加载env source
	cfg = xerror.ExitErr(config.NewConfig()).(config.Config)
	xerror.Exit(cfg.Load(
		mEnv.NewSource(mEnv.WithStrippedPrefix(Domain)),
	))
}

func GetCfg() config.Config {
	return cfg
}
