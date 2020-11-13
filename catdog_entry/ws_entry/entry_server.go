package ws_entry

import (
	"context"
	"fmt"
	"github.com/asim/nitro/v3/metadata"
	"github.com/asim/nitro/v3/server"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
	"github.com/valyala/fasthttp"
	"net/http"
	"reflect"
	"strings"
)

type entryServerWrapper struct {
	server.Server
	router fiber.Router
}

func (t *entryServerWrapper) Start() (err error) {
	defer xerror.RespErr(&err)
	return nil
}

func (t *entryServerWrapper) wsHandle(ctx context.Context, _ server.Request, rsp interface{}) (err error) {
	defer xerror.RespErr(&err)

	var upgrader = websocket.FastHTTPUpgrader{
		HandshakeTimeout:  cfg.HandshakeTimeout,
		Subprotocols:      cfg.Subprotocols,
		ReadBufferSize:    cfg.ReadBufferSize,
		WriteBufferSize:   cfg.WriteBufferSize,
		EnableCompression: cfg.EnableCompression,
		CheckOrigin: func(fctx *fasthttp.RequestCtx) bool {
			if cfg.Origins[0] == "*" {
				return true
			}
			origin := utils.GetString(fctx.Request.Header.Peek("Origin"))
			for i := range cfg.Origins {
				if cfg.Origins[i] == origin {
					return true
				}
			}
			return false
		},
	}

	view := rsp.(*fiber.Ctx)
	xerror.Panic(upgrader.Upgrade(view.Context(), func(conn *websocket.Conn) {
		c := rsp.(*websocket.Conn)
		cancel := xprocess.GoLoop(func(_ context.Context) (err error) {
			defer xerror.RespErr(&err)
			mt, msg, err := c.ReadMessage()
			xerror.Panic(err)
			return xerror.Wrap(c.WriteMessage(mt, msg))
		})

		c.SetCloseHandler(func(code int, text string) error {
			xlog.Debugf("%d, %s", code, text)
			return xerror.Wrap(cancel())
		})
	}))

	return nil
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

func (t *entryServerWrapper) Stop() (err error) {
	defer xerror.RespErr(&err)
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

			var upgrader = websocket.FastHTTPUpgrader{
				HandshakeTimeout:  cfg.HandshakeTimeout,
				Subprotocols:      cfg.Subprotocols,
				ReadBufferSize:    cfg.ReadBufferSize,
				WriteBufferSize:   cfg.WriteBufferSize,
				EnableCompression: cfg.EnableCompression,
				CheckOrigin: func(fctx *fasthttp.RequestCtx) bool {
					if cfg.Origins[0] == "*" {
						return true
					}
					origin := utils.GetString(fctx.Request.Header.Peek("Origin"))
					for i := range cfg.Origins {
						if cfg.Origins[i] == origin {
							return true
						}
					}
					return false
				},
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

type fastHttpRequest struct{}
type mth struct{}
