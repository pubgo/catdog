package catdog_version

import (
	"encoding/json"
	"fmt"

	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
)

type Plugin struct {
	catdog_plugin.Plugin
}

func (p Plugin) Commands() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "catdog version",
		Aliases: []string{"v"},
		Run: func(cmd *cobra.Command, args []string) {
			for name, v := range List() {
				fmt.Println(name, string(xerror.PanicBytes(json.MarshalIndent(v, "", "  "))))
			}
		},
	}
}

func init() {
	p := &Plugin{Plugin: catdog_plugin.NewBase("catdog_version")}
	xerror.Exit(catdog_plugin.Register(p))
}
