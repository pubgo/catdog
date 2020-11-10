// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: example/proto/login/code.proto
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

// Client API for Code service
type CodeService interface {
	SendCode(ctx context.Context, in *SendCodeRequest, opts ...client.CallOption) (*SendCodeResponse, error)
	Verify(ctx context.Context, in *VerifyRequest, opts ...client.CallOption) (*VerifyResponse, error)
	IsCheckImageCode(ctx context.Context, in *IsCheckImageCodeRequest, opts ...client.CallOption) (*IsCheckImageCodeResponse, error)
	VerifyImageCode(ctx context.Context, in *VerifyImageCodeRequest, opts ...client.CallOption) (*VerifyImageCodeResponse, error)
	GetSendStatus(ctx context.Context, in *GetSendStatusRequest, opts ...client.CallOption) (*GetSendStatusResponse, error)
}

type codeService struct {
	c    client.Client
	name string
}

func NewCodeService(name string, c client.Client) CodeService {
	return &codeService{
		c:    c,
		name: name,
	}
}
func (c *codeService) SendCode(ctx context.Context, in *SendCodeRequest, opts ...client.CallOption) (*SendCodeResponse, error) {

	req := c.c.NewRequest(c.name, "Code.SendCode", in)
	out := new(SendCodeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Code_SendCodeService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type CodeSendCode struct {
	stream client.Stream
}

func (x *CodeSendCode) Close() error {
	return x.stream.Close()
}

func (x *CodeSendCode) Context() context.Context {
	return x.stream.Context()
}

func (x *CodeSendCode) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *CodeSendCode) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *codeService) Verify(ctx context.Context, in *VerifyRequest, opts ...client.CallOption) (*VerifyResponse, error) {

	req := c.c.NewRequest(c.name, "Code.Verify", in)
	out := new(VerifyResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Code_VerifyService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type CodeVerify struct {
	stream client.Stream
}

func (x *CodeVerify) Close() error {
	return x.stream.Close()
}

func (x *CodeVerify) Context() context.Context {
	return x.stream.Context()
}

func (x *CodeVerify) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *CodeVerify) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *codeService) IsCheckImageCode(ctx context.Context, in *IsCheckImageCodeRequest, opts ...client.CallOption) (*IsCheckImageCodeResponse, error) {

	req := c.c.NewRequest(c.name, "Code.IsCheckImageCode", in)
	out := new(IsCheckImageCodeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Code_IsCheckImageCodeService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type CodeIsCheckImageCode struct {
	stream client.Stream
}

func (x *CodeIsCheckImageCode) Close() error {
	return x.stream.Close()
}

func (x *CodeIsCheckImageCode) Context() context.Context {
	return x.stream.Context()
}

func (x *CodeIsCheckImageCode) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *CodeIsCheckImageCode) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *codeService) VerifyImageCode(ctx context.Context, in *VerifyImageCodeRequest, opts ...client.CallOption) (*VerifyImageCodeResponse, error) {

	req := c.c.NewRequest(c.name, "Code.VerifyImageCode", in)
	out := new(VerifyImageCodeResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Code_VerifyImageCodeService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type CodeVerifyImageCode struct {
	stream client.Stream
}

func (x *CodeVerifyImageCode) Close() error {
	return x.stream.Close()
}

func (x *CodeVerifyImageCode) Context() context.Context {
	return x.stream.Context()
}

func (x *CodeVerifyImageCode) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *CodeVerifyImageCode) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *codeService) GetSendStatus(ctx context.Context, in *GetSendStatusRequest, opts ...client.CallOption) (*GetSendStatusResponse, error) {

	req := c.c.NewRequest(c.name, "Code.GetSendStatus", in)
	out := new(GetSendStatusResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Code_GetSendStatusService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type CodeGetSendStatus struct {
	stream client.Stream
}

func (x *CodeGetSendStatus) Close() error {
	return x.stream.Close()
}

func (x *CodeGetSendStatus) Context() context.Context {
	return x.stream.Context()
}

func (x *CodeGetSendStatus) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *CodeGetSendStatus) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

// Server API for Code service
type CodeHandler interface {
	SendCode(context.Context, *SendCodeRequest, *SendCodeResponse) error
	Verify(context.Context, *VerifyRequest, *VerifyResponse) error
	IsCheckImageCode(context.Context, *IsCheckImageCodeRequest, *IsCheckImageCodeResponse) error
	VerifyImageCode(context.Context, *VerifyImageCodeRequest, *VerifyImageCodeResponse) error
	GetSendStatus(context.Context, *GetSendStatusRequest, *GetSendStatusResponse) error
}

func RegisterCodeHandler(s server.Server, hdlr CodeHandler, opts ...server.HandlerOption) error {
	type code interface {
		SendCode(ctx context.Context, in *SendCodeRequest, out *SendCodeResponse) error
		Verify(ctx context.Context, in *VerifyRequest, out *VerifyResponse) error
		IsCheckImageCode(ctx context.Context, in *IsCheckImageCodeRequest, out *IsCheckImageCodeResponse) error
		VerifyImageCode(ctx context.Context, in *VerifyImageCodeRequest, out *VerifyImageCodeResponse) error
		GetSendStatus(ctx context.Context, in *GetSendStatusRequest, out *GetSendStatusResponse) error
	}

	type Code struct {
		code
	}
	h := &codeHandler{hdlr}
	opts = append(opts, server.EndpointMetadata("SendCode", map[string]string{"POST": "/user/code/send-code"}))
	opts = append(opts, server.EndpointMetadata("Verify", map[string]string{"POST": "/user/code/verify"}))
	opts = append(opts, server.EndpointMetadata("IsCheckImageCode", map[string]string{"POST": "/user/code/is-check-image-code"}))
	opts = append(opts, server.EndpointMetadata("VerifyImageCode", map[string]string{"POST": "/user/code/verify-image-code"}))
	opts = append(opts, server.EndpointMetadata("GetSendStatus", map[string]string{"POST": "/user/code/get-send-status"}))
	return s.Handle(s.NewHandler(&Code{h}, opts...))
}

func init() { catdog_data.Add("RegisterCodeHandler", RegisterCodeHandler) }

type codeHandler struct {
	CodeHandler
}

func (h *codeHandler) SendCode(ctx context.Context, in *SendCodeRequest, out *SendCodeResponse) error {
	return h.CodeHandler.SendCode(ctx, in, out)
}

func (h *codeHandler) Verify(ctx context.Context, in *VerifyRequest, out *VerifyResponse) error {
	return h.CodeHandler.Verify(ctx, in, out)
}

func (h *codeHandler) IsCheckImageCode(ctx context.Context, in *IsCheckImageCodeRequest, out *IsCheckImageCodeResponse) error {
	return h.CodeHandler.IsCheckImageCode(ctx, in, out)
}

func (h *codeHandler) VerifyImageCode(ctx context.Context, in *VerifyImageCodeRequest, out *VerifyImageCodeResponse) error {
	return h.CodeHandler.VerifyImageCode(ctx, in, out)
}

func (h *codeHandler) GetSendStatus(ctx context.Context, in *GetSendStatusRequest, out *GetSendStatusResponse) error {
	return h.CodeHandler.GetSendStatus(ctx, in, out)
}
