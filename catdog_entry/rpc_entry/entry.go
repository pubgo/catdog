package rpc_entry

import (
	"context"
	"github.com/pubgo/catdog/catdog_entry/base_entry"
	"time"

	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
	"github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber/middleware"
	"github.com/pubgo/catdog/catdog_entry"
)

type entry struct {
	catdog_entry.Entry
	c client.Client
}

func (t *entry) middleware() []interface{} {
	return []interface{}{
		middleware.Recover(),
		middleware.Logger(middleware.LoggerConfig{
			Format:     "#${pid} - ${time} ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: time.RFC3339,
		}),
	}
}

func newEntry(name string) *entry {
	ent := &entry{
		Entry: base_entry.New(
			name,
			grpc.NewServer(server.Context(context.Background()), server.Name(name)),
		),
		c: grpcC.NewClient(),
	}
	ent.Options().App.Use(ent.middleware()...)

	return ent
}

func New(name string) catdog_entry.Entry {
	return newEntry(name)
}
