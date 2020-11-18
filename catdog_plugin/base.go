package catdog_plugin

import (
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ Plugin = (*Base)(nil)

type Base struct {
	Name       string
	Enabled    bool `yaml:"enabled"`
	OnInit     func(ent catdog_entry.Entry)
	OnWatch    func(r reader.Value)
	OnCommands func(cmd *cobra.Command)
	OnHandler  func() *catdog_handler.Handler
	OnFlags    func(flags *pflag.FlagSet)
}

func (p *Base) Init(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)
	r, err := catdog_config.Load(p.Name)
	xerror.Panic(err)
	xerror.Panic(r.Scan(p))

	var status = "disabled"
	if p.Enabled {
		status = "enabled"
	}

	xlog.Debugf("plugin [%s] init, status: %s", p.Name, status)

	if !p.Enabled {
		return nil
	}

	if p.OnInit != nil {
		xlog.Debugf("[%s] start init", p.Name)
		p.OnInit(ent)
	}
	return nil
}

func (p *Base) Watch(r reader.Value) (err error) {
	defer xerror.RespErr(&err)

	if !p.Enabled {
		return nil
	}

	if p.OnWatch != nil {
		xlog.Debugf("[%s] start watch", p.Name)
		p.OnWatch(r)
	}
	return nil
}

func (p *Base) Commands() *cobra.Command {
	if !p.Enabled {
		return nil
	}

	if p.OnCommands != nil {
		cmd := &cobra.Command{Use: p.Name}
		p.OnCommands(cmd)
		return cmd
	}
	return nil
}

func (p *Base) Handler() *catdog_handler.Handler {
	if !p.Enabled {
		return nil
	}

	if p.OnHandler != nil {
		return p.OnHandler()
	}
	return nil
}

func (p *Base) String() string {
	return p.Name
}

func (p *Base) Flags() *pflag.FlagSet {
	if !p.Enabled {
		return nil
	}

	flags := pflag.NewFlagSet(p.Name, pflag.PanicOnError)
	if p.OnFlags != nil {
		p.OnFlags(flags)
	}
	return flags
}
