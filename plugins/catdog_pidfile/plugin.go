package catdog_pidfile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pubgo/catdog/catdog_config"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Exit(catdog_plugin.Register(&catdog_plugin.Base{
		Name: "pidfile",
		OnInit: func() {
			// 检查存放pid的目录是否存在, 不存在就创建
			xerror.Panic(catdog_abc.WithBeforeStart(func() {
				pidPath := filepath.Dir(GetPidPath())
				if !catdog_util.PathExist(pidPath) {
					xerror.Exit(os.MkdirAll(pidPath, pidPerm))
				}
			}))

			// 保存pid到文件当中
			xerror.Panic(catdog_abc.WithAfterStart(func() {
				xerror.Panic(SavePid())

				if catdog_config.Trace {
					pid, err := GetPid()
					xerror.Panic(err)
					xlog.Debugf("path, pid trace")
					fmt.Println(GetPidPath(), pid)
				}
			}))
		},
	}))
}
