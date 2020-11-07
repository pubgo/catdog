package catdog_watcher

import (
	"context"
	"strings"

	"github.com/asim/nitro/v3/config"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
)

func Watch(name string, watcher func(r reader.Value) error) error {
	if name == "" {
		return xerror.Fmt("[name] should not be empty")
	}

	if watcher == nil {
		return xerror.Fmt("[watcher] should not be nil")
	}

	return xerror.Wrap(dix.WithBeforeStart(func() {
		key := strings.Join([]string{catdog_config.Project, name}, ".")
		resp := xerror.PanicErr(catdog_config.Load(key)).(reader.Value)
		if resp.Bytes() != nil {
			xerror.Panic(watcher(resp))
		}

		xlog.Debugf("Start Watch Config, Key: %s", key)
		w := xerror.PanicErr(catdog_config.GetCfg().Watch(key)).(config.Watcher)

		// 开启监听配置
		cancel := xprocess.Go(func(ctx context.Context) (err error) {
			defer xerror.RespErr(&err)
			defer func() {
				xlog.Debugf("Stop Watch Config, Key: %s", key)
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
		xerror.Exit(dix.WithAfterStart(func() { xerror.Exit(cancel()) }))
	}))
}
