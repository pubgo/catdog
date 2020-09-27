package catdog_abc

import (
	"github.com/micro/go-micro/v3/broker"
	"github.com/micro/go-micro/v3/registry"
	"github.com/pubgo/catdog/catdog_broker"
	"github.com/pubgo/catdog/catdog_client"
	"github.com/pubgo/catdog/catdog_model"
	"github.com/pubgo/catdog/catdog_registry"
	"github.com/pubgo/catdog/catdog_server"
	"github.com/pubgo/dix"
	"github.com/pubgo/xerror"
	"os"

	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/model"
	"github.com/micro/go-micro/v3/server"
)

// CatDog is an interface for a micro catdog_service
type CatDog interface {
	// The catdog_service name
	Name() string
	// initCatDog initialises options
	Init(...Option)
	// Options returns the current options
	Options() Options
	// Client is used to call services
	Client() client.Client
	// Server is for handling requests and events
	Server() server.Server
	// Model is used to access data
	Model() model.Model
	// Run the catdog_service
	Run() error
	// The catdog_service implementation
	String() string
}

type Option func(*Options)
type Options struct {
	IsSignal bool
	Signal   os.Signal

	Broker   broker.Broker
	Client   client.Client
	Model    model.Model
	Registry registry.Registry
	Server   server.Server

	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error
}

func NewOption(opts ...Option) Options {
	options := Options{
		IsSignal: true,
		Broker:   catdog_broker.Default,
		Client:   catdog_client.Default,
		Model:    catdog_model.Default,
		Registry: catdog_registry.Default,
		Server:   catdog_server.Default,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return options
}

func Watch(fn func(cat CatDog) error) error {
	return xerror.Wrap(dix.Dix(fn))
}
