package redis

import (
	"context"
	"github.com/pubgo/catdog/internal/tracing"
	"github.com/pubgo/catdog/plugins/catdog_redis"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	DbType                  = "redis"
	SpanKind                = ext.SpanKindEnum("redis-client")
	MaxPipelineNameCmdCount = 3
	DefaultRWTimeout        = time.Second
)

type Option func(options *redis.Options)

// WithReadTimeout 设置读超时
func WithReadTimeout(d time.Duration) Option {
	return func(o *redis.Options) {
		o.ReadTimeout = d
	}
}

// WithReadTimeout 设置写超时
func WithWriteTimeout(d time.Duration) Option {
	return func(o *redis.Options) {
		o.WriteTimeout = d
	}
}

// GetRedis get a redis client with open tracing context
func GetRedis(ctx context.Context, prefix string, options ...Option) (*redis.Client, error) {
	client, err := catdog_redis.PickupRedisClient(prefix)
	if err != nil {
		return nil, err
	}

	cc := client.WithContext(ctx)
	opts := cc.Options()

	// 默认的读写超时时间为 1s
	opts.WriteTimeout = DefaultRWTimeout
	opts.ReadTimeout = DefaultRWTimeout

	// 处理外部进来的参数配置
	for _, o := range options {
		o(opts)
	}

	cc.WrapProcess(wrapper(ctx, opts))
	cc.WrapProcessPipeline(wrapperPipeline(ctx, opts))

	return cc, nil
}

// =============================================================================================================
// Wrap the redis command process func
// GoRedisProcessFunc is an alias of cmd process func
type GoRedisProcessFunc = func(cmd redis.Cmder) error

// GoRedisWrapProcessFunc is an alias of wrapper that wrap process
type GoRedisWrapProcessFunc = func(oldProcess GoRedisProcessFunc) GoRedisProcessFunc

func setTag(span opentracing.Span, opts *redis.Options, method, key string) {
	ext.DBType.Set(span, DbType)
	ext.PeerAddress.Set(span, opts.Addr)
	ext.SpanKind.Set(span, SpanKind)

	// add redis command
	span.SetTag("db.method", method)
	span.SetTag("db.key", key)
}

func wrapper(ctx context.Context, opts *redis.Options) GoRedisWrapProcessFunc {
	return func(oldProcess GoRedisProcessFunc) GoRedisProcessFunc {
		return func(cmd redis.Cmder) (err error) {
			_, span, err := tracing.StartSpanFromContext(ctx, opentracing.GlobalTracer(), "Redis.Client")
			if err != nil {
				// Maybe there will be many logs, annotate it.
				// log.Warn("[wrapper] start span error. ", err)
				_ = err
			}

			// Didn't handler the error. So need to compare with nil.
			if span != nil {
				defer func() {
					method, key := wrapperName(cmd)
					setTag(span, opts, method, key)
					tracing.SetIfError(span, err)
					span.Finish()
				}()
			}

			return oldProcess(cmd)
		}
	}
}

func wrapperName(cmd redis.Cmder) (name string, key string) {
	name = cmd.Name()
	args := cmd.Args()
	if len(args) > 1 {
		k, ok := args[1].(string)
		if ok {
			key = k
		}
	}
	return
}

// =============================================================================================================
// Wrap the redis command process pipeline func
// GoRedisProcessPipelineFunc is an alias of process pipeline func
type GoRedisProcessPipelineFunc = func([]redis.Cmder) error

// GoRedisWrapProcessPipelineFunc is an alias of wrapper that wrap pipeline
type GoRedisWrapProcessPipelineFunc = func(oldProcess GoRedisProcessPipelineFunc) GoRedisProcessPipelineFunc

func wrapperPipeline(ctx context.Context, opts *redis.Options) GoRedisWrapProcessPipelineFunc {
	f := func(oldProcess GoRedisProcessPipelineFunc) GoRedisProcessPipelineFunc {
		return func(cmders []redis.Cmder) (err error) {
			_, span, err := tracing.StartSpanFromContext(ctx, opentracing.GlobalTracer(), "Redis.Client.Pipeline")
			if err != nil {
				// Maybe there will be many logs, annotate it.
				// log.Warn("[wrapperPipeline] start span error. ", err)
				_ = err
			}

			// Didn't handler the error. So need to compare with nil.
			if span != nil {
				defer func() {
					method, key := wrapperPipelineName(cmders)
					setTag(span, opts, method, key)
					tracing.SetIfError(span, err)
					span.Finish()
				}()
			}

			return oldProcess(cmders)
		}
	}

	return f
}

func wrapperPipelineName(cmders []redis.Cmder) (name string, key string) {
	names := make([]string, len(cmders))
	keys := make([]string, len(cmders))
	for i, cmd := range cmders {
		// keep MaxPipelineNameCmdCount cmd name
		if i >= MaxPipelineNameCmdCount {
			break
		}
		names[i], keys[i] = wrapperName(cmd)
	}

	name = strings.Join(names, "->")
	key = strings.Join(keys, "->")

	return
}
