package cfg

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/muchlist/moneymagnet/pkg/env"
)

type Config struct {
	App       App
	DB        DbConfig
	Telemetry Telemetry
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		// dont panic, because in prod we will not use this.
		// use os env instead
		log.Print(".env notfound [env still can be read from os env variable]")
	}

	return &Config{
		App: App{
			Name:      env.Get("APP_NAME", "money-magnet"),
			Port:      env.Get("APP_PORT", 8081),
			DebugPort: env.Get("APP_DEBUG_PORT", 4000),
			Env:       env.Get("APP_ENV", "dev"),
			Secret:    env.Get("APP_SECRET", "xoxoxoxo"),
		},
		DB: DbConfig{
			DSN:         env.Get("DB_DSN", "postgres://postgres:postgres@localhost:5432/money_magnet?sslmode=disable"),
			MaxOpenCons: env.Get("DB_MAX_CONN", 100),
			MinOpenCons: env.Get("DB_MIN_CONN", 2),
		},
		Telemetry: Telemetry{
			URL:      env.Get("OTEL_URL", "localhost:4317"),
			Key:      env.Get("OTEL_KEY", "example-api-key"),
			Insecure: env.Get("OTEL_INSECURE", true),
		},
	}

}
