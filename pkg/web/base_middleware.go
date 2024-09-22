package web

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/muchlist/moneymagnet/pkg/global"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// midLogger ...
func midLogger(l mlogger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			clientID := r.Header.Get("X-Client-Id")
			ipAddress := r.Header.Get("X-Forwarded-For")

			path := fmt.Sprintf("%s_%s", r.Method, r.URL.Path)

			// log incomming request
			l.Info("request started",
				zap.String("path", path),
				zap.String("request_id", ReadRequestID(r.Context())),
				zap.String("trace_id", ReadTraceID(r.Context())),
				zap.String("client_id", clientID),
				zap.String("source_ip", ipAddress),
			)

			// log request end
			defer func() {
				l.Info("request completed",
					zap.String("path", path),
					zap.String("request_id", ReadRequestID(r.Context())),
					zap.String("trace_id", ReadTraceID(r.Context())),
					zap.String("client_id", clientID),
					zap.String("source_ip", ipAddress),
					zap.String("latency", fmt.Sprint(time.Since(t1))),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

var RequestIDHeader = "X-Request-Id"

// requestID read header with key X-Request-Id, if exist that value used to traceID
// if not, generate uuid for traceID
func requestAndTraceID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get RequestID from header
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Get TraceID from otel tracer
		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID().String()

		// set requestID and traceID to context
		ctx = context.WithValue(ctx, global.RequestIDKey, requestID)
		ctx = context.WithValue(ctx, global.TraceIDKey, traceID)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// response writted when got panic
func panicRecovery(l mlogger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 2048)
					n := runtime.Stack(buf, false)
					buf = buf[:n]

					l.Info(fmt.Sprintf("recovering from err %v\n %s", err, buf))
					ServerErrorResponse(w, r, fmt.Errorf("%v", err))
				}
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
