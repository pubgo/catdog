package ctl_entry

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/asim/nitro/v3/server"
	ver "github.com/hashicorp/go-version"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var _ catdog_entry.Entry = (*ctlEntry)(nil)

type ctlEntry struct {
	opts    catdog_entry.Options
	handler func() error
}

func (t *ctlEntry) Handler(handler interface{}, opts ...server.HandlerOption) error {
	if handler == nil {
		xlog.Warn("[handler] should not be nil")
		return nil
	}

	vh := reflect.ValueOf(handler)

	switch vh.Kind() {
	case reflect.Func:
		t.handler = handler.(func() error)
	case reflect.Struct, reflect.Ptr:
		t.handler = handler.(func() error)
	default:
		return xerror.Fmt("[handler] type error, type:%#v", handler)
	}

	return nil
}

func (t *ctlEntry) Init() (err error) {
	defer xerror.RespErr(&err)

	catdog_config.IsBlock = false
	catdog_config.Project = t.Options().Name
	t.opts.Initialized = true

	return nil
}

func (t *ctlEntry) Start() (err error) {
	defer xerror.RespErr(&err)

	return xerror.Wrap(t.handler())
}

func (t *ctlEntry) Stop() error {
	return nil
}

func (t *ctlEntry) Options() catdog_entry.Options {
	return t.opts
}

func (t *ctlEntry) Flags(fn func(flags *pflag.FlagSet)) (err error) {
	defer xerror.RespErr(&err)
	fn(t.opts.Command.PersistentFlags())
	return nil
}

func (t *ctlEntry) Description(description ...string) error {
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

func (t *ctlEntry) Version(v string) error {
	t.opts.Version = strings.TrimSpace(v)
	if t.opts.Version == "" {
		return xerror.New("[version] should not be null")
	}

	t.opts.Command.Version = v
	_, err := ver.NewVersion(v)
	return xerror.WrapF(err, "[v] version format error")
}

func (t *ctlEntry) Commands(commands ...*cobra.Command) error {
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

func newEntry(name string) *ctlEntry {
	name = strings.TrimSpace(name)
	if name == "" {
		xerror.Panic(xerror.New("the [name] parameter should not be empty"))
	}

	rootCmd := &cobra.Command{Use: name}
	runCmd := &cobra.Command{Use: "run", Short: "run as a service"}
	rootCmd.AddCommand(runCmd)

	ent := &ctlEntry{
		opts: catdog_entry.Options{
			Name:       name,
			RunCommand: runCmd,
			Command:    rootCmd,
		},
	}

	//xerror.Panic(ent.Flags(func(flags *pflag.FlagSet) {
	//	flags.StringVar()
	//}))
	return ent
}

func New(name string) *ctlEntry {
	return newEntry(name)
}
