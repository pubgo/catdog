package rpc_entry

import (
	"context"
	"fmt"
	"github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/catdog/plugins/catdog_broker"
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
	"github.com/pubgo/catdog/plugins/catdog_server"
)

const defaultContentType = "application/json"

var _ catdog_entry.Entry = (*rpcEntry)(nil)

type rpcEntry struct {
	s        *wrapper
	opts     catdog_entry.Options
	mux      sync.Mutex
	app      *fiber.App
	addr     string
	gwPrefix string
}

func (r *rpcEntry) Start() (err error) {
	defer xerror.RespErr(&err)
	catdog_server.Default.Server = r.s
	return xerror.Wrap(catdog_server.Default.Start())
}

func (r *rpcEntry) Stop() (err error) {
	defer xerror.RespErr(&err)
	return xerror.Wrap(catdog_server.Default.Stop())
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

func (r *rpcEntry) Options() catdog_entry.Options {
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

func newEntry() *rpcEntry {
	ent := &rpcEntry{
		s: &wrapper{Server: grpc.NewServer(server.Context(context.Background()))},
		opts: catdog_entry.Options{
			RunCommand: &cobra.Command{Use: "run"},
			Command:    &cobra.Command{},
		},
		app:      fiber.New(),
		addr:     ":8080",
		gwPrefix: "api",
	}
	ent.opts.Command.AddCommand(ent.opts.RunCommand)
	ent.app.Use(ent.middleware()...)
	catdog_server.Default.Server = grpc.NewServer(server.Context(context.Background()))

	xerror.Exit(ent.Flags(func(flags *pflag.FlagSet) {
		flags.StringVar(&ent.addr, "gw_addr", ent.addr, "gateway address")
	}))

	xerror.Exit(catdog_abc.WithBeforeStart(func() {
		xerror.Exit(ent.s.Init(server.Broker(catdog_broker.Default.Broker)))

		g := ent.app.Group(ent.gwPrefix)
		for i := range ent.s.handlers {
			xerror.Panic(ent.s.handlers[i](g))
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

func New() catdog_entry.Entry {
	return newEntry()
}
