package catdog_abc

import (
	"github.com/pubgo/dix"
)

func WithBeforeStart(fn func()) error { return dix.WithBeforeStart(fn) }
func WithAfterStart(fn func()) error  { return dix.WithAfterStart(fn) }
func WithBeforeStop(fn func()) error  { return dix.WithBeforeStop(fn) }
func WithAfterStop(fn func()) error   { return dix.WithAfterStop(fn) }
