package monitor

import (
	"runtime"
)

type SimpleMonUsage struct {
	AllocMB      uint64
	TotalAllocMB uint64
	SysMB        uint64
	NumGC        uint32
	NumGoroutine int
}

func GetMemUsage() SimpleMonUsage {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	return SimpleMonUsage{
		AllocMB:      bToMb(m.Alloc),
		TotalAllocMB: bToMb(m.TotalAlloc),
		SysMB:        bToMb(m.Sys),
		NumGC:        m.NumGC,
		NumGoroutine: runtime.NumGoroutine(),
	}
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
