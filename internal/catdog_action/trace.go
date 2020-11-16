package catdog_action

import (
	"bytes"
	"fmt"

	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/xerror"
)

func Trace() {
	var buf = bytes.NewBuffer(nil)
	buf.WriteString("beforeStart trace\n")
	for _, fn := range GetBeforeStart() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		buf.WriteRune('\n')
	}

	buf.WriteString("afterStart trace\n")
	for _, fn := range GetAfterStop() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		buf.WriteRune('\n')
	}

	buf.WriteString("beforeStop trace\n")
	for _, fn := range GetBeforeStop() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		buf.WriteRune('\n')
	}

	buf.WriteString("afterStop trace\n")
	for _, fn := range GetAfterStop() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
	}
	fmt.Println(buf.String())
}
