package catdog_client

import (
	"github.com/micro/go-micro/v3/client/grpc"
)

var Default = &catdogClient{Client: grpc.NewClient()}
