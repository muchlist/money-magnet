package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	keyTable      = "pockets"
	keyID         = "id"
	keyOwnerID    = "owner_id"
	keyEditorID   = "editor_id"
	keyWatcherID  = "watcher_id"
	keyPocketName = "pocket_name"
	keyBalance    = "balance"
	keyCurrency   = "currency"
	keyIcon       = "icon"
	keyLevel      = "level"
	keyCreatedAt  = "created_at"
	keyUpdatedAt  = "updated_at"
	keyVersion    = "version"
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
func (r Repo) Insert(pctx context.Context, pocket *model.Pocket) error {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-Insert")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pocket.Sanitize()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyPocketName,
			keyCurrency,
			keyOwnerID,
			keyEditorID,
			keyWatcherID,
			keyVersion,
			keyIcon,
			keyLevel,
			keyUpdatedAt,
			keyCreatedAt,
		).
		Values(
			pocket.PocketName,
			pocket.Currency,
			pocket.OwnerID,
			pocket.EditorID,
			pocket.WatcherID,
			pocket.Version,
			pocket.Icon,
			pocket.Level,
			pocket.CreatedAt,
			pocket.UpdatedAt).
		Suffix(db.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query insert pocket: %w", err)
	}

	err = r.mod(ctx).QueryRow(ctx, sqlStatement, args...).Scan(&pocket.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(pctx context.Context, pocket *model.Pocket) error {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-Edit")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pocket.Sanitize()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyPocketName: pocket.PocketName,
			keyCurrency:   pocket.Currency,
			keyOwnerID:    pocket.OwnerID,
			keyEditorID:   pocket.EditorID,
			keyWatcherID:  pocket.WatcherID,
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

	err = r.mod(ctx).QueryRow(ctx, sqlStatement, args...).Scan(&pocket.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// UpdateBalance ...
func (r Repo) UpdateBalance(pctx context.Context, pocketID uuid.UUID, balance int64, isSetOperaton bool) (int64, error) {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-UpdateBalance")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	setMapValue := sq.Eq{}
	if isSetOperaton {
		setMapValue[keyBalance] = balance
	} else {
		operation := "+" // when minus operation convert to + -20000 while alse understood by postgres
		setMapValue[keyBalance] = sq.Expr(fmt.Sprintf("(pockets.balance %s %d)", operation, balance))
	}

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(setMapValue).
		Where(sq.Eq{keyID: pocketID}).
		Suffix(db.Returning(keyBalance)).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query update pocket balance: %w", err)
	}

	var newBalance int64
	err = r.mod(ctx).QueryRow(ctx, sqlStatement, args...).Scan(&newBalance)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return 0, db.ParseError(err)
	}

	return newBalance, nil
}

// Delete ...
func (r Repo) Delete(pctx context.Context, id uuid.UUID) error {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-Delete")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.Eq{keyID: id}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete pocket: %w", err)
	}

	res, err := r.mod(ctx).Exec(ctx, sqlStatement, args...)
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

// GetByID get one pocket by id
func (r Repo) GetByID(pctx context.Context, id uuid.UUID) (model.Pocket, error) {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-GetByID")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyOwnerID,
		keyEditorID,
		keyWatcherID,
		keyPocketName,
		keyBalance,
		keyCurrency,
		keyIcon,
		keyLevel,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyID: id}).ToSql()

	if err != nil {
		return model.Pocket{}, fmt.Errorf("build query get pocket by id: %w", err)
	}

	var pocket model.Pocket
	err = r.mod(ctx).QueryRow(ctx, sqlStatement, args...).
		Scan(
			&pocket.ID,
			&pocket.OwnerID,
			&pocket.EditorID,
			&pocket.WatcherID,
			&pocket.PocketName,
			&pocket.Balance,
			&pocket.Currency,
			&pocket.Icon,
			&pocket.Level,
			&pocket.CreatedAt,
			&pocket.UpdatedAt,
			&pocket.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return model.Pocket{}, db.ParseError(err)
	}

	return pocket, nil
}

// Find get all pocket
func (r Repo) Find(pctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error) {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-Find")
	defer span.End()

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
		keyOwnerID,
		keyEditorID,
		keyWatcherID,
		keyPocketName,
		keyBalance,
		keyCurrency,
		keyIcon,
		keyLevel,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).
		From(keyTable).
		Where(sq.Eq{keyOwnerID: owner}).
		OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("build query find pocket: %w", err)
	}

	rows, err := r.mod(ctx).Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, data.Metadata{}, db.ParseError(err)
	}
	defer rows.Close()

	totalRecords := 0
	pockets := make([]model.Pocket, 0)
	for rows.Next() {
		var pocket model.Pocket
		err := rows.Scan(
			&totalRecords,
			&pocket.ID,
			&pocket.OwnerID,
			&pocket.EditorID,
			&pocket.WatcherID,
			&pocket.PocketName,
			&pocket.Balance,
			&pocket.Currency,
			&pocket.Icon,
			&pocket.Level,
			&pocket.CreatedAt,
			&pocket.UpdatedAt,
			&pocket.Version)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
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
// DEPRECATED
// func (r Repo) FindUserPockets(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error) {

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
// 		keyOwnerID,
// 		keyEditorID,
// 		keyWatcherID,
// 		keyPocketName,
// 		keyBalance,
// 		keyCurrency,
// 		keyIcon,
// 		keyLevel,
// 		keyCreatedAt,
// 		keyUpdatedAt,
// 		keyVersion,
// 	).
// 		From(keyTable).
// 		Where(fmt.Sprintf("'%s' = ANY(%s)", owner.String(), keyWatcherID)).
// 		OrderBy(filter.SortColumnDirection()).
// 		Limit(uint64(filter.Limit())).
// 		Offset(uint64(filter.Offset())).
// 		ToSql()

// 	if err != nil {
// 		return nil, data.Metadata{}, fmt.Errorf("build query find user pocket: %w", err)
// 	}

// 	rows, err := r.mod(ctx).Query(ctx, sqlStatement, args...)
// 	if err != nil {
// 		r.log.InfoT(ctx, err.Error())
// 		return nil, data.Metadata{}, db.ParseError(err)
// 	}
// 	defer rows.Close()

// 	totalRecords := 0
// 	pockets := make([]model.Pocket, 0)
// 	for rows.Next() {
// 		var pocket model.Pocket
// 		err := rows.Scan(
// 			&totalRecords,
// 			&pocket.ID,
// 			&pocket.OwnerID,
// 			&pocket.EditorID,
// 			&pocket.WatcherID,
// 			&pocket.PocketName,
// 			&pocket.Balance,
// 			&pocket.Currency,
// 			&pocket.Icon,
// 			&pocket.Level,
// 			&pocket.CreatedAt,
// 			&pocket.UpdatedAt,
// 			&pocket.Version)
// 		if err != nil {
// 			r.log.InfoT(ctx, err.Error())
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

// FindUserPockets get all pocket user has uuid in it by relation constrain
func (r Repo) FindUserPocketsByRelation(pctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error) {
	ctx, span := observ.GetTracer().Start(pctx, "pocket-repo-FindUserPocketsByRelation")
	defer span.End()

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
		db.A(keyOwnerID),
		db.A(keyEditorID),
		db.A(keyWatcherID),
		db.A(keyPocketName),
		db.A(keyBalance),
		db.A(keyCurrency),
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

	rows, err := r.mod(ctx).Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, data.Metadata{}, db.ParseError(err)
	}
	defer rows.Close()

	totalRecords := 0
	pockets := make([]model.Pocket, 0)
	for rows.Next() {
		var pocket model.Pocket
		err := rows.Scan(
			&totalRecords,
			&pocket.ID,
			&pocket.OwnerID,
			&pocket.EditorID,
			&pocket.WatcherID,
			&pocket.PocketName,
			&pocket.Balance,
			&pocket.Currency,
			&pocket.Icon,
			&pocket.Level,
			&pocket.CreatedAt,
			&pocket.UpdatedAt,
			&pocket.Version)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
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
