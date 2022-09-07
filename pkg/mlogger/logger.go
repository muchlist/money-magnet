package mlogger

import "context"

type Logger interface {
	Debug(msg string, tags ...Field)
	Info(msg string, tags ...Field)
	InfoT(ctx context.Context, msg string, tags ...Field)
	Warn(msg string, err error, tags ...Field)
	WarnT(ctx context.Context, msg string, err error, tags ...Field)
	Error(msg string, err error, tags ...Field)
	ErrorT(ctx context.Context, msg string, err error, tags ...Field)
	Printf(format string, v ...interface{})
	Print(format string, v ...interface{})
}
