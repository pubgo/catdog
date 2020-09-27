package catdog_debug_plugin

import (
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_debug_plugin/handler"
	debug "github.com/pubgo/catdog/plugins/catdog_debug_plugin/proto"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
}

func (p *Plugin) Flags() *pflag.FlagSet {
	return nil
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) String() string {
	return p.name
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return catdog_handler.Register(
		debug.RegisterDebugHandler,
		handler.NewHandler(),
		//api.WithEndpoint(&api.Endpoint{
		//	name:    "TestApi.Version",
		//	Path:    []string{"/v1/example/version"},
		//	Method:  []string{"POST"},
		//	Body:    "*",
		//	handler: "rpc",
		//}),
	)
}

func (p *Plugin) Init(cat catdog_abc.CatDog) error {
	return xerror.Wrap(dix.Dix(p))
}

func NewPlugin() *Plugin {
	return &Plugin{
		name: "catdog_debug_plugin",
	}
}
