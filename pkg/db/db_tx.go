package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	Begin(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type PGStore struct {
	NonTX *pgxpool.Pool
	Tx    pgx.Tx
}

// NewPGStore return interface can execute TX and pgx.Pool
func NewPGStore(pool *pgxpool.Pool, tx pgx.Tx) DBTX {
	var pgstore PGStore
	if tx != nil {
		pgstore.Tx = tx
		return &pgstore
	}
	pgstore.NonTX = pool
	return &pgstore
}

// Begin implements DBTX
func (p *PGStore) Begin(ctx context.Context) (pgx.Tx, error) {
	if p.Tx != nil {
		return nil, errors.New("cannot begin inside running transaction")
	}
	return p.NonTX.Begin(ctx)
}

// Commit implements DBTX
func (p *PGStore) Commit(ctx context.Context) error {
	if p.Tx != nil {
		return p.Tx.Commit(ctx)
	}
	return errors.New("cannot commit: nil tx value")
}

// Rollback implements DBTX
func (p *PGStore) Rollback(ctx context.Context) error {
	if p.Tx != nil {
		return p.Tx.Rollback(ctx)
	}
	return errors.New("cannot roleback: nil tx value")
}

// Exec implements DBTX
func (p *PGStore) Exec(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error) {
	if p.Tx != nil {
		return p.Tx.Exec(ctx, sql, arguments...)
	}
	return p.NonTX.Exec(ctx, sql, arguments...)
}

// Prepare implements DBTX
func (p *PGStore) Prepare(ctx context.Context, name string, sql string) (*pgconn.StatementDescription, error) {
	if p.Tx != nil {
		return p.Tx.Prepare(ctx, name, sql)
	}
	return nil, errors.New("cannot prefare: pool does not have prefare method")
}

// Query implements DBTX
func (p *PGStore) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if p.Tx != nil {
		return p.Tx.Query(ctx, sql, args...)
	}
	return p.NonTX.Query(ctx, sql, args...)
}

// QueryRow implements DBTX
func (p *PGStore) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if p.Tx != nil {
		return p.Tx.QueryRow(ctx, sql, args...)
	}
	return p.NonTX.QueryRow(ctx, sql, args...)
}
