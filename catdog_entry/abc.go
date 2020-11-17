package catdog_entry

import (
	"github.com/asim/nitro/v3/server"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Entry interface {
	Init() error
	Start() error
	Stop() error
	Description(description ...string) error
	Version(v string) error
	Flags(fn func(flags *pflag.FlagSet)) error
	Commands(commands ...*cobra.Command) error
	Options() Options
	Handler(handler interface{}, opts ...server.HandlerOption) error
}

type Option func(o *Options)
type Options struct {
	App         *fiber.App
	Initialized bool
	RestAddr    string
	Name        string
	Version     string
	RunCommand  *cobra.Command
	Command     *cobra.Command
}
