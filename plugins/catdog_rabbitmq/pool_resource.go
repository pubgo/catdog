package catdog_rabbitmq

import (
	"context"
	"errors"
	"github.com/pubgo/xlog"
	"time"

	"github.com/streadway/amqp"

	"vitess.io/vitess/go/pools"
)

const (
	DefaultMQTimeout   = time.Second * 2
	DefaultIdleTime    = 0
	DefaultHeartbeat   = time.Second * 2
	DefaultCapacity    = 10
	MaxCapacity        = 50
	prefillParallelism = 10
)

var (
	ErrorOutOfCapacity = errors.New("rabbitMQ resource pool has got the Max Capacity. ")
)

type ResourcePool struct {
	*pools.ResourcePool
}

// NewResourcePool create a resource pool
func NewResourcePool(config *RbmqConfig) (*ResourcePool, error) {
	// Set the catdog_rabbitmq_plugin pool by DefaultCapacity
	resourcePool := pools.NewResourcePool(newResource(config), DefaultCapacity, MaxCapacity, DefaultIdleTime, prefillParallelism, func(t time.Time) {
		return
	})

	// Check the connect is ok
	ctx, cancel := context.WithTimeout(context.TODO(), DefaultMQTimeout)
	defer cancel()

	conn, err := resourcePool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return &ResourcePool{
		ResourcePool: resourcePool,
	}, nil
}

func (rp *ResourcePool) Get() (*Resource, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), DefaultMQTimeout)
	defer cancel()

	r, err := rp.ResourcePool.Get(ctx)
	if err != nil {
		if err == pools.ErrTimeout {
			cp := rp.Capacity()
			if cp == MaxCapacity {
				return nil, ErrorOutOfCapacity
			}

			cp += DefaultCapacity
			if cp >= MaxCapacity {
				cp = MaxCapacity
			}
			if err := rp.ResourcePool.SetCapacity(int(cp)); err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	return r.(*Resource), nil
}

// Resource adapts a catdog_rabbitmq_plugin connection to a Vitess Resource.
type Resource struct {
	*amqp.Connection
	config *RbmqConfig
}

// newResource return an closure for create a catdog_rabbitmq_plugin connection
func newResource(config *RbmqConfig) pools.Factory {
	return func(ctx context.Context) (pools.Resource, error) {
		c := amqp.Config{
			Heartbeat: DefaultHeartbeat,
		}
		conn, err := amqp.DialConfig(config.URL, c)

		if err != nil {
			return nil, err
		}

		return &Resource{conn, config}, nil
	}
}

// Close is put the conn to the poll
func (r *Resource) Close() {
	if err := r.Connection.Close(); err != nil {
		xlog.Error(err)
	}
}
