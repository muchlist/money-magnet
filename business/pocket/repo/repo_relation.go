package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/observ"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	keyTableUP  = "user_pocket"
	keyIDUP     = "id"
	keyUserUP   = "user_id"
	keyPocketUP = "pocket_id"
)

// InsertPocketUser ...
func (r Repo) InsertPocketUser(pctx context.Context, userIDs []uuid.UUID, pocketID uuid.UUID) error {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-CreatePocket")
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

	_, err = r.mod(ctx).Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db.ParseError(err)
	}

	return nil
}

// DeletePocketUser ...
func (r Repo) DeletePocketUser(pctx context.Context, userID uuid.UUID, pocketID uuid.UUID) error {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-DeletePocketUser")
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

	res, err := r.mod(ctx).Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return db.ErrDBNotFound
	}

	return nil
}
