package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/business/request/model"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
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
	keyTable       = "requests"
	keyID          = "id"
	keyRequesterID = "requester_id"
	keyApproverID  = "approver_id"
	keyPocketID    = "pocket_id"
	keyPocketName  = "pocket_name"
	keyIsApproved  = "is_approved"
	keyIsRejected  = "is_rejected"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
)

// Repo manages the set of APIs for pocket access.
type Repo struct {
	db  *pgxpool.Pool
	log mlogger.Logger
	sb  sq.StatementBuilderType
}

// NewRepo constructs a data for api access..
func NewRepo(sqlDB *pgxpool.Pool, log mlogger.Logger) Repo {
	return Repo{
		db:  sqlDB,
		log: log,
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// =========================================================================
// MANIPULATOR

// Insert ...
func (r Repo) Insert(ctx context.Context, request *model.RequestPocket) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyRequesterID,
			keyApproverID,
			keyPocketID,
			keyPocketName,
			keyCreatedAt,
			keyUpdatedAt,
		).
		Values(
			request.RequesterID,
			request.ApproverID,
			request.PocketID,
			request.PocketName,
			request.CreatedAt,
			request.UpdatedAt).
		Suffix(db.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query insert request: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&request.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// UpdateStatus update approver, is_approved and udpdated_at
func (r Repo) UpdateStatus(ctx context.Context, request *model.RequestPocket) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyIsApproved: request.IsApproved,
			keyIsRejected: request.IsRejected,
			keyUpdatedAt:  time.Now(),
		}).
		Suffix(db.Returning(keyID,
			keyRequesterID,
			keyApproverID,
			keyPocketID,
			keyPocketName,
			keyIsApproved,
			keyIsRejected,
			keyCreatedAt,
			keyUpdatedAt,
		)).ToSql()

	if err != nil {
		return fmt.Errorf("build query update request: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(
		&request.ID,
		&request.RequesterID,
		&request.ApproverID,
		&request.PocketID,
		&request.PocketName,
		&request.IsApproved,
		&request.IsRejected,
		&request.CreatedAt,
		&request.UpdatedAt,
	)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// =========================================================================
// GETTER
// GetByID get one pocket by email
func (r Repo) GetByID(ctx context.Context, id uint64) (model.RequestPocket, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyRequesterID,
		keyApproverID,
		keyPocketID,
		keyPocketName,
		keyIsApproved,
		keyIsRejected,
		keyCreatedAt,
		keyUpdatedAt,
	).From(keyTable).Where(sq.Eq{keyID: id}).ToSql()

	if err != nil {
		return model.RequestPocket{}, fmt.Errorf("build query get request by id: %w", err)
	}

	var request model.RequestPocket
	err = r.db.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&request.ID,
			&request.RequesterID,
			&request.ApproverID,
			&request.PocketID,
			&request.PocketName,
			&request.IsApproved,
			&request.IsRejected,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return model.RequestPocket{}, db.ParseError(err)
	}

	return request, nil
}

// Find get all request by FIND model
func (r Repo) Find(ctx context.Context, findBy model.FindBy, filter data.Filters) ([]model.RequestPocket, data.Metadata, error) {

	// Validation filter
	filter.SortSafelist = []string{"updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, data.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// where builder
	var orBuilder sq.Or
	var inputCount int
	if findBy.ApproverID != "" {
		orBuilder = append(orBuilder, sq.Eq{keyApproverID: findBy.ApproverID})
		inputCount++
	}
	if len(findBy.PocketIDs) != 0 {
		orBuilder = append(orBuilder, sq.Eq{keyPocketID: findBy.PocketIDs})
		inputCount++
	}
	if findBy.RequesterID != "" {
		orBuilder = append(orBuilder, sq.Eq{keyRequesterID: findBy.RequesterID})
		inputCount++
	}

	if inputCount == 0 {
		return []model.RequestPocket{}, data.Metadata{}, nil
	}

	sqlStatement, args, err := r.sb.Select(
		"count(*) OVER()",
		keyID,
		keyRequesterID,
		keyApproverID,
		keyPocketID,
		keyPocketName,
		keyIsApproved,
		keyIsRejected,
		keyCreatedAt,
		keyUpdatedAt,
	).
		From(keyTable).
		Where(orBuilder).
		OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("build query find request: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, data.Metadata{}, db.ParseError(err)
	}
	defer rows.Close()

	totalRecords := 0
	requests := make([]model.RequestPocket, 0)
	for rows.Next() {
		var request model.RequestPocket
		err := rows.Scan(
			&totalRecords,
			&request.ID,
			&request.RequesterID,
			&request.ApproverID,
			&request.PocketID,
			&request.PocketName,
			&request.IsApproved,
			&request.IsRejected,
			&request.CreatedAt,
			&request.UpdatedAt,
		)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, data.Metadata{}, db.ParseError(err)
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return requests, metadata, nil
}
