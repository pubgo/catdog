package catdog_mongo

import (
	"context"
	"errors"
	"github.com/pubgo/catdog/plugins/catdog_tracing/tracing"
	"github.com/pubgo/xlog"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	traceLog "github.com/opentracing/opentracing-go/log"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const default_timeout = time.Second * 5

var initInjector = sync.Map{}

func GetMongoClient(prefix string) (*mongo.Client, error) {
	_, loaded := initInjector.LoadOrStore(prefix, 1)
	if !loaded {
		if err := injectMonitor(prefix); err != nil {
			return nil, err
		}
	}

	clientAndOpts, err := PickupMongoClient(prefix)
	if err != nil {
		return nil, err
	}

	return clientAndOpts.Client, nil
}

func injectMonitor(prefix string) error {
	// Inject catdog_mongo_plugin tracing
	opts := options.Client()
	opts.SetMonitor(&event.CommandMonitor{
		Started:   commandStarted,
		Succeeded: commandSucceeded,
		Failed:    commandFailed,
	})
	// catdog_mongo_plugin pool event, open in dev mod
	// opts.SetPoolMonitor(&event.PoolMonitor{
	// 	Event: poolEvent,
	// })
	opts.SetConnectTimeout(default_timeout)
	opts.SetServerSelectionTimeout(default_timeout)

	if err := RefreshClient(prefix, opts); err != nil {
		xlog.Errorf("Tracing.catdog_mongo_plugin refresh error:%s", err)
		return err
	}
	return nil
}

func commandStarted(ctx context.Context, e *event.CommandStartedEvent) {
	ctx, span, _ := tracing.StartSpanFromContext(ctx, opentracing.GlobalTracer(), "Mongo.RunCommand.Started")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "MongoDB")
			span.SetTag("ConnectionID", e.ConnectionID)
			span.SetTag("RequestID", e.RequestID)
			span.SetTag("RunCommand.name", e.CommandName)
			span.SetTag("DatabaseName", e.DatabaseName)

			span.LogFields(traceLog.String("RunCommand", e.Command.String()))
			span.Finish()
		}()
	}
}

func commandSucceeded(ctx context.Context, e *event.CommandSucceededEvent) {
	_, span, _ := tracing.StartSpanFromContext(ctx, opentracing.GlobalTracer(), "Mongo.RunCommand.Succeeded")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "MongoDB")
			span.SetTag("ConnectionID", e.ConnectionID)
			span.SetTag("RequestID", e.RequestID)
			span.SetTag("RunCommand.name", e.CommandName)
			span.SetTag("DurationNanos", e.DurationNanos)

			span.LogFields(traceLog.String("Reply", e.Reply.String()))
			span.Finish()
		}()
	}
}

func commandFailed(ctx context.Context, e *event.CommandFailedEvent) {
	_, span, _ := tracing.StartSpanFromContext(ctx, opentracing.GlobalTracer(), "Mongo.RunCommand.Failed")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "MongoDB")
			span.SetTag("ConnectionID", e.ConnectionID)
			span.SetTag("RequestID", e.RequestID)
			span.SetTag("RunCommand.name", e.CommandName)
			span.SetTag("DurationNanos", e.DurationNanos)
			tracing.SetIfError(span, errors.New(e.Failure))
			span.Finish()
		}()
	}
}

func poolEvent(e *event.PoolEvent) {
	span := opentracing.GlobalTracer().StartSpan("Mongo.Pool.Event")
	if span != nil {
		ext.SpanKindRPCClient.Set(span)
		span.SetTag("DBType", "MongoDB")
		span.SetTag("ConnectionID", e.ConnectionID)
		span.SetTag("Address", e.Address)
		span.SetTag("EventType", e.Type)
		span.SetTag("EventReason", e.Reason)
		if e.PoolOptions != nil {
			span.SetTag("Options.MaxPoolSize", e.PoolOptions.MaxPoolSize)
			span.SetTag("Options.MinPoolSize", e.PoolOptions.MinPoolSize)
			span.SetTag("Options.WaitQueueTimeoutMS", e.PoolOptions.WaitQueueTimeoutMS)
		}

		span.Finish()
	}
}
