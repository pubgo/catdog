package rpc_entry

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/asim/nitro/v3/metadata"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber"
	"github.com/pubgo/xerror"
)

type wrapper struct {
	server.Server
	handlers []func(g fiber.Router) error
}

func (r *wrapper) httpHandler(httpMethod, path string, handlers ...fiber.Handler) {
	r.handlers = append(r.handlers, func(router fiber.Router) error {
		if router == nil {
			return xerror.New("[router] should not be nil")
		}

		router.Add(strings.ToUpper(httpMethod), path, handlers...)
		return nil
	})
}

func (r *wrapper) Handle(handler server.Handler) (err error) {
	defer xerror.RespErr(&err)
	if handler == nil {
		return xerror.New("[handler] should not be nil")
	}

	xerror.Panic(r.Server.Handle(handler))

	hdr := reflect.ValueOf(handler.Handler())
	for mthName, httpRule := range handler.Options().Metadata {
		if httpRule == nil || len(httpRule) == 0 || mthName == "" {
			continue
		}

		mth := hdr.MethodByName(mthName)
		mthInType := mth.Type().In(1)
		mthOutType := mth.Type().In(2)

		handler := func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
			defer xerror.RespErr(&err)

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

		var httpMethod, httpPath string
		for k, v := range httpRule {
			httpMethod = k
			httpPath = v
			break
		}

		r.httpHandler(httpMethod, httpPath, func(view *fiber.Ctx) {
			defer xerror.Resp(func(err xerror.XErr) {
				_ = view.
					Status(http.StatusInternalServerError).
					JSON(err)
			})
			xerror.Panic(view.Error())

			hd := r.Options().HdlrWrappers
			for i := len(hd); i > 0; i-- {
				handler = hd[i-1](handler)
			}

			// 处理 metadata fastHttpRequest
			headers := make(metadata.Metadata)
			view.Fasthttp.Request.Header.VisitAll(func(key, value []byte) {
				headers[strings.ToLower(string(key))] = string(value)
			})

			serviceName, mthName := apiRoute(string(view.Fasthttp.RequestURI()))
			request := &httpRequest{
				service:     fmt.Sprintf("%s.%s", r.Options().Name, serviceName),
				contentType: defaultContentType,
				method:      mthName,
				body:        view.Fasthttp.Request.Body(),
				header:      headers,
			}

			xerror.Panic(handler(context.WithValue(view.Context(), fastHttpRequest{}, view), request, view))
		})
	}

	return
}
