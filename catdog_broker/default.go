package catdog_broker

import (
	mBroker "github.com/asim/nitro/v3/broker/memory"
)

var (
	Default = &wrapper{Broker: mBroker.NewBroker()}
)
