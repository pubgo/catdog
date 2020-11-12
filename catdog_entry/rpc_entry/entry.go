package rpc_entry

import (
	"context"
	"fmt"
	"github.com/pubgo/xlog"
	"net/http"
	"strings"
	"sync"
	"time"

	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
	"github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	ver "github.com/hashicorp/go-version"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xprocess"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/internal/catdog_abc"
)

const defaultContentType = "application/json"

var _ catdog_entry.Entry = (*rpcEntry)(nil)

type rpcEntry struct {
	s        *serverWrapper
	c        client.Client
	opts     catdog_entry.Options
	mux      sync.Mutex
	app      *fiber.App
	addr     string
	gwPrefix string
}

func (r *rpcEntry) Server() server.Server {
	return r.s.Server
}

func (r *rpcEntry) Client() client.Client {
	return r.c
}

func (r *rpcEntry) Start() (err error) {
	defer xerror.RespErr(&err)

	g := r.app.Group(r.gwPrefix)
	for i := range r.s.handlers {
		xerror.Panic(r.s.handlers[i](g))
	}

	cancel := xprocess.Go(func(ctx context.Context) (err error) {
		defer xerror.RespErr(&err)
		log.Infof("Server [http] Listening on http://%s", r.addr)
		xerror.Exit(r.app.Listen(r.addr))
		log.Infof("Server [http] Closed OK")
		return nil
	})

	xerror.Exit(catdog_abc.WithBeforeStop(func() {
		xerror.Panic(cancel())
		if err := r.app.Shutdown(); err != nil && err != http.ErrServerClosed {
			fmt.Println(xerror.Parse(err).Println())
		}
	}))

	return xerror.Wrap(r.s.Start())
}

func (r *rpcEntry) Stop() (err error) {
	defer xerror.RespErr(&err)
	return xerror.Wrap(r.s.Stop())
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
	return xerror.Wrap(catdog_handler.Register(r.s, hdlr, opts...))
}

func newEntry() *rpcEntry {
	ent := &rpcEntry{
		c: grpcC.NewClient(),
		s: &serverWrapper{Server: grpc.NewServer(server.Context(context.Background()))},
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

	xerror.Exit(ent.Flags(func(flags *pflag.FlagSet) {
		flags.StringVar(&ent.addr, "gw_addr", ent.addr, "rpc gateway address")
	}))

	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !catdog_config.Trace {
			return
		}

		xlog.Debug("rpc entry trace")
		for _, stacks := range ent.app.Stack() {
			for _, stack := range stacks {
				if stack.Path == "/" {
					continue
				}

				log.Debugf("%s %s", stack.Method, stack.Path)
			}
		}
		fmt.Println()
	}))

	return ent
}

func New() catdog_entry.Entry {
	return newEntry()
}
