package catdog_tracing_plugin

import (
	"context"
	"fmt"
	"github.com/pubgo/catdog/internal/tracing"

	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/server"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"

)

// HandlerWrap is a handler wrapper, look at micro's handler wrapper option
func HandlerWrap(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) (err error) {
		// create
		ctx = NewContextWithOld(ctx)

		name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
		ctx, span, err := tracing.StartSpanFromContext(ctx, opentracing.GlobalTracer(), name)
		if err != nil {
			// Maybe there will be many logs, annotate it.
			// log.Warn("[HandlerWrap] start span error. ", err)
			_ = err
		}

		// Didn't handler the error. So need to compare with nil.
		if span != nil {
			defer func() {
				setHandlerSpanTag(span, req)
				// handler wrap means the remote caller must record the parent span.
				// So we just need to record when remote is nil.
				tracing.SetIfError(span, err) // remote error
				// If context was canceled or timeout
				// tracing.SetIfContextError(span, ctx)
				span.Finish()
			}()
		}

		return fn(ctx, req, rsp)
	}
}

// setHandlerSpanTag a child span which saved the necessary tags for server handler
func setHandlerSpanTag(span opentracing.Span, req server.Request) {
	ext.SpanKindRPCServer.Set(span)
	ext.PeerHostname.Set(span, req.Service())
	ext.PeerAddress.Set(span, req.Endpoint())
}

// =======================================================================================================
// tracingWrapper client tracing
type tracingWrapper struct {
	client.Client
}

// Call tracing wrapper implement the interface of micro's client wrapper
func (tw *tracingWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) (err error) {
	name := fmt.Sprintf("%s.%s", req.Service(), req.Endpoint())
	ctx, span, err := tracing.StartSpanFromContext(ctx, opentracing.GlobalTracer(), name)
	if err != nil {
		// Maybe there will be many logs, annotate it.
		// log.Warn("[WrapperClient.Call] start span error. ", err)
		_ = err
	}

	// Didn't handler the error. So need to compare with nil.
	if span != nil {
		defer func() {
			setClientSpanTag(span, req)
			// Call wrap means the develop must use tracing.WrapChild func to record the call.
			// So we just need to record when caller is nil.
			tracing.SetIfError(span, err) // remote error
			// tracing.SetIfContextError(span, ctx) // If context was canceled or timeout
			span.Finish()
		}()
	}

	return tw.Client.Call(ctx, req, rsp, opts...)
}

// ClientWrap is a client wrapper func, look at the micro's client wrapper option
func ClientWrap(c client.Client) client.Client {
	return &tracingWrapper{c}
}

// setClientSpanTag return a child span which saved the necessary tags for client caller
func setClientSpanTag(span opentracing.Span, req client.Request) {
	ext.SpanKindRPCClient.Set(span)
	ext.PeerHostname.Set(span, req.Service())
	ext.PeerAddress.Set(span, req.Endpoint())
}
