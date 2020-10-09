package catdog_server

import (
	"context"
	"github.com/micro/go-micro/v3/server"
	"github.com/micro/go-micro/v3/server/grpc"
)

var Default = &wrapper{Server: grpc.NewServer(server.Context(context.Background()))}
