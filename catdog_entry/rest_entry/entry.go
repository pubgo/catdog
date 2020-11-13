package rest_entry

import (
	"context"
	"fmt"
	"github.com/pubgo/catdog/catdog_entry/base_entry"
	"github.com/pubgo/xprocess"
	"net/http"
	"time"

	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/catdog/internal/plugins/server/server_http"
	"github.com/pubgo/xerror"
)

type entry struct {
	catdog_entry.Entry
	c    client.Client
	app  *fiber.App
	addr string
}

func (t *entry) Start() (err error) {
	defer xerror.RespErr(&err)

	cancel := xprocess.Go(func(ctx context.Context) (err error) {
		defer xerror.RespErr(&err)
		log.Infof("Server [http] Listening on http://%s", t.addr)
		xerror.Exit(t.app.Listen(t.addr))
		log.Infof("Server [http] Closed OK")
		return nil
	})

	xerror.Exit(catdog_abc.WithBeforeStop(func() {
		xerror.Panic(cancel())
		if err := t.app.Shutdown(); err != nil && err != http.ErrServerClosed {
			fmt.Println(xerror.Parse(err).Println())
		}
	}))

	return xerror.Wrap(t.Entry.Start())
}

func (t *entry) Stop() (err error) {
	defer xerror.RespErr(&err)
	return xerror.Wrap(t.Entry.Stop())
}

func (t *entry) Group(relativePath string, handlers ...fiber.Handler) fiber.Router {
	return t.app.Group(relativePath, handlers...)
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
	app := fiber.New()

	ent := &entry{
		Entry: base_entry.New(name, &entryServerWrapper{Server: server_http.NewServer(server.Context(context.Background()))}, app),
		app:   app,
		addr:  ":8080",
	}
	ent.app.Use(ent.middleware()...)

	return ent
}

func New(name string) catdog_entry.Entry {
	return newEntry(name)
}
