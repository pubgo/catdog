package base_entry

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type fastHttpRequest struct{}

func RequestFromCtx(ctx context.Context) *fiber.Ctx {
	ret := ctx.Value(fastHttpRequest{})
	if ret == nil {
		return nil
	}
	return ret.(*fiber.Ctx)
}

func RequestBackFromCtx(ctx context.Context, fn func(*fiber.Ctx)) {
	ret := ctx.Value(fastHttpRequest{})
	if ret == nil {
		return
	}
	fn(ret.(*fiber.Ctx))
}
