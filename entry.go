package catdog

import (
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/catdog_entry/catdog_rpc_entry"
)

type Entry = catdog_entry.Entry

func NewEntry() catdog_entry.Entry {
	return catdog_rpc_entry.New()
}
