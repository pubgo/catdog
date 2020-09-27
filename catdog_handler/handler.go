package catdog_handler

import (
	"github.com/micro/go-micro/v3/server"
)

type Handler struct {
	Register interface{}
	Handler  interface{}
	Opts     []server.HandlerOption
}

func Register(register, hdlr interface{}, opts ...server.HandlerOption) *Handler {
	return &Handler{
		Register: register,
		Handler:  hdlr,
		Opts:     opts,
	}
}
