package repo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/business/category/port"
	"github.com/muchlist/moneymagnet/constant"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
)

const (
	keyTable            = "categories"
	keyID               = "id"
	keyCategoryName     = "category_name"
	keyCategoryIcon     = "category_icon"
	keyPocketID         = "pocket_id"
	keyIsIncome         = "is_income"
	keyDefaultSpendType = "default_spend_type"
	keyCreatedAt        = "created_at"
	keyUpdatedAt        = "updated_at"
)

// make sure the implementation satisfies the interface
var _ port.CategoryStorer = (*Repo)(nil)

// Repo manages the set of APIs for pocket access.
type Repo struct {
	db  *pgxpool.Pool
	log mlogger.Logger
	sb  sq.StatementBuilderType
}

// NewRepo constructs a data for api access..
func NewRepo(sqlDB *pgxpool.Pool, log mlogger.Logger) *Repo {
	return &Repo{
		db:  sqlDB,
		log: log,
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// =========================================================================
// MANIPULATOR

// Insert ...
func (r *Repo) Insert(ctx context.Context, category *model.Category) error {
	ctx, span := observ.GetTracer().Start(ctx, "category-repo-Insert")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyID,
			keyCategoryName,
			keyCategoryIcon,
			keyPocketID,
			keyIsIncome,
			keyDefaultSpendType,
			keyUpdatedAt,
			keyCreatedAt,
		).
		Values(
			category.ID,
			category.CategoryName,
			category.CategoryIcon,
			category.PocketID,
			category.IsIncome,
			category.DefaultSpendType,
			category.CreatedAt,
			category.UpdatedAt).
		Suffix(db.Returning(keyID)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query insert category: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&category.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Insert Many...
func (r *Repo) InsertMany(ctx context.Context, categories []model.Category) error {
	ctx, span := observ.GetTracer().Start(ctx, "category-repo-Insert")
	defer span.End()

	if len(categories) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, 6*time.Second)
	defer cancel()

	sqBuilder := r.sb.Insert(keyTable).
		Columns(
			keyID,
			keyPocketID,
			keyCategoryName,
			keyCategoryIcon,
			keyIsIncome,
			keyDefaultSpendType,
			keyUpdatedAt,
			keyCreatedAt,
		)

	for _, category := range categories {
		sqBuilder = sqBuilder.Values(
			category.ID,
			category.PocketID,
			category.CategoryName,
			category.CategoryIcon,
			category.IsIncome,
			category.DefaultSpendType,
			category.CreatedAt,
			category.UpdatedAt,
		)
	}

	sqlStatement, args, err := sqBuilder.ToSql()

	if err != nil {
		return fmt.Errorf("build query insert many category: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	_, err = dbtx.Exec(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r *Repo) Edit(ctx context.Context, category *model.Category) error {
	ctx, span := observ.GetTracer().Start(ctx, "category-repo-Edit")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyCategoryName:     category.CategoryName,
			keyCategoryIcon:     category.CategoryIcon,
			keyDefaultSpendType: category.DefaultSpendType,
			keyUpdatedAt:        time.Now(),
		}).
		Where(sq.Eq{keyID: category.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit category: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	res, err := dbtx.Exec(ctx, sqlStatement, args...)
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
func (r *Repo) Delete(ctx context.Context, id string) error {
	ctx, span := observ.GetTracer().Start(ctx, "category-repo-Delete")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.Eq{keyID: id}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete category: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	res, err := dbtx.Exec(ctx, sqlStatement, args...)
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

// GetByID get one category by id
func (r *Repo) GetByID(ctx context.Context, id string) (model.Category, error) {
	ctx, span := observ.GetTracer().Start(ctx, "category-repo-GetByID")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyCategoryName,
		keyCategoryIcon,
		keyIsIncome,
		keyDefaultSpendType,
		keyPocketID,
		keyCreatedAt,
		keyUpdatedAt,
	).From(keyTable).Where(sq.Eq{keyID: id}).ToSql()

	if err != nil {
		return model.Category{}, fmt.Errorf("build query get category by id: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	var cat model.Category
	err = dbtx.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&cat.ID,
			&cat.CategoryName,
			&cat.CategoryIcon,
			&cat.IsIncome,
			&cat.DefaultSpendType,
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

// Find get all category within pocketID
func (r *Repo) Find(ctx context.Context, pocketID string, filter paging.Filters) ([]model.Category, paging.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "category-repo-Find")
	defer span.End()

	// Validation filter
	filter.SortSafelist = []string{"category_name", "-category_name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, paging.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	where := sq.Or{
		sq.Eq{keyPocketID: pocketID},
		sq.Eq{keyPocketID: constant.POCK_MAIN_ID}, // OR 00000000000000000000000000
	}

	sqlStatement, args, err := r.sb.Select(
		"count(*) OVER()",
		keyID,
		keyCategoryName,
		keyCategoryIcon,
		keyIsIncome,
		keyDefaultSpendType,
		keyPocketID,
		keyCreatedAt,
		keyUpdatedAt,
	).
		From(keyTable).
		Where(where).
		OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, paging.Metadata{}, fmt.Errorf("build query find category: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, paging.Metadata{}, db.ParseError(err)
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
			&cat.CategoryIcon,
			&cat.IsIncome,
			&cat.DefaultSpendType,
			&cat.PocketID,
			&cat.CreatedAt,
			&cat.UpdatedAt)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, paging.Metadata{}, db.ParseError(err)
		}
		cats = append(cats, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, paging.Metadata{}, err
	}

	metadata := paging.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return cats, metadata, nil
}
