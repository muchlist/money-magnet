package main

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

func debugMux(db *pgxpool.Pool) *http.ServeMux {
	mux := http.NewServeMux()

	expvar.NewString("api_version").Set(version)
	expvar.Publish("api_timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("database", expvar.Func(func() interface{} {
		stat := db.Stat()
		return map[string]interface{}{
			"conn_max":       stat.MaxConns(),
			"conn_idle":      stat.IdleConns(),
			"conn_in_use":    stat.TotalConns(),
			"acquire_total":  stat.AcquireCount(),
			"acquire_cancel": stat.CanceledAcquireCount(),
		}
	}))

	// Register all the standard library debug endpoints.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
