package catdog_log

import (
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xlog/xlog_config"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_plugin"
)

func initLog(cfg xlog_config.Config) (err error) {
	defer xerror.RespErr(&err)

	zapL := xerror.PanicErr(xlog_config.NewZapLoggerFromConfig(cfg)).(*zap.Logger)
	log := xlog.New(zapL.WithOptions(xlog.AddCaller(), xlog.AddCallerSkip(1)))
	xerror.Panic(xlog.SetDefault(log.Named(catdog_config.Domain, xlog.AddCallerSkip(1))))
	xerror.Panic(dix.Dix(log.Named(catdog_config.Domain)))

	trace(cfg)
	return nil
}

func init() {
	var config = xlog_config.NewDevConfig()
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "log",
		OnFlags: func(flags *pflag.FlagSet) {
			flags.StringVar(&config.Level, "level", config.Level, "log level")
		},
		OnInit: func(r reader.Value) {
			xerror.Panic(r.Scan(&config))
			xerror.Panic(initLog(config))
		},
		OnWatch: func(r reader.Value) {
			xerror.Panic(r.Scan(&config))
			xerror.Panic(initLog(config))
		},
	}))
}
