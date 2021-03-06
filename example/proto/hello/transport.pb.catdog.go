// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: example/proto/hello/transport.proto
package hello

import (
	context "context"
	fmt "fmt"
	math "math"

	"github.com/pubgo/catdog/catdog_data"

	client "github.com/asim/nitro/v3/client"
	server "github.com/asim/nitro/v3/server"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Transport service
type TransportService interface {
	TestStream(ctx context.Context, opts ...client.CallOption) (Transport_TestStreamService, error)
	TestStream1(ctx context.Context, opts ...client.CallOption) (Transport_TestStream1Service, error)
	TestStream2(ctx context.Context, in *Message, opts ...client.CallOption) (Transport_TestStream2Service, error)
	TestStream3(ctx context.Context, in *Message, opts ...client.CallOption) (*Message, error)
}

type transportService struct {
	c    client.Client
	name string
}

func NewTransportService(name string, c client.Client) TransportService {
	return &transportService{
		c:    c,
		name: name,
	}
}
func (c *transportService) TestStream(ctx context.Context, opts ...client.CallOption) (Transport_TestStreamService, error) {

	req := c.c.NewRequest(c.name, "Transport.TestStream", &Message{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return Transport.TestStream{stream}, nil
}

// Stream auxiliary types and methods.
type Transport_TestStreamService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error

	Send(*Message) error

	Recv(*Message) error
}
type TransportTestStream struct {
	stream client.Stream
}

func (x *TransportTestStream) Close() error {
	return x.stream.Close()
}

func (x *TransportTestStream) Context() context.Context {
	return x.stream.Context()
}

func (x *TransportTestStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *TransportTestStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *TransportTestStream) Send(m *Message) error {
	return x.stream.Send(m)
}

func (x *TransportTestStream) Recv() (*Message, error) {
	m := new(Message)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transportService) TestStream1(ctx context.Context, opts ...client.CallOption) (Transport_TestStream1Service, error) {

	req := c.c.NewRequest(c.name, "Transport.TestStream1", &Message{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return Transport.TestStream1{stream}, nil
}

// Stream auxiliary types and methods.
type Transport_TestStream1Service interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error

	Send(*Message) error
}
type TransportTestStream1 struct {
	stream client.Stream
}

func (x *TransportTestStream1) Close() error {
	return x.stream.Close()
}

func (x *TransportTestStream1) Context() context.Context {
	return x.stream.Context()
}

func (x *TransportTestStream1) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *TransportTestStream1) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *TransportTestStream1) Send(m *Message) error {
	return x.stream.Send(m)
}

func (c *transportService) TestStream2(ctx context.Context, in *Message, opts ...client.CallOption) (Transport_TestStream2Service, error) {

	req := c.c.NewRequest(c.name, "Transport.TestStream2", &Message{})
	stream, err := c.c.Stream(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	if err := stream.Send(in); err != nil {
		return nil, err
	}
	return Transport.TestStream2{stream}, nil
}

// Stream auxiliary types and methods.
type Transport_TestStream2Service interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error

	Recv(*Message) error
}
type TransportTestStream2 struct {
	stream client.Stream
}

func (x *TransportTestStream2) Close() error {
	return x.stream.Close()
}

func (x *TransportTestStream2) Context() context.Context {
	return x.stream.Context()
}

func (x *TransportTestStream2) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *TransportTestStream2) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *TransportTestStream2) Recv() (*Message, error) {
	m := new(Message)
	err := x.stream.Recv(m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *transportService) TestStream3(ctx context.Context, in *Message, opts ...client.CallOption) (*Message, error) {

	req := c.c.NewRequest(c.name, "Transport.TestStream3", in)
	out := new(Message)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Transport_TestStream3Service interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type TransportTestStream3 struct {
	stream client.Stream
}

func (x *TransportTestStream3) Close() error {
	return x.stream.Close()
}

func (x *TransportTestStream3) Context() context.Context {
	return x.stream.Context()
}

func (x *TransportTestStream3) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *TransportTestStream3) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

// Server API for Transport service
type TransportHandler interface {
	TestStream(context.Context, Transport_TestStreamStream) error
	TestStream1(context.Context, Transport_TestStream1Stream) error
	TestStream2(context.Context, *Message, Transport_TestStream2Stream) error
	TestStream3(context.Context, *Message, *Message) error
}

func RegisterTransportHandler(s server.Server, hdlr TransportHandler, opts ...server.HandlerOption) error {
	type transport interface {
		TestStream(ctx context.Context, stream server.Stream) error
		TestStream1(ctx context.Context, stream server.Stream) error
		TestStream2(ctx context.Context, stream server.Stream) error
		TestStream3(ctx context.Context, in *Message, out *Message) error
	}

	type Transport struct {
		transport
	}
	h := &transportHandler{hdlr}
	opts = append(opts, server.EndpointMetadata("TestStream", map[string]string{"POST": "hello_transport/test_stream"}))
	opts = append(opts, server.EndpointMetadata("TestStream1", map[string]string{"POST": "hello_transport/test_stream1"}))
	opts = append(opts, server.EndpointMetadata("TestStream2", map[string]string{"POST": "hello_transport/test_stream2"}))
	opts = append(opts, server.EndpointMetadata("TestStream3", map[string]string{"POST": "hello_transport/test_stream3"}))
	return s.Handle(s.NewHandler(&Transport{h}, opts...))
}

func init() { catdog_data.Add("RegisterTransportHandler", RegisterTransportHandler) }

type transportHandler struct {
	TransportHandler
}

func (h *transportHandler) TestStream(ctx context.Context, stream server.Stream) error {

	return h.TransportHandler.TestStream(ctx, &transportTestStreamStream{stream})

}

type Transport_TestStreamStream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error

	Send(*Message) error

	Recv() (*Message, error)
}

type transportTestStreamStream struct {
	stream server.Stream
}

func (x *transportTestStreamStream) Close() error {
	return x.stream.Close()
}

func (x *transportTestStreamStream) Context() context.Context {
	return x.stream.Context()
}

func (x *transportTestStreamStream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *transportTestStreamStream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *transportTestStreamStream) Send(*Message) error {
	return x.stream.Send(m)
}

func (x *transportTestStreamStream) Recv() (*Message, error) {
	m := new(Message)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *transportHandler) TestStream1(ctx context.Context, stream server.Stream) error {

	return h.TransportHandler.TestStream1(ctx, &transportTestStream1Stream{stream})

}

type Transport_TestStream1Stream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error

	Recv() (*Message, error)
}

type transportTestStream1Stream struct {
	stream server.Stream
}

func (x *transportTestStream1Stream) Close() error {
	return x.stream.Close()
}

func (x *transportTestStream1Stream) Context() context.Context {
	return x.stream.Context()
}

func (x *transportTestStream1Stream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *transportTestStream1Stream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *transportTestStream1Stream) Recv() (*Message, error) {
	m := new(Message)
	if err := x.stream.Recv(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *transportHandler) TestStream2(ctx context.Context, stream server.Stream) error {

	m := new(Message)
	if err := stream.Recv(m); err != nil {
		return err
	}
	return h.TransportHandler.TestStream2(ctx, m, &transportTestStream2Stream{stream})

}

type Transport_TestStream2Stream interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error

	Send(*Message) error
}

type transportTestStream2Stream struct {
	stream server.Stream
}

func (x *transportTestStream2Stream) Close() error {
	return x.stream.Close()
}

func (x *transportTestStream2Stream) Context() context.Context {
	return x.stream.Context()
}

func (x *transportTestStream2Stream) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *transportTestStream2Stream) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (x *transportTestStream2Stream) Send(*Message) error {
	return x.stream.Send(m)
}

func (h *transportHandler) TestStream3(ctx context.Context, in *Message, out *Message) error {
	return h.TransportHandler.TestStream3(ctx, in, out)
}
