package catdog_mongo_plugin

import (
	"github.com/pubgo/catdog/catdog_abc"
	"github.com/pubgo/catdog/catdog_log"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/spf13/pflag"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	catdog_plugin.Plugin
	log xlog.XLog
}

func (p *Plugin) GetClient(name string) {

}

func (p *Plugin) Init(cat catdog_abc.CatDog) (rErr error) {
	defer xerror.RespErr(&rErr)

	cat.Init()

	return xerror.Wrap(dix.Dix(p))
}

func (p *Plugin) Flags() *pflag.FlagSet {
	flags := p.Plugin.Flags()
	return flags
}

func NewPlugin() *Plugin {
	return &Plugin{
		Plugin: catdog_plugin.New("mongo"),
		log:    catdog_log.GetLog().Named("catdog.plugin.mongo"),
	}
}
