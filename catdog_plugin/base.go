package catdog_plugin

import (
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ Plugin = (*Base)(nil)

type Base struct {
	Name       string
	OnInit     func()
	OnWatch    func(r reader.Value) error
	OnCommands func(cmd *cobra.Command)
	OnHandler  func() *catdog_handler.Handler
	OnFlags    func(flags *pflag.FlagSet)
}

func (p *Base) Init() {
	if p.OnInit != nil {
		p.OnInit()
	}
}

func (p *Base) Watch(r reader.Value) error {
	if p.OnWatch != nil {
		return p.OnWatch(r)
	}
	return nil
}

func (p *Base) Commands() *cobra.Command {
	if p.OnCommands != nil {
		cmd := &cobra.Command{Use: p.Name}
		p.OnCommands(cmd)
		return cmd
	}
	return nil
}

func (p *Base) Handler() *catdog_handler.Handler {
	if p.OnHandler != nil {
		return p.OnHandler()
	}
	return nil
}

func (p *Base) String() string {
	return p.Name
}

func (p *Base) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(p.Name, pflag.PanicOnError)
	if p.OnFlags != nil {
		p.OnFlags(flags)
	}
	return flags
}
