package main

import (
	"github.com/pubgo/catdog"
	"github.com/pubgo/catdog/example/hello/entry2"
	"github.com/pubgo/catdog/example/hello/entry3"
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/example/hello/entry"
	"github.com/pubgo/catdog/example/hello/entry1"
)

func main() {
	xerror.Exit(catdog.Init())
	xerror.Exit(catdog.Run(
		entry.GetEntry(),
		entry1.GetEntry(),
		entry2.GetEntry(),
		entry3.GetEntry(),
	))
}
