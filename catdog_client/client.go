package catdog_client

import (
	"github.com/micro/go-micro/v3/client"
)

type wrapper struct {
	client.Client
}
