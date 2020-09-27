package catdog_tracing_plugin

import (
	"github.com/pubgo/catdog/catdog_handler"
	"github.com/pubgo/catdog/catdog_plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"strings"
)

var tracerCloser io.Closer

var _ catdog_plugin.Plugin = (*Plugin)(nil)

type Plugin struct {
	name string
}

func (p *Plugin) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(p.name, pflag.PanicOnError)

	/*
		cli.StringFlag{
			Name:   "jaeger_service_name",
			Value:  "",
			Usage:  "The service name",
			EnvVar: "JAEGER_SERVICE_NAME",
		}
			cli.StringFlag{
				Name:   "jaeger_disabled",
				Usage:  "Whether the Plugin is disabled or not. If true, the default opentracing.NoopTracer is used.",
				EnvVar: "JAEGER_DISABLED",
			}
			cli.StringFlag{
				Name:   "jaeger_RPC_metrics",
				Usage:  "Whether to store RPC metrics",
				EnvVar: "JAEGER_RPC_METRICS",
			}
			cli.StringFlag{
				Name:   "jaeger_tags",
				Usage:  "A comma separated list of name = value Plugin level tags, which get added to all reported spans.",
				EnvVar: "JAEGER_TAGS",
			}
			cli.StringFlag{
				Name:   "jaeger_sampler_type",
				Value:  "const",
				Usage:  "The sampler type",
				EnvVar: "JAEGER_SAMPLER_TYPE",
			}
			cli.StringFlag{
				Name:   "jaeger_sampler_param",
				Value:  "1",
				Usage:  "The sampler parameter (number)",
				EnvVar: "JAEGER_SAMPLER_PARAM",
			}
			cli.StringFlag{
				Name:   "jaeger_sampler_manager_host_port",
				Usage:  "The HTTP endpoint when using the remote sampler, i.e. http://jaeger-agent:5778/sampling",
				EnvVar: "JAEGER_SAMPLER_MANAGER_HOST_PORT",
			}
			cli.StringFlag{
				Name:   "jaeger_sampler_max_operations",
				Usage:  "The maximum number of operations that the sampler will keep track of",
				EnvVar: "JAEGER_SAMPLER_MAX_OPERATIONS",
			}
			cli.StringFlag{
				Name:   "jaeger_sampler_refresh_interval",
				Usage:  "How often the remotely controlled sampler will poll jaeger-agent for the appropriate sampling strategy, with units",
				EnvVar: "JAEGER_SAMPLER_REFRESH_INTERVAL",
			}
			cli.StringFlag{
				Name:   "jaeger_reporter_max_queue_size",
				Usage:  "The reporter's maximum queue size",
				EnvVar: "JAEGER_REPORTER_MAX_QUEUE_SIZE",
			}
			cli.StringFlag{
				Name:   "jaeger_reporter_flush_interval",
				Usage:  "The reporter's flush interval, with units, e.g. 500ms or 2s (valid units)",
				EnvVar: "JAEGER_REPORTER_FLUSH_INTERVAL",
			}
			cli.StringFlag{
				Name:   "jaeger_reporter_log_spans",
				Usage:  "Whether the reporter should also log the spans",
				EnvVar: "JAEGER_REPORTER_LOG_SPANS",
			}
			cli.StringFlag{
				Name:   "jaeger_endpoint",
				Usage:  "The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces",
				EnvVar: "JAEGER_ENDPOINT",
			}
			cli.StringFlag{
				Name:   "jaeger_user",
				Usage:  "Username to send as part of Basic authentication to the collector endpoint",
				EnvVar: "JAEGER_USER",
			}
			cli.StringFlag{
				Name:   "jaeger_password",
				Usage:  "Password to send as part of Basic authentication to the collector endpoint",
				EnvVar: "JAEGER_PASSWORD",
			}
			cli.StringFlag{
				Name:   "jaeger_agent_host",
				Usage:  "The hostname for communicating with agent via UDP",
				EnvVar: "JAEGER_AGENT_HOST",
			}
			cli.StringFlag{
				Name:   "jaeger_agent_port",
				Usage:  "The port for communicating with agent via UDP",
				EnvVar: "JAEGER_AGENT_PORT",
			}
			cli.StringFlag{
				Name:   "jaeger_save_as_file",
				Value:  "1",
				Usage:  "Use the self-implemented transport. Save the spans into a log file with http reporter's rules. Default is true.",
				EnvVar: "JAEGER_SAVE_AS_FILE",
			}
	*/

	return flags
}

func (p *Plugin) Commands() *cobra.Command {
	return nil
}

func (p *Plugin) Handler() *catdog_handler.Handler {
	return nil
}

func (p *Plugin) Init1() {
	/*
		for _, v := range p.Flags() {
			sflag, ok := v.(cli.StringFlag)

			if !ok {
				log.Error("Tracer flag set failed. %v", v.GetName())
				continue
			}

			name := sflag.GetName()
			value := ctx.String(name)
			if value == "" {
				continue
			}

			if err := os.Setenv(sflag.EnvVar, value); err != nil {
				return err
			}
		}

		// Tags is special, jaeger config will panic
		envTags := ctx.String("jaeger_tags")
		if envTags != "" {
			if !isTagsvalid(envTags) {
				return fmt.Errorf("tags is not valid, get tag %v, expected 'key=value' ", envTags)
			}
		}

		// * jaeger service name must set before Plugin init.
		if ctx.String("jaeger_service_name") == "" {
			if err := setDefaultJaegerServiceName(ctx); err != nil {
				return err
			}
		}

		// initCatDog Plugin
		cfg, err := jaegercfg.FromEnv()
		if err != nil {
			return err
		}

		var tracer opentracing.Tracer
		var opts []jaegercfg.Option
		opts = append(opts, jaegercfg.Logger(NewLogger()))

		// Default is save to file
		if ctx.String("jaeger_endpoint") == "" {
			if e := ctx.String("jaeger_save_as_file"); e != "" {
				if value, err := strconv.ParseBool(e); err == nil && value {
					r := jaeger.NewRemoteReporter(transport.NewIOTransport(NewLogger(), 1)) // Using logging transport
					opts = append(opts, jaegercfg.Reporter(r))
				}
			}
		}

		tracer, tracerCloser, err = cfg.NewTracer(opts...)
		if err != nil {
			return err
		}

		opentracing.SetGlobalTracer(tracer)
		return err
	*/
}

func (p *Plugin) String() string {
	return p.name
}

func isTagsvalid(tags string) bool {
	pairs := strings.Split(tags, ",")
	for _, p := range pairs {
		if !strings.Contains(p, "=") {
			return false
		}
	}
	return true
}

func NewPlugin() catdog_plugin.Plugin {
	p := &Plugin{name: "tracing"}
	return p
}

func Close() error {
	if tracerCloser == nil {
		return nil
	}

	return tracerCloser.Close()
}
