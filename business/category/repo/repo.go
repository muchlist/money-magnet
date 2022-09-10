package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
)

const (
	keyTable        = "categories"
	keyID           = "id"
	keyCategoryName = "category_name"
	keyPocketID     = "pocket_id"
	keyIsIncome     = "is_income"
	keyCreatedAt    = "created_at"
	keyUpdatedAt    = "updated_at"
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
func (r Repo) Insert(ctx context.Context, category *model.Category) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyID,
			keyCategoryName,
			keyPocketID,
			keyIsIncome,
			keyUpdatedAt,
			keyCreatedAt,
		).
		Values(
			category.ID,
			category.CategoryName,
			category.PocketID,
			category.IsIncome,
			category.CreatedAt,
			category.UpdatedAt).
		Suffix(db.Returning(keyID)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query insert category: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&category.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(ctx context.Context, category *model.Category) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyCategoryName: category.CategoryName,
			keyUpdatedAt:    time.Now(),
		}).
		Where(sq.Eq{keyID: category.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit category: %w", err)
	}

	res, err := r.db.Exec(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}
	if res.RowsAffected() == 0 {
		return db.ErrDBNotFound
	}

	return nil
}

// Delete ...
func (r Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.Eq{keyID: id}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete category: %w", err)
	}

	res, err := r.db.Exec(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return db.ErrDBNotFound
	}

	return nil
}

// =========================================================================
// GETTER

// GetByID get one category by email
func (r Repo) GetByID(ctx context.Context, id uuid.UUID) (model.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyCategoryName,
		keyIsIncome,
		keyPocketID,
		keyCreatedAt,
		keyUpdatedAt,
	).From(keyTable).Where(sq.Eq{keyID: id}).ToSql()

	if err != nil {
		return model.Category{}, fmt.Errorf("build query get category by id: %w", err)
	}

	var cat model.Category
	err = r.db.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&cat.ID,
			&cat.CategoryName,
			&cat.IsIncome,
			&cat.PocketID,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return model.Category{}, db.ParseError(err)
	}

	return cat, nil
}

// Find get all category within user
func (r Repo) Find(ctx context.Context, pocketID uuid.UUID, filter data.Filters) ([]model.Category, data.Metadata, error) {

	// Validation filter
	filter.SortSafelist = []string{"category_name", "-category_name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, data.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		"count(*) OVER()",
		keyID,
		keyCategoryName,
		keyIsIncome,
		keyPocketID,
		keyCreatedAt,
		keyUpdatedAt,
	).
		From(keyTable).
		Where(sq.Eq{keyPocketID: pocketID}).
		OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("build query find category: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, data.Metadata{}, db.ParseError(err)
	}
	defer rows.Close()

	totalRecords := 0
	cats := make([]model.Category, 0)
	for rows.Next() {
		var cat model.Category
		err := rows.Scan(
			&totalRecords,
			&cat.ID,
			&cat.CategoryName,
			&cat.IsIncome,
			&cat.PocketID,
			&cat.CreatedAt,
			&cat.UpdatedAt)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, data.Metadata{}, db.ParseError(err)
		}
		cats = append(cats, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return cats, metadata, nil
}
