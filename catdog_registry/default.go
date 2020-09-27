package catdog_registry

import (
	"github.com/micro/go-micro/v3/registry/mdns"
)

var (
	Default = &catdogRegistry{Registry: mdns.NewRegistry()}
)
