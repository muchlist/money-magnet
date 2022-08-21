package mlogger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mlog struct {
	zap *zap.Logger
}

// Field Type alias dari zap.Field,
// sehingga core app tidak memerlukan pemanggilan zap core
type Field = zap.Field

var (
	Binary     = zap.Binary
	Bool       = zap.Bool
	ByteString = zap.ByteString
	Float64    = zap.Float64
	Float32    = zap.Float32
	Int        = zap.Int
	Int64      = zap.Int64
	Int32      = zap.Int32
	String     = zap.String
	Stringp    = zap.Stringp
	Stack      = zap.Stack
	StackSkip  = zap.StackSkip
	Durationp  = zap.Durationp
	Any        = zap.Any
)

func New(level string, output string) *mlog {
	logConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLevel(level)),
		OutputPaths: []string{getOutput(output)},
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

	log := mlog{}
	var err error
	if log.zap, err = logConfig.Build(); err != nil {
		panic(err)
	}

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

func (l *mlog) InfoT(traceID string, msg string, tags ...Field) {
	tags = append(tags, zap.String("trace_id", traceID))
	l.zap.Info(msg, tags...)
}

func (l *mlog) Warn(msg string, tags ...Field) {
	l.zap.Warn(msg, tags...)
}

func (l *mlog) WarnT(traceID string, msg string, tags ...Field) {
	tags = append(tags, zap.String("trace_id", traceID))
	l.zap.Info(msg, tags...)
}

func (l *mlog) Error(msg string, err error, tags ...Field) {
	tags = append(tags, zap.NamedError("error", err), zap.StackSkip("stacktrace", 1))
	l.zap.Error(msg, tags...)
}

func (l *mlog) ErrorT(traceID string, msg string, err error, tags ...Field) {
	tags = append(tags, zap.String("trace_id", traceID), zap.NamedError("error", err), zap.StackSkip("stacktrace", 1))
	l.zap.Error(msg, tags...)
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
	case "error":
		return zap.ErrorLevel
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
