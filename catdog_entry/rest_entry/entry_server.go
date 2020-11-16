package rest_entry

import (
	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/xerror"
)

type entryServerWrapper struct {
	server.Server
}

func (t *entryServerWrapper) Start() (err error) {
	defer xerror.RespErr(&err)
	return nil
}

func (t *entryServerWrapper) Stop() (err error) {
	defer xerror.RespErr(&err)
	return nil
}
