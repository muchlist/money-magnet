package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
)

type TxManager interface {
	WithAtomic(ctx context.Context, tFunc func(ctx context.Context) error) error
}

type txManager struct {
	db  *pgxpool.Pool
	log mlogger.Logger
}

func NewTxManager(sqlDB *pgxpool.Pool, log mlogger.Logger) TxManager {
	return &txManager{
		db:  sqlDB,
		log: log,
	}
}

// =========================================================================
// TRANSACTION

// WithAtomic runs function within transaction
// The transaction commits when function were finished without error
func (r *txManager) WithAtomic(ctx context.Context, tFunc func(ctx context.Context) error) error {

	// begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// run callback
	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			r.log.Error("rollback transaction", errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(ctx); errCommit != nil {
		return fmt.Errorf("commit transaction: %w", errCommit)
	}
	return nil
}
