package catdog_broker

import (
	"github.com/micro/go-micro/v3/broker/http"
)

var Default = &wrapper{Broker: http.NewBroker()}
