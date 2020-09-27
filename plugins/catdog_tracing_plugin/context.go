package catdog_tracing_plugin

import (
	"context"
	"github.com/pubgo/catdog/internal/tracing"
)

// NewContextWithOld 为请求生成不依赖于父一级服务超时设置的请求
func NewContextWithOld(ctx context.Context) context.Context {
	return tracing.NewContextWithOld(ctx)
}

// GetRequestIDFromContext return request ID or call it unique_id
func GetRequestIDFromContext(ctx context.Context) string {
	return tracing.GetRequestIDFromContext(ctx)
}
