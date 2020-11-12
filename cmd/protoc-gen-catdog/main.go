package main

import (
	"log"

	"github.com/pubgo/xerror"
	"github.com/pubgo/xprotogen/gen"
)

func main() {
	m := gen.New("catdog")
	m.Parameter(func(key, value string) {
		log.Println("params:", key, "=", value)
	})

	xerror.Exit(m.Init(func(fd *gen.FileDescriptor) {
		header(fd)
		for _, ss := range fd.GetService() {
			service(ss)
		}
	}))
}
