package catdog_action

import (
	"github.com/pubgo/xerror"
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

var mu sync.Mutex

func mutex(fn func()) {
	mu.Lock()
	defer mu.Unlock()
	fn()
}

func WithBeforeStart(fn func()) (err error) {
	defer xerror.RespErr(&err)
	mutex(func() { beforeStart = append(beforeStart, fn) })
	return nil
}

func WithAfterStart(fn func()) (err error) {
	defer xerror.RespErr(&err)
	mutex(func() { afterStart = append(afterStart, fn) })
	return nil
}

func WithBeforeStop(fn func()) (err error) {
	defer xerror.RespErr(&err)
	mutex(func() { beforeStop = append(beforeStop, fn) })
	return nil
}

func WithAfterStop(fn func()) (err error) {
	defer xerror.RespErr(&err)
	mutex(func() { afterStop = append(afterStop, fn) })
	return nil
}
