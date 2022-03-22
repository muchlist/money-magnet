package mlogger

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, tags ...zap.Field)
	Error(msg string, err error, tags ...zap.Field)
	InfoT(traceID string, msg string, tags ...zap.Field)
	ErrorT(traceID string, msg string, err error, tags ...zap.Field)
	Printf(format string, v ...interface{})
	Print(format string, v ...interface{})
}
