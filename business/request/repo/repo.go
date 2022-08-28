package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/business/request/model"
	"github.com/muchlist/moneymagnet/pkg/db"
)

/*
CREATE TABLE IF NOT EXISTS "requests" (
  "id" BIGSERIAL PRIMARY KEY,
  "requester" uuid,
  "pocket" bigint,
  "pocket_name" varchar(100) NOT NULL,
  "is_approved" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);
*/

const (
	keyTable      = "requests"
	keyID         = "id"
	keyRequester  = "requester"
	keyApprover   = "approver"
	keyPocket     = "pocket"
	keyPocketName = "pocket_name"
	keyIsApproved = "is_approved"
	keyCreatedAt  = "created_at"
	keyUpdatedAt  = "updated_at"
)

// Repo manages the set of APIs for pocket access.
type Repo struct {
	db *pgxpool.Pool
	sb sq.StatementBuilderType
}

// NewRepo constructs a data for api access..
func NewRepo(sqlDB *pgxpool.Pool) Repo {
	return Repo{
		db: sqlDB,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// =========================================================================
// MANIPULATOR

// Insert ...
func (r Repo) Insert(ctx context.Context, pocket *model.RequestPocket) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyRequester,
			keyPocket,
			keyPocketName,
			keyCreatedAt,
			keyUpdatedAt,
		).
		Values(
			pocket.Requester,
			pocket.Pocket,
			pocket.PocketName,
			pocket.CreatedAt,
			pocket.UpdatedAt).
		Suffix(db.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query register request: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&pocket.ID)
	if err != nil {
		return db.ParseError(err)
	}

	return nil
}

// UpdateStatus update approver, is_approved and udpdated_at
func (r Repo) UpdateStatus(ctx context.Context, pocket *model.RequestPocket) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyApprover,
			keyIsApproved,
			keyUpdatedAt,
		).
		Values(
			pocket.Approver,
			pocket.IsApproved,
			time.Now(),
		).
		Suffix(db.Returning(keyID,
			keyRequester,
			keyApprover,
			keyPocket,
			keyPocketName,
			keyIsApproved,
			keyCreatedAt,
			keyUpdatedAt,
		)).ToSql()

	if err != nil {
		return fmt.Errorf("build query update request: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(
		&pocket.ID,
		&pocket.Requester,
		&pocket.Approver,
		&pocket.Pocket,
		&pocket.PocketName,
		&pocket.IsApproved,
		&pocket.CreatedAt,
		&pocket.UpdatedAt,
	)
	if err != nil {
		return db.ParseError(err)
	}

	return nil
}
