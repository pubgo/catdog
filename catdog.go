package catdog

import (
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/model"
	"github.com/micro/go-micro/v3/server"
	"github.com/pubgo/catdog/plugins/catdog_debug_plugin"
	"github.com/pubgo/catdog/plugins/catdog_pidfile_plugin"
	"github.com/pubgo/catdog/plugins/catdog_recovery_plugin"
	"github.com/pubgo/catdog/plugins/catdog_version_plugin"
	"github.com/pubgo/xlog"
	"os"
	"os/signal"
	"syscall"

	"github.com/pubgo/xerror"

	signalutil "github.com/micro/go-micro/v3/util/signal"

	"github.com/pubgo/catdog/catdog_abc"

	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/catdog/plugins/catdog_broker_plugin"
	"github.com/pubgo/catdog/plugins/catdog_client_plugin"
	"github.com/pubgo/catdog/plugins/catdog_config_plugin"
	"github.com/pubgo/catdog/plugins/catdog_log_plugin"
	"github.com/pubgo/catdog/plugins/catdog_model_plugin"
	"github.com/pubgo/catdog/plugins/catdog_registry_plugin"
	"github.com/pubgo/catdog/plugins/catdog_server_plugin"
)

type catDog struct {
	opts catdog_abc.Options
}

func (t *catDog) Name() string {
	return t.Server().Options().Name
}

func (t *catDog) loadDefaultPlugin() {
	defer xerror.RespExit()
	var register = func(plugin catdog_plugin.Plugin) {
		xerror.Panic(catdog_plugin.Register(plugin, catdog_plugin.Module(_globalPlugin)))
	}
	register(catdog_config_plugin.NewPlugin())
	register(catdog_log_plugin.NewPlugin())
	register(catdog_model_plugin.NewPlugin())
	register(catdog_registry_plugin.NewPlugin())
	register(catdog_broker_plugin.NewPlugin())
	register(catdog_server_plugin.NewPlugin())
	register(catdog_client_plugin.NewPlugin())
	register(catdog_version_plugin.NewPlugin())
	register(catdog_pidfile_plugin.NewPlugin())
	register(catdog_debug_plugin.NewPlugin())
	register(catdog_recovery_plugin.NewPlugin())
}

func (t *catDog) Init(opts ...catdog_abc.Option) {
	for _, o := range opts {
		o(&t.opts)
	}
}

func (t *catDog) Options() catdog_abc.Options {
	return t.opts
}

func (t *catDog) Client() client.Client {
	return t.opts.Client
}

func (t *catDog) Server() server.Server {
	return t.opts.Server
}

func (t *catDog) Model() model.Model {
	return t.opts.Model
}

func (t *catDog) String() string {
	return t.Server().Options().Name
}

func (t *catDog) Start() (err error) {
	defer xerror.RespErr(&err)

	for _, fn := range t.opts.BeforeStart {
		xlog.DebugF("BeforeStart: %s", xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		xerror.Panic(fn())
	}

	xerror.Panic(t.Server().Start())

	for _, fn := range t.opts.AfterStart {
		xlog.DebugF("AfterStart: %s", xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		xerror.Panic(fn())
	}

	return
}

func (t *catDog) Stop() (err error) {
	defer xerror.RespErr(&err)

	var errs []error
	for _, fn := range t.opts.BeforeStop {
		xlog.DebugF("BeforeStop: %s", xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		if err := fn(); err != nil {
			errs = append(errs, xerror.Wrap(err))
		}
	}

	if err := t.Server().Stop(); err != nil {
		return xerror.Combine(append(errs, xerror.Wrap(err))...)
	}

	for _, fn := range t.opts.AfterStop {
		xlog.DebugF("AfterStop: %s", xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		if err := fn(); err != nil {
			errs = append(errs, xerror.Wrap(err))
		}
	}

	return xerror.Combine(errs...)
}

// Run catdog service
func (t *catDog) Run() (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(t.Start())

	if t.opts.IsSignal {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, append(signalutil.Shutdown(), syscall.SIGHUP)...)
		t.opts.Signal = <-ch
	}

	xerror.Panic(t.Stop())
	return
}

func newCatDog(opts ...catdog_abc.Option) *catDog {
	return &catDog{
		opts: catdog_abc.NewOption(opts...),
	}
}

func New(opts ...catdog_abc.Option) *catDog {
	return newCatDog(opts...)
}
