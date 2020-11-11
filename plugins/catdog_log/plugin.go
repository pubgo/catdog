package catdog_log

import (
	"fmt"
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xlog/xlog_config"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

func init() {
	var config = xlog_config.NewDevConfig()
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "log",
		OnFlags: func(flags *pflag.FlagSet) {
			flags.StringVar(&config.Level, "log_level", config.Level, "log level")
		},
		OnInit: func() {
			xerror.Exit(catdog_abc.WithAfterStart(func() {
				if !catdog_config.Trace {
					return
				}

				xlog.Debug("log trace")
				fmt.Printf("%v\n", config)
			}))
		},
		OnWatch: func(r reader.Value) {
			xerror.Panic(r.Scan(&config))

			zapL := xerror.PanicErr(xlog_config.NewZapLoggerFromConfig(config)).(*zap.Logger)
			log := xlog.New(zapL.WithOptions(xlog.AddCaller(), xlog.AddCallerSkip(1)))
			xerror.Panic(xlog.SetDefault(log.Named(catdog_config.Domain, xlog.AddCallerSkip(1))))
			xerror.Exit(dix.Dix(log.Named(catdog_config.Domain)))
		},
	}))
}
