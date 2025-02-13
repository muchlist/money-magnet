package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type commonValue interface {
	~int | ~bool | ~string | ~int64
}

func Get[T commonValue](key string, def T) T {
	var result any = def
	valueStr := strings.TrimSpace(os.Getenv(key))
	if valueStr == "" {
		return def
	}

	switch any(def).(type) {
	case string:
		result = valueStr

	case bool:
		val, err := strconv.ParseBool(strings.ToLower(valueStr))
		if err != nil {
			return def
		}
		result = val

	case int:
		val, err := strconv.Atoi(valueStr)
		if err != nil {
			return def
		}
		result = val

	case time.Duration:
		duration, err := time.ParseDuration(valueStr)
		if err != nil {
			return def
		}
		result = duration
	}

	return result.(T)
}
