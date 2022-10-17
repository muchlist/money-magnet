package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	// Make sure to import this so the instrumented driver is registered.
	_ "github.com/signalfx/splunk-otel-go/instrumentation/github.com/jackc/pgx/splunkpgx"
)

type Config struct {
	DSN          string
	MaxOpenConns int
	MinOpenConns int
}

// OpenDB init open database pool
func OpenDB(cfg Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(cfg.MaxOpenConns)
	config.MinConns = int32(cfg.MinOpenConns)

	db, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
