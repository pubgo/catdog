package catdog_entry

import (
	"context"
	"fmt"
	"github.com/pubgo/catdog/catdog_data"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xprocess"
	"github.com/spf13/cobra"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"

	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_server"
	"github.com/pubgo/xerror"
	"github.com/spf13/pflag"
)

var _ Entry = (*rpcEntry)(nil)

type rpcEntry struct {
	opts     Options
	mux      sync.Mutex
	app      *fiber.App
	addr     string
	gwPrefix string
}

func (r *rpcEntry) Task(name string, handler interface{}, opts ...server.SubscriberOption) error {
	panic("implement me")
}

func (r *rpcEntry) Group(relativePath string, handlers ...fiber.Handler) fiber.Router {
	return r.app.Group(relativePath, handlers...)
}

func (r *rpcEntry) initFlags() error {
	return xerror.Wrap(r.Flags(func(flags *pflag.FlagSet) {
		flags.StringVar(&r.addr, "gw_addr", r.addr, "gateway address")
	}))
}

func (r *rpcEntry) stopService() error {
	if err := r.app.Shutdown(); err != nil && err != http.ErrServerClosed {
		return xerror.Wrap(err)
	}
	return nil
}

func (r *rpcEntry) pathRouterTrace() error {
	if catdog_config.Trace {
		for _, stacks := range r.app.Stack() {
			for _, stack := range stacks {
				log.DebugF("%s %s", stack.Method, stack.Path)
			}
		}
	}
	return nil
}

func (r *rpcEntry) startService() error {
	xprocess.Go(func(ctx context.Context) (err error) {
		defer xerror.RespErr(&err)
		log.InfoF("Server [http] Listening on http://%s", r.addr)
		xerror.Exit(r.app.Listen(r.addr))
		log.InfoF("Server [http] Closed OK")
		return nil
	})
	return nil
}

func (r *rpcEntry) middleware() []interface{} {
	return []interface{}{
		middleware.Recover(),
		middleware.Logger(middleware.LoggerConfig{
			Format:     "#${pid} - ${time} ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: time.RFC3339,
		}),
	}

}

func (r *rpcEntry) initCatDog(cat catdog_app.CatDog) (err error) {
	xerror.RespErr(&err)

	opts := r.Options()
	cat.Init(catdog_app.Name(opts.Name), catdog_app.Version(opts.Version))
	cat.Init(opts.Options...)
	cat.Init(
		catdog_app.BeforeStart(r.startService),
		catdog_app.BeforeStop(r.stopService),
		catdog_app.AfterStart(r.pathRouterTrace),
	)

	g := r.app.Group(r.gwPrefix)
	handlers := catdog_server.Default.Handlers()
	for i := range handlers {
		xerror.Panic(handlers[i](g))
	}

	return nil
}

func (r *rpcEntry) Options() Options {
	return r.opts
}

func (r *rpcEntry) Flags(fn func(flag *pflag.FlagSet)) (err error) {
	defer xerror.RespErr(&err)
	fn(r.opts.Command.PersistentFlags())
	return nil
}

func (r *rpcEntry) Name(name string, description ...string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return xerror.New("the name parameter should not be empty")
	}
	r.opts.Command.Use = name
	r.opts.Name = r.opts.Command.Name()
	r.opts.Command.Short = fmt.Sprintf("This is a %s service", r.opts.Name)

	if len(description) > 0 {
		r.opts.Command.Short = description[0]
	}

	return nil
}

func (r *rpcEntry) Version(v string) error {
	r.opts.Version = strings.TrimSpace(v)
	if r.opts.Version == "" {
		return xerror.New("version should not be null")
	}
	return nil
}

func (r *rpcEntry) Commands(commands ...*cobra.Command) error {
	for _, command := range commands {
		if command == nil {
			continue
		}
		r.opts.Command.AddCommand(command)
	}
	return nil
}

// func(s server.Server, handle TestHandler, opts ...server.HandlerOption) error
func (r *rpcEntry) Handler(hdlr interface{}, opts ...server.HandlerOption) error {
	return xerror.Wrap(catdog_handler.Register(hdlr,opts...))
}

func (r *rpcEntry) Plugins(pgs ...catdog_plugin.Plugin) (err error) {
	defer xerror.RespErr(&err)

	for _, pg := range pgs {
		xerror.PanicF(catdog_plugin.Register(pg,
			catdog_plugin.Module(r.opts.Name)), "Plugin [%s] Register error", pg.String())
		xerror.Panic(r.Flags(func(flag *pflag.FlagSet) { flag.AddFlagSet(pg.Flags()) }))
		xerror.Panic(r.Commands(pg.Commands()))
	}

	return
}

func newEntry() *rpcEntry {
	ent := &rpcEntry{
		opts: Options{
			RunCommand: &cobra.Command{Use: "run"},
			Command:    &cobra.Command{},
		},
		app:      fiber.New(),
		addr:     ":8080",
		gwPrefix: "api",
	}
	ent.opts.Command.AddCommand(ent.opts.RunCommand)
	ent.app.Use(ent.middleware()...)
	xerror.Exit(ent.initFlags())
	xerror.Exit(catdog_app.Watch(ent.initCatDog))
	return ent
}

func New() Entry {
	return newEntry()
}

func unWrapType(tye reflect.Type) reflect.Type {
	for isElem(tye) {
		tye = tye.Elem()
	}
	return tye
}

func isElem(tye reflect.Type) bool {
	switch tye.Kind() {
	case reflect.Chan, reflect.Map, reflect.Ptr, reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}
