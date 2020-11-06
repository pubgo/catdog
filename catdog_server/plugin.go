package catdog_server

import (
	"crypto/tls"
	grpcS "github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/pubgo/xerror"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"time"

	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
)

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
	opts server.Options
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *Plugin) String() string {
	return p.name
}

func (p *Plugin) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(p.name, pflag.PanicOnError)
	flags.StringVar(&p.opts.Address, "server_addr", p.opts.Address, "server address")
	flags.StringVar(&p.opts.Name, "server_name", p.opts.Name, "server name")
	return flags
}

func NewPlugin() *Plugin {
	p := &Plugin{
		opts: Default.Options(),
	}

	return p
}

// WithTLS sets the TLS config for the catdog_service
func WithTLS(t *tls.Config) Option {
	return func(o *Options) {
		xerror.Exit(Default.Server.Init(grpcS.AuthTLS(t)))
	}
}

//sets the address of the internal_catdog_server
func Address(addr string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Address(addr)))
	}
}

// name of the catdog_service
func Name(n string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Name(n)))
	}
}

// Version of the catdog_service
func Version(v string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Version(v)))
	}
}

// Metadata associated with the catdog_service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.Metadata(md)))
	}
}

// RegisterTTL specifies the TTL to use when registering the catdog_service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.RegisterTTL(t)))
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		xerror.Exit(o.Server.Init(server.RegisterInterval(t)))
	}
}
