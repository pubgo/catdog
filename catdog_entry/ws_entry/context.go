package ws_entry

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type fastHttpRequest1 struct{}

func RequestFromCtx(ctx context.Context) *fiber.Ctx {
	ret := ctx.Value(fastHttpRequest{})
	if ret == nil {
		return nil
	}
	return ret.(*fiber.Ctx)
}
