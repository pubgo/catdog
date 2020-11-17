package catdog_app

import (
	"os"
	"os/signal"
	"syscall"

	signalUtil "github.com/asim/nitro/v3/util/signal"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_version"
	"github.com/pubgo/dix/dix_run"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
)

func Start(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(ent.Init())

	// 启动配置, 初始化组件, 初始化插件
	plugins := catdog_plugin.List(catdog_plugin.Module(ent.Options().Name))
	for _, pg := range append(catdog_plugin.List(), plugins...) {
		key := pg.String()
		r, err := catdog_config.Load(key)
		xerror.PanicF(err, "plugin [%s] load error", key)
		xerror.PanicF(pg.Init(r), "plugin [%s] init error", key)

		hdlr := pg.Handler()
		if hdlr != nil {
			xerror.Panic(ent.Handler(hdlr.Handler, hdlr.Opts...))
		}
	}

	xerror.Panic(dix_run.BeforeStart())
	xerror.Panic(ent.Start())
	xerror.Panic(dix_run.AfterStart())

	return
}

func Stop(ent catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(dix_run.BeforeStop())
	xerror.Panic(ent.Stop())
	xerror.Panic(dix_run.AfterStop())

	return nil
}

func Run(entries ...catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	if len(entries) == 0 {
		return xerror.New("[entries] should not be zero")
	}

	for _, ent := range entries {
		if ent == nil {
			return xerror.New("[ent] should not be nil")
		}

		opt := ent.Options()
		if opt.Name == "" || opt.Version == "" {
			return xerror.New("neither [name] nor [version] can be empty")
		}
	}

	var rootCmd = &cobra.Command{Use: catdog_config.Domain, Version: catdog_version.Version}
	rootCmd.PersistentFlags().AddFlagSet(catdog_config.DefaultFlags())
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error { return xerror.Wrap(cmd.Help()) }

	for _, ent := range entries {
		ent := ent
		cmd := ent.Options().Command

		// 检查Command是否注册
		for _, c := range rootCmd.Commands() {
			if c.Name() == cmd.Name() {
				return xerror.Fmt("command(%s) already exists", cmd.Name())
			}
		}

		// 注册plugin的command和flags
		entPlugins := catdog_plugin.List(catdog_plugin.Module(ent.Options().Name))
		for _, pl := range append(catdog_plugin.List(), entPlugins...) {
			cmd.PersistentFlags().AddFlagSet(pl.Flags())
			xerror.Panic(ent.Commands(pl.Commands()))
		}

		runCmd := ent.Options().RunCommand
		runCmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			xerror.Panic(Start(ent))

			if catdog_config.IsBlock {
				ch := make(chan os.Signal, 1)
				signal.Notify(ch, append(signalUtil.Shutdown(), syscall.SIGHUP)...)
				<-ch
			}

			xerror.Panic(Stop(ent))
			return nil
		}
		rootCmd.AddCommand(cmd)
	}
	return xerror.Wrap(rootCmd.Execute())
}
