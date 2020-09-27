package main

import (
	"github.com/pubgo/catdog"
	"github.com/pubgo/catdog/example/hello/entry"
	"github.com/pubgo/xerror"
)

func main() {
	xerror.Exit(catdog.Run(
		entry.GetEntry(),
	))
}
