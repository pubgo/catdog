package catdog

import (
	"github.com/pubgo/catdog/catdog_entry"
)

type Entry = catdog_entry.Entry

func NewEntry() catdog_entry.Entry {
	return catdog_entry.New()
}
