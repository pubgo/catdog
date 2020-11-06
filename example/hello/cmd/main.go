package main

import (
	"github.com/pubgo/catdog"
	"github.com/pubgo/xerror"

	"github.com/pubgo/catdog/example/hello"
)

func main() {
	xerror.Exit(catdog.Init())
	xerror.Exit(catdog.Run(hello.GetEntry()))
}
