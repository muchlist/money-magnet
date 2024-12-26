package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/xulid"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	keyTable      = "spends"
	keyID         = "id"
	keyUserID     = "user_id"
	keyPocketID   = "pocket_id"
	keyCategoryID = "category_id"
	keyName       = "name"
	keyPrice      = "price"
	keyBalance    = "balance_snapshoot"
	keyIsIncome   = "is_income"
	keyType       = "type"
	keyDate       = "date"
	keyCreatedAt  = "created_at"
	keyUpdatedAt  = "updated_at"
	keyVersion    = "version"
)

// Repo manages the set of APIs for spend access.
type Repo struct {
	db  *pgxpool.Pool
	log mlogger.Logger
	sb  sq.StatementBuilderType
}

// NewRepo constructs a data for api access..
func NewRepo(sqlDB *pgxpool.Pool, logger mlogger.Logger) *Repo {
	return &Repo{
		db:  sqlDB,
		log: logger,
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// =========================================================================
// MANIPULATOR

// Insert ...
func (r *Repo) Insert(ctx context.Context, spend *model.Spend) error {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-Insert")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyID,
			keyUserID,
			keyPocketID,
			keyCategoryID,
			keyName,
			keyPrice,
			keyBalance,
			keyIsIncome,
			keyType,
			keyDate,
			keyCreatedAt,
			keyUpdatedAt,
			keyVersion,
		).
		Values(
			&spend.ID,
			&spend.UserID,
			&spend.PocketID,
			&spend.CategoryID,
			&spend.Name,
			&spend.Price,
			&spend.BalanceSnapshoot,
			&spend.IsIncome,
			&spend.SpendType,
			&spend.Date,
			&spend.CreatedAt,
			&spend.UpdatedAt,
			&spend.Version,
		).
		Suffix(db.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query insert spend: %w,", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&spend.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r *Repo) Edit(ctx context.Context, spend *model.Spend) error {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-Edit")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyUserID:     spend.UserID,
			keyPocketID:   spend.PocketID,
			keyCategoryID: spend.CategoryID,
			keyName:       spend.Name,
			keyPrice:      spend.Price,
			keyBalance:    spend.BalanceSnapshoot,
			keyIsIncome:   spend.IsIncome,
			keyType:       spend.SpendType,
			keyDate:       spend.Date,
			keyCreatedAt:  spend.CreatedAt,
			keyUpdatedAt:  spend.UpdatedAt,
			keyVersion:    spend.Version + 1,
		}).
		Where(sq.Eq{keyID: spend.ID}).
		Suffix(db.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit spend: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&spend.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Delete ...
func (r *Repo) Delete(ctx context.Context, id xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-Delete")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.Eq{keyID: id}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete spend: %w", err)
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

// GetByID get one spend by email
func (r *Repo) GetByID(ctx context.Context, id xulid.ULID) (model.Spend, error) {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-GetByID")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		db.A(keyID),
		db.A(keyUserID),
		db.A(keyPocketID),
		db.A(keyCategoryID),
		db.A(keyName),
		db.A(keyPrice),
		db.A(keyBalance),
		db.A(keyIsIncome),
		db.A(keyType),
		db.A(keyDate),
		db.A(keyCreatedAt),
		db.A(keyUpdatedAt),
		db.A(keyVersion),
		db.B("name"),
		db.C("pocket_name"),
		db.CoalesceString(db.D("category_name"), ""),
		db.CoalesceInt(db.D("category_icon"), 0),
	).
		From(keyTable + " A").
		LeftJoin("users B ON A.user_id = B.id").
		LeftJoin("pockets C ON A.pocket_id = C.id").
		LeftJoin("categories D ON A.category_id = D.id").
		Where(sq.Eq{"A.id": id}).ToSql()

	if err != nil {
		return model.Spend{}, fmt.Errorf("build query get spend by id: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	var spend model.Spend
	err = dbtx.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&spend.ID,
			&spend.UserID,
			&spend.PocketID,
			&spend.CategoryID,
			&spend.Name,
			&spend.Price,
			&spend.BalanceSnapshoot,
			&spend.IsIncome,
			&spend.SpendType,
			&spend.Date,
			&spend.CreatedAt,
			&spend.UpdatedAt,
			&spend.Version,
			&spend.UserName,
			&spend.PocketName,
			&spend.CategoryName,
			&spend.CategoryIcon,
		)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return model.Spend{}, db.ParseError(err)
	}

	return spend, nil
}

// Find get all spend
func (r *Repo) Find(ctx context.Context, spendFilter model.SpendFilter, filter paging.Filters) ([]model.Spend, paging.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-Find")
	defer span.End()

	// Validation filter
	filter.SortSafelist = []string{"-date", "date", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, paging.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	//  := r.sb.Select(
	query := r.sb.Select(
		"count(*) OVER()",
		db.A(keyID),
		db.A(keyUserID),
		db.A(keyPocketID),
		db.A(keyCategoryID),
		db.A(keyName),
		db.A(keyPrice),
		db.A(keyBalance),
		db.A(keyIsIncome),
		db.A(keyType),
		db.A(keyDate),
		db.A(keyCreatedAt),
		db.A(keyUpdatedAt),
		db.A(keyVersion),
		db.B("name"),
		db.C("pocket_name"),
		db.CoalesceString(db.D("category_name"), ""),
		db.CoalesceInt(db.D("category_icon"), 0),
	).
		From("spends A").
		LeftJoin("users B ON A.user_id = B.id").
		LeftJoin("pockets C ON A.pocket_id = C.id").
		LeftJoin("categories D ON A.category_id = D.id")

	// WHERE builder
	// mapping where filter equal
	whereMap := sq.Eq{db.A(keyPocketID): spendFilter.PocketID.ULID}
	if spendFilter.User.Valid {
		whereMap[db.A(keyUserID)] = spendFilter.User.ULID
	}
	if spendFilter.IsIncome != nil {
		whereMap[db.A(keyIsIncome)] = *spendFilter.IsIncome
	}
	if len(spendFilter.Type) != 0 {
		whereMap[db.A(keyType)] = spendFilter.Type
	}

	// building where clause
	query = query.Where(whereMap)
	if spendFilter.Category.Valid {
		query = query.Where(
			sq.Eq{db.A(keyCategoryID): spendFilter.Category.ULID},
		)
	}
	if spendFilter.DateStart != nil {
		query = query.Where(sq.GtOrEq{db.A(keyDate): *spendFilter.DateStart})
	}
	if spendFilter.DateEnd != nil {
		query = query.Where(sq.Lt{db.A(keyDate): *spendFilter.DateEnd})
	}
	// searchable name
	if spendFilter.Name != "" {
		// query = query.Where(sq.Expr("A.name % ?", spendFilter.Name))
		query = query.Where(
			sq.Or{
				sq.Expr("A.name % ?", spendFilter.Name),
				sq.Like{"A.name": fmt.Sprint("%", spendFilter.Name, "%")},
			},
		)
	}

	sqlStatement, args, err := query.OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, paging.Metadata{}, fmt.Errorf("build query find spend: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, paging.Metadata{}, db.ParseError(err)
	}
	defer rows.Close()

	totalRecords := 0
	spends := make([]model.Spend, 0)
	for rows.Next() {
		var spend model.Spend
		err := rows.Scan(
			&totalRecords,
			&spend.ID,
			&spend.UserID,
			&spend.PocketID,
			&spend.CategoryID,
			&spend.Name,
			&spend.Price,
			&spend.BalanceSnapshoot,
			&spend.IsIncome,
			&spend.SpendType,
			&spend.Date,
			&spend.CreatedAt,
			&spend.UpdatedAt,
			&spend.Version,
			&spend.UserName,
			&spend.PocketName,
			&spend.CategoryName,
			&spend.CategoryIcon,
		)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, paging.Metadata{}, db.ParseError(err)
		}
		spends = append(spends, spend)
	}

	if err := rows.Err(); err != nil {
		return nil, paging.Metadata{}, err
	}

	metadata := paging.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return spends, metadata, nil
}

// Find With Cursor Pagination
func (r *Repo) FindWithCursor(ctx context.Context, spendFilter model.SpendFilter, filter paging.Cursor) ([]model.Spend, error) {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-FindWithCursor")
	defer span.End()

	// Validation filter
	filter.SetCursorList([]string{"-date", "date", "updated_at", "-updated_at"})
	if err := filter.Validate(); err != nil {
		return nil, db.ErrDBInvalidCursorType
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := r.sb.Select(
		db.A(keyID),
		db.A(keyUserID),
		db.A(keyPocketID),
		db.A(keyCategoryID),
		db.A(keyName),
		db.A(keyPrice),
		db.A(keyBalance),
		db.A(keyIsIncome),
		db.A(keyType),
		db.A(keyDate),
		db.A(keyCreatedAt),
		db.A(keyUpdatedAt),
		db.A(keyVersion),
		db.B("name"),
		db.C("pocket_name"),
		db.CoalesceString(db.D("category_name"), ""),
		db.CoalesceInt(db.D("category_icon"), 0),
	).
		From("spends A").
		LeftJoin("users B ON A.user_id = B.id").
		LeftJoin("pockets C ON A.pocket_id = C.id").
		LeftJoin("categories D ON A.category_id = D.id")

	// WHERE builder
	// mapping where filter equal
	whereMap := sq.Eq{db.A(keyPocketID): spendFilter.PocketID.ULID}
	if spendFilter.User.Valid {
		whereMap[db.A(keyUserID)] = spendFilter.User.ULID
	}
	if spendFilter.IsIncome != nil {
		whereMap[db.A(keyIsIncome)] = *spendFilter.IsIncome
	}
	if len(spendFilter.Type) != 0 {
		whereMap[db.A(keyType)] = spendFilter.Type
	}

	// building where clause
	query = query.Where(whereMap)

	// apply cursor value
	if filter.GetCursor() != "" {
		cursorColumn, _ := filter.GetCursorColumn()
		direction := filter.GetDirection()

		if direction == ">" {
			query = query.Where(sq.Gt{cursorColumn: filter.GetCursor()})
		} else {
			query = query.Where(sq.Lt{cursorColumn: filter.GetCursor()})
		}
	}

	if spendFilter.Category.Valid {
		query = query.Where(
			sq.Eq{db.A(keyCategoryID): spendFilter.Category.ULID},
		)
	}
	if spendFilter.DateStart != nil {
		query = query.Where(sq.GtOrEq{db.A(keyDate): *spendFilter.DateStart})
	}
	if spendFilter.DateEnd != nil {
		query = query.Where(sq.Lt{db.A(keyDate): *spendFilter.DateEnd})
	}

	// searchable name
	if spendFilter.Name != "" {
		query = query.Where(
			sq.Or{
				sq.Expr("A.name % ?", spendFilter.Name),
				sq.Like{"A.name": fmt.Sprint("%", spendFilter.Name, "%")},
			},
		)
	}

	// apply order by
	orderByStr, err := filter.GetSortColumnDirection()
	if err != nil {
		return nil, db.ErrDBSortFilter
	}

	sqlStatement, args, err := query.OrderBy(orderByStr).
		Limit(uint64(filter.GetPageSizePlusOne())).
		ToSql()

	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, fmt.Errorf("build query find spend: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, db.ParseError(err)
	}
	defer rows.Close()

	spends := make([]model.Spend, 0)
	for rows.Next() {
		var spend model.Spend
		err := rows.Scan(
			&spend.ID,
			&spend.UserID,
			&spend.PocketID,
			&spend.CategoryID,
			&spend.Name,
			&spend.Price,
			&spend.BalanceSnapshoot,
			&spend.IsIncome,
			&spend.SpendType,
			&spend.Date,
			&spend.CreatedAt,
			&spend.UpdatedAt,
			&spend.Version,
			&spend.UserName,
			&spend.PocketName,
			&spend.CategoryName,
			&spend.CategoryIcon,
		)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, db.ParseError(err)
		}
		spends = append(spends, spend)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spends, nil
}

// FindWithCursorMultiPockets With Cursor Pagination
func (r *Repo) FindWithCursorMultiPockets(ctx context.Context, spendFilter model.SpendFilterMultiPocket, filter paging.Cursor) ([]model.Spend, error) {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-FindWithCursorMultiPockets")
	defer span.End()

	// Validation filter
	filter.SetCursorList([]string{"-date", "date", "updated_at", "-updated_at"})
	if err := filter.Validate(); err != nil {
		return nil, db.ErrDBInvalidCursorType
	}

	// Must have at least one pocketID
	if len(spendFilter.Pockets) == 0 {
		return nil, db.ErrDBInvalidInput
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := r.sb.Select(
		db.A(keyID),
		db.A(keyUserID),
		db.A(keyPocketID),
		db.A(keyCategoryID),
		db.A(keyName),
		db.A(keyPrice),
		db.A(keyBalance),
		db.A(keyIsIncome),
		db.A(keyType),
		db.A(keyDate),
		db.A(keyCreatedAt),
		db.A(keyUpdatedAt),
		db.A(keyVersion),
		db.B("name"),
		db.C("pocket_name"),
		db.CoalesceString(db.D("category_name"), ""),
		db.CoalesceInt(db.D("category_icon"), 0),
	).
		From("spends A").
		LeftJoin("users B ON A.user_id = B.id").
		LeftJoin("pockets C ON A.pocket_id = C.id").
		LeftJoin("categories D ON A.category_id = D.id")

	// WHERE builder
	// mapping where filter equal

	whereMap := sq.Eq{db.A(keyPocketID): spendFilter.Pockets}

	if len(spendFilter.Users) != 0 {
		whereMap[db.A(keyUserID)] = spendFilter.Users
	}
	if len(spendFilter.Categories) != 0 {
		whereMap[db.A(keyCategoryID)] = spendFilter.Categories
	}
	if spendFilter.IsIncome != nil {
		whereMap[db.A(keyIsIncome)] = *spendFilter.IsIncome
	}
	if len(spendFilter.Types) != 0 {
		whereMap[db.A(keyType)] = spendFilter.Types
	}

	// building where clause
	query = query.Where(whereMap)

	// apply cursor value
	if filter.GetCursor() != "" {
		cursorColumn, _ := filter.GetCursorColumn()
		direction := filter.GetDirection()

		if direction == ">" {
			query = query.Where(sq.Gt{cursorColumn: filter.GetCursor()})
		} else {
			query = query.Where(sq.Lt{cursorColumn: filter.GetCursor()})
		}
	}
	if spendFilter.DateStart != nil {
		query = query.Where(sq.GtOrEq{db.A(keyDate): *spendFilter.DateStart})
	}
	if spendFilter.DateEnd != nil {
		query = query.Where(sq.Lt{db.A(keyDate): *spendFilter.DateEnd})
	}

	// searchable name
	if spendFilter.Name != "" {
		query = query.Where(
			sq.Or{
				sq.Expr("A.name % ?", spendFilter.Name),
				sq.Like{"A.name": fmt.Sprint("%", spendFilter.Name, "%")},
			},
		)
	}

	// apply order by
	orderByStr, err := filter.GetSortColumnDirection()
	if err != nil {
		return nil, db.ErrDBSortFilter
	}

	sqlStatement, args, err := query.OrderBy(orderByStr).
		Limit(uint64(filter.GetPageSizePlusOne())).
		ToSql()

	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, fmt.Errorf("build query find spend: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, db.ParseError(err)
	}
	defer rows.Close()

	spends := make([]model.Spend, 0)
	for rows.Next() {
		var spend model.Spend
		err := rows.Scan(
			&spend.ID,
			&spend.UserID,
			&spend.PocketID,
			&spend.CategoryID,
			&spend.Name,
			&spend.Price,
			&spend.BalanceSnapshoot,
			&spend.IsIncome,
			&spend.SpendType,
			&spend.Date,
			&spend.CreatedAt,
			&spend.UpdatedAt,
			&spend.Version,
			&spend.UserName,
			&spend.PocketName,
			&spend.CategoryName,
			&spend.CategoryIcon,
		)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, db.ParseError(err)
		}
		spends = append(spends, spend)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return spends, nil
}

// Count All Price
func (r *Repo) CountAllPrice(ctx context.Context, pocketID xulid.ULID) (int64, error) {
	ctx, span := observ.GetTracer().Start(ctx, "spend-repo-CountAllPrice")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select("sum(price)").
		From(keyTable).
		Where(sq.Eq{"pocket_id": pocketID}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("build query find spend: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	var balance int64
	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&balance)
	if err != nil {
		return 0, db.ParseError(err)
	}

	return balance, nil
}
