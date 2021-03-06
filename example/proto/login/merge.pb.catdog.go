// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: example/proto/login/merge.proto
package login

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

// Client API for Merge service
type MergeService interface {
	Telephone(ctx context.Context, in *TelephoneRequest, opts ...client.CallOption) (*Reply, error)
	TelephoneCheck(ctx context.Context, in *TelephoneRequest, opts ...client.CallOption) (*Reply, error)
	WeChat(ctx context.Context, in *WeChatRequest, opts ...client.CallOption) (*Reply, error)
	WeChatCheck(ctx context.Context, in *WeChatRequest, opts ...client.CallOption) (*Reply, error)
	WeChatUnMerge(ctx context.Context, in *WeChatUnMergeRequest, opts ...client.CallOption) (*Reply, error)
}

type mergeService struct {
	c    client.Client
	name string
}

func NewMergeService(name string, c client.Client) MergeService {
	return &mergeService{
		c:    c,
		name: name,
	}
}
func (c *mergeService) Telephone(ctx context.Context, in *TelephoneRequest, opts ...client.CallOption) (*Reply, error) {

	req := c.c.NewRequest(c.name, "Merge.Telephone", in)
	out := new(Reply)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Merge_TelephoneService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type MergeTelephone struct {
	stream client.Stream
}

func (x *MergeTelephone) Close() error {
	return x.stream.Close()
}

func (x *MergeTelephone) Context() context.Context {
	return x.stream.Context()
}

func (x *MergeTelephone) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *MergeTelephone) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *mergeService) TelephoneCheck(ctx context.Context, in *TelephoneRequest, opts ...client.CallOption) (*Reply, error) {

	req := c.c.NewRequest(c.name, "Merge.TelephoneCheck", in)
	out := new(Reply)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Merge_TelephoneCheckService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type MergeTelephoneCheck struct {
	stream client.Stream
}

func (x *MergeTelephoneCheck) Close() error {
	return x.stream.Close()
}

func (x *MergeTelephoneCheck) Context() context.Context {
	return x.stream.Context()
}

func (x *MergeTelephoneCheck) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *MergeTelephoneCheck) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *mergeService) WeChat(ctx context.Context, in *WeChatRequest, opts ...client.CallOption) (*Reply, error) {

	req := c.c.NewRequest(c.name, "Merge.WeChat", in)
	out := new(Reply)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Merge_WeChatService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type MergeWeChat struct {
	stream client.Stream
}

func (x *MergeWeChat) Close() error {
	return x.stream.Close()
}

func (x *MergeWeChat) Context() context.Context {
	return x.stream.Context()
}

func (x *MergeWeChat) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *MergeWeChat) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *mergeService) WeChatCheck(ctx context.Context, in *WeChatRequest, opts ...client.CallOption) (*Reply, error) {

	req := c.c.NewRequest(c.name, "Merge.WeChatCheck", in)
	out := new(Reply)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Merge_WeChatCheckService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type MergeWeChatCheck struct {
	stream client.Stream
}

func (x *MergeWeChatCheck) Close() error {
	return x.stream.Close()
}

func (x *MergeWeChatCheck) Context() context.Context {
	return x.stream.Context()
}

func (x *MergeWeChatCheck) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *MergeWeChatCheck) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *mergeService) WeChatUnMerge(ctx context.Context, in *WeChatUnMergeRequest, opts ...client.CallOption) (*Reply, error) {

	req := c.c.NewRequest(c.name, "Merge.WeChatUnMerge", in)
	out := new(Reply)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Merge_WeChatUnMergeService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type MergeWeChatUnMerge struct {
	stream client.Stream
}

func (x *MergeWeChatUnMerge) Close() error {
	return x.stream.Close()
}

func (x *MergeWeChatUnMerge) Context() context.Context {
	return x.stream.Context()
}

func (x *MergeWeChatUnMerge) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *MergeWeChatUnMerge) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

// Server API for Merge service
type MergeHandler interface {
	Telephone(context.Context, *TelephoneRequest, *Reply) error
	TelephoneCheck(context.Context, *TelephoneRequest, *Reply) error
	WeChat(context.Context, *WeChatRequest, *Reply) error
	WeChatCheck(context.Context, *WeChatRequest, *Reply) error
	WeChatUnMerge(context.Context, *WeChatUnMergeRequest, *Reply) error
}

func RegisterMergeHandler(s server.Server, hdlr MergeHandler, opts ...server.HandlerOption) error {
	type merge interface {
		Telephone(ctx context.Context, in *TelephoneRequest, out *Reply) error
		TelephoneCheck(ctx context.Context, in *TelephoneRequest, out *Reply) error
		WeChat(ctx context.Context, in *WeChatRequest, out *Reply) error
		WeChatCheck(ctx context.Context, in *WeChatRequest, out *Reply) error
		WeChatUnMerge(ctx context.Context, in *WeChatUnMergeRequest, out *Reply) error
	}

	type Merge struct {
		merge
	}
	h := &mergeHandler{hdlr}
	opts = append(opts, server.EndpointMetadata("Telephone", map[string]string{"POST": "/user/merge/telephone"}))
	opts = append(opts, server.EndpointMetadata("TelephoneCheck", map[string]string{"POST": "/user/merge/telephone-check"}))
	opts = append(opts, server.EndpointMetadata("WeChat", map[string]string{"POST": "/user/merge/we-chat"}))
	opts = append(opts, server.EndpointMetadata("WeChatCheck", map[string]string{"POST": "/user/merge/we-chat-check"}))
	opts = append(opts, server.EndpointMetadata("WeChatUnMerge", map[string]string{"POST": "/user/merge/we-chat-un-merge"}))
	return s.Handle(s.NewHandler(&Merge{h}, opts...))
}

func init() { catdog_data.Add("RegisterMergeHandler", RegisterMergeHandler) }

type mergeHandler struct {
	MergeHandler
}

func (h *mergeHandler) Telephone(ctx context.Context, in *TelephoneRequest, out *Reply) error {
	return h.MergeHandler.Telephone(ctx, in, out)
}

func (h *mergeHandler) TelephoneCheck(ctx context.Context, in *TelephoneRequest, out *Reply) error {
	return h.MergeHandler.TelephoneCheck(ctx, in, out)
}

func (h *mergeHandler) WeChat(ctx context.Context, in *WeChatRequest, out *Reply) error {
	return h.MergeHandler.WeChat(ctx, in, out)
}

func (h *mergeHandler) WeChatCheck(ctx context.Context, in *WeChatRequest, out *Reply) error {
	return h.MergeHandler.WeChatCheck(ctx, in, out)
}

func (h *mergeHandler) WeChatUnMerge(ctx context.Context, in *WeChatUnMergeRequest, out *Reply) error {
	return h.MergeHandler.WeChatUnMerge(ctx, in, out)
}
