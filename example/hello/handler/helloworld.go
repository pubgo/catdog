package handler

import (
	"context"

	helloworld "github.com/pubgo/catdog/example/hello/proto"
)

func NewHelloworld() helloworld.HelloworldHandler {
	return &Helloworld{}
}

type Helloworld struct{}

// Call is a single request handler called via catdog_client_plugin.Call or the generated catdog_client_plugin code
func (e *Helloworld) Call(ctx context.Context, req *helloworld.Request, rsp *helloworld.Response) error {
	log.Infof("Received Helloworld.Call request, name: %s", req.Name)
	rsp.Msg = req.Name
	return nil
}
