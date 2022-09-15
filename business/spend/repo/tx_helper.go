package repo

import (
	"context"
	"fmt"

	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/observ"
)

// =========================================================================
// TRANSACTION

// WithinTransaction runs function within transaction
// The transaction commits when function were finished without error
func (r Repo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-WithinTransaction")
	defer span.End()

	// begin transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	// run callback
	err = tFunc(db.InjectTx(ctx, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			r.log.Error("rollback transaction", errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(ctx); errCommit != nil {
		r.log.Error("commit transaction", errCommit)
	}
	return nil
}

// mod returns query model with context with or without transaction extracted from context
func (r Repo) mod(ctx context.Context) db.DBTX {
	tx := db.ExtractTx(ctx)
	if tx != nil {
		return db.NewPGStore(nil, tx)
	}
	return db.NewPGStore(r.db, nil)
}
