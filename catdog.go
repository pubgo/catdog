package catdog

import (
	"github.com/pubgo/catdog/catdog_app"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/plugins/catdog_pidfile"
	"github.com/pubgo/xerror"
)

func Run(entries ...Entry) error {
	return catdog_app.Run(entries...)
}

func Init() (err error) {
	defer xerror.RespErr(&err)
	catdog_pidfile.Debug()
	return nil
}

type Entry = catdog_entry.Entry

func NewEntry() catdog_entry.Entry {
	return catdog_entry.New()
}
