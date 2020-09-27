package catdog_rpc_entry

import (
	"github.com/pubgo/catdog/catdog_log"
	"github.com/pubgo/xlog"
	"net/http"
	"sync"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"

	"github.com/micro/go-micro/v3/server"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_server"
	"github.com/pubgo/xerror"
	"github.com/spf13/pflag"
)

var _ catdog_entry.Entry = (*rpcEntry)(nil)

const Name = "entry.rpc"

var log xlog.XLog

func init() {
	xerror.Exit(catdog_log.Watch(Name, &log))
}

type rpcEntry struct {
	catdog_entry.Entry

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
	return xerror.Wrap(r.Entry.Flags(func(flags *pflag.FlagSet) {
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
	go func() {
		defer xerror.RespGoroutine(r.Options().Name)
		log.InfoF("Server [http] Listening on http://%s", r.addr)
		xerror.Exit(r.app.Listen(r.addr))
		log.InfoF("Server [http] Closed OK")
	}()
	return nil
}

func (r *rpcEntry) middleware() []interface{} {
	return []interface{}{
		middleware.Logger("${method} - ${path}\n"),
	}
}

func (r *rpcEntry) initCatDog(cat catdog_abc.CatDog) (err error) {
	xerror.RespErr(&err)

	opts := r.Options()
	cat.Init(catdog_abc.Name(opts.Name), catdog_abc.Version(opts.Version))
	cat.Init(append(opts.Options,
		catdog_abc.BeforeStart(r.startService),
		catdog_abc.BeforeStop(r.stopService),
		catdog_abc.AfterStart(r.pathRouterTrace),
	)...)

	g := r.app.Group(r.gwPrefix)
	handlers := catdog_server.Default.Handlers()
	for i := range handlers {
		xerror.Panic(handlers[i](g))
	}

	return nil
}

func newEntry() *rpcEntry {
	ent := &rpcEntry{
		Entry:    catdog_entry.NewBase(),
		app:      fiber.New(),
		addr:     ":8080",
		gwPrefix: "",
	}
	ent.app.Use(ent.middleware()...)
	xerror.Exit(ent.initFlags())
	xerror.Exit(catdog_abc.Watch(ent.initCatDog))
	return ent
}

func New() catdog_entry.Entry {
	return newEntry()
}
