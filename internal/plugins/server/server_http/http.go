// Package http implements a go-micro.Server
package server_http

import (
	"errors"
	"fmt"
	"github.com/pubgo/xerror"
	"net"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/asim/nitro/v3/broker"
	"github.com/asim/nitro/v3/codec"
	"github.com/asim/nitro/v3/codec/jsonrpc"
	"github.com/asim/nitro/v3/codec/protorpc"
	"github.com/asim/nitro/v3/registry"
	"github.com/asim/nitro/v3/server"
)

var (
	defaultCodecs = map[string]codec.NewCodec{
		"application/json":         jsonrpc.NewCodec,
		"application/json-rpc":     jsonrpc.NewCodec,
		"application/protobuf":     protorpc.NewCodec,
		"application/proto-rpc":    protorpc.NewCodec,
		"application/octet-stream": protorpc.NewCodec,
	}
)

type httpServer struct {
	sync.Mutex
	opts         server.Options
	hd           server.Handler
	exit         chan chan error
	registerOnce sync.Once
	subscribers  map[*httpSubscriber][]broker.Subscriber
	// used for first registration
	registered bool
}

func (t *httpServer) newCodec(contentType string) (codec.NewCodec, error) {
	if cf, ok := t.opts.Codecs[contentType]; ok {
		return cf, nil
	}
	if cf, ok := defaultCodecs[contentType]; ok {
		return cf, nil
	}
	return nil, fmt.Errorf("Unsupported Content-Type: %s", contentType)
}

func (t *httpServer) Options() server.Options {
	t.Lock()
	opts := t.opts
	t.Unlock()
	return opts
}

func (t *httpServer) Init(opts ...server.Option) error {
	t.Lock()
	for _, o := range opts {
		o(&t.opts)
	}
	t.Unlock()
	return nil
}

func (t *httpServer) Handle(handler server.Handler) error {
	if _, ok := handler.Handler().(http.Handler); !ok {
		return errors.New("Handler requires http.Handler")
	}
	t.Lock()
	t.hd = handler
	t.Unlock()
	return nil
}

func (t *httpServer) NewHandler(handler interface{}, opts ...server.HandlerOption) server.Handler {
	options := server.HandlerOptions{
		Metadata: make(map[string]map[string]string),
	}

	for _, o := range opts {
		o(&options)
	}

	var eps []*registry.Endpoint

	if !options.Internal {
		for name, metadata := range options.Metadata {
			eps = append(eps, &registry.Endpoint{
				Name:     name,
				Metadata: metadata,
			})
		}
	}

	return &httpHandler{
		eps:  eps,
		hd:   handler,
		opts: options,
	}
}

func (t *httpServer) NewSubscriber(topic string, handler interface{}, opts ...server.SubscriberOption) server.Subscriber {
	return newSubscriber(topic, handler, opts...)
}

func (t *httpServer) Subscribe(sb server.Subscriber) error {
	sub, ok := sb.(*httpSubscriber)
	if !ok {
		return fmt.Errorf("invalid subscriber: expected *httpSubscriber")
	}
	if len(sub.handlers) == 0 {
		return fmt.Errorf("invalid subscriber: no handler functions")
	}

	if err := validateSubscriber(sb); err != nil {
		return err
	}

	t.Lock()
	defer t.Unlock()
	_, ok = t.subscribers[sub]
	if ok {
		return fmt.Errorf("subscriber %v already exists", t)
	}
	t.subscribers[sub] = nil
	return nil
}

func (t *httpServer) Register() error {
	t.Lock()
	opts := t.opts
	eps := t.hd.Endpoints()
	t.Unlock()

	service := serviceDef(opts)
	service.Endpoints = eps

	t.Lock()
	var subscriberList []*httpSubscriber
	for e := range t.subscribers {
		// Only advertise non internal subscribers
		if !e.Options().Internal {
			subscriberList = append(subscriberList, e)
		}
	}
	sort.Slice(subscriberList, func(i, j int) bool {
		return subscriberList[i].topic > subscriberList[j].topic
	})
	for _, e := range subscriberList {
		service.Endpoints = append(service.Endpoints, e.Endpoints()...)
	}
	t.Unlock()

	rOpts := []registry.RegisterOption{
		registry.RegisterTTL(opts.RegisterTTL),
	}

	t.registerOnce.Do(func() {
		log.Infof("Registering node: %s", opts.Name+"-"+opts.Id)
	})

	if err := opts.Registry.Register(service, rOpts...); err != nil {
		return err
	}

	t.Lock()
	defer t.Unlock()

	if t.registered {
		return nil
	}
	t.registered = true

	for sb, _ := range t.subscribers {
		handler := t.createSubHandler(sb, opts)
		var subOpts []broker.SubscribeOption
		if queue := sb.Options().Queue; len(queue) > 0 {
			subOpts = append(subOpts, broker.Queue(queue))
		}
		sub, err := opts.Broker.Subscribe(sb.Topic(), handler, subOpts...)
		if err != nil {
			return err
		}
		t.subscribers[sb] = []broker.Subscriber{sub}
	}
	return nil
}

func (t *httpServer) Deregister() error {
	t.Lock()
	opts := t.opts
	t.Unlock()

	log.Infof("Deregistering node: %s", opts.Name+"-"+opts.Id)

	service := serviceDef(opts)
	if err := opts.Registry.Deregister(service); err != nil {
		return err
	}

	t.Lock()
	if !t.registered {
		t.Unlock()
		return nil
	}
	t.registered = false

	for sb, subs := range t.subscribers {
		for _, sub := range subs {
			log.Infof("Unsubscribing from topic: %s", sub.Topic())
			sub.Unsubscribe()
		}
		t.subscribers[sb] = nil
	}
	t.Unlock()
	return nil
}

func (t *httpServer) Start() (err error) {
	defer xerror.RespErr(&err)

	t.Lock()
	opts := t.opts
	hd := t.hd
	t.Unlock()

	ln, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return err
	}

	log.Infof("Listening on %s", ln.Addr().String())

	t.Lock()
	t.opts.Address = ln.Addr().String()
	t.Unlock()

	handler, ok := hd.Handler().(http.Handler)
	if !ok {
		return errors.New("Server required http.Handler")
	}

	// register
	xerror.Panic(t.Register())

	go http.Serve(ln, handler)

	go func() {
		tk := new(time.Ticker)

		// only process if it exists
		if opts.RegisterInterval > time.Duration(0) {
			// new ticker
			tk = time.NewTicker(opts.RegisterInterval)
		}

		// return error chan
		var ch chan error

	Loop:
		for {
			select {
			// register self on interval
			case <-tk.C:
				if err := t.Register(); err != nil {
					log.Infof("Server register error: %v", err)
				}
			// wait for exit
			case ch = <-t.exit:
				break Loop
			}
		}

		ch <- ln.Close()

		// deregister
		t.Deregister()

		opts.Broker.Disconnect()
	}()

	return opts.Broker.Connect()
}

func (t *httpServer) Stop() error {
	ch := make(chan error)
	t.exit <- ch
	return <-ch
}

func (t *httpServer) String() string {
	return "http"
}

func newServer(opts ...server.Option) server.Server {
	return &httpServer{
		opts:        newOptions(opts...),
		exit:        make(chan chan error),
		subscribers: make(map[*httpSubscriber][]broker.Subscriber),
	}
}

func NewServer(opts ...server.Option) server.Server {
	return newServer(opts...)
}
