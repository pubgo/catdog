package catdog_log

import (
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xlog/xlog_config"
	"go.uber.org/zap"

	"github.com/pubgo/catdog/catdog_config"
)

func GetDevLog() xlog.XLog {
	zl, err := xlog_config.NewZapLoggerFromConfig(xlog_config.NewDevConfig())
	if err != nil {
		xerror.Panic(err)
	}

	zl = zl.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).Named(catdog_config.Project)
	return xlog.New(zl)
}

func Watch(fn func(logs xlog.XLog)) error {
	fn(GetDevLog())
	return xerror.Wrap(dix.Dix(func(logs xlog.XLog) {
		fn(logs)
	}))
}
