package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type KeyTransaction string

const TXKey KeyTransaction = "moneymag-transaction"

// ExtractTx extract transaction from context and transform database into db.DBTX
func ExtractTx(ctx context.Context, defaultPool *pgxpool.Pool) DBTX {
	tx, ok := ctx.Value(TXKey).(pgx.Tx)
	if !ok || tx == nil {
		return NewPGStore(defaultPool, nil)
	}
	return NewPGStore(nil, tx)
}

// injectTx injects transaction to context
func injectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TXKey, tx)
}
