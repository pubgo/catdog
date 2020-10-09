package catdog_model

import (
	"github.com/micro/go-micro/v3/model/mud"
)

var (
	Default = &wrapper{Model: mud.NewModel()}
)
