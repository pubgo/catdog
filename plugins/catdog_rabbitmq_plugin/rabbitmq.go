package catdog_rabbitmq_plugin

import (
	"errors"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)


type ExchangeKind string

const (
	// ExchangeKindFanout ExchangeKind = "fanout"
	ExchangeKindDirect ExchangeKind = "direct"
	// ExchangeKindTopic  ExchangeKind = "topic"

	contentType     = "application/json"
	contentEncoding = "UTF-8"

	retryCount = 2
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

	var conf *RbmqConfig
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

func store(prefix string, conf *RbmqConfig) error {
	resourcePool, err := NewResourcePool(conf)
	if err != nil {
		return fmt.Errorf("catdog_rabbitmq_plugin(%+v) 连接失败, error=%+v", prefix, err)
	}

	initRbmqMap.Store(prefix, resourcePool)
	return nil
}

func NewRbmqWatcher() watcher.Watcher {
	// 注册prefix == "etcd/watch/redis"
	// 获取注册的配置中心名称，默认etcd
	// 获取prefix节点，遍历服务项目名称，得到列表
	return &rbmqwather{
		WatchKey: "catdog_rabbitmq_plugin",
	}
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
				_ = grpool.Submit(func() {
					resourcePool.Close()
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
