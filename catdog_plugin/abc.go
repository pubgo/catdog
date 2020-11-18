package catdog_plugin

import (
	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Manager is the internal_plugin manager which stores plugins and allows them to be retrieved.
// This is used by all the components of micro.
type Manager interface {
	Plugins(...ManagerOption) []Plugin
	Register(Plugin, ...ManagerOption) error
}

type ManagerOption func(o *ManagerOptions)
type ManagerOptions struct {
	Module string
}

// Plugin is the interface for plugins to micro. It differs from go-micro in that it's for
// the micro API, Web, Sidecar, CLI. It's a method of building middleware for the HTTP side.
type Plugin interface {
	Watch(r reader.Value) error
	Init(ent catdog_entry.Entry) error
	Flags() *pflag.FlagSet
	Commands() *cobra.Command
	String() string
}

type Option func(o *Options)
type Options struct {
	Name     string
	Flags    *pflag.FlagSet
	Commands *cobra.Command
}
