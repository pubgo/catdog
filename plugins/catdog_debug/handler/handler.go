// Package handler implements catdog_service catdog_debug_plugin handler embedded in nitro services
package handler

import (
	"context"
	"time"

	"github.com/asim/nitro/v3/debug/log"
	"github.com/asim/nitro/v3/debug/stats"
	"github.com/asim/nitro/v3/debug/trace"

	pb "github.com/pubgo/catdog/plugins/catdog_debug/proto"

	memLog "github.com/asim/nitro/v3/debug/log/memory"
	"github.com/asim/nitro/v3/debug/profile"
	memStats "github.com/asim/nitro/v3/debug/stats/memory"
	memTrace "github.com/asim/nitro/v3/debug/trace/memory"
)

var (
	DefaultLog      log.Log         = memLog.NewLog()
	DefaultTracer   trace.Tracer    = memTrace.NewTracer()
	DefaultStats    stats.Stats     = memStats.NewStats()
	DefaultProfiler profile.Profile = nil
)

// NewHandler returns an instance of the Debug handler
func NewHandler() pb.DebugHandler {
	return &Debug{
		log:   DefaultLog,
		stats: DefaultStats,
		trace: DefaultTracer,
	}
}

type Debug struct {
	// the logger for retrieving logs
	log log.Log
	// the stats collector
	stats stats.Stats
	// the tracer
	trace trace.Tracer
}

func (d *Debug) Health(ctx context.Context, req *pb.HealthRequest, rsp *pb.HealthResponse) error {
	rsp.Status = "ok"
	return nil
}

func (d *Debug) Stats(ctx context.Context, req *pb.StatsRequest, rsp *pb.StatsResponse) error {
	stats, err := d.stats.Read()
	if err != nil {
		return err
	}

	if len(stats) == 0 {
		return nil
	}

	// write the response values
	rsp.Timestamp = uint64(stats[0].Timestamp)
	rsp.Started = uint64(stats[0].Started)
	rsp.Uptime = uint64(stats[0].Uptime)
	rsp.Memory = stats[0].Memory
	rsp.Gc = stats[0].GC
	rsp.Threads = stats[0].Threads
	rsp.Requests = stats[0].Requests
	rsp.Errors = stats[0].Errors

	return nil
}

func (d *Debug) Trace(ctx context.Context, req *pb.TraceRequest, rsp *pb.TraceResponse) error {
	traces, err := d.trace.Read(trace.ReadTrace(req.Id))
	if err != nil {
		return err
	}

	for _, t := range traces {
		var typ pb.SpanType
		switch t.Type {
		case trace.SpanTypeRequestInbound:
			typ = pb.SpanType_INBOUND
		case trace.SpanTypeRequestOutbound:
			typ = pb.SpanType_OUTBOUND
		}
		rsp.Spans = append(rsp.Spans, &pb.Span{
			Trace:    t.Trace,
			Id:       t.Id,
			Parent:   t.Parent,
			Name:     t.Name,
			Started:  uint64(t.Started.UnixNano()),
			Duration: uint64(t.Duration.Nanoseconds()),
			Type:     typ,
			Metadata: t.Metadata,
		})
	}

	return nil
}

// Log returns some log lines
func (d *Debug) Log(ctx context.Context, req *pb.LogRequest, rsp *pb.LogResponse) error {
	var options []log.ReadOption

	since := time.Unix(req.Since, 0)
	if !since.IsZero() {
		options = append(options, log.Since(since))
	}

	count := int(req.Count)
	if count > 0 {
		options = append(options, log.Count(count))
	}

	// get the log records
	records, err := d.log.Read(options...)
	if err != nil {
		return err
	}

	for _, record := range records {
		// copy metadata
		metadata := make(map[string]string)
		for k, v := range record.Metadata {
			metadata[k] = v
		}
		// send record
		rsp.Records = append(rsp.Records, &pb.Record{
			Timestamp: record.Timestamp.Unix(),
			Message:   record.Message.(string),
			Metadata:  metadata,
		})
	}

	return nil
}
