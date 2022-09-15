package global

type KeyRequestIDType string

// RequestIDKey used for key context to set-get requestID value
// use global variable bacause ctx value must be passing to different lib
const RequestIDKey KeyRequestIDType = "Request-Id"

// TraceIDKey used for key context to set-get traceID from otel tracer
// use global variable bacause ctx value must be passing to different lib
const TraceIDKey KeyRequestIDType = "Trace-Id"
