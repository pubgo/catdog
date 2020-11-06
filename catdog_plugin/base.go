package catdog_plugin

import (
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ Plugin = (*basePlugin)(nil)

type basePlugin struct {
	name string
}

func (p *basePlugin) Watch(r reader.Value) error {
	return nil
}

func (p *basePlugin) Commands() *cobra.Command {
	return nil
}

func (p *basePlugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *basePlugin) String() string {
	return p.name
}

func (p *basePlugin) Flags() *pflag.FlagSet {
	return nil
}

func NewBase(name string) Plugin {
	return &basePlugin{name: name}
}
