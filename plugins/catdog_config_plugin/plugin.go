package catdog_config_plugin

import (
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v3/config"
	"github.com/micro/go-micro/v3/config/encoder/yaml"
	"github.com/micro/go-micro/v3/config/source"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"strings"

	mFile "github.com/micro/go-micro/v3/config/source/file"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	ss   []source.Source
	cfg  config.Config
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *Plugin) String() string {
	return p.name
}

func (p *Plugin) GetCfg() config.Config {
	return p.cfg
}

func (p *Plugin) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(p.name, pflag.PanicOnError)
	flags.StringVarP(&catdog_config.Mode, "mode", "m", catdog_config.Mode, "running mode(dev|test|stag|prod|release)")
	flags.StringVarP(&catdog_config.CfgDir, "cfg_dir", "", catdog_config.CfgDir, "project config dir")
	flags.StringVarP(&catdog_config.CfgPath, "cfg_path", "", catdog_config.CfgPath, "project config path")
	flags.BoolVar(&catdog_config.Debug, "debug", catdog_config.Debug, "enable debug")
	flags.BoolVar(&catdog_config.Trace, "trace", catdog_config.Trace, "enable trace")
	flags.StringVar(&catdog_config.Project, "name", catdog_config.Project, "project name")
	return flags
}

func (p *Plugin) catDogWatcher(cat catdog_abc.CatDog) (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(catdog_config.CheckRunMode())

	cat.Init(catdog_abc.AfterStart(func() error {
		if catdog_config.Trace {
			fmt.Println("plugins", dix.Graph())
			fmt.Println("config", string(xerror.PanicBytes(json.MarshalIndent(catdog_config.GetCfg().Map(), "", "  "))))
			fmt.Println("plugins", catdog_plugin.String())
		}
		return nil
	}))

	xerror.Exit(p.cfg.Load(
		mFile.NewSource(
			mFile.WithPath(filepath.Join(catdog_config.CfgDir, "config", "config.yaml")),
			source.WithEncoder(yaml.NewEncoder()),
		),
	))

	return xerror.Wrap(catdog_config.WatchStart())

	//xerror.Exit(p.cfg.Load(
	//	mFile.NewSource(
	//		mFile.WithPath("path"),
	//		mFile.WithPath("/i/do/not/exists.json"),
	//		source.WithEncoder(yaml.NewEncoder()),
	//	),
	//	mEtcd.NewSource(
	//		mEtcd.WithAddress("10.0.0.10:8500"),
	//		mEtcd.WithPrefix("/my/prefix"),
	//		mEtcd.StripPrefix(true),
	//	),
	//))
}

func New() *Plugin {
	p := &Plugin{
		name: "config",
		cfg:  catdog_config.GetCfg(),
	}
	xerror.Exit(catdog_abc.Watch(p.catDogWatcher))
	return p
}

func SetConfigPath(path string) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}

	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}

	if fi.IsDir() {
		return
	}
}
