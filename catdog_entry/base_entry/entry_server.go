package base_entry

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

type entryServerWrapper struct {
	server.Server
	router fiber.Router
}

func (t *entryServerWrapper) httpHandler(httpMethod, path string, handlers ...fiber.Handler) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if t.router != nil {
		t.router.Add(strings.ToUpper(httpMethod), path, handlers...)
	}

	return nil
}

func (t *entryServerWrapper) mthHandle(ctx context.Context, req server.Request, rsp interface{}) (err error) {
	defer xerror.RespErr(&err)

	mth := ctx.Value("mth").(reflect.Value)
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

	xerror.Panic(t.Server.Handle(handler))

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

		xerror.Panic(t.httpHandler(httpMethod, httpPath, func(view *fiber.Ctx) {
			defer xerror.Resp(func(err xerror.XErr) {
				_ = view.Status(http.StatusInternalServerError).JSON(err)
			})
			xerror.Panic(view.Error())

			// 处理 metadata fastHttpRequest
			headers := make(metadata.Metadata)
			view.Fasthttp.Request.Header.VisitAll(func(key, value []byte) {
				headers[strings.ToLower(string(key))] = string(value)
			})

			request := &httpRequest{
				service:     fmt.Sprintf("%s.%s", t.Options().Name, handler.Name()),
				contentType: defaultContentType,
				method:      mthName,
				body:        view.Fasthttp.Request.Body(),
				header:      headers,
			}

			handle := t.mthHandle
			hd := t.Server.Options().HdlrWrappers
			for i := len(hd); i > 0; i-- {
				handle = hd[i-1](handle)
			}

			ctx := context.WithValue(view.Context(), fastHttpRequest{}, view)
			ctx = context.WithValue(ctx, "mth", hdlr.MethodByName(mthName))
			xerror.Panic(handle(ctx, request, view))
		}))
	}

	return
}
