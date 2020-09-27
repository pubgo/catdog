package tracing

import (
	"context"
	"time"

	"github.com/micro/go-micro/v3/metadata"
)

const (
	UniqueIDKey = "unique_id"
)

// NewContextWithOld 为请求生成不依赖于父一级服务超时设置的请求
func NewContextWithOld(ctx context.Context) context.Context {
	// nil
	if ctx == nil {
		return nil
	}

	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(metadata.Metadata)
	}
	md = metadata.Copy(md)

	var deadline time.Time
	if d, ok := ctx.Deadline(); ok {
		deadline = d
	} else {
		deadline = time.Now().Add(time.Second * 2)
	}

	nctx, cancel := context.WithDeadline(context.Background(), deadline)
	_ = cancel
	nctx = metadata.NewContext(nctx, md)

	return nctx
}

// GetRequestIDFromContext return request ID or call it unique_id
func GetRequestIDFromContext(ctx context.Context) string {
	var uniqueID string

	// nil
	if ctx == nil {
		return uniqueID
	}

	var md, ok = metadata.FromContext(ctx)
	if ok {
		if v, suc := md[UniqueIDKey]; suc {
			uniqueID = v
		}
	}

	return uniqueID
}
