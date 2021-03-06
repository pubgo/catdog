// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: example/proto/login/login.proto
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

// Client API for Login service
type LoginService interface {
	Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error)
	Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...client.CallOption) (*AuthenticateResponse, error)
}

type loginService struct {
	c    client.Client
	name string
}

func NewLoginService(name string, c client.Client) LoginService {
	return &loginService{
		c:    c,
		name: name,
	}
}
func (c *loginService) Login(ctx context.Context, in *LoginRequest, opts ...client.CallOption) (*LoginResponse, error) {

	req := c.c.NewRequest(c.name, "Login.Login", in)
	out := new(LoginResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Login_LoginService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type LoginLogin struct {
	stream client.Stream
}

func (x *LoginLogin) Close() error {
	return x.stream.Close()
}

func (x *LoginLogin) Context() context.Context {
	return x.stream.Context()
}

func (x *LoginLogin) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *LoginLogin) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

func (c *loginService) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...client.CallOption) (*AuthenticateResponse, error) {

	req := c.c.NewRequest(c.name, "Login.Authenticate", in)
	out := new(AuthenticateResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Stream auxiliary types and methods.
type Login_AuthenticateService interface {
	Context() context.Context
	SendMsg(interface{}) error
	RecvMsg(interface{}) error
	Close() error
}
type LoginAuthenticate struct {
	stream client.Stream
}

func (x *LoginAuthenticate) Close() error {
	return x.stream.Close()
}

func (x *LoginAuthenticate) Context() context.Context {
	return x.stream.Context()
}

func (x *LoginAuthenticate) SendMsg(m interface{}) error {
	return x.stream.Send(m)
}

func (x *LoginAuthenticate) RecvMsg(m interface{}) error {
	return x.stream.Recv(m)
}

// Server API for Login service
type LoginHandler interface {
	Login(context.Context, *LoginRequest, *LoginResponse) error
	Authenticate(context.Context, *AuthenticateRequest, *AuthenticateResponse) error
}

func RegisterLoginHandler(s server.Server, hdlr LoginHandler, opts ...server.HandlerOption) error {
	type login interface {
		Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error
		Authenticate(ctx context.Context, in *AuthenticateRequest, out *AuthenticateResponse) error
	}

	type Login struct {
		login
	}
	h := &loginHandler{hdlr}
	opts = append(opts, server.EndpointMetadata("Login", map[string]string{"POST": "/user/login/login"}))
	opts = append(opts, server.EndpointMetadata("Authenticate", map[string]string{"POST": "/user/login/authenticate"}))
	return s.Handle(s.NewHandler(&Login{h}, opts...))
}

func init() { catdog_data.Add("RegisterLoginHandler", RegisterLoginHandler) }

type loginHandler struct {
	LoginHandler
}

func (h *loginHandler) Login(ctx context.Context, in *LoginRequest, out *LoginResponse) error {
	return h.LoginHandler.Login(ctx, in, out)
}

func (h *loginHandler) Authenticate(ctx context.Context, in *AuthenticateRequest, out *AuthenticateResponse) error {
	return h.LoginHandler.Authenticate(ctx, in, out)
}
