package ptrepo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
)

const (
	keyTableUP  = "user_pocket"
	keyIDUP     = "id"
	keyUserUP   = "user_id"
	keyPocketUP = "pocket_id"
)

// InsertPocketUser ...
func (r Repo) InsertPocketUser(ctx context.Context, userID uuid.UUID, pocketID uint64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTableUP).
		Columns(
			keyUserUP,
			keyPocketUP,
		).
		Values(
			userID,
			pocketID,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query insert pocket user relation: %w", err)
	}

	fmt.Println(sqlStatement)
	fmt.Println(args)

	_, err = r.db.Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db.ParseError(err)
	}

	return nil
}

// DeletePocketUser ...
func (r Repo) DeletePocketUser(ctx context.Context, userID uuid.UUID, pocketID uint64) error {
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

	res, err := r.db.Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return db.ErrDBNotFound
	}

	return nil
}
