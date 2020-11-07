package catdog_redis

import (
	"encoding/json"

	"github.com/pubgo/xerror"

	"github.com/asim/nitro/v3/config/reader"
	"github.com/pubgo/catdog/catdog_plugin"
)

type Plugin struct {
	catdog_plugin.Plugin
	name string
}

func (p *Plugin) Watch(r reader.Value) error {
	var cfg config
	xerror.Exit(json.Unmarshal(r.Bytes(), &cfg))

	return nil
}

func init() {
	p := &Plugin{Plugin: catdog_plugin.NewBase("redis")}
	xerror.Exit(catdog_plugin.Register(p))
}
