package catdog_handler

import (
	"reflect"

	"github.com/asim/nitro/v3/server"
	"github.com/pubgo/catdog/catdog_data"
	"github.com/pubgo/xerror"
)

type Handler struct {
	Handler interface{}
	Opts    []server.HandlerOption
}

func New(hdlr interface{}, opts ...server.HandlerOption) *Handler {
	return &Handler{
		Handler: hdlr,
		Opts:    opts,
	}
}

func Register(s server.Server, hdlr interface{}, opts ...server.HandlerOption) (err error) {
	defer xerror.RespErr(&err)

	if hdlr == nil {
		return xerror.New("[params] should not be nil")
	}
	if s == nil {
		return xerror.New("[server] should not be nil")
	}

	var vRegister reflect.Value
	hd := reflect.New(reflect.Indirect(reflect.ValueOf(hdlr)).Type()).Type()
	for _, v := range catdog_data.List() {
		v1 := reflect.TypeOf(v)
		if v1.NumIn() < 2 {
			continue
		}

		if hd.Implements(reflect.TypeOf(v).In(1)) {
			vRegister = reflect.ValueOf(v)
			break
		}
	}

	if !vRegister.IsValid() || vRegister.IsNil() {
		return xerror.Fmt("[%s] 没有找到匹配的interface", hd.Name())
	}

	vHandler := reflect.ValueOf(hdlr)
	if vRegister.Kind() != reflect.Func ||
		vRegister.Type().NumIn() < 2 ||
		vRegister.Type().In(0).String() != "server.Server" {
		return xerror.New("the first parameter should be <func(s server.Server, hdlr handler, opts ...server.HandlerOption) error> type")
	}

	if !vHandler.Type().Implements(vRegister.Type().In(1)) {
		return xerror.Fmt("the second parameter type does not match")
	}

	var sOpts = []reflect.Value{reflect.ValueOf(s), vHandler}
	for _, opt := range opts {
		sOpts = append(sOpts, reflect.ValueOf(opt))
	}

	if ret := vRegister.Call(sOpts); !ret[0].IsNil() {
		return xerror.WrapF(ret[0].Interface().(error), "%v, %v", vHandler.Type(), vRegister.Type())
	}
	return
}
