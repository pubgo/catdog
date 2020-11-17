package ws_entry

import (
	"github.com/asim/nitro/v3/client"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_entry/entry"
	"github.com/pubgo/catdog/internal/plugins/server/server_http"
)

type wsEntry struct {
	catdog_entry.Entry
	c client.Client
}

func newEntry(name string) *wsEntry {
	ent := &wsEntry{
		Entry: entry.New(name, &entryServerWrapper{Server: server_http.NewServer()}),
	}

	return ent
}

func New(name string) catdog_entry.Entry {
	return newEntry(name)
}
