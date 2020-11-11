package catdog_rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/pubgo/catdog/plugins/catdog_tracing/tracing"
	"github.com/pubgo/xprocess"
	"github.com/uber/jaeger-client-go"
	"sync"
	"vitess.io/vitess/go/vt/log"

	"github.com/streadway/amqp"
)

type ExchangeKind string

const (
	ExchangeKindFanout ExchangeKind = "fanout"
	ExchangeKindDirect ExchangeKind = "direct"
	ExchangeKindTopic  ExchangeKind = "topic"

	contentType     = "application/json"
	contentEncoding = "UTF-8"

	retryCount             = 2
	DeliveryModeTransient  = 1
	DeliveryModePersistent = 2
	prefetchCount          = 200
)

var initRbmqMap sync.Map

type rbmqwather struct {
	WatchKey string
}

type RabbitChan struct {
	*amqp.Channel
}

/*
 定义exchange
 参数：
	name exchange名称
	kind exchange种类
	durable 是否持久化
	autoDelete 当所有绑定的队列都与交换器解绑后，交换器会自动删除
 返回值：
	error 操作期间产生的错误
*/
// Deprecated: Do not use it! TODO::This function will remove in v0.2.0
func (r *RabbitChan) DeclareExchange(name string, kind ExchangeKind, durable, autoDelete bool) error {
	return r.Channel.ExchangeDeclare(name, string(kind), durable, autoDelete, false, false, nil)
}

/*
 定义队列
 参数：
	name 队列的名称
	durable 是否持久化
	autoDelete 当所有消费者都断开时，队列会自动删除
 返回值：
	error 操作期间产生的错误
*/
// Deprecated: Do not use it! TODO::This function will remove in v0.2.0
func (r *RabbitChan) DeclareQueue(name string, durable, autoDelete bool) (amqp.Queue, error) {
	return r.Channel.QueueDeclare(name, durable, autoDelete, false, false, nil)
}

/*
 exchange和queue绑定
 参数：
	queue 队列名称
	bindkey 绑定的key
	exchange 交换器名称
 返回值：
	error 操作期间产生的错误
*/
// Deprecated: Do not use it! TODO::This function will remove in v0.2.0
func (r *RabbitChan) Bind(queue, bindKey, exchange string) error {
	return r.Channel.QueueBind(queue, bindKey, exchange, false, nil)
}

/*
 发送消息到指定的队列。通过默认exchange把消息发送到指定的消息队列
 参数：
 	queue 队列名称
 	msg 消息内容
 返回值：
	error 操作期间产生的错误
*/
// Deprecated: Do not use it! TODO::This function will remove in v0.2.0
func (r *RabbitChan) Send(queue string, msg []byte, deliveryMode uint8) error {
	return r.Publish("", queue, msg, deliveryMode)
}

/*
 发送消息到指定的exchange
 参数：
	routingKey 路由key
 	msg 消息内容
 返回值：
	error 操作期间产生的错误
*/
// Deprecated: Do not use it! TODO::This function will remove in v0.2.0
func (r *RabbitChan) Publish(exchange, routingKey string, msg []byte, deliveryMode uint8) (err error) {
	defer func() {
		e := r.close()
		if e != nil {
			err = e
		}
		if err != nil {
			log.Error(err)
		}
	}()
	err = r.Channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		ContentType:     contentType,
		ContentEncoding: contentEncoding,
		DeliveryMode:    deliveryMode, // 持久化
		Body:            msg,
	})
	return err
}

/*
 消费消息队列
 参数：
	queue 队列名称
	autoAck	是否自动回复ack
 返回值：
	Delivery 传递消息的单向通道，可以通过读取该通道获取接收到的消息
	error 操作期间产生的错误
*/
// Deprecated: Do not use it! TODO::This function will remove in v0.2.0
func (r *RabbitChan) Consume(queue string, autoAck bool) (<-chan amqp.Delivery, error) {
	defer func() {
		err := r.close()
		if err != nil {
			log.Error(err)
		}
	}()
	return r.Channel.Consume(queue, "", autoAck, false, false, false, nil)
}

/*
 关闭
 返回值：
	error 操作期间产生的错误
*/
// Deprecated: Do not use it! TODO::This function will remove in v0.2.0
func (r *RabbitChan) close() error {
	if r.Channel != nil {
		return r.Channel.Close()
	}

	return nil
}

func (rw *rbmqwather) GetKey() string {
	return rw.WatchKey
}

func (rw *rbmqwather) Watch(key string, value []byte) error {
	err := buildRbmqData(key, value)
	if err != nil {
		log.Error(err)
	}
	return err
}

func (rw *rbmqwather) Name() string {
	return "catdog_rabbitmq_plugin"
}

// 组装转换数据
func buildRbmqData(prefix string, value []byte) error {
	if len(value) == 0 {
		return errors.New("watch时读取json出错,prefix= " + prefix)
	}

	var conf *rabbitConfig
	var err error
	if conf, err = Parse(value); err != nil {
		return err
	}

	if err = store(prefix, conf); err != nil {
		return err
	}

	log.Info("rebuild rbmq pool done - ", prefix)

	return nil
}

func store(prefix string, conf *rabbitConfig) error {
	resourcePool, err := NewResourcePool(conf)
	if err != nil {
		return fmt.Errorf("catdog_rabbitmq_plugin(%+v) 连接失败, error=%+v", prefix, err)
	}

	initRbmqMap.Store(prefix, resourcePool)
	return nil
}

func PickupRabiitMQClient(prefix string) (*RabbitChan, error) {
	for i := 0; i < retryCount; i++ {
		v, ok := initRbmqMap.Load(prefix)
		if !ok {
			log.Errorf("can not get rbmq client ,prefix=" + prefix)
			return nil, errors.New("can not get catdog_rabbitmq_plugin client ,prefix=" + prefix)
		}

		resourcePool, ok := v.(*ResourcePool)
		if !ok {
			return nil, fmt.Errorf("convert RabbitMQ fail, InitRbmqMap.Load(%s) = %+v", prefix, resourcePool)
		}

		resource, err := resourcePool.Get()
		if err != nil {
			return nil, fmt.Errorf("get RabbitMQ from the pool failed. %+v ", err)
		}

		ch, err := resource.Channel()
		if err != nil {
			if err == amqp.ErrClosed {
				if err := store(prefix, resource.config); err != nil {
					return nil, fmt.Errorf("catdog_rabbitmq_plugin(%+v) reconnect error: %+v", prefix, err)
				}

				// 如果已经重连成功了，则把老的连接池关闭，释放对象
				_ = xprocess.Go(func(ctx context.Context) error {
					resourcePool.Close()
					return nil
				})
			}

			log.Error(err, " prefix=", prefix)
			continue
		}

		// Put back if the connection is ok
		resourcePool.Put(resource)

		return &RabbitChan{
			Channel: ch,
		}, nil
	}

	return nil, fmt.Errorf("catdog_rabbitmq_plugin(%+v) failed when retry %+v times. ", prefix, retryCount)
}

type RabbitChanWithContext struct {
	rabbitChan *RabbitChan
	ctx        context.Context
}

// Caution!!! This function return a channel, you should close it after use.
func GetRbmqClient(ctx context.Context, prefix string) (*RabbitChanWithContext, error) {
	rabbitmqChan, err := PickupRabiitMQClient(prefix)
	if err != nil || rabbitmqChan == nil {
		return nil, err
	}

	// We Should Set Qos to [100-300] due to https://www.rabbitmq.com/confirms.html#channel-qos-prefetch
	if err := rabbitmqChan.Qos(prefetchCount, 0, false); err != nil {
		return nil, err
	}

	rabbitChanWithContext := &RabbitChanWithContext{
		rabbitChan: rabbitmqChan,
		ctx:        ctx,
	}

	return rabbitChanWithContext, nil
}

func (r *RabbitChanWithContext) GetContext() context.Context {
	return r.ctx
}

func (r *RabbitChanWithContext) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *RabbitChanWithContext) newSpan(opeartionName string) opentracing.Span {
	_, span, err := tracing.StartSpanFromContext(r.ctx, opentracing.GlobalTracer(), opeartionName)
	if err != nil {
		// Maybe there will be many logs, annotate it.
		// log.Warn("[RabbitChanWithContext.newSpan] start span error. ", err)
		_ = err
	}

	return span
}

/*
 定义exchange
 参数：
	name exchange名称
	kind exchange种类
	durable 是否持久化
	autoDelete 当所有绑定的队列都与交换器解绑后，交换器会自动删除
 返回值：
	error 操作期间产生的错误
*/
func (r *RabbitChanWithContext) DeclareExchange(name string, kind ExchangeKind, durable, autoDelete bool) (err error) {
	span := r.newSpan("Rabbitmq.DeclareExchange")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "Rabbitmq")
			span.SetTag("Method", "DeclareExchang")
			span.SetTag("Exchange.name", name)
			span.SetTag("Exchange.Kind", kind)
			tracing.SetIfError(span, err)
			span.Finish()
		}()
	}

	return r.rabbitChan.ExchangeDeclare(name, string(kind), durable, autoDelete, false, false, nil)
}

/*
 定义队列
 参数：
	name 队列的名称
	durable 是否持久化
	autoDelete 当所有消费者都断开时，队列会自动删除
 返回值：
	error 操作期间产生的错误
*/
func (r *RabbitChanWithContext) DeclareQueue(name string, durable, autoDelete bool) (q amqp.Queue, err error) {
	span := r.newSpan("Rabbitmq.DeclareQueue")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "Rabbitmq")
			span.SetTag("Method", "DeclareQueue")
			span.SetTag("Queue.name", name)
			tracing.SetIfError(span, err)
			span.Finish()
		}()
	}

	return r.rabbitChan.QueueDeclare(name, durable, autoDelete, false, false, nil)
}

/*
 exchange和queue绑定
 参数：
	queue 队列名称
	bindkey 绑定的key
	exchange 交换器名称
 返回值：
	error 操作期间产生的错误
*/
func (r *RabbitChanWithContext) Bind(queue, bindKey, exchange string) (err error) {
	span := r.newSpan("Rabbitmq.Bind")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "Rabbitmq")
			span.SetTag("Method", "Bind")
			span.SetTag("Bind.queue", queue)
			span.SetTag("Bind.bindKey", bindKey)
			span.SetTag("Bind.exchange", exchange)
			tracing.SetIfError(span, err)
			span.Finish()
		}()
	}

	return r.rabbitChan.QueueBind(queue, bindKey, exchange, false, nil)
}

/*
 发送消息到指定的队列。通过默认exchange把消息发送到指定的消息队列
 参数：
 	queue 队列名称
 	msg 消息内容
 返回值：
	error 操作期间产生的错误
*/
func (r *RabbitChanWithContext) Send(queue string, msg []byte, deliveryMode uint8) (err error) {
	//TODO:: consider about this span
	span := r.newSpan("Rabbitmq.Send")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "Rabbitmq")
			span.SetTag("Method", "Send")
			span.SetTag("Send.queue", queue)
			span.SetTag("Send.msg", msg)
			tracing.SetIfError(span, err)
			span.Finish()
		}()
	}

	return r.Publish("", queue, msg, deliveryMode)
}

/*
 发送消息到指定的exchange
 参数：
	routingKey 路由key
 	msg 消息内容
 返回值：
	error 操作期间产生的错误
*/
func (r *RabbitChanWithContext) Publish(exchange, routingKey string, msg []byte, deliveryMode uint8) (err error) {
	span := r.newSpan("Rabbitmq.Publish")
	headers := amqp.Table{}

	if span != nil {
		defer func() {
			if exchange == "" {
				//When param exchange is empty, catdog_rabbitmq_plugin will use default exchange.
				exchange = "default"
			}
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "Rabbitmq")
			span.SetTag("Method", "Publish")
			span.SetTag("Publish.exchange", exchange)
			span.SetTag("Publish.routingKey", routingKey)
			span.SetTag("Publish.msg", msg)
			tracing.SetIfError(span, err)
			span.Finish()
		}()

		carrier := opentracing.TextMapCarrier{}
		if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, &carrier); err != nil {
			log.Error("Tracing:catdog_rabbitmq_plugin tracerInjectError:", err)
		}

		err = carrier.ForeachKey(func(key, val string) error {
			headers[key] = val
			// carrier.ForeachKey() will return true forever
			return nil
		})
		if err != nil {
			log.Error("Tracing:catdog_rabbitmq_plugin buildRabbitHeaderError:", err)
		}
	}

	//TODO:: delete these codes in v0.2.0
	//defer func() {
	//	if r.rabbitChan.Channel != nil {
	//		err := r.rabbitChan.Close()
	//		if err != nil {
	//			log.Error("Tracing:catdog_rabbitmq_plugin chanCloseError:", err)
	//		}
	//	}
	//}()

	return r.rabbitChan.Channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		Headers:         headers,
		ContentType:     contentType,
		ContentEncoding: contentEncoding,
		DeliveryMode:    deliveryMode, // 持久化
		Body:            msg,
	})
}

/*
 消费消息队列
 参数：
	queue 队列名称
	autoAck	是否自动回复ack
 返回值：
	Delivery 传递消息的单向通道，可以通过读取该通道获取接收到的消息
	error 操作期间产生的错误
*/
func (r *RabbitChanWithContext) Consume(queue string, autoAck bool) (ch <-chan amqp.Delivery, err error) {
	span := r.newSpan("Rabbitmq.Consume")
	if span != nil {
		defer func() {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "Rabbitmq")
			span.SetTag("Method", "Consume")
			span.SetTag("Consume.queue", queue)
			tracing.SetIfError(span, err)
			span.Finish()
		}()
	}

	//TODO:: delete these codes in v0.2.0
	//defer func() {
	//	if r.rabbitChan.Channel != nil {
	//		err := r.rabbitChan.Close()
	//		if err != nil {
	//			log.Error("Defer error:", err)
	//		}
	//	}
	//}()
	return r.rabbitChan.Channel.Consume(queue, "", autoAck, false, false, false, nil)
}

func (r *RabbitChanWithContext) Close() error {
	return r.rabbitChan.Close()
}

//This function decode the header from consumed message. Use tracing.StartSpanFromSpanContext() to build next span.
func SpanContextFromDelivery(d amqp.Delivery) (opentracing.SpanContext, error) {
	spanStr, _ := d.Headers[jaeger.TraceContextHeaderName].(string)
	//spanStr will be empty string when can not find in header
	return jaeger.ContextFromString(spanStr)
}

func TraceDelivery(d amqp.Delivery) {
	spanCtx, err := SpanContextFromDelivery(d)
	if err == nil {
		_, span, err := tracing.StartSpanFromSpanContext(spanCtx, opentracing.GlobalTracer(), "Rabbitmq.Received")
		if err == nil {
			ext.SpanKindRPCClient.Set(span)
			span.SetTag("DBType", "Rabbitmq")
			span.SetTag("Method", "Receive")
			span.SetTag("Delivery", d)
			tracing.SetIfError(span, err)
			go span.Finish()
		}
	}
}
