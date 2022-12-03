package mlogger

import "go.uber.org/zap"

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

// Level const helper
const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// Options for init Log
type Options struct {
	Level        string
	Output       string
	ContextField map[string]any
	SkipCaller   int
}
