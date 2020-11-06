package catdog_server

import (
	"context"
	"github.com/asim/nitro/v3/server"
	grpc "github.com/asim/nitro-plugins/server/grpc/v3"
)

var Default = &wrapper{Server: grpc.NewServer(server.Context(context.Background()))}
