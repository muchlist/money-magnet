package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/business/pocket/port"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/xulid"

	sq "github.com/Masterminds/squirrel"
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

// make sure the implementation satisfies the interface
var _ port.PocketStorer = (*Repo)(nil)

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
func (r Repo) Insert(ctx context.Context, pocket *model.Pocket) error {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-Insert")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	pocket.Sanitize()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyID,
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
			pocket.ID,
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

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&pocket.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(ctx context.Context, pocket *model.Pocket) error {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-Edit")
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

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&pocket.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// UpdateBalance ...
func (r Repo) UpdateBalance(ctx context.Context, pocketID xulid.ULID, balance int64, isSetOperaton bool) (int64, error) {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-UpdateBalance")
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

	dbtx := db.ExtractTx(ctx, r.db)

	var newBalance int64
	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&newBalance)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return 0, db.ParseError(err)
	}

	return newBalance, nil
}

// Delete ...
func (r Repo) Delete(ctx context.Context, id xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-Delete")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.Eq{keyID: id}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete pocket: %w", err)
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

// GetByID get one pocket by id
func (r Repo) GetByID(ctx context.Context, id xulid.ULID) (model.Pocket, error) {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-GetByID")
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

	dbtx := db.ExtractTx(ctx, r.db)

	var pocket model.Pocket
	err = dbtx.QueryRow(ctx, sqlStatement, args...).
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
func (r Repo) Find(ctx context.Context, owner xulid.ULID, filter paging.Filters) ([]model.Pocket, paging.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-Find")
	defer span.End()

	// Validation filter
	filter.SortSafelist = []string{"name", "-name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, paging.Metadata{}, db.ErrDBSortFilter
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
		return nil, paging.Metadata{}, fmt.Errorf("build query find pocket: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, paging.Metadata{}, db.ParseError(err)
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
			return nil, paging.Metadata{}, db.ParseError(err)
		}
		pockets = append(pockets, pocket)
	}

	if err := rows.Err(); err != nil {
		return nil, paging.Metadata{}, err
	}

	metadata := paging.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return pockets, metadata, nil
}

// FindUserPockets get all pocket user has uuid in it by relation constrain
func (r Repo) FindUserPocketsByRelation(ctx context.Context, owner xulid.ULID, filter paging.Filters) ([]model.Pocket, paging.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "pocket-repo-FindUserPocketsByRelation")
	defer span.End()

	// Validation filter
	filter.SortSafelist = []string{"pocket_name", "-pocket_name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, paging.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

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
		return nil, paging.Metadata{}, fmt.Errorf("build query find user pocket: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, paging.Metadata{}, db.ParseError(err)
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
			return nil, paging.Metadata{}, db.ParseError(err)
		}
		pockets = append(pockets, pocket)
	}

	if err := rows.Err(); err != nil {
		return nil, paging.Metadata{}, err
	}

	metadata := paging.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return pockets, metadata, nil
}
