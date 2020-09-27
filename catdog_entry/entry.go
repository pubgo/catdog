package catdog_entry

import (
	"github.com/gofiber/fiber"
	"github.com/micro/go-micro/v3/server"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Entry interface {
	Name(name string, description ...string) error
	Version(v string) error
	Flags(fn func(flags *pflag.FlagSet)) error
	Commands(commands ...*cobra.Command) error
	Handler(register interface{}, hdlr interface{}, opts ...server.HandlerOption) error
	Init(opts ...catdog_abc.Option) Entry
	Plugins(pgs ...catdog_plugin.Plugin) error
	Group(relativePath string, handlers ...fiber.Handler) fiber.Router
	Task(name string, handler interface{}, opts ...server.SubscriberOption) error
	Options() Options
}

type Option func(o *Options)
type Options struct {
	Name       string
	Version    string
	RunCommand *cobra.Command
	Command    *cobra.Command
	Options    []catdog_abc.Option
}
