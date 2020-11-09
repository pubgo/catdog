package catdog_client

import (
	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
)

var Default = &wrapper{Client: grpcC.NewClient()}
