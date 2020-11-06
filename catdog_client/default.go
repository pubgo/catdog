package catdog_client

import (
	"crypto/tls"
	grpc "github.com/asim/nitro-plugins/client/grpc/v3"
	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/pubgo/xerror"
)

var Default = &wrapper{Client: grpc.NewClient()}

// WithTLS sets the TLS config for the catdog_service
func WithTLS(t *tls.Config) Option {
	return func(o *Options) {
		xerror.Exit(o.Client.Init(grpcC.AuthTLS(t)))
		xerror.Exit(o.Server.Init(grpcS.AuthTLS(t)))
	}
}
