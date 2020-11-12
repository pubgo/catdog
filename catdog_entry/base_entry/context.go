package base_entry

import (
	"context"
	"github.com/gofiber/fiber"
)

type fastHttpRequest struct{}

func HttpRequestFromCtx(ctx context.Context) *fiber.Ctx {
	ret := ctx.Value(fastHttpRequest{})
	if ret == nil {
		return nil
	}
	return ret.(*fiber.Ctx)
}
