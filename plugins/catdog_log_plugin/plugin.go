package catdog_log_plugin

import (
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xlog/xlog_config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name   string
	config xlog_config.Config
	log    xlog.XLog
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *Plugin) String() string {
	return p.name
}

func (p *Plugin) catDogWatcher(cat catdog_abc.CatDog) (rErr error) {
	defer xerror.RespErr(&rErr)

	zapL := xerror.PanicErr(xlog_config.NewZapLoggerFromConfig(p.config)).(*zap.Logger)
	log := xlog.New(zapL.WithOptions(xlog.AddCaller(), xlog.AddCallerSkip(1)))
	p.log = log.Named(catdog_config.Domain)

	xerror.Panic(xlog.SetLog(log.Named(catdog_config.Domain, xlog.AddCallerSkip(1))))
	return xerror.Wrap(dix.Dix(p.log))
}

func (p *Plugin) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(p.name, pflag.PanicOnError)
	flags.StringVar(&p.config.Level, "log_level", p.config.Level, "log level")
	return flags
}

func New() *Plugin {
	p := &Plugin{
		name:   "log",
		config: xlog_config.NewDevConfig(),
	}
	xerror.Exit(catdog_abc.Watch(p.catDogWatcher))

	return p
}
