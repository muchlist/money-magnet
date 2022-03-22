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

func (l *mlog) Info(msg string, tags ...zap.Field) {
	l.zap.Info(msg, tags...)
	_ = l.zap.Sync()
}

func (l *mlog) InfoT(traceID string, msg string, tags ...zap.Field) {
	tags = append(tags, zap.String("trace_id", traceID))
	l.zap.Info(msg, tags...)
	_ = l.zap.Sync()
}

func (l *mlog) Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	l.zap.Error(msg, tags...)
	_ = l.zap.Sync()
}

func (l *mlog) ErrorT(traceID string, msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.String("trace_id", traceID), zap.NamedError("error", err))
	l.zap.Error(msg, tags...)
	_ = l.zap.Sync()
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
