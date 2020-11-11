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
	"github.com/pubgo/catdog/version"
)

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

		cmd := ent.Options().Command
		cmd.Version = ent.Options().Version
		cmd.Use = ent.Options().Name
		if len(ent.Options().Description) > 0 {
			cmd.Short = ent.Options().Description[0]
		}
		if len(ent.Options().Description) > 1 {
			cmd.Long = ent.Options().Description[1]
		}
		if len(ent.Options().Description) > 2 {
			cmd.Example = ent.Options().Description[2]
		}
	}

	var rootCmd = &cobra.Command{Use: catdog_config.Domain, Version: version.Version}
	rootCmd.PersistentFlags().AddFlagSet(catdog_config.DefaultFlags())
	rootCmd.RunE = func(cmd *cobra.Command, args []string) error { return xerror.Wrap(cmd.Help()) }

	for _, ent := range entries {
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
