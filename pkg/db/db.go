package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	DSN          string
	MaxOpenConns int32
	MinOpenConns int32
}

// OpenDB init open database pool
func OpenDB(cfg Config) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	config.MaxConns = cfg.MaxOpenConns
	config.MinConns = cfg.MinOpenConns

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
