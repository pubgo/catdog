package catdog_config

import (
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/config/reader"
	"github.com/pubgo/dix"
	"github.com/pubgo/xlog"
	"os"
	"strings"
	"time"

	"github.com/pubgo/xerror"
)

// RunMode 项目运行模式
var RunMode = struct {
	Dev     string
	Test    string
	Stag    string
	Prod    string
	Release string
}{
	Dev:     "dev",
	Test:    "test",
	Stag:    "stag",
	Prod:    "prod",
	Release: "release",
}

func Env(env *string, names ...string) {
	getEnv(env, names...)
}

func SysEnv(env *string, names ...string) {
	getSysEnv(env, names...)
}

func getSysEnv(val *string, names ...string) {
	for _, name := range names {
		env, ok := os.LookupEnv(strings.ToUpper(name))
		env = strings.TrimSpace(env)
		if ok && env != "" {
			*val = env
		}
	}
}

func getEnv(val *string, names ...string) {
	for _, name := range names {

		if Domain == "" {
			name = strings.ToUpper(name)
		} else {
			name = strings.ToUpper(strings.Join([]string{Domain, name}, "_"))
		}

		env, ok := os.LookupEnv(strings.ToUpper(name))
		env = strings.TrimSpace(env)
		if ok && env != "" {
			*val = env
		}
	}
}

func CheckRunMode() error {
	// 运行环境检查
	switch Mode {
	case RunMode.Dev, RunMode.Stag, RunMode.Prod, RunMode.Test, RunMode.Release:
	default:
		return xerror.Fmt("running mode does not match, mode: %s", Mode)
	}
	return nil
}

func PluginWrap(names ...string) []string {
	return append([]string{Domain, Project, "plugins"}, names...)
}

type watchCtx struct {
	time.Time
}

func WatchStart() error {
	return xerror.Wrap(dix.Dix(&watchCtx{time.Now()}))
}

func Watch(name string, watcher func(r reader.Value) error) {
	xerror.Exit(dix.Dix(func(ctx *watchCtx) {
		nameWrap := strings.Join(PluginWrap(name), "/")

		if watcher == nil {
			return
		}

		cfg := GetCfg()
		xlog.DebugF("Start Config Watch, Key: %s", nameWrap)
		w := xerror.PanicErr(cfg.Watch(nameWrap)).(config.Watcher)

		go func() {
			defer xerror.RespGoroutine("plugin_" + name)
			defer func() {
				xlog.DebugF("Stop Config Watch, Key: %s", nameWrap)
				xerror.Panic(w.Stop())
			}()

			for {
				r, err := w.Next()
				if err != nil && strings.Contains(err.Error(), "stopped") {
					return
				}
				xerror.Panic(err)
				xerror.Panic(watcher(r))
			}
		}()
	}))
}
