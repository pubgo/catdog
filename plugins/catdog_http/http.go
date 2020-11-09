package catdog_http

import (
	"context"
	"github.com/pubgo/catdog/plugins/catdog_tracing/tracing"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

	"github.com/gojektech/heimdall/hystrix"
)

const (
	MaxConcurrentHTTPRequest = 5000
	DefaultTimeout           = time.Second * 2
)

type httpClient struct {
	client *http.Client
	ctx    context.Context
	name   string
}

func setTag(span opentracing.Span, req *http.Request, rsp *http.Response, name string) {
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())

	var code uint16 = 500
	if rsp != nil {
		code = uint16(rsp.StatusCode)
	}
	ext.HTTPStatusCode.Set(span, code)
	span.SetTag("http-client-name", name)
	span.SetTag(tracing.UniqueIDKey, req.Header.Get("unique_id"))
}

func (c *httpClient) Do(req *http.Request) (rsp *http.Response, err error) {
	if req.Header == nil {
		req.Header = make(http.Header)
	}

	uniqueID := tracing.GetRequestIDFromContext(c.ctx)
	req.Header.Add(tracing.UniqueIDKey, uniqueID)

	_, span, err := tracing.StartSpanFromContext(c.ctx, opentracing.GlobalTracer(), "HTTP.Client")
	if err != nil {
		// Maybe there will be many logs, annotate it.
		// log.Warn("[httpClient.Do] start span error. ", err)
		_ = err
	}

	// Didn't handler the error. So need to compare with nil.
	if span != nil {
		defer func() {
			// need rsp info
			setTag(span, req, rsp, c.name)
			tracing.SetIfError(span, err)
			span.Finish()
		}()
	}

	return c.client.Do(req)
}

// New 返回一个带 tracing 的 client 实例
// 超时时间从 ctx 读取，hystrix.Option 的超时时间对返回的 client 不生效
func New(ctx context.Context, opts ...hystrix.Option) *hystrix.Client {
	timeout := DefaultTimeout
	if ctx != nil {
		if deadline, ok := ctx.Deadline(); ok {
			timeout = time.Until(deadline)
		}
	}

	name := "http.client"
	realc := &httpClient{
		ctx:  ctx,
		name: name,
		client: &http.Client{
			Timeout: timeout,
		},
	}

	var options []hystrix.Option
	options = append(options,
		hystrix.WithCommandName(name),
		hystrix.WithHTTPClient(realc),
		hystrix.WithMaxConcurrentRequests(MaxConcurrentHTTPRequest),
	)
	options = append(options, opts...)
	c := hystrix.NewClient(options...)

	return c
}
