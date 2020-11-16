package catdog_abc

import (
	"bytes"
	"fmt"
	"github.com/pubgo/catdog/catdog_util"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xlog"
	"sync"
)

var beforeStart []func()
var afterStart []func()
var beforeStop []func()
var afterStop []func()

func GetBeforeStart() []func() {
	var data = make([]func(), len(beforeStart))
	mutex(func() { copy(data, beforeStart) })
	return data
}

func GetAfterStart() []func() {
	var data = make([]func(), len(afterStart))
	mutex(func() { copy(data, afterStart) })
	return data
}

func GetBeforeStop() []func() {
	var data = make([]func(), len(beforeStop))
	mutex(func() { copy(data, beforeStop) })
	return data
}

func GetAfterStop() []func() {
	var data = make([]func(), len(afterStop))
	mutex(func() { copy(data, afterStop) })
	return data
}

func Trace() {
	var buf = bytes.NewBuffer(nil)
	xlog.Debug("beforeStart trace")
	for _, fn := range GetBeforeStart() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		buf.WriteRune('\n')
	}
	xlog.Debug("afterStart trace")
	for _, fn := range GetAfterStop() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		buf.WriteRune('\n')
	}
	xlog.Debug("beforeStop trace")
	for _, fn := range GetBeforeStop() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
		buf.WriteRune('\n')
	}
	xlog.Debug("afterStop trace")
	for _, fn := range GetAfterStop() {
		buf.WriteString(xerror.PanicStr(catdog_util.CallerWithFunc(fn)))
	}
	fmt.Println(buf.String())
}

var mu sync.Mutex

func mutex(fn func()) {
	mu.Lock()
	defer mu.Unlock()
	fn()
}

func WithBeforeStart(fn func()) error {
	mutex(func() { beforeStart = append(beforeStart, fn) })
	return nil
}

func WithAfterStart(fn func()) error {
	mutex(func() { afterStart = append(afterStart, fn) })
	return nil
}

func WithBeforeStop(fn func()) error {
	mutex(func() { beforeStop = append(beforeStop, fn) })
	return nil
}

func WithAfterStop(fn func()) error {
	mutex(func() { afterStop = append(afterStop, fn) })
	return nil
}
