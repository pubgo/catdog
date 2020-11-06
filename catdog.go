package catdog

import (
	"fmt"
	"path/filepath"

	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"

	"github.com/asim/nitro/v3/config"
	"github.com/asim/nitro/v3/config/source"
	mFile "github.com/asim/nitro/v3/config/source/file"

	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/plugins/config/encoder/yaml"
)

func Run(entries ...Entry) error {
	return catdog_app.Run(entries...)
}

func Init(project string) error {
	// 加载本地配置文件
	catdog_config.Init(
		config.WithSource(mFile.NewSource(
			mFile.WithPath(filepath.Join(catdog_config.Home, "config", "config.yaml")),
			source.WithEncoder(yaml.NewEncoder()),
		)),
	)

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

	xerror.Exit(dix.WithBeforeStart(func() {
		if catdog_config.Trace {
			fmt.Println("deps", dix.Graph())
			fmt.Println("config", string(catdog_config.LoadBytes()))
			fmt.Println("plugins", catdog_plugin.String())
		}
	}))
}

type Entry = catdog_entry.Entry

func NewEntry() catdog_entry.Entry {
	return catdog_entry.New()
}
