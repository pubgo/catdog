package entry

import (
	"github.com/pubgo/catdog"
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/example/hello/handler"
	"github.com/pubgo/catdog/version"
)

func GetEntry() catdog.Entry {
	//catdog_config.Domain = "dev"

	ent := catdog.NewRpcEntry("hello")
	xerror.Exit(ent.Description("hello grpc 服务"))
	xerror.Exit(ent.Version(version.Version))

	xerror.Exit(ent.Handler(handler.NewHelloworld()))
	xerror.Exit(ent.Handler(handler.NewTestAPIHandler()))
	//xerror.Exit(ent.Handler(handler.NewTestAPIHandler1()))
	return ent
}
