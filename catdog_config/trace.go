package catdog_config

import (
	"encoding/json"
	"fmt"

	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
)

func init() {
	// debug and trace
	xerror.Exit(catdog_abc.WithAfterStart(func() {
		if !Trace {
			return
		}

		var data = make(map[string]interface{})
		xerror.Panic(json.Unmarshal(LoadBytes(), &data))
		xlog.Debug("config trace")
		fmt.Println(catdog_util.MarshalIndent(data))

		xlog.Debug("deps trace")
		fmt.Println(dix.Graph())
	}))
}
