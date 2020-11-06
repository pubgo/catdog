package catdog_config

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/asim/nitro/v3/config"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
	"github.com/spf13/pflag"
)

func DefaultFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("app", pflag.PanicOnError)
	flags.StringVarP(&Mode, "mode", "m", Mode, "running mode(dev|test|stag|prod|release)")
	flags.StringVarP(&Home, "home", "c", Home, "project config home dir")
	flags.BoolVarP(&Debug, "debug", "d", Debug, "enable debug")
	flags.BoolVarP(&Trace, "trace", "t", Trace, "enable trace")
	flags.StringVarP(&Project, "project", "p", Project, "project name")
	return flags
}

type Config struct {
	config.Config
}

// 默认的全局配置
var (
	Domain  = "catdog"
	Project = "catdog"
	Debug   = true
	Trace   = false
	Mode    = "dev"
	Home    = xerror.PanicStr(filepath.Abs(filepath.Dir("")))
	cfg     *Config
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

func Init(opts ...config.Option) {
	xerror.Exit(cfg.Init(opts...))
}

func Load(path ...string) (reader.Value, error) {
	return cfg.Load(path...)
}

func LoadBytes() []byte {
	return xerror.PanicErr(cfg.Load()).(reader.Value).Bytes()
}

func Watch(name string, watcher func(r reader.Value) error) error {
	if name == "" {
		return xerror.Fmt("[name] should not be empty")
	}

	if watcher == nil {
		return xerror.Fmt("[watcher] should not be nil")
	}

	return xerror.Wrap(dix.WithBeforeStart(func() {
		key := strings.Join([]string{Project, name}, ".")
		resp := xerror.PanicErr(cfg.Load(key)).(reader.Value)
		if resp.Bytes() != nil {
			xerror.Panic(watcher(resp))
		}

		xlog.Debugf("Start Watch Config, Key: %s", key)
		w := xerror.PanicErr(cfg.Watch(key)).(config.Watcher)

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
