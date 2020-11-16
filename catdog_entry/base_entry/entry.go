package base_entry

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
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/catdog/plugins/catdog_server"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_entry.Entry = (*entry)(nil)

type entry struct {
	s    *entryServerWrapper
	opts catdog_entry.Options
}

func (t *entry) Init() (err error) {
	defer xerror.RespErr(&err)

	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !catdog_config.Trace || !t.opts.Initialized {
			return
		}

		xlog.Debug("entry rest trace")
		for _, stacks := range t.opts.App.Stack() {
			for _, stack := range stacks {
				if stack.Path == "/" {
					continue
				}

				log.Debugf("%s %s", stack.Method, stack.Path)
			}
		}
		fmt.Println()
	}))

	t.opts.Initialized = true
	catdog_config.Project = t.Options().Name
	catdog_server.Default.Server = t.s.Server

	return nil
}

func (t *entry) Start() (err error) {
	defer xerror.RespErr(&err)

	cancel := xprocess.Go(func(ctx context.Context) (err error) {
		defer xerror.RespErr(&err)

		addr := t.Options().RestAddr
		log.Infof("Server [http] Listening on http://%s", addr)
		xerror.Exit(t.opts.App.Listen(addr))
		log.Infof("Server [http] Closed OK")

		return nil
	})

	xerror.Exit(catdog_abc.WithBeforeStop(func() {
		xerror.Panic(cancel())
		if err := t.opts.App.Shutdown(); err != nil && err != http.ErrServerClosed {
			fmt.Println(xerror.Parse(err).Println())
		}
	}))

	return xerror.Wrap(t.s.Start())
}

func (t *entry) Stop() (err error) {
	defer xerror.RespErr(&err)
	return xerror.Wrap(t.s.Stop())
}

func (t *entry) Options() catdog_entry.Options {
	return t.opts
}

func (t *entry) Flags(fn func(flags *pflag.FlagSet)) (err error) {
	defer xerror.RespErr(&err)
	fn(t.opts.Command.PersistentFlags())
	return nil
}

func (t *entry) Description(description ...string) error {
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

func (t *entry) Version(v string) error {
	t.opts.Version = strings.TrimSpace(v)
	if t.opts.Version == "" {
		return xerror.New("[version] should not be null")
	}

	t.opts.Command.Version = v
	_, err := ver.NewVersion(v)
	return xerror.WrapF(err, "[v] version format error")
}

func (t *entry) Commands(commands ...*cobra.Command) error {
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
func (t *entry) Handler(hdlr interface{}, opts ...server.HandlerOption) error {
	return xerror.Wrap(catdog_handler.Register(t.s, hdlr, opts...))
}

func newEntry(name string, srv server.Server) *entry {
	name = strings.TrimSpace(name)
	if name == "" {
		xerror.Panic(xerror.New("the [name] parameter should not be empty"))
	}

	rootCmd := &cobra.Command{Use: name}
	runCmd := &cobra.Command{Use: "run", Short: "run as a service"}
	rootCmd.AddCommand(runCmd)

	xerror.Panic(srv.Init(
		server.Name(name),
		server.Context(context.Background()),
	))

	app := fiber.New()
	ent := &entry{
		s: &entryServerWrapper{Server: srv, router: app.Group(name)},
		opts: catdog_entry.Options{
			App:        app,
			Name:       Name,
			RestAddr:   ":8080",
			RunCommand: runCmd,
			Command:    rootCmd,
		},
	}

	xerror.Panic(ent.Flags(func(flags *pflag.FlagSet) {
		flags.StringVar(&ent.opts.RestAddr, "rest_addr", ent.opts.RestAddr, "the http server address")
	}))

	return ent
}

func New(name string, srv server.Server) catdog_entry.Entry {
	return newEntry(name, srv)
}
