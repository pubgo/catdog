package catdog_config

import (
	"context"
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/config/reader"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
	"os"
	"strings"
)

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

func PluginPrefix(names ...string) []string {
	return append([]string{Domain, Project, "plugins"}, names...)
}

func ProjectPrefix() []string {
	return []string{Domain, Project}
}

func WatchStart() error {
	return xerror.Wrap(dix.Start())
}

func WatchStop() error {
	return xerror.Wrap(dix.Stop())
}

func Watch(name string, watcher func(r reader.Value) error) {
	xerror.Exit(dix.Dix(func(ctx *dix.StartCtx) {
		if name == "" {
			return
		}
		if watcher == nil {
			return
		}

		wKey := strings.Join(PluginPrefix(name), "/")
		xlog.DebugF("Start Config Watch, Key: %s", wKey)
		w := xerror.PanicErr(cfg.Watch(wKey)).(config.Watcher)

		// 开启监听配置变化
		cancel := xprocess.Go(func(ctx context.Context) (err error) {
			defer xerror.RespErr(&err)
			defer func() {
				xlog.DebugF("Stop Config Watch, Key: %s", wKey)
				xerror.Panic(w.Stop())
			}()

			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					r, err := w.Next()
					if err != nil && strings.Contains(err.Error(), "stopped") {
						return nil
					}
					xerror.Panic(err)
					xerror.Panic(watcher(r))
				}
			}
		})

		// 关闭监听配置变化
		xerror.Exit(dix.Dix(func(stopCtx *dix.StopCtx) error {
			return xerror.Wrap(cancel())
		}))
	}))
}
