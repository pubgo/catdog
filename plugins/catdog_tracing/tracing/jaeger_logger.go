package tracing

import (
	"github.com/pubgo/xlog"
	"io"

	jaegerlog "github.com/uber/jaeger-client-go/log"
)

var (
	_ io.Writer        = (*logger)(nil)
	_ jaegerlog.Logger = (*logger)(nil)
)

type logger struct{}

func NewLogger() *logger {
	return &logger{}
}

func (l *logger) Error(msg string) {
	xlog.Errorf("ERROR: %s", msg)
}

// Infof logs a message at info priority
func (l *logger) Infof(msg string, args ...interface{}) {
	xlog.Infof(msg, args...)
}

func (l *logger) Write(b []byte) (n int, err error) {
	xlog.With(xlog.String("type", "tracing")).Info(string(b))
	return len(b), nil
}
