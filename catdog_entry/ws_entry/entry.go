package ws_entry

import (
	"github.com/asim/nitro/v3/client"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_entry/entry"
	"github.com/pubgo/catdog/internal/plugins/server/server_http"
)

type wsEntry struct {
	*entry.BaseEntry
	c client.Client
}

func newEntry(name string) *wsEntry {
	s := &entryServerWrapper{Server: server_http.NewServer()}
	base := entry.New(name, s)
	s.base = base

	ent := &wsEntry{BaseEntry: base}
	return ent
}

func New(name string) catdog_entry.Entry {
	return newEntry(name)
}
