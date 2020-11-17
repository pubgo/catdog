package catdog_server

import (
	"context"

	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/server"
)

var Default = &wrapper{Server: grpcS.NewServer(server.Context(context.Background()))}
