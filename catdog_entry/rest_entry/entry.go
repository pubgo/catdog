package rest_entry

import (
	"context"
	"github.com/pubgo/catdog/catdog_entry/base_entry"
	"time"

	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/internal/plugins/server/server_http"
)

type entry struct {
	catdog_entry.Entry
	c client.Client
}

func (t *entry) Group(relativePath string, handlers ...fiber.Handler) fiber.Router {
	return t.Options().App.Group(relativePath, handlers...)
}

func (t *entry) middleware() []interface{} {
	return []interface{}{
		middleware.Recover(),
		middleware.Logger(middleware.LoggerConfig{
			Format:     "#${pid} - ${time} ${status} - ${latency} ${method} ${path}\n",
			TimeFormat: time.RFC3339,
		}),
	}
}

func newEntry(name string) *entry {
	ent := &entry{
		Entry: base_entry.New(name, &entryServerWrapper{Server: server_http.NewServer(server.Context(context.Background()))}),
	}
	ent.Options().App.Use(ent.middleware()...)

	return ent
}

func New(name string) catdog_entry.Entry {
	return newEntry(name)
}
