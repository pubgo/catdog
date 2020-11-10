package catdog_version

import (
	"encoding/json"
	"fmt"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
)

func init() {
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "catdog_version",
		OnCommands: func(cmd *cobra.Command) {
			cmd.Use = "version"
			cmd.Short = "catdog version"
			cmd.Aliases = []string{"v"}
			cmd.Run = func(cmd *cobra.Command, args []string) {
				for name, v := range List() {
					fmt.Println(name, string(xerror.PanicBytes(json.MarshalIndent(v, "", "  "))))
				}
			}
		},
	}))
}
