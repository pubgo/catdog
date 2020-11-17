package entry

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/asim/nitro/v3/metadata"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/xerror"
)

type mth struct{}

const defaultContentType = "application/json"

type entryServerWrapper struct {
	server.Server
	routers []func(r fiber.Router)
}

var httpMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodConnect: {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

func (t *entryServerWrapper) httpHandler(httpMethod, path string, handlers ...fiber.Handler) error {
	httpMethod = strings.ToUpper(httpMethod)
	if _, ok := httpMethods[httpMethod]; !ok {
		return nil
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	t.routers = append(t.routers, func(r fiber.Router) {
		r.Add(strings.ToUpper(httpMethod), path, handlers...)
	})

	return nil
}

func (t *entryServerWrapper) httpHandle(ctx context.Context, _ server.Request, rsp interface{}) (err error) {
	defer xerror.RespErr(&err)

	mth := ctx.Value(mth{}).(reflect.Value)
	mthInType := mth.Type().In(1)
	mthOutType := mth.Type().In(2)

	view := rsp.(*fiber.Ctx)

	mthIn := reflect.New(mthInType.Elem())
	ret := reflect.ValueOf(view.BodyParser).Call([]reflect.Value{mthIn})
	if !ret[0].IsNil() {
		return xerror.Wrap(ret[0].Interface().(error))
	}

	mthOut := reflect.New(mthOutType.Elem())
	ret = mth.Call([]reflect.Value{reflect.ValueOf(ctx), mthIn, mthOut})
	if !ret[0].IsNil() {
		return xerror.Wrap(ret[0].Interface().(error))
	}

	return xerror.Wrap(view.JSON(mthOut.Interface()))
}

func (t *entryServerWrapper) Handle(handler server.Handler) (err error) {
	defer xerror.RespErr(&err)
	if handler == nil {
		return xerror.New("[handler] should not be nil")
	}

	xerror.PanicF(t.Server.Handle(handler), handler.Name())

	hdlr := reflect.ValueOf(handler.Handler())
	for mthName, httpRule := range handler.Options().Metadata {
		if mthName == "" || httpRule == nil || len(httpRule) == 0 {
			continue
		}

		var httpMethod, httpPath string
		for k, v := range httpRule {
			if k == "" || v == "" {
				continue
			}

			httpMethod = k
			httpPath = v
			break
		}

		mthHandle := hdlr.MethodByName(mthName)
		xerror.Panic(t.httpHandler(httpMethod, httpPath, func(view *fiber.Ctx) (err error) {
			defer xerror.RespErr(&err)

			handle := t.httpHandle
			hd := t.Server.Options().HdlrWrappers
			for i := len(hd); i > 0; i-- {
				handle = hd[i-1](handle)
			}

			// 处理 metadata fastHttpRequest
			headers := make(metadata.Metadata)
			view.Request().Header.VisitAll(func(key, value []byte) {
				headers[strings.ToLower(string(key))] = string(value)
			})

			request := &httpRequest{
				service:     fmt.Sprintf("%s.%s", t.Options().Name, handler.Name()),
				contentType: defaultContentType,
				method:      mthName,
				body:        view.Body(),
				header:      headers,
			}

			ctx := metadata.NewContext(view.Context(), headers)
			ctx = context.WithValue(ctx, fastHttpRequest{}, view)
			ctx = context.WithValue(ctx, mth{}, mthHandle)
			return xerror.Wrap(handle(ctx, request, view))
		}))
	}

	return
}
