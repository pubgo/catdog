package catdog_log

import (
	"fmt"
	"github.com/pubgo/xerror"
	"sync"

	"github.com/micro/go-micro/v3/logger"
	"github.com/pubgo/xlog"
)

func init() {
	xerror.Exit(Watch(func(logs xlog.XLog) {
		logger.DefaultLogger = newMicroLogger(logs)
	}))
}

func newMicroLogger(xl xlog.XLog) logger.Logger {
	return &microLogger{log: xl.Named("micro", xlog.AddCallerSkip(2))}
}

type microLogger struct {
	sync.RWMutex
	log  xlog.XLog
	opts logger.Options
}

func (h *microLogger) Init(opts ...logger.Option) error {
	h.Lock()
	defer h.Unlock()

	var hOpts = h.opts
	for _, o := range opts {
		o(&hOpts)
	}

	if h.opts.Fields != nil {
		var fields []xlog.Field
		for k, v := range h.opts.Fields {
			fields = append(fields, xlog.Any(k, v))
		}
		if len(fields) > 0 {
			h.log = h.log.With(fields...)
		}
	}

	return nil
}

func (h *microLogger) Options() logger.Options {
	return h.opts
}

func (h *microLogger) Fields(fields map[string]interface{}) logger.Logger {
	h.Lock()
	defer h.Unlock()

	if fields != nil {
		opts := h.opts
		if opts.Fields == nil {
			opts.Fields = make(map[string]interface{})
		}

		var logFields []xlog.Field
		for k, v := range fields {
			logFields = append(logFields, xlog.Any(k, v))
			opts.Fields[k] = v
		}

		if len(logFields) > 0 {
			return &microLogger{log: h.log.With(logFields...), opts: opts}
		}
	}

	return h
}

func (h *microLogger) Log(level logger.Level, v ...interface{}) {
	log := h.log
	switch level {
	case logger.TraceLevel:
		log.Debug(fmt.Sprint(v...))
	case logger.DebugLevel:
		log.Debug(fmt.Sprint(v...))
	case logger.InfoLevel:
		log.Info(fmt.Sprint(v...))
	case logger.WarnLevel:
		log.Warn(fmt.Sprint(v...))
	case logger.ErrorLevel:
		log.Error(fmt.Sprint(v...))
	case logger.FatalLevel:
		log.Fatal(fmt.Sprint(v...))
	}
}

func (h *microLogger) Logf(level logger.Level, format string, v ...interface{}) {
	log := h.log
	switch level {
	case logger.TraceLevel:
		log.Debug(fmt.Sprintf(format, v...))
	case logger.DebugLevel:
		log.Debug(fmt.Sprintf(format, v...))
	case logger.InfoLevel:
		log.Info(fmt.Sprintf(format, v...))
	case logger.WarnLevel:
		log.Warn(fmt.Sprintf(format, v...))
	case logger.ErrorLevel:
		log.Error(fmt.Sprintf(format, v...))
	case logger.FatalLevel:
		log.Fatal(fmt.Sprintf(format, v...))
	}
}

func (h *microLogger) String() string {
	return "micro_logger"
}
