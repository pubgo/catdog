package catdog_config_plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v3/config/encoder/yaml"
	"github.com/micro/go-micro/v3/config/source"
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	mFile "github.com/micro/go-micro/v3/config/source/file"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type config struct {
	Address  []string `json:"address"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	Timeout  int32    `json:"timeout"`
}

type Plugin struct {
	name string
	ss   []source.Source
	cfg  *catdog_config.Config
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

	cat.Init(catdog_abc.AfterStart(func() error {
		if catdog_config.Trace {
			fmt.Println("deps", dix.Graph())
			fmt.Println("config", catdog_util.MarshalIndent(catdog_config.GetCfg().Map()))
			fmt.Println("plugins", catdog_plugin.String())
		}
		return nil
	}), catdog_abc.BeforeStart(func() error {
		return xerror.Wrap(catdog_config.WatchStart())
	}), catdog_abc.AfterStop(func() error {
		return xerror.Wrap(catdog_config.WatchStop())
	}))

	// 加载本地配置文件
	xerror.Exit(p.cfg.Load(
		mFile.NewSource(
			mFile.WithPath(filepath.Join(catdog_config.CfgDir, "config", "config.yaml")),
			source.WithEncoder(yaml.NewEncoder()),
		),
	))

	//var cfg config
	//xerror.Exit(p.cfg.Get("plugins", "config", "watcher").Scan(&cfg))
	//
	//xerror.Exit(p.cfg.Load(
	//	mEtcd.NewSource(
	//		mEtcd.WithAddress(cfg.Address...),
	//		mEtcd.WithPrefix(strings.Join(catdog_config.ProjectPrefix(), "/")),
	//		mEtcd.StripPrefix(true),
	//	),
	//))

	return nil
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
