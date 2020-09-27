package catdog_tracing_plugin

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/micro/go-micro/v3/metadata"
	"github.com/smartystreets/goconvey/convey"
)

func TestNewContextWithOld(t *testing.T) {
	convey.Convey("New context from old one", t, func() {
		convey.Convey("nil context", func() {
			var old context.Context
			ctx := NewContextWithOld(old)
			convey.So(ctx, convey.ShouldBeNil)
		})

		convey.Convey("metadata && with no deadline", func() {
			var ctx context.Context
			convey.Convey("no metadata", func() {
				ctx = NewContextWithOld(context.Background())
				convey.So(ctx, convey.ShouldNotBeNil)
				md, ok := metadata.FromContext(ctx)
				convey.So(ok, convey.ShouldBeTrue)
				convey.So(md, convey.ShouldNotBeNil)
				convey.So(len(md), convey.ShouldEqual, 0)
			})

			convey.Convey("with metadata", func() {
				md := make(metadata.Metadata)
				md["key"] = "value"
				old := metadata.NewContext(context.Background(), md)
				ctx = NewContextWithOld(old)
				convey.So(ctx, convey.ShouldNotBeNil)
				md, ok := metadata.FromContext(ctx)
				convey.So(ok, convey.ShouldBeTrue)
				convey.So(md, convey.ShouldNotBeNil)
				convey.So(len(md), convey.ShouldEqual, 1)
				val, ok := md["key"]
				convey.So(ok, convey.ShouldBeTrue)
				convey.So(val, convey.ShouldEqual, "value")
			})

			var (
				d  time.Time
				ok bool
			)

			convey.So(func() { d, ok = ctx.Deadline() }, convey.ShouldNotPanic)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(math.Abs(time.Since(d).Seconds()), convey.ShouldBeBetween, 1, 2)
		})

		convey.Convey("with deadline", func() {
			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*10))
			defer cancel()
			ctx = NewContextWithOld(ctx)
			convey.So(ctx, convey.ShouldNotBeNil)
			var (
				d  time.Time
				ok bool
			)
			convey.So(func() { d, ok = ctx.Deadline() }, convey.ShouldNotPanic)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(math.Abs(time.Since(d).Seconds()), convey.ShouldBeBetween, 9, 10)
		})

		convey.Convey("unique ID to request ID", func() {
			convey.Convey("with no unique_id", func() {
				ctx := NewContextWithOld(context.Background())
				convey.So(ctx, convey.ShouldNotBeNil)
				reqID := GetRequestIDFromContext(ctx)
				convey.So(reqID, convey.ShouldBeEmpty)
			})

			convey.Convey("with unique_id", func() {
				md := metadata.Metadata{"unique_id": "id123"}
				old := metadata.NewContext(context.Background(), md)
				ctx := NewContextWithOld(old)
				convey.So(ctx, convey.ShouldNotBeNil)
				reqID := GetRequestIDFromContext(ctx)
				convey.So(reqID, convey.ShouldEqual, "id123")
			})
		})
	})
}
