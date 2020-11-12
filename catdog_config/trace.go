package catdog_config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"github.com/pubgo/xprocess"
)

func init() {
	// debug and trace
	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !Trace {
			return
		}

		//fmt.Println(tag)
		//fmt.Println()

		var data = make(map[string]interface{})
		xerror.Panic(json.Unmarshal(LoadBytes(), &data))
		xlog.Debug("config trace")
		fmt.Println(catdog_util.MarshalIndent(data))
		fmt.Println()

		xlog.Debug("deps trace")
		fmt.Println(dix.Graph())
		fmt.Println()

		xlog.Debug("goroutine trace")
		data = make(map[string]interface{})
		xerror.Panic(json.Unmarshal([]byte(xprocess.Stack()), &data))
		fmt.Println(catdog_util.MarshalIndent(data))
		fmt.Println()

		xlog.Debug("env trace")
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, Domain) {
				fmt.Println(env)
			}
		}
		fmt.Println()
	}))
}
