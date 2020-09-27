package catdog_log

import (
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xlog/xlog_config"
	"go.uber.org/zap"
)

func GetDevLog() xlog.XLog {
	zl, err := xlog_config.NewZapLoggerFromConfig(xlog_config.NewDevConfig())
	if err != nil {
		xerror.Exit(err)
	}

	zl = zl.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Named("debug")
	return xlog.New(zl)
}

func Watch(name string, log *xlog.XLog) error {
	*log = GetDevLog().Named(name)
	return xerror.Wrap(dix.Dix(func(logs xlog.XLog) {
		*log = logs.Named(name)
	}))
}
