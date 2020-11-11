package entry

import (
	"github.com/pubgo/catdog"
	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/example/hello/handler"
	"github.com/pubgo/catdog/version"
	"github.com/pubgo/xerror"
)

func GetEntry() catdog.Entry {
	//catdog_config.Domain = "dev"
	catdog_config.Project = "hello"

	ent := catdog.NewEntry()
	xerror.Exit(ent.Name(catdog_config.Project, "hello 服务"))
	xerror.Exit(ent.Version(version.Version))

	xerror.Exit(ent.Handler(handler.NewHelloworld()))
	xerror.Exit(ent.Handler(handler.NewTestAPIHandler()))
	//xerror.Exit(ent.Handler(handler.NewTestAPIHandler1()))
	return ent
}
