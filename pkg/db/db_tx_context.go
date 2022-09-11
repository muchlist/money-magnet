package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type KeyTransaction string

const TXKey KeyTransaction = "moneymag-transaction"

// InjectTx injects transaction to context
func InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TXKey, tx)
}

// ExtractTx extracts transaction from context
func ExtractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(TXKey).(pgx.Tx); ok {
		return tx
	}
	return nil
}
