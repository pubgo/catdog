package catdog_version_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_version"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
}

func (p Plugin) Flags() *pflag.FlagSet {
	return nil
}

func (p Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p Plugin) String() string {
	return p.name
}

func (p Plugin) Commands() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Short:   "catdog version",
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			defer xerror.RespErr(&err)
			for name, v := range catdog_version.List() {
				fmt.Println(name, string(xerror.PanicBytes(json.MarshalIndent(v, "", "  "))))
			}
			return nil
		},
	}
}

func NewPlugin() *Plugin {
	return &Plugin{
		name: "catdog_version",
	}
}
