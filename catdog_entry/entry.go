package catdog_entry

import (
	"context"
	"fmt"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	ver "github.com/hashicorp/go-version"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xprocess"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/catdog_server"
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

func (r *rpcEntry) middleware() []interface{} {
	return []interface{}{
		middleware.Recover(),
		middleware.Logger(middleware.LoggerConfig{
			Format:     "#${pid} - ${time} ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: time.RFC3339,
		}),
	}

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
	r.opts.Description = description

	return nil
}

func (r *rpcEntry) Version(v string) error {
	r.opts.Version = strings.TrimSpace(v)
	if r.opts.Version == "" {
		return xerror.New("[version] should not be null")
	}

	_, err := ver.NewVersion(v)
	return xerror.WrapF(err, "[v] version format error")
}

func (r *rpcEntry) Commands(commands ...*cobra.Command) error {
	rootCmd := r.opts.Command
	for _, cmd := range commands {
		if cmd == nil {
			continue
		}

		if rootCmd.Name() == cmd.Name() {
			return xerror.Fmt("command(%s) already exists", cmd.Name())
		}

		rootCmd.AddCommand(cmd)
	}
	return nil
}

// func(s server.Server, handle TestHandler, opts ...server.HandlerOption) error
func (r *rpcEntry) Handler(hdlr interface{}, opts ...server.HandlerOption) error {
	return xerror.Wrap(catdog_handler.Register(catdog_server.Default, hdlr, opts...))
}

func (r *rpcEntry) Plugins(pgs ...catdog_plugin.Plugin) (err error) {
	defer xerror.RespErr(&err)

	for _, pg := range pgs {
		xerror.PanicF(catdog_plugin.Register(pg, catdog_plugin.Module(r.opts.Name)), "Plugin [%s] Register error", pg.String())
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

	xerror.Exit(ent.Flags(func(flags *pflag.FlagSet) {
		flags.StringVar(&ent.addr, "gw_addr", ent.addr, "gateway address")
	}))

	xerror.Exit(catdog_abc.WithBeforeStart(func() {
		g := ent.app.Group(ent.gwPrefix)
		handlers := catdog_server.Default.Handlers()
		for i := range handlers {
			xerror.Panic(handlers[i](g))
		}

		cancel := xprocess.Go(func(ctx context.Context) (err error) {
			defer xerror.RespErr(&err)
			log.Infof("Server [http] Listening on http://%s", ent.addr)
			xerror.Exit(ent.app.Listen(ent.addr))
			log.Infof("Server [http] Closed OK")
			return nil
		})

		xerror.Exit(catdog_abc.WithBeforeStop(func() {
			xerror.Panic(cancel())
			if err := ent.app.Shutdown(); err != nil && err != http.ErrServerClosed {
				fmt.Println(xerror.Parse(err).Println())
			}
		}))
	}))

	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !catdog_config.Trace {
			return
		}

		for _, stacks := range ent.app.Stack() {
			for _, stack := range stacks {
				log.Debugf("%s %s", stack.Method, stack.Path)
			}
		}
	}))

	return ent
}

func New() Entry {
	return newEntry()
}
