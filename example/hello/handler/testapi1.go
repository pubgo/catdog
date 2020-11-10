package handler

import (
	"context"
	helloworld "github.com/pubgo/catdog/example/hello/proto"
)

func NewTestAPIHandler1() helloworld.TestApiV2Handler {
	return &testapiHandler1{}
}

type testapiHandler1 struct {
}

func (h *testapiHandler1) Version(ctx context.Context, in *helloworld.TestReq, out *helloworld.TestApiOutput) error {
	log.Infof("Received Helloworld.Call request, name: %s", in.Input)
	out.Msg = in.Input
	return nil
}

func (h *testapiHandler1) VersionTest(ctx context.Context, in *helloworld.TestReq, out *helloworld.TestApiOutput) error {
	out.Msg = in.Input + "_test"
	return nil
}
