package catdog_tracing_plugin

import (
	"context"
	"fmt"
	"github.com/pubgo/catdog/internal/tracing/http"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/codec"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/log"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport"
)

func init() {
	tracer, _ := jaeger.NewTracer(
		"micro.plugins.test",
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(transport.NewHTTPTransport("http://10.3.7.11:14268/api/traces?format=jaeger.thrift")))

	opentracing.SetGlobalTracer(tracer)
}

func normalHandler(ctx context.Context, req server.Request, rsp interface{}) error {
	id := GetRequestIDFromContext(ctx)
	log.Info("request_id ==> ", id)
	c := http.New(ctx)
	r, err := c.Get("http://w.weipaitang.com", nil)
	if err != nil {
		return err
	}

	if r != nil && r.Body != nil {
		defer func() {
			_ = r.Body.Close()
		}()
	}

	return nil
}

type req struct {
}

// Service name requested
func (r *req) Service() string {
	return "test"
}

// The action requested
func (r *req) Method() string {
	return "test-api"
}

// Endpoint name requested
func (r *req) Endpoint() string {
	return "123.123.123.123"
}

// Content type provided
func (r *req) ContentType() string {
	return "application/json"
}

// Header of the request
func (r *req) Header() map[string]string {
	return nil
}

// Body is the initial decoded value
func (r *req) Body() interface{} {
	return nil
}

// Read the undecoded request body
func (r *req) Read() ([]byte, error) {
	return nil, nil
}

// The encoded message stream
func (r *req) Codec() codec.Reader {
	return nil
}

// Indicates whether its a stream
func (r *req) Stream() bool {
	return false
}

func TestHandlerWrap(t *testing.T) {
	span := opentracing.StartSpan("/test/handler/wrapper")
	spanContext := opentracing.ContextWithSpan(context.Background(), span)

	md := make(metadata.Metadata)
	md[jaeger.TraceContextHeaderName] = fmt.Sprintf("%s", span)
	traceIDContext := metadata.NewContext(context.Background(), md)
	defer func() {
		span.Finish()
		time.Sleep(time.Second * 2)
	}()

	type args struct {
		fn  server.HandlerFunc
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Default Context",
			args: args{
				fn:  normalHandler,
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Remote Span Context",
			args: args{
				fn:  normalHandler,
				ctx: spanContext,
			},
			wantErr: false,
		},
		{
			name: "Span From Trace ID",
			args: args{
				fn:  normalHandler,
				ctx: traceIDContext,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HandlerWrap(tt.args.fn)
			if got == nil {
				t.Errorf("HandlerWrap() = %v ", got)
				return
			}

			if err := got(tt.args.ctx, &req{}, nil); err != nil {
				t.Error(err)
			}
		})
	}
}

type req2 struct {
}

// The service to call
func (r *req2) Service() string {
	return "client"
}

// The action to take
func (r *req2) Method() string {
	return "POST"
}

// The endpoint to invoke
func (r *req2) Endpoint() string {
	return "localhost"
}

// The content type
func (r *req2) ContentType() string {
	return "application/json"
}

// The unencoded request body
func (r *req2) Body() interface{} {
	return nil
}

// Write to the encoded request writer. This is nil before a call is made
func (r *req2) Codec() codec.Writer {
	return nil
}

// indicates whether the request will be a streaming one rather than unary
func (r *req2) Stream() bool {
	return false
}

func TestClientWrap(t *testing.T) {
	span := opentracing.StartSpan("/test/client/wrapper")
	spanContext := opentracing.ContextWithSpan(context.Background(), span)

	md := make(metadata.Metadata)
	md[jaeger.TraceContextHeaderName] = fmt.Sprintf("%s", span)
	traceIDContext := metadata.NewContext(context.Background(), md)
	defer func() {
		span.Finish()
		time.Sleep(time.Second * 2)
	}()

	type args struct {
		c   client.Client
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Default Context",
			args: args{
				c:   client.NewClient(),
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Remote Span Context",
			args: args{
				c:   client.NewClient(),
				ctx: spanContext,
			},
			wantErr: false,
		},
		{
			name: "Span From Trace ID",
			args: args{
				c:   client.NewClient(),
				ctx: traceIDContext,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClientWrap(tt.args.c)
			if got == nil {
				t.Errorf("ClientWrap() = %v", got)
				return
			}

			_ = got.Call(tt.args.ctx, &req2{}, nil)
		})
	}
}
