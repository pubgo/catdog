package catdog_app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"

	signalUtil "github.com/asim/nitro/v3/util/signal"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
)

func Run(entries ...catdog_entry.Entry) (err error) {
	defer xerror.RespErr(&err)

	var rootCmd = &cobra.Command{Use: catdog_config.Domain}
	rootCmd.PersistentFlags().AddFlagSet(catdog_config.DefaultFlags())
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error { return xerror.Wrap(cmd.Help()) }

	if len(entries) == 0 {
		return xerror.New("[entries] should not be zero")
	}

	for _, ent := range entries {
		cmd := ent.Options().Command
		runCmd := ent.Options().RunCommand

		// 检查Command是否注册
		for _, c := range rootCmd.Commands() {
			if c.Name() == cmd.Name() {
				return xerror.Fmt("command(%s) already exists", cmd.Name())
			}
		}

		// 注册插件
		nameModule := catdog_plugin.Module(ent.Options().Name)
		plugins := append(catdog_plugin.List(), catdog_plugin.List(nameModule)...)
		for _, pl := range plugins {
			cmd.PersistentFlags().AddFlagSet(pl.Flags())

			if pl.Commands() != nil {
				cmd.AddCommand(pl.Commands())
			}
		}

		runCmd.RunE = func(cmd *cobra.Command, args []string) (err error) {
			defer xerror.RespErr(&err)

			xerror.Panic(Start(ent))

			ch := make(chan os.Signal, 1)
			signal.Notify(ch, append(signalUtil.Shutdown(), syscall.SIGHUP)...)

			xerror.Panic(Stop())
			return nil
		}
		rootCmd.AddCommand(cmd)
	}
	return xerror.Wrap(rootCmd.Execute())
}
