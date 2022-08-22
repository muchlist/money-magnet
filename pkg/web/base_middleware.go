package web

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"go.uber.org/zap"
)

// midLogger ...
func midLogger(l mlogger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			clientID := getFirstValueFromHeader(r, "X-Client-Id")
			ipAddress := getFirstValueFromHeader(r, "X-Forwarded-For")

			// log incomming request
			l.Info("request started",
				zap.String("path", r.URL.Path),
				zap.String("trace_id", ReadTraceID(r.Context())),
				zap.String("client_id", clientID),
				zap.String("source_ip", ipAddress),
			)

			// log request end
			defer func() {
				l.Info("request completed",
					zap.String("path", r.URL.Path),
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

type ctxKeyRequestID int

const RequestIDKey ctxKeyRequestID = 0

var RequestIDHeader = "X-Request-Id"

// requestID read header with key X-Request-Id, if exist that value used to traceID
// if not, generate uuid for traceID
func requestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
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

func getFirstValueFromHeader(req *http.Request, key string) string {
	vs, ok := req.Header[key]
	if ok {
		if len(vs) != 0 {
			return vs[0]
		}
	}
	return ""
}