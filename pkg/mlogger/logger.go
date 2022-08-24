package mlogger

type Logger interface {
	Debug(msg string, tags ...Field)
	Info(msg string, tags ...Field)
	InfoT(traceID string, msg string, tags ...Field)
	Warn(msg string, err error, tags ...Field)
	WarnT(traceID string, msg string, err error, tags ...Field)
	Error(msg string, err error, tags ...Field)
	ErrorT(traceID string, msg string, err error, tags ...Field)
	Printf(format string, v ...interface{})
	Print(format string, v ...interface{})
}
