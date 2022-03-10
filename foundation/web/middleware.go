package web

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
	"go.uber.org/zap"
)

func midLogger(l mlogger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				l.Info("served",
					zap.String("path", r.URL.Path),
					zap.String("trace_id", middleware.GetReqID(r.Context())),
					zap.Duration("latency", time.Since(t1)),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
