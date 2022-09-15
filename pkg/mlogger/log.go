package mlogger

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mlog struct {
	zap          *zap.Logger
	contextField map[string]any
}

func New(opt Options) *mlog {
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLevel(opt.Level)),
		OutputPaths: []string{getOutput(opt.Output)},
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:     "lvl",
			TimeKey:      "time",
			MessageKey:   "msg",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var log mlog
	var err error
	if log.zap, err = logConfig.Build(); err != nil {
		panic(err)
	}
	log.contextField = opt.ContextField

	return &log
}

func (l *mlog) Sync() error {
	return l.zap.Sync()
}

func (l *mlog) Debug(msg string, tags ...Field) {
	l.zap.Debug(msg, tags...)
}

func (l *mlog) Info(msg string, tags ...Field) {
	l.zap.Info(msg, tags...)
}

func (l *mlog) InfoT(ctx context.Context, msg string, tags ...Field) {
	fields := l.getFieldFromContext(ctx)
	fields = append(fields, tags...)
	l.zap.Info(msg, fields...)
}

func (l *mlog) Warn(msg string, err error, tags ...Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.zap.Warn(msg, tags...)
}

func (l *mlog) WarnT(ctx context.Context, msg string, err error, tags ...Field) {
	fields := l.getFieldFromContext(ctx)
	fields = append(fields, zap.NamedError("error", err))
	fields = append(fields, tags...)
	l.zap.Warn(msg, fields...)
}

func (l *mlog) Error(msg string, err error, tags ...Field) {
	tags = append(tags, zap.NamedError("error", err), zap.StackSkip("stacktrace", 1))
	l.zap.Error(msg, tags...)
}

func (l *mlog) ErrorT(ctx context.Context, msg string, err error, tags ...Field) {
	// send error to otel
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	fields := l.getFieldFromContext(ctx)
	fields = append(fields, zap.NamedError("error", err), zap.StackSkip("stacktrace", 1))
	fields = append(fields, tags...)
	l.zap.Error(msg, fields...)
}

// Printf used to mimic setLogger interface on other lib, ex : ElasticSearch
func (l *mlog) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		l.Info(format)
	} else {
		l.Info(fmt.Sprintf(format, v...))
	}
}

func (l *mlog) Print(format string, v ...interface{}) {
	if len(v) == 0 {
		l.Info(format)
	} else {
		l.Info(fmt.Sprintf(format, v...))
	}
}

func getLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "panic":
		return zap.PanicLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput(output string) string {
	out := strings.TrimSpace(output)
	if out == "" {
		return "stdout"
	}
	return out
}

func (l *mlog) getFieldFromContext(ctx context.Context) []zapcore.Field {
	if ctx == nil {
		return nil
	}
	fields := make([]zapcore.Field, 0, len(l.contextField))
	for key, v := range l.contextField {
		if ctxValue, ok := ctx.Value(v).(string); ok {
			fields = append(fields, zap.String(key, ctxValue))
		}
	}
	return fields
}
