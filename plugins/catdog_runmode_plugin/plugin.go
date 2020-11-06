package catdog_runmode_plugin

import (
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
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

func (p *Plugin) Flags() *pflag.FlagSet {
	return nil
}

func (p *Plugin) catDogWatcher(cat catdog_app.CatDog) error {
	cat.Init(catdog_app.BeforeStart(func() error {
		return xerror.Wrap(checkRunMode())
	}))
	return xerror.Wrap(dix.Dix(p))
}

func New() *Plugin {
	p := &Plugin{
		name: "runmode",
	}
	xerror.Exit(catdog_app.Watch(p.catDogWatcher))
	return p
}

// checkRunMode 运行环境检查
func checkRunMode() error {
	var runMode = catdog_config.RunMode

	switch catdog_config.Mode {
	case runMode.Dev, runMode.Stag, runMode.Prod, runMode.Test, runMode.Release:
	default:
		return xerror.Fmt("running mode does not match, mode: %s", catdog_config.Mode)
	}

	return nil
}
