package tracing

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"

	"github.com/micro/go-micro/v3/metadata"
	"github.com/opentracing/opentracing-go"
)

type referHandler func(sc opentracing.SpanContext) opentracing.SpanReference

/*
X-W-Traceid
X-W-Spanid
X-W-ParentspanId
X-W-Sampled
*/
const (
	PHPRequestTraceID      = "x-w-traceid"
	PHPRequestSpanID       = "x-w-spanid"
	PHPRequestParentSpanID = "x-w-parentspanid"
	PHPRequestSampleID     = "x-w-sampled"

	// error keys
	KeyErrorMessage        = "error_msg"
	KeyContextErrorMessage = "context_error_msg"
)

var (
	ErrorContextNil   = errors.New("got a nil context. ")
	ErrorNotFoundSpan = errors.New("not found span in context. ")
)

// spanFromContext 从 context 生成 span
func spanFromContext(
	handler referHandler,
	ctx context.Context,
	tracer opentracing.Tracer,
	name string,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span, error) {
	parentSpanCtx, err := GetParentSpanContext(ctx, tracer)
	if err != nil {
		return ctx, nil, err
	}

	opts = append(opts, handler(parentSpanCtx))
	sp := tracer.StartSpan(name, opts...)

	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	// copy the metadata to prevent race
	md = metadata.Copy(md)
	if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(md)); err != nil {
		return ctx, nil, err
	}

	ctx = opentracing.ContextWithSpan(ctx, sp)
	ctx = metadata.NewContext(ctx, md)
	return ctx, sp, nil
}

func spanFromPHPRequest(md metadata.Metadata) (jaeger.SpanContext, error) {
	var emptyContext jaeger.SpanContext
	var err error

	sample, ok := md[PHPRequestSampleID]
	if !ok {
		return emptyContext, errors.New("no tracing info. ")
	}

	n, err := strconv.Atoi(sample)
	if err != nil {
		return emptyContext, err
	}
	isSampled := n == 1
	if !isSampled {
		return emptyContext, errors.New("sample is not allowed. ")
	}

	trace, ok := md[PHPRequestTraceID]
	if !ok {
		return emptyContext, errors.New("didn't find trace ID. ")
	}
	traceID, err := jaeger.TraceIDFromString(trace)
	if err != nil {
		return emptyContext, err
	}

	span, ok := md[PHPRequestSpanID]
	if !ok {
		span = trace
	}
	spanID, err := jaeger.SpanIDFromString(span)
	if err != nil {
		return emptyContext, err
	}

	parentSpan, ok := md[PHPRequestParentSpanID]
	if !ok {
		parentSpan = span
	}
	parentSpanID, err := jaeger.SpanIDFromString(parentSpan)
	if err != nil {
		return emptyContext, err
	}

	ctx := jaeger.NewSpanContext(traceID, spanID, parentSpanID, isSampled, nil)
	return ctx, nil
}

// StartSpanFromContext returns a new span with the given operation name and options. If a span
// is found in the context, it will be used as the parent of the resulting span.
func StartSpanFromContext(
	ctx context.Context,
	tracer opentracing.Tracer,
	name string,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span, error) {

	return spanFromContext(opentracing.ChildOf, ctx, tracer, name, opts...)
}

// StartFollowSpanFromContext returns a new span with the given operation name and options. If a span
// is found in the context, it will be used as the parent of the resulting span.
func StartFollowSpanFromContext(
	ctx context.Context,
	tracer opentracing.Tracer,
	name string,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span, error) {

	return spanFromContext(opentracing.FollowsFrom, ctx, tracer, name, opts...)
}

func StartSpanFromSpanContext(
	sctx opentracing.SpanContext,
	tracer opentracing.Tracer,
	name string,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span, error) {
	opts = append(opts, opentracing.ChildOf(sctx))
	sp := tracer.StartSpan(name, opts...)
	return opentracing.ContextWithSpan(context.TODO(), sp), sp, nil
}

func StartFollowSpanFromSpanContext(
	sctx opentracing.SpanContext,
	tracer opentracing.Tracer,
	name string,
	opts ...opentracing.StartSpanOption,
) (context.Context, opentracing.Span, error) {
	opts = append(opts, opentracing.FollowsFrom(sctx))
	sp := tracer.StartSpan(name, opts...)
	return opentracing.ContextWithSpan(context.TODO(), sp), sp, nil
}

// GetParentSpanContext get parent span-context from context
func GetParentSpanContext(
	ctx context.Context,
	tracer opentracing.Tracer,
) (opentracing.SpanContext, error) {
	var parentSpanCtx opentracing.SpanContext

	if ctx == nil {
		return parentSpanCtx, ErrorContextNil
	}

	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = make(map[string]string)
	}

	// copy the metadata to prevent race
	md = metadata.Copy(md)

	if tracer == nil {
		tracer = opentracing.GlobalTracer()
	}

	var errRet error
	// Find parent span.
	// First try to get span within current service boundary.
	// If there doesn't exist, try to get it from go-micro metadata(which is cross boundary)
	if span := opentracing.SpanFromContext(ctx); span != nil {
		parentSpanCtx = span.Context()
	} else if spanCtx, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(md)); err == nil {
		parentSpanCtx = spanCtx
	} else if reqCtx, err := spanFromPHPRequest(md); err == nil {
		parentSpanCtx = reqCtx
	} else {
		errRet = ErrorNotFoundSpan
	}

	return parentSpanCtx, errRet
}

type TraceInfo struct {
	SpanID    string
	TraceID   string
	ParentID  string
	IsSampled bool
}

// GetFullTraceInfo parse the trace from context
func GetFullTraceInfo(ctx context.Context, tracer opentracing.Tracer) (*TraceInfo, error) {
	parentSpanCtx, err := GetParentSpanContext(ctx, tracer)

	if err != nil {
		return nil, err
	}

	jSpanCtx, err := jaeger.ContextFromString(fmt.Sprintf("%v", parentSpanCtx))
	if err != nil {
		return nil, err
	}

	ret := &TraceInfo{
		SpanID:    jSpanCtx.SpanID().String(),
		TraceID:   jSpanCtx.TraceID().String(),
		ParentID:  jSpanCtx.ParentID().String(),
		IsSampled: jSpanCtx.IsSampled(),
	}

	return ret, nil
}

// GetTraceID parse the trace id from context
func GetTraceID(ctx context.Context, tracer opentracing.Tracer) (string, error) {
	traceInfo, err := GetFullTraceInfo(ctx, tracer)
	if err != nil {
		return "", err
	}

	return traceInfo.TraceID, nil
}

// SetIfError add error info and flag for error
func SetIfError(span opentracing.Span, err error) {
	if span == nil {
		return
	}

	if err != nil {
		ext.Error.Set(span, true)
		span.SetTag(KeyErrorMessage, err.Error())
	}
}

// SetIfContextError record error
func SetIfContextError(span opentracing.Span, ctx context.Context) {
	if span == nil {
		return
	}

	if ctx != nil {
		if err := ctx.Err(); err != nil {
			ext.Error.Set(span, true)
			span.SetTag(KeyContextErrorMessage, err.Error())
		}
	}
}
