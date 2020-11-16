package catdog_request

import (
	"context"
	"fmt"
	"strings"

	"github.com/asim/nitro/v3/metadata"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber/v2"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/xerror"
)

const defaultContentType = "application/json"

func HttpMiddleware(fn server.HandlerFunc) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		defer xerror.Resp(func(err xerror.XErr) { ctx.Next(err) })

		headers := make(metadata.Metadata)
		ctx.Fasthttp.Request.Header.VisitAll(func(key, value []byte) {
			headers[strings.ToLower(string(key))] = string(value)
		})

		s, m := apiRoute(ctx.OriginalURL())
		request := &Request{
			service:     fmt.Sprintf("%s.%s", catdog_config.Project, s),
			contentType: defaultContentType,
			method:      m,
			body:        ctx.Fasthttp.Request.Body(),
			header:      headers,
		}

		ctx.Next(fn(ctx.Fasthttp, request, ctx))
	}
}

func HttpMiddleware1(fn server.HandlerWrapper) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		defer xerror.Resp(func(err xerror.XErr) { ctx.Next(err) })

		headers := make(metadata.Metadata)
		ctx.Fasthttp.Request.Header.VisitAll(func(key, value []byte) {
			headers[strings.ToLower(string(key))] = string(value)
		})

		s, m := apiRoute(ctx.OriginalURL())
		request := &Request{
			service:     fmt.Sprintf("%s.%s", catdog_config.Project, s),
			contentType: defaultContentType,
			method:      m,
			body:        ctx.Fasthttp.Request.Body(),
			header:      headers,
		}

		xerror.Panic(fn(func(ctx1 context.Context, req server.Request, rsp interface{}) (err error) {
			defer xerror.RespErr(&err)
			rsp.(*fiber.Ctx).Next()
			return nil
		})(ctx.Fasthttp, request, ctx))
	}
}
