package server_http

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/asim/nitro/v3/broker"
	"github.com/asim/nitro/v3/codec"
	"github.com/asim/nitro/v3/registry"
	"github.com/asim/nitro/v3/server"
)

const (
	subSig = "func(context.Context, interface{}) error"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type handler struct {
	method  reflect.Value
	reqType reflect.Type
	ctxType reflect.Type
}

type httpSubscriber struct {
	topic      string
	rcvr       reflect.Value
	typ        reflect.Type
	subscriber interface{}
	handlers   []*handler
	endpoints  []*registry.Endpoint
	opts       server.SubscriberOptions
}

// Is this an exported - upper case - name?
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

func newSubscriber(topic string, sub interface{}, opts ...server.SubscriberOption) server.Subscriber {
	var options server.SubscriberOptions
	for _, o := range opts {
		o(&options)
	}

	var endpoints []*registry.Endpoint
	var handlers []*handler

	if typ := reflect.TypeOf(sub); typ.Kind() == reflect.Func {
		h := &handler{
			method: reflect.ValueOf(sub),
		}

		switch typ.NumIn() {
		case 1:
			h.reqType = typ.In(0)
		case 2:
			h.ctxType = typ.In(0)
			h.reqType = typ.In(1)
		}

		handlers = append(handlers, h)

		endpoints = append(endpoints, &registry.Endpoint{
			Name:    "Func",
			Request: extractSubValue(typ),
			Metadata: map[string]string{
				"topic":      topic,
				"subscriber": "true",
			},
		})
	} else {
		hdlr := reflect.ValueOf(sub)
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			h := &handler{
				method: method.Func,
			}

			switch method.Type.NumIn() {
			case 2:
				h.reqType = method.Type.In(1)
			case 3:
				h.ctxType = method.Type.In(1)
				h.reqType = method.Type.In(2)
			}

			handlers = append(handlers, h)

			endpoints = append(endpoints, &registry.Endpoint{
				Name:    name + "." + method.Name,
				Request: extractSubValue(method.Type),
				Metadata: map[string]string{
					"topic":      topic,
					"subscriber": "true",
				},
			})
		}
	}

	return &httpSubscriber{
		rcvr:       reflect.ValueOf(sub),
		typ:        reflect.TypeOf(sub),
		topic:      topic,
		subscriber: sub,
		handlers:   handlers,
		endpoints:  endpoints,
		opts:       options,
	}
}

func validateSubscriber(sub server.Subscriber) error {
	typ := reflect.TypeOf(sub.Subscriber())
	var argType reflect.Type

	if typ.Kind() == reflect.Func {
		name := "Func"
		switch typ.NumIn() {
		case 2:
			argType = typ.In(1)
		default:
			return fmt.Errorf("subscriber %v takes wrong number of args: %v required signature %s", name, typ.NumIn(), subSig)
		}
		if !isExportedOrBuiltinType(argType) {
			return fmt.Errorf("subscriber %v argument type not exported: %v", name, argType)
		}
		if typ.NumOut() != 1 {
			return fmt.Errorf("subscriber %v has wrong number of outs: %v require signature %s",
				name, typ.NumOut(), subSig)
		}
		if returnType := typ.Out(0); returnType != typeOfError {
			return fmt.Errorf("subscriber %v returns %v not error", name, returnType.String())
		}
	} else {
		hdlr := reflect.ValueOf(sub.Subscriber())
		name := reflect.Indirect(hdlr).Type().Name()

		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)

			switch method.Type.NumIn() {
			case 3:
				argType = method.Type.In(2)
			default:
				return fmt.Errorf("subscriber %v.%v takes wrong number of args: %v required signature %s",
					name, method.Name, method.Type.NumIn(), subSig)
			}

			if !isExportedOrBuiltinType(argType) {
				return fmt.Errorf("%v argument type not exported: %v", name, argType)
			}
			if method.Type.NumOut() != 1 {
				return fmt.Errorf(
					"subscriber %v.%v has wrong number of outs: %v require signature %s",
					name, method.Name, method.Type.NumOut(), subSig)
			}
			if returnType := method.Type.Out(0); returnType != typeOfError {
				return fmt.Errorf("subscriber %v.%v returns %v not error", name, method.Name, returnType.String())
			}
		}
	}

	return nil
}

func (t *httpServer) createSubHandler(sb *httpSubscriber, opts server.Options) broker.Handler {
	return func(msg *broker.Message) error {
		ct := msg.Header["Content-Type"]
		cf, err := t.newCodec(ct)
		if err != nil {
			return err
		}

		hdr := make(map[string]string)
		for k, v := range msg.Header {
			hdr[k] = v
		}
		delete(hdr, "Content-Type")
		//ctx := metadata.NewContext(context.Background(), hdr)

		results := make(chan error, len(sb.handlers))

		for i := 0; i < len(sb.handlers); i++ {
			handler := sb.handlers[i]

			var isVal bool
			var req reflect.Value

			if handler.reqType.Kind() == reflect.Ptr {
				req = reflect.New(handler.reqType.Elem())
			} else {
				req = reflect.New(handler.reqType)
				isVal = true
			}
			if isVal {
				req = req.Elem()
			}

			b := &buffer{bytes.NewBuffer(msg.Body)}
			co := cf(b)
			defer co.Close()

			if err := co.ReadHeader(&codec.Message{}, codec.Event); err != nil {
				return err
			}

			if err := co.ReadBody(req.Interface()); err != nil {
				return err
			}

			fn := func(ctx context.Context, msg server.Message) error {
				var vals []reflect.Value
				if sb.typ.Kind() != reflect.Func {
					vals = append(vals, sb.rcvr)
				}
				if handler.ctxType != nil {
					vals = append(vals, reflect.ValueOf(ctx))
				}

				vals = append(vals, reflect.ValueOf(msg.Payload()))

				returnValues := handler.method.Call(vals)
				if err := returnValues[0].Interface(); err != nil {
					return err.(error)
				}
				return nil
			}

			for i := len(opts.SubWrappers); i > 0; i-- {
				fn = opts.SubWrappers[i-1](fn)
			}

			go func() {
				//results <- fn(ctx, &httpMessage{
				//	topic:       sb.topic,
				//	contentType: ct,
				//	payload:     req.Interface(),
				//})
			}()
		}

		var errors []string

		for i := 0; i < len(sb.handlers); i++ {
			if err := <-results; err != nil {
				errors = append(errors, err.Error())
			}
		}

		if len(errors) > 0 {
			return fmt.Errorf("subscriber error: %s", strings.Join(errors, "\n"))
		}

		return nil
	}
}

func (s *httpSubscriber) Topic() string {
	return s.topic
}

func (s *httpSubscriber) Subscriber() interface{} {
	return s.subscriber
}

func (s *httpSubscriber) Endpoints() []*registry.Endpoint {
	return s.endpoints
}

func (s *httpSubscriber) Options() server.SubscriberOptions {
	return s.opts
}