package catdog_handler

import (
	"github.com/micro/go-micro/v3/server"
	"github.com/pubgo/catdog/catdog_server"
	"github.com/pubgo/xerror"
	"reflect"
)

type Handler struct {
	Register interface{}
	Handler  interface{}
	Opts     []server.HandlerOption
}

func New(register, hdlr interface{}, opts ...server.HandlerOption) *Handler {
	return &Handler{
		Register: register,
		Handler:  hdlr,
		Opts:     opts,
	}
}

func Register(register interface{}, hdlr interface{}, opts ...server.HandlerOption) (err error) {
	defer xerror.RespErr(&err)

	if register == nil || hdlr == nil {
		return xerror.New("params should not be nil")
	}

	vRegister := reflect.ValueOf(register)
	vHandler := reflect.ValueOf(hdlr)

	if vRegister.Kind() != reflect.Func ||
		vRegister.Type().NumIn() < 2 ||
		vRegister.Type().In(0).String() != "server.Server" {
		return xerror.New("the first parameter should be <func(s server.Server, hdlr handler, opts ...server.HandlerOption) error> type")
	}

	if !vHandler.Type().Implements(vRegister.Type().In(1)) {
		return xerror.Fmt("the second parameter type does not match")
	}

	var sOpts = []reflect.Value{
		reflect.ValueOf(catdog_server.Default),
		vHandler,
	}
	for _, opt := range opts {
		sOpts = append(sOpts, reflect.ValueOf(opt))
	}
	if ret := vRegister.Call(sOpts); !ret[0].IsNil() {
		return xerror.WrapF(ret[0].Interface().(error), "%v, %v", vHandler.Type(), vRegister.Type())
	}
	return
}
