package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/xulid"

	sq "github.com/Masterminds/squirrel"
)

const (
	keyTableUP  = "user_pocket"
	keyIDUP     = "id"
	keyUserUP   = "user_id"
	keyPocketUP = "pocket_id"
)

// InsertPocketUser ...
func (r Repo) InsertPocketUser(ctx context.Context, userIDs []string, pocketID xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-CreatePocket")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := r.sb.Insert(keyTableUP).
		Columns(
			keyUserUP,
			keyPocketUP,
		)

	for _, userID := range userIDs {
		query = query.Values(
			userID,
			pocketID,
		)
	}
	sqlStatement, args, err := query.Suffix("ON CONFLICT DO NOTHING").ToSql()

	if err != nil {
		return fmt.Errorf("build query insert pocket user relation: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	_, err = dbtx.Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db.ParseError(err)
	}

	return nil
}

// DeletePocketUser ...
func (r Repo) DeletePocketUser(ctx context.Context, userID xulid.ULID, pocketID xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-DeletePocketUser")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTableUP).
		Where(sq.And{
			sq.Eq{keyUserUP: userID},
			sq.Eq{keyPocketUP: pocketID},
		}).ToSql()

	if err != nil {
		return fmt.Errorf("build query delete pocket user relation: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	res, err := dbtx.Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return db.ErrDBNotFound
	}

	return nil
}
