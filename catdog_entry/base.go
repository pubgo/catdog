package catdog_entry

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber"
	"github.com/micro/go-micro/v3/server"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_server"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ Entry = (*baseEntry)(nil)

type baseEntry struct {
	opts Options
}

func (b *baseEntry) Init(opts ...catdog_abc.Option) Entry {
	b.opts.Options = append(b.opts.Options, opts...)
	return b
}

func (b *baseEntry) Options() Options {
	return b.opts
}

func (b *baseEntry) Flags(fn func(flag *pflag.FlagSet)) (err error) {
	defer xerror.RespErr(&err)
	fn(b.opts.Command.PersistentFlags())
	return nil
}

func (b *baseEntry) Name(name string, description ...string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return xerror.New("the name parameter should not be empty")
	}
	b.opts.Command.Use = name
	b.opts.Name = b.opts.Command.Name()
	b.opts.Command.Short = fmt.Sprintf("This is a %s service", b.opts.Name)

	if len(description) > 0 {
		b.opts.Command.Short = description[0]
	}

	return nil
}

func (b *baseEntry) Version(v string) error {
	b.opts.Version = strings.TrimSpace(v)
	if b.opts.Version == "" {
		return xerror.New("version should not be null")
	}
	return nil
}

func (b *baseEntry) Commands(commands ...*cobra.Command) error {
	for _, command := range commands {
		if command == nil {
			continue
		}
		b.opts.Command.AddCommand(command)
	}
	return nil
}

// func(s server.Server, handle TestHandler, opts ...server.HandlerOption) error
func (b *baseEntry) Handler(register interface{}, hdlr interface{}, opts ...server.HandlerOption) error {
	return xerror.Wrap(catdog_server.RegHandler(register, hdlr, opts...))
}

func (b *baseEntry) Plugins(pgs ...catdog_plugin.Plugin) (err error) {
	defer xerror.RespErr(&err)

	for _, pg := range pgs {
		xerror.PanicF(catdog_plugin.Register(pg,
			catdog_plugin.Module(b.opts.Name)), "Plugin [%s] Register error", pg.String())
		xerror.Panic(b.Flags(func(flag *pflag.FlagSet) { flag.AddFlagSet(pg.Flags()) }))
		xerror.Panic(b.Commands(pg.Commands()))
	}

	return
}

func (b *baseEntry) Group(relativePath string, handlers ...fiber.Handler) fiber.Router {
	panic("implement me")
}

func (b *baseEntry) Task(name string, handler interface{}, opts ...server.SubscriberOption) error {
	panic("implement me")
}

func NewBase() Entry {
	base := &baseEntry{
		opts: Options{
			RunCommand: &cobra.Command{Use: "run"},
			Command:    &cobra.Command{},
		},
	}
	base.opts.Command.AddCommand(base.opts.RunCommand)
	return base
}
