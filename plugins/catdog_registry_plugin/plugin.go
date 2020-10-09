package catdog_registry_plugin

import (
	"github.com/micro/go-micro/v3/registry"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_registry"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	registry.Options
}

func (p *Plugin) Flags() *pflag.FlagSet {
	return nil
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

func (p *Plugin) Init(cat catdog_abc.CatDog) error {

	//Default.initCatDog()

	return xerror.Wrap(dix.Dix(p))
}

func New() *Plugin {
	return &Plugin{
		name:    "registry",
		Options: catdog_registry.Default.Options(),
	}
}
