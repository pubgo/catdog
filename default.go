package catdog

import (
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/version"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
)

const _globalPlugin = "__global_default_plugins"

func init() {
	// 加载系统默认插件
	defaultCatDog.loadDefaultPlugin()

	// 初始化version
	version.Init()
}

var _ catdog_abc.CatDog = (*catDog)(nil)

var defaultCatDog = newCatDog()

func Run(entries ...catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	if len(entries) == 0 {
		return xerror.New("parameters should not be zero")
	}

	var rootCmd = &cobra.Command{}
	rootCmd.Use = catdog_config.Domain
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return xerror.Wrap(cmd.Help())
	}

	for _, ent := range entries {
		cmd := ent.Options().RunCommand

		// 检查服务是否注册
		for _, c := range rootCmd.Commands() {
			if c.Name() == cmd.Parent().Name() {
				return xerror.Fmt("command(%s) already exists", cmd.Parent().Name())
			}
		}

		// 系统默认注册的插件, 用户全局注册的插件
		pls := append(catdog_plugin.List(catdog_plugin.Module(_globalPlugin)), catdog_plugin.List()...)

		// 加载全局plugin
		for _, pg := range pls {
			cmd.Parent().PersistentFlags().AddFlagSet(pg.Flags())

			if pg.Commands() != nil {
				cmd.Parent().AddCommand(pg.Commands())
			}

			if pg.Handler() != nil {
				hdlr := pg.Handler()
				xerror.Panic(catdog_handler.Register(hdlr.Register, hdlr.Handler, hdlr.Opts...))
			}
		}

		plugins := catdog_plugin.List(catdog_plugin.Module(ent.Options().Name))
		cmd.Parent().PersistentPreRunE = func(cmd *cobra.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			for _, pl := range plugins {
				if pl.Handler() != nil {
					hdlr := pl.Handler()
					xerror.Panic(catdog_handler.Register(hdlr.Register, hdlr.Handler, hdlr.Opts...))
				}
			}
			return xerror.Wrap(dix.Dix(defaultCatDog))
		}
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			return xerror.Wrap(defaultCatDog.Run())
		}
		rootCmd.AddCommand(cmd.Parent())
	}

	return xerror.Wrap(rootCmd.Execute())
}
