package ptrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/bussines/pocket/ptmodel"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	keyTable      = "pockets"
	keyID         = "id"
	keyOwner      = "owner"
	keyEditor     = "editor"
	keyWatcher    = "watcher"
	keyPocketName = "pocket_name"
	keyIcon       = "icon"
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
			keyEditor,
			keyWatcher,
			keyVersion,
			keyIcon,
			keyLevel,
			keyUpdatedAt,
			keyCreatedAt,
		).
		Values(
			pocket.PocketName,
			pocket.Owner,
			pocket.Editor,
			pocket.Watcher,
			pocket.Version,
			pocket.Icon,
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
			keyEditor:     pocket.Editor,
			keyWatcher:    pocket.Watcher,
			keyIcon:       pocket.Icon,
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
		keyEditor,
		keyWatcher,
		keyPocketName,
		keyIcon,
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
			&pocket.Editor,
			&pocket.Watcher,
			&pocket.PocketName,
			&pocket.Icon,
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
		keyEditor,
		keyWatcher,
		keyPocketName,
		keyIcon,
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
			&pocket.Editor,
			&pocket.Watcher,
			&pocket.PocketName,
			&pocket.Icon,
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

// FindUserPockets get all pocket user has uuid in it
func (r Repo) FindUserPockets(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error) {

	// Validation filter
	filter.SortSafelist = []string{"pocket_name", "-pocket_name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, data.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// SELECT count(*) OVER(), id, owner, editor, watcher, pocket_name, icon, level, created_at, updated_at, version
	// FROM pockets
	// WHERE 'a502f2bf-f813-40e2-b39a-bec07374076f'=ANY(watcher)
	// ORDER BY pocket_name ASC LIMIT 50 OFFSET 0
	sqlStatement, args, err := r.sb.Select(
		"count(*) OVER()",
		keyID,
		keyOwner,
		keyEditor,
		keyWatcher,
		keyPocketName,
		keyIcon,
		keyLevel,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).
		From(keyTable).
		Where(fmt.Sprintf("'%s' = ANY(%s)", owner.String(), keyWatcher)).
		OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("build query find user pocket: %w", err)
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
			&pocket.Editor,
			&pocket.Watcher,
			&pocket.PocketName,
			&pocket.Icon,
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

// FindUserPockets get all pocket user has uuid in it by relation constrain
func (r Repo) FindUserPocketsByRelation(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error) {

	// Validation filter
	filter.SortSafelist = []string{"pocket_name", "-pocket_name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, data.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// SELECT count(*) OVER(), id, owner, editor, watcher, pocket_name, icon, level, created_at, updated_at, version
	// FROM pockets
	// WHERE 'a502f2bf-f813-40e2-b39a-bec07374076f'=ANY(watcher)
	// ORDER BY pocket_name ASC LIMIT 50 OFFSET 0
	sqlStatement, args, err := r.sb.Select(
		"count(*) OVER()",
		db.A(keyID),
		db.A(keyOwner),
		db.A(keyEditor),
		db.A(keyWatcher),
		db.A(keyPocketName),
		db.A(keyIcon),
		db.A(keyLevel),
		db.A(keyCreatedAt),
		db.A(keyUpdatedAt),
		db.A(keyVersion),
	).
		From("pockets A").
		Join("user_pocket B ON A.id = B.pocket_id").
		Join("users C ON B.user_id = C.id").
		Where(sq.Eq{"C.id": owner}).
		OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("build query find user pocket: %w", err)
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
			&pocket.Editor,
			&pocket.Watcher,
			&pocket.PocketName,
			&pocket.Icon,
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
