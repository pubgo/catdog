package catdog_registry

import (
	"github.com/asim/nitro-plugins/registry/mdns"
)

var (
	Default = &wrapper{Registry: mdns.NewRegistry()}
)
