package ws_entry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asim/nitro/v3/metadata"
	"github.com/asim/nitro/v3/server"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
	"github.com/valyala/fasthttp"
	"net/http"
	"reflect"
	"strings"
)

type fastHttpRequest struct{}

type mth struct{}

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

type entryServerWrapper struct {
	server.Server
	router fiber.Router
}

func (t *entryServerWrapper) Start() (err error) {
	defer xerror.RespErr(&err)
	return nil
}

func (t *entryServerWrapper) Stop() (err error) {
	defer xerror.RespErr(&err)
	return nil
}

func (t *entryServerWrapper) wsHandle(ctx context.Context, _ server.Request, rsp interface{}) (err error) {
	defer xerror.RespErr(&err)

	var upgrade = websocket.FastHTTPUpgrader{
		HandshakeTimeout:  0,
		Subprotocols:      nil,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		EnableCompression: true,
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			return true
		},
	}

	mth := ctx.Value(mth{}).(reflect.Value)
	mthInType := mth.Type().In(1)
	mthOutType := mth.Type().In(2)
	view := rsp.(*fiber.Ctx)
	err = upgrade.Upgrade(view.Context(), func(conn *websocket.Conn) {
		//defer conn.Close()
		//conn.SetReadLimit(maxMessageSize)
		//conn.SetReadDeadline(time.Now().Add(pongWait))
		//conn.SetWriteDeadline(time.Now().Add(writeWait))
		//conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		c := rsp.(*websocket.Conn)
		cancel := xprocess.GoLoop(func(_ context.Context) (err error) {
			defer xerror.RespErr(&err)
			mt, msg, err := c.ReadMessage()
			xerror.Panic(err)

			mthIn := reflect.New(mthInType.Elem())
			ret := reflect.ValueOf(json.Unmarshal).Call([]reflect.Value{reflect.ValueOf(msg), mthIn})
			if !ret[0].IsNil() {
				return xerror.Wrap(ret[0].Interface().(error))
			}

			mthOut := reflect.New(mthOutType.Elem())
			ret = mth.Call([]reflect.Value{reflect.ValueOf(ctx), mthIn, mthOut})
			if !ret[0].IsNil() {
				return xerror.Wrap(ret[0].Interface().(error))
			}

			dt, _err := json.Marshal(mthOut.Interface())
			if err != nil {
				return xerror.Wrap(_err)
			}

			return xerror.Wrap(c.WriteMessage(mt, dt))
		})

		c.SetCloseHandler(func(code int, text string) error {
			xlog.Debugf("%d, %s", code, text)
			return xerror.Wrap(cancel())
		})
	})

	if err != nil {
		if err == websocket.ErrBadHandshake {
			log.Errorf("%#v", err)
		}
		return
	}

	return nil
}

func (t *entryServerWrapper) wsHandler(httpMethod, path string, handlers fiber.Handler) error {
	httpMethod = strings.ToUpper(httpMethod)
	if _, ok := httpMethods[httpMethod]; !ok {
		return nil
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	t.router.Add(httpMethod, path, func(ctx *fiber.Ctx) error {
		return xerror.Wrap(handlers(ctx))
	})
	return nil
}

const defaultContentType = "application/json"

func (t *entryServerWrapper) Handle(handler server.Handler) (err error) {
	if handler == nil {
		return xerror.New("[handler] should not be nil")
	}

	defer xerror.RespErr(&err)

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

		xerror.Panic(t.wsHandler(httpMethod, httpPath, func(view *fiber.Ctx) (err error) {
			defer xerror.RespErr(&err)

			// 处理 metadata fastHttpRequest
			headers := make(metadata.Metadata)
			view.Request().Header.VisitAll(func(key, value []byte) {
				headers[strings.ToLower(string(key))] = string(value)
			})

			request := &httpRequest{
				service:     fmt.Sprintf("%s.%s", t.Options().Name, handler.Name()),
				contentType: defaultContentType,
				method:      mthName,
				header:      headers,
			}

			handle := t.wsHandle
			hd := t.Server.Options().HdlrWrappers
			for i := len(hd); i > 0; i-- {
				handle = hd[i-1](handle)
			}

			ctx := context.WithValue(view.Context(), fastHttpRequest{}, view)
			ctx = context.WithValue(ctx, mth{}, hdlr.MethodByName(mthName))
			xerror.Panic(handle(ctx, request, ctx))
			return nil
		}))
	}

	return
}
