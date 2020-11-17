package entry3

import (
	"fmt"
	"github.com/pubgo/catdog"
	"github.com/pubgo/catdog/version"
	"github.com/pubgo/xerror"
)

func GetEntry() catdog.Entry {
	ent := catdog.NewCtlEntry("hello3")
	xerror.Exit(ent.Description("hello3 ctl 服务"))
	xerror.Exit(ent.Version(version.Version))

	xerror.Exit(ent.Handler(func() error {
		fmt.Println("ok")
		return nil
	}))

	xerror.Exit(ent.Handler(func() error {
		fmt.Println("ok")
		return nil
	}))

	return ent
}
