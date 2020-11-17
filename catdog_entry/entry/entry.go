package entry

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber/v2"
	ver "github.com/hashicorp/go-version"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xprocess"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_entry.Entry = (*BaseEntry)(nil)

type BaseEntry struct {
	app  *fiber.App
	s    *entryServerWrapper
	opts catdog_entry.Options
}

func (t *BaseEntry) Router(fn func(r fiber.Router)) {
	t.s.routers = append(t.s.routers, fn)
}

func (t *BaseEntry) Init() (err error) {
	defer xerror.RespErr(&err)

	t.opts.Initialized = true
	catdog_config.Project = t.Options().Name
	catdog_server.Default.Server = t.s.Server

	// 初始化routes
	t.app = fiber.New(t.opts.RestCfg)
	r := t.app.Group(t.opts.Name)
	for i := range t.s.routers {
		t.s.routers[i](r)
	}

	return nil
}

func (t *BaseEntry) Start() (err error) {
	defer xerror.RespErr(&err)

	cancel := xprocess.Go(func(ctx context.Context) (err error) {
		defer xerror.RespErr(&err)

		addr := t.Options().RestAddr
		log.Infof("Server [http] Listening on http://%s", addr)
		xerror.Panic(t.app.Listen(addr))
		log.Infof("Server [http] Closed OK")

		return nil
	})

	xerror.Panic(dix_run.WithBeforeStop(func(ctx *dix_run.BeforeStopCtx) {
		xerror.Panic(cancel())
		if err := t.app.Shutdown(); err != nil && err != http.ErrServerClosed {
			fmt.Println(xerror.Parse(err).Println())
		}
	}))

	return xerror.Wrap(t.s.Start())
}

func (t *BaseEntry) Stop() (err error) {
	defer xerror.RespErr(&err)
	return xerror.Wrap(t.s.Stop())
}

func (t *BaseEntry) Options() catdog_entry.Options {
	return t.opts
}

func (t *BaseEntry) initCfg() {
	r, err := catdog_config.Load("http_server")
	xerror.Panic(err)
	if r != nil {
		xerror.Panic(r.Scan(&t.opts.RestCfg))
	}
}

func (t *BaseEntry) initFlags() {
	xerror.Panic(t.Flags(func(flags *pflag.FlagSet) {
		flags.StringVar(&t.opts.RestAddr, "http_addr", t.opts.RestAddr, "the http server address")
		flags.BoolVar(&t.opts.RestCfg.DisableStartupMessage, "disable_startup_message", t.opts.RestCfg.DisableStartupMessage, "print out the http server art and listening address")
	}))
}

func (t *BaseEntry) Flags(fn func(flags *pflag.FlagSet)) (err error) {
	defer xerror.RespErr(&err)
	fn(t.opts.Command.PersistentFlags())
	return nil
}

func (t *BaseEntry) Description(description ...string) error {
	t.opts.Command.Short = fmt.Sprintf("This is a %s service", t.opts.Name)

	if len(description) > 0 {
		t.opts.Command.Short = description[0]
	}
	if len(description) > 1 {
		t.opts.Command.Long = description[1]
	}
	if len(description) > 2 {
		t.opts.Command.Example = description[2]
	}

	return nil
}

func (t *BaseEntry) Version(v string) error {
	t.opts.Version = strings.TrimSpace(v)
	if t.opts.Version == "" {
		return xerror.New("[version] should not be null")
	}

	t.opts.Command.Version = v
	_, err := ver.NewVersion(v)
	return xerror.WrapF(err, "[v] version format error")
}

func (t *BaseEntry) Commands(commands ...*cobra.Command) error {
	rootCmd := t.opts.Command
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
func (t *BaseEntry) Handler(hdlr interface{}, opts ...server.HandlerOption) error {
	return xerror.Wrap(catdog_handler.Register(t.s, hdlr, opts...))
}

func newEntry(name string, srv server.Server) *BaseEntry {
	name = strings.TrimSpace(name)
	if name == "" {
		xerror.Panic(xerror.New("the [name] parameter should not be empty"))
	}

	catdog_server.Default.Server = srv
	xerror.Panic(srv.Init(
		server.Name(name),
		server.Context(context.Background()),
	))

	rootCmd := &cobra.Command{Use: name}
	runCmd := &cobra.Command{Use: "run", Short: "run as a service"}
	rootCmd.AddCommand(runCmd)

	ent := &BaseEntry{
		s: &entryServerWrapper{Server: srv},
		opts: catdog_entry.Options{
			RestCfg:    fiber.New().Config(),
			Name:       name,
			RestAddr:   ":8080",
			RunCommand: runCmd,
			Command:    rootCmd,
		},
	}
	ent.initFlags()
	ent.initCfg()
	ent.trace()

	return ent
}

func New(name string, srv server.Server) *BaseEntry {
	return newEntry(name, srv)
}

//"#${pid} - ${time} ${status} - ${latency} ${method} ${path}\n"
