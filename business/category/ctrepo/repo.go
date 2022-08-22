package ctrepo

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/business/category/ctmodel"
	"github.com/muchlist/moneymagnet/pkg/db"
)

const (
	keyTable        = "categories"
	keyID           = "id"
	keyCategoryName = "category_name"
	keyPocket       = "pocket"
	keyIsIncome     = "is_income"
	keyCreatedAt    = "created_at"
	keyUpdatedAt    = "updated_at"
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
func (r Repo) Insert(ctx context.Context, category *ctmodel.Category) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyCategoryName,
			keyPocket,
			keyIsIncome,
			keyUpdatedAt,
			keyCreatedAt,
		).
		Values(
			category.CategoryName,
			category.Pocket,
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
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(ctx context.Context, category *ctmodel.Category) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyCategoryName: category.CategoryName,
			keyIsIncome:     category.IsIncome,
			keyUpdatedAt:    time.Now(),
		}).
		Where(sq.Eq{keyID: category.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit category: %w", err)
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

// Delete ...
func (r Repo) Delete(ctx context.Context, id uint64) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.Eq{keyID: id}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete category: %w", err)
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

// GetByID get one category by email
func (r Repo) GetByID(ctx context.Context, id uuid.UUID) (ctmodel.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyCategoryName,
		keyIsIncome,
		keyPocket,
		keyCreatedAt,
		keyUpdatedAt,
	).From(keyTable).Where(sq.Eq{keyID: id}).ToSql()

	if err != nil {
		return ctmodel.Category{}, fmt.Errorf("build query get category by id: %w", err)
	}

	var cat ctmodel.Category
	err = r.db.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&cat.ID,
			&cat.CategoryName,
			&cat.IsIncome,
			&cat.Pocket,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
	if err != nil {
		return ctmodel.Category{}, db.ParseError(err)
	}

	return cat, nil
}

// // Find get all pocket
// func (r Repo) Find(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error) {

// 	// Validation filter
// 	filter.SortSafelist = []string{"name", "-name", "updated_at", "-updated_at"}
// 	if err := filter.Validate(); err != nil {
// 		return nil, data.Metadata{}, db.ErrDBSortFilter
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 	defer cancel()

// 	sqlStatement, args, err := r.sb.Select(
// 		"count(*) OVER()",
// 		keyID,
// 		keyOwner,
// 		keyEditor,
// 		keyWatcher,
// 		keyPocketName,
// 		keyIcon,
// 		keyLevel,
// 		keyCreatedAt,
// 		keyUpdatedAt,
// 		keyVersion,
// 	).
// 		From(keyTable).
// 		Where(sq.Eq{keyOwner: owner}).
// 		OrderBy(filter.SortColumnDirection()).
// 		Limit(uint64(filter.Limit())).
// 		Offset(uint64(filter.Offset())).
// 		ToSql()

// 	if err != nil {
// 		return nil, data.Metadata{}, fmt.Errorf("build query find pocket: %w", err)
// 	}

// 	rows, err := r.db.Query(ctx, sqlStatement, args...)
// 	if err != nil {
// 		return nil, data.Metadata{}, db.ParseError(err)
// 	}
// 	defer rows.Close()

// 	totalRecords := 0
// 	pockets := make([]ptmodel.Pocket, 0)
// 	for rows.Next() {
// 		var pocket ptmodel.Pocket
// 		err := rows.Scan(
// 			&totalRecords,
// 			&pocket.ID,
// 			&pocket.Owner,
// 			&pocket.Editor,
// 			&pocket.Watcher,
// 			&pocket.PocketName,
// 			&pocket.Icon,
// 			&pocket.Level,
// 			&pocket.CreatedAt,
// 			&pocket.UpdatedAt,
// 			&pocket.Version)
// 		if err != nil {
// 			return nil, data.Metadata{}, db.ParseError(err)
// 		}
// 		pockets = append(pockets, pocket)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, data.Metadata{}, err
// 	}

// 	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

// 	return pockets, metadata, nil
// }

// // FindUserPockets get all pocket user has uuid in it
// func (r Repo) FindUserPockets(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error) {

// 	// Validation filter
// 	filter.SortSafelist = []string{"pocket_name", "-pocket_name", "updated_at", "-updated_at"}
// 	if err := filter.Validate(); err != nil {
// 		return nil, data.Metadata{}, db.ErrDBSortFilter
// 	}

// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 	defer cancel()

// 	// SELECT count(*) OVER(), id, owner, editor, watcher, pocket_name, icon, level, created_at, updated_at, version
// 	// FROM pockets
// 	// WHERE 'a502f2bf-f813-40e2-b39a-bec07374076f'=ANY(watcher)
// 	// ORDER BY pocket_name ASC LIMIT 50 OFFSET 0
// 	sqlStatement, args, err := r.sb.Select(
// 		"count(*) OVER()",
// 		keyID,
// 		keyOwner,
// 		keyEditor,
// 		keyWatcher,
// 		keyPocketName,
// 		keyIcon,
// 		keyLevel,
// 		keyCreatedAt,
// 		keyUpdatedAt,
// 		keyVersion,
// 	).
// 		From(keyTable).
// 		Where(fmt.Sprintf("'%s' = ANY(%s)", owner.String(), keyWatcher)).
// 		OrderBy(filter.SortColumnDirection()).
// 		Limit(uint64(filter.Limit())).
// 		Offset(uint64(filter.Offset())).
// 		ToSql()

// 	if err != nil {
// 		return nil, data.Metadata{}, fmt.Errorf("build query find user pocket: %w", err)
// 	}

// 	rows, err := r.db.Query(ctx, sqlStatement, args...)
// 	if err != nil {
// 		return nil, data.Metadata{}, db.ParseError(err)
// 	}
// 	defer rows.Close()

// 	totalRecords := 0
// 	pockets := make([]ptmodel.Pocket, 0)
// 	for rows.Next() {
// 		var pocket ptmodel.Pocket
// 		err := rows.Scan(
// 			&totalRecords,
// 			&pocket.ID,
// 			&pocket.Owner,
// 			&pocket.Editor,
// 			&pocket.Watcher,
// 			&pocket.PocketName,
// 			&pocket.Icon,
// 			&pocket.Level,
// 			&pocket.CreatedAt,
// 			&pocket.UpdatedAt,
// 			&pocket.Version)
// 		if err != nil {
// 			return nil, data.Metadata{}, db.ParseError(err)
// 		}
// 		pockets = append(pockets, pocket)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, data.Metadata{}, err
// 	}

// 	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

// 	return pockets, metadata, nil
// }
