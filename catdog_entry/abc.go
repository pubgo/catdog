package catdog_entry

import (
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type RpcEntry interface {
	Entry
}

type Entry interface {
	Start() error
	Stop() error
	Server() server.Server
	Client() client.Client
	Name(name string, description ...string) error
	Version(v string) error
	Flags(fn func(flags *pflag.FlagSet)) error
	Commands(commands ...*cobra.Command) error
	Options() Options
	Handler(handler interface{}, opts ...server.HandlerOption) error
}

type Option func(o *Options)
type Options struct {
	Name        string
	Description []string
	Version     string
	RunCommand  *cobra.Command
	Command     *cobra.Command
}
