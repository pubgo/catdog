package catdog_server

import (
	"context"

	grpc "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/server"
)

var Default = &wrapper{Server: grpc.NewServer(server.Context(context.Background()))}
