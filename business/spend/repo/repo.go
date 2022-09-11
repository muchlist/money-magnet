package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	keyTable       = "spends"
	keyID          = "id"
	keyUserID      = "user_id"
	keyPocketID    = "pocket_id"
	keyCategoryID  = "category_id"
	keyCategoryID2 = "category_id_2"
	keyName        = "name"
	keyPrice       = "price"
	keyBalance     = "balance_snapshoot"
	keyIsIncome    = "is_income"
	keyType        = "type"
	keyDate        = "date"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
	keyVersion     = "version"
)

// Repo manages the set of APIs for spend access.
type Repo struct {
	db  *pgxpool.Pool
	log mlogger.Logger
	sb  sq.StatementBuilderType
}

// NewRepo constructs a data for api access..
func NewRepo(sqlDB *pgxpool.Pool, logger mlogger.Logger) Repo {
	return Repo{
		db:  sqlDB,
		log: logger,
		sb:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// =========================================================================
// MANIPULATOR

// Insert ...
func (r Repo) Insert(ctx context.Context, spend *model.Spend) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyID,
			keyUserID,
			keyPocketID,
			keyCategoryID,
			keyCategoryID2,
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
			&spend.CategoryID2,
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

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&spend.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(ctx context.Context, spend *model.Spend) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyUserID:      spend.UserID,
			keyPocketID:    spend.PocketID,
			keyCategoryID:  spend.CategoryID,
			keyCategoryID2: spend.CategoryID2,
			keyName:        spend.Name,
			keyPrice:       spend.Price,
			keyBalance:     spend.BalanceSnapshoot,
			keyIsIncome:    spend.IsIncome,
			keyType:        spend.SpendType,
			keyDate:        spend.Date,
			keyCreatedAt:   spend.CreatedAt,
			keyUpdatedAt:   spend.UpdatedAt,
			keyVersion:     spend.Version + 1,
		}).
		Where(sq.Eq{keyID: spend.ID}).
		Suffix(db.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit spend: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&spend.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
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
		return fmt.Errorf("build query delete spend: %w", err)
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

// GetByID get one spend by email
func (r Repo) GetByID(ctx context.Context, id uuid.UUID) (model.Spend, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		db.A(keyID),
		db.A(keyUserID),
		db.A(keyPocketID),
		db.A(keyCategoryID),
		db.A(keyCategoryID2),
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
		db.CoalesceString(db.E("category_name"), ""),
	).
		From(keyTable + " A").
		LeftJoin("users B ON A.user_id = B.id").
		LeftJoin("pockets C ON A.pocket_id = C.id").
		LeftJoin("categories D ON A.category_id = D.id").
		LeftJoin("categories E ON A.category_id_2 = E.id").
		Where(sq.Eq{"A.id": id}).ToSql()

	if err != nil {
		return model.Spend{}, fmt.Errorf("build query get spend by id: %w", err)
	}

	var spend model.Spend
	err = r.db.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&spend.ID,
			&spend.UserID,
			&spend.PocketID,
			&spend.CategoryID,
			&spend.CategoryID2,
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
			&spend.CategoryName2,
		)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return model.Spend{}, db.ParseError(err)
	}

	return spend, nil
}

// Find get all spend
func (r Repo) Find(ctx context.Context, spendFilter model.SpendFilter, filter data.Filters) ([]model.Spend, data.Metadata, error) {

	// Validation filter
	filter.SortSafelist = []string{"-date", "date", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, data.Metadata{}, db.ErrDBSortFilter
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
		db.A(keyCategoryID2),
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
		db.CoalesceString(db.E("category_name"), ""),
	).
		From(keyTable + " A").
		LeftJoin("users B ON A.user_id = B.id").
		LeftJoin("pockets C ON A.pocket_id = C.id").
		LeftJoin("categories D ON A.category_id = D.id").
		LeftJoin("categories E ON A.category_id_2 = E.id")

	// WHERE builder
	// mapping where filter
	whereMap := sq.Eq{db.A(keyPocketID): spendFilter.PocketID.UUID}
	if spendFilter.User.Valid {
		whereMap[db.A(keyUserID)] = spendFilter.User.UUID
	}
	// if spendFilter.Category.Valid {
	// 	whereMap[db.A(keyCategoryID)] = spendFilter.Category.UUID
	// }
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
			sq.Or{
				sq.Eq{db.A(keyCategoryID): spendFilter.Category.UUID},
				sq.Eq{db.A(keyCategoryID2): spendFilter.Category.UUID},
			},
		)
	}
	if spendFilter.DateStart != nil {
		query = query.Where(sq.GtOrEq{db.A(keyDate): *spendFilter.DateStart})
	}
	if spendFilter.DateEnd != nil {
		query = query.Where(sq.Lt{db.A(keyDate): *spendFilter.DateEnd})
	}

	sqlStatement, args, err := query.OrderBy(filter.SortColumnDirection()).
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, data.Metadata{}, fmt.Errorf("build query find spend: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, data.Metadata{}, db.ParseError(err)
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
			&spend.CategoryID2,
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
			&spend.CategoryName2,
		)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, data.Metadata{}, db.ParseError(err)
		}
		spends = append(spends, spend)
	}

	if err := rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return spends, metadata, nil
}

// Count All Price
func (r Repo) CountAllPrice(ctx context.Context, pocketID uuid.UUID) (int64, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select("sum(price)").
		From(keyTable).
		Where(sq.Eq{"pocket_id": pocketID}).
		ToSql()

	if err != nil {
		return 0, fmt.Errorf("build query find spend: %w", err)
	}

	var balance int64
	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&balance)
	if err != nil {
		return 0, db.ParseError(err)
	}

	return balance, nil
}
