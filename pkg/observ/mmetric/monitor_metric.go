package mmetric

import (
	"context"
	"sync"

	"github.com/muchlist/moneymagnet/pkg/monitor"
	"go.opentelemetry.io/otel/metric"
)

var (
	doOnce     sync.Once
	alloc      metric.Int64ObservableGauge
	totalAlloc metric.Int64ObservableGauge
	sys        metric.Int64ObservableGauge
	gc         metric.Int64ObservableGauge
	goroutine  metric.Int64ObservableGauge
)

func init() {
	var err error
	alloc, err = meter.Int64ObservableGauge(
		"mem.alloc",
		metric.WithDescription("allocated heap objects in MB"),
		metric.WithUnit("MB"),
	)
	if err != nil {
		panic(err)
	}

	totalAlloc, err = meter.Int64ObservableGauge(
		"mem.total_alloc",
		metric.WithDescription("cumulative MB allocated for heap objects"),
		metric.WithUnit("MB"),
	)
	if err != nil {
		panic(err)
	}

	sys, err = meter.Int64ObservableGauge(
		"mem.sys",
		metric.WithDescription("total MB of memory obtained from the OS"),
		metric.WithUnit("MB"),
	)
	if err != nil {
		panic(err)
	}

	gc, err = meter.Int64ObservableGauge(
		"num.gc",
		metric.WithDescription("number of garbage collectors finished"),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(err)
	}

	goroutine, err = meter.Int64ObservableGauge(
		"goroutine",
		metric.WithDescription("number of goroutine active"),
		metric.WithUnit("1"),
	)
	if err != nil {
		panic(err)
	}
}

func RegisterMonitorMetric(ctx context.Context) {
	doOnce.Do(func() {
		_, err := meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
			data := monitor.GetMemUsage()

			o.ObserveInt64(alloc, int64(data.AllocMB), metric.WithAttributes(uniquePerNodeID))
			o.ObserveInt64(totalAlloc, int64(data.TotalAllocMB), metric.WithAttributes(uniquePerNodeID))
			o.ObserveInt64(sys, int64(data.SysMB), metric.WithAttributes(uniquePerNodeID))
			o.ObserveInt64(gc, int64(data.NumGC), metric.WithAttributes(uniquePerNodeID))
			o.ObserveInt64(goroutine, int64(data.NumGoroutine), metric.WithAttributes(uniquePerNodeID))
			return nil
		}, alloc, totalAlloc, sys, gc, goroutine)

		if err != nil {
			panic(err)
		}
	})
}
