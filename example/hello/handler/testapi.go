package handler

import (
	"context"
	helloworld "github.com/pubgo/catdog/example/hello/proto"
)

func NewTestAPIHandler() helloworld.TestApiHandler {
	return &testapiHandler{}
}

type testapiHandler struct {
}

func (h *testapiHandler) Version(ctx context.Context, in *helloworld.TestReq, out *helloworld.TestApiOutput) error {
	log.Infof("Received Helloworld.Call request, name: %s", in.Input)
	out.Msg = in.Input
	return nil
}

func (h *testapiHandler) VersionTest(ctx context.Context, in *helloworld.TestReq, out *helloworld.TestApiOutput) error {
	out.Msg = in.Input + "_test"
	return nil
}
