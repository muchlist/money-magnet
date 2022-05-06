package ptrepo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptmodel"
	"github.com/muchlist/moneymagnet/bussines/sys/data"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
)

const (
	keyTable      = "pockets"
	keyID         = "id"
	keyOwner      = "owner"
	keyWathcer    = "watcher"
	keyPocketName = "pocket_name"
	keyLevel      = "level"
	keyCreatedAt  = "created_at"
	keyUpdatedAt  = "updated_at"
	keyVersion    = "version"
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
func (r Repo) Insert(ctx context.Context, pocket *ptmodel.Pocket) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyPocketName,
			keyOwner,
			keyWathcer,
			keyVersion,
			keyLevel,
			keyUpdatedAt,
			keyCreatedAt,
		).
		Values(
			pocket.PocketName,
			pocket.Owner,
			pocket.Watcher,
			pocket.Version,
			pocket.Level,
			pocket.CreatedAt,
			pocket.UpdatedAt).
		Suffix(db.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query insert pocket: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&pocket.ID)
	if err != nil {
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(ctx context.Context, pocket *ptmodel.Pocket) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyPocketName: pocket.PocketName,
			keyOwner:      pocket.Owner,
			keyWathcer:    pocket.Watcher,
			keyLevel:      pocket.Level,
			keyUpdatedAt:  time.Now(),
			keyVersion:    pocket.Version + 1,
		}).
		Where(sq.Eq{keyID: pocket.ID}).
		Suffix(db.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit pocket: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&pocket.Version)
	if err != nil {
		return db.ParseError(err)
	}

	return nil
}

// Delete ...
func (r Repo) Delete(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.Eq{keyID: id}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete pocket: %w", err)
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

// =========================================================================
// GETTER

// GetByID get one pocket by email
func (r Repo) GetByID(ctx context.Context, id uint64) (ptmodel.Pocket, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyOwner,
		keyWathcer,
		keyPocketName,
		keyLevel,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyID: id}).ToSql()

	if err != nil {
		return ptmodel.Pocket{}, fmt.Errorf("build query get pocket by id: %w", err)
	}

	var pocket ptmodel.Pocket
	err = r.db.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&pocket.ID,
			&pocket.Owner,
			&pocket.Watcher,
			&pocket.PocketName,
			&pocket.Level,
			&pocket.CreatedAt,
			&pocket.UpdatedAt,
			&pocket.Version)
	if err != nil {
		return ptmodel.Pocket{}, db.ParseError(err)
	}

	return pocket, nil
}

// Find get all pocket
func (r Repo) Find(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error) {

	// Validation filter
	filter.SortSafelist = []string{"name", "-name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, data.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		"count(*) OVER()",
		keyID,
		keyOwner,
		keyWathcer,
		keyPocketName,
		keyLevel,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).
		From(keyTable).
		Where(sq.Eq{keyOwner: owner}).
		OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("build query find pocket: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, data.Metadata{}, db.ParseError(err)
	}
	defer rows.Close()

	totalRecords := 0
	pockets := make([]ptmodel.Pocket, 0)
	for rows.Next() {
		var pocket ptmodel.Pocket
		err := rows.Scan(
			&totalRecords,
			&pocket.ID,
			&pocket.Owner,
			&pocket.Watcher,
			&pocket.PocketName,
			&pocket.Level,
			&pocket.CreatedAt,
			&pocket.UpdatedAt,
			&pocket.Version)
		if err != nil {
			return nil, data.Metadata{}, db.ParseError(err)
		}
		pockets = append(pockets, pocket)
	}

	if err := rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return pockets, metadata, nil
}
