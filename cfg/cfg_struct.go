package cfg

import "time"

type App struct {
	Name         string
	Port         int
	DebugPort    int
	Env          string
	Secret       string
	LoggerOutput string
}

type DbConfig struct {
	DSN         string
	MaxOpenCons int
	MinOpenCons int
}

type RedisConfig struct {
	RedisURL         string
	RedisPass        string
	RedisDB          int
	RedisDefDuration time.Duration
}

type GoogleConfig struct {
	CredentialLocation string
}

type Telemetry struct {
	URL      string
	Key      string
	Insecure bool
}

type Toggle struct {
	TraceON  bool
	MetricON bool
	CacheON  bool
}
