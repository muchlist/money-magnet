package mmetric

import (
	"context"
	"sync"

	"github.com/muchlist/moneymagnet/pkg/monitor"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
)

var doOnce sync.Once

var alloc, _ = meter.AsyncInt64().Gauge(
	"mem.alloc",
	instrument.WithUnit(unit.Dimensionless),
	instrument.WithDescription("allocated heap objects in MB"),
)

var totalAlloc, _ = meter.AsyncInt64().Gauge(
	"mem.total_alloc",
	instrument.WithUnit(unit.Dimensionless),
	instrument.WithDescription("cumulative MB allocated for heap objects"),
)

var sys, _ = meter.AsyncInt64().Gauge(
	"mem.sys",
	instrument.WithUnit(unit.Dimensionless),
	instrument.WithDescription("total MB of memory obtained from the OS"),
)

var gc, _ = meter.AsyncInt64().Gauge(
	"num.gc",
	instrument.WithUnit(unit.Dimensionless),
	instrument.WithDescription("number of garbage collectore finished"),
)

var goroutine, _ = meter.AsyncInt64().Gauge(
	"goroutine",
	instrument.WithUnit(unit.Dimensionless),
	instrument.WithDescription("number of goroutine active"),
)

func RegisterMonitorMetric(ctx context.Context) {
	doOnce.Do(func() {
		if err := meter.RegisterCallback(
			[]instrument.Asynchronous{
				alloc,
				totalAlloc,
				sys,
				gc,
				goroutine,
			},
			func(ctx context.Context) {
				data := monitor.GetMemUsage()

				alloc.Observe(ctx, int64(data.AllocMB), uniquePerNodeID)
				totalAlloc.Observe(ctx, int64(data.TotalAllocMB), uniquePerNodeID)
				sys.Observe(ctx, int64(data.SysMB), uniquePerNodeID)
				gc.Observe(ctx, int64(data.NumGC), uniquePerNodeID)
				goroutine.Observe(ctx, int64(data.NumGoroutine), uniquePerNodeID)
			},
		); err != nil {
			panic(err)
		}
	})
}
