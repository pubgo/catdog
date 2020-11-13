package rpc_entry

import (
	"context"
	"fmt"
	"github.com/pubgo/catdog/catdog_entry/base_entry"
	"net/http"
	"time"

	grpcC "github.com/asim/nitro-plugins/client/grpc/v3"
	"github.com/asim/nitro-plugins/server/grpc/v3"
	"github.com/asim/nitro/v3/client"
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/pubgo/catdog/catdog_entry"
	"github.com/pubgo/catdog/internal/catdog_abc"
	"github.com/pubgo/xerror"
	"github.com/pubgo/xprocess"
)

type entry struct {
	catdog_entry.Entry
	c   client.Client
	app *fiber.App
}

func (t *entry) Init() (err error) {
	defer xerror.RespErr(&err)

	xerror.Panic(t.Entry.Init())
	return nil
}

func (t *entry) Start() (err error) {
	defer xerror.RespErr(&err)

	cancel := xprocess.Go(func(ctx context.Context) (err error) {
		defer xerror.RespErr(&err)

		addr := t.Options().RestAddr
		log.Infof("Server [http] Listening on http://%s", addr)
		xerror.Exit(t.app.Listen(addr))
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
		Entry: base_entry.New(
			name,
			grpc.NewServer(server.Context(context.Background()), server.Name(name)),
			app,
		),
		c:   grpcC.NewClient(),
		app: app,
	}
	ent.app.Use(ent.middleware()...)

	return ent
}

func New(name string) catdog_entry.Entry {
	return newEntry(name)
}
