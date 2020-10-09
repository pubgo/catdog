package catdog_broker

import (
	"github.com/micro/go-micro/v3/broker"
)

type wrapper struct {
	broker.Broker
}
