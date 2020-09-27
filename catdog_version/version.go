package catdog_version

import (
	"github.com/imdario/mergo"
	"github.com/pubgo/xerror"
	"sync"
)

var versions sync.Map

type M = map[string]string

func Register(name string, data M) {
	versions.Store(name, data)
}

func Get(name string) (v M) {
	m, ok := versions.Load(name)
	if ok {
		xerror.Exit(mergo.Map(&v, m))
		return
	}
	return
}

func List() map[string]M {
	ms := make(map[string]M)
	versions.Range(func(key, value interface{}) bool {
		var v M
		xerror.Exit(mergo.Map(&v, value))
		ms[key.(string)] = v
		return true
	})
	return ms
}
