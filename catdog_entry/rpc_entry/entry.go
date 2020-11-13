package rpc_entry

import (
	"context"
	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
	"github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_entry/base_entry"
)

type entry struct {
	catdog_entry.Entry
	c client.Client
}

//"#${pid} - ${time} ${status} - ${latency} ${method} ${path}\n"

func newEntry(name string) *entry {
	ent := &entry{
		Entry: base_entry.New(
			name,
			grpc.NewServer(server.Context(context.Background()), server.Name(name)),
		),
		c: grpcC.NewClient(),
	}

	return ent
}

func New(name string) catdog_entry.Entry {
	return newEntry(name)
}
