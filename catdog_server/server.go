package catdog_server

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/micro/go-micro/v3/api"
	"github.com/micro/go-micro/v3/metadata"
	"github.com/micro/go-micro/v3/server"
	"github.com/pubgo/xerror"
	"net/http"
	"reflect"
	"strings"
)

const defaultContentType = "application/json"

type wrapper struct {
	server.Server
	handlers []func(g fiber.Router) error
}

func (r *wrapper) Handlers() []func(g fiber.Router) error {
	return r.handlers
}

func (r *wrapper) httpHandler(httpMethod, relativePath string, handlers ...fiber.Handler) {
	r.handlers = append(r.handlers, func(router fiber.Router) error {
		if router == nil {
			return xerror.New("please init router group")
		}
		router.Add(httpMethod, relativePath, handlers...)
		return nil
	})
}

func (r *wrapper) Handle(handler server.Handler) (err error) {
	defer xerror.RespErr(&err)
	xerror.Panic(r.Server.Handle(handler))

	hdr := reflect.ValueOf(handler.Handler())
	for _, e := range handler.Endpoints() {
		endpoint := api.Decode(e.Metadata)
		if len(endpoint.Method) == 0 || endpoint.Method[0] == "" || len(endpoint.Path) == 0 || endpoint.Path[0] == "" {
			continue
		}

		mthS := strings.Split(e.Name, ".")
		mth := hdr.MethodByName(mthS[len(mthS)-1])
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

		r.httpHandler(endpoint.Method[0], endpoint.Path[0], func(view *fiber.Ctx) {
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

			ctx := context.WithValue(view.Context(), fastHttpRequest{}, view)
			xerror.Panic(handler(ctx, request, view))
		})
	}

	return
}
