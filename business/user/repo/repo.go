package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/xulid"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	keyTable     = "users"
	keyID        = "id"
	keyEmail     = "email"
	keyName      = "name"
	keyPassword  = "password"
	keyRoles     = "roles"
	keyFCM       = "fcm"
	keyCreatedAt = "created_at"
	keyUpdatedAt = "updated_at"
	keyVersion   = "version"
)

// Repo manages the set of APIs for user access.
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
func (r Repo) Insert(ctx context.Context, user *model.User) error {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-Insert")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyID,
			keyName,
			keyEmail,
			keyPassword,
			keyRoles,
			keyFCM,
			keyCreatedAt,
			keyUpdatedAt).
		Values(
			user.ID.String(),
			user.Name,
			user.Email,
			user.Password,
			user.Roles,
			user.Fcm,
			user.CreatedAt,
			user.UpdatedAt).
		Suffix(db.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query insert user: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&user.ID)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(ctx context.Context, user *model.User) error {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-Edit")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyName:      user.Name,
			keyEmail:     user.Email,
			keyRoles:     user.Roles,
			keyFCM:       user.Fcm,
			keyUpdatedAt: time.Now(),
			keyVersion:   user.Version + 1,
		}).
		Where(sq.Eq{keyID: user.ID}).
		Suffix(db.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit user: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&user.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

func (r Repo) EditFCM(ctx context.Context, id xulid.ULID, fcm string) error {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-EditFCM")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyFCM:       fcm,
			keyUpdatedAt: time.Now(),
		}).
		Where(sq.Eq{keyID: id}).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query update fcm user: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	_, err = dbtx.Exec(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// Delete ...
func (r Repo) Delete(ctx context.Context, id xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-Delete")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.And{
			sq.Eq{keyID: id},
		}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete user: %w", err)
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

// ChangePassword ...
func (r Repo) ChangePassword(ctx context.Context, user *model.User) error {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-ChangePassword")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyName:      user.Password,
			keyUpdatedAt: time.Now(),
			keyVersion:   user.Version + 1,
		}).
		Where(sq.Eq{keyID: user.ID}).
		Suffix(db.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query change password user: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	err = dbtx.QueryRow(ctx, sqlStatement, args...).Scan(&user.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return db.ParseError(err)
	}

	return nil
}

// =========================================================================
// GETTER

// GetByID get one user by ulid
func (r Repo) GetByID(ctx context.Context, ulid xulid.ULID) (model.User, error) {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-GetByID")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyName,
		keyEmail,
		keyPassword,
		keyRoles,
		keyFCM,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyID: ulid}).ToSql()

	if err != nil {
		return model.User{}, fmt.Errorf("build query get user by id: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	var user model.User
	err = dbtx.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Roles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return model.User{}, db.ParseError(err)
	}

	return user, nil
}

// GetByIDs get many user by []uuid
func (r Repo) GetByIDs(ctx context.Context, ulids []string) ([]model.User, error) {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-GetByIDs")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyName,
		keyEmail,
		keyRoles,
		keyFCM,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyID: ulids}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query get user by ids: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, db.ParseError(err)
	}
	defer rows.Close()

	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Roles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, db.ParseError(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByEmail get one user by email
func (r Repo) GetByEmail(ctx context.Context, email string) (model.User, error) {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-GetByEmail")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyName,
		keyEmail,
		keyPassword,
		keyRoles,
		keyFCM,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyEmail: email}).ToSql()

	if err != nil {
		return model.User{}, fmt.Errorf("build query get user by email: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	var user model.User
	err = dbtx.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Roles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return model.User{}, db.ParseError(err)
	}

	return user, nil
}

// Find get all user
func (r Repo) Find(ctx context.Context, name string, filter paging.Filters) ([]model.User, paging.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "user-repo-Find")
	defer span.End()

	// Validation filter
	filter.SortSafelist = []string{"name", "-name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, paging.Metadata{}, db.ErrDBSortFilter
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlFrom := r.sb.Select(
		"count(*) OVER()",
		keyID,
		keyName,
		keyEmail,
		keyPassword,
		keyRoles,
		keyFCM,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable)

	if len(name) > 0 {
		sqlFrom = sqlFrom.Where(sq.ILike{keyName: fmt.Sprint("%", name, "%")})
	}

	sqlStatement, args, err := sqlFrom.OrderBy(filter.SortColumnDirection(), keyCreatedAt+" ASC").
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return nil, paging.Metadata{}, fmt.Errorf("build query find user: %w", err)
	}

	dbtx := db.ExtractTx(ctx, r.db)

	rows, err := dbtx.Query(ctx, sqlStatement, args...)
	if err != nil {
		r.log.InfoT(ctx, err.Error())
		return nil, paging.Metadata{}, db.ParseError(err)
	}
	defer rows.Close()

	totalRecords := 0
	users := make([]model.User, 0)
	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Roles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
		if err != nil {
			r.log.InfoT(ctx, err.Error())
			return nil, paging.Metadata{}, db.ParseError(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, paging.Metadata{}, err
	}

	metadata := paging.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return users, metadata, nil
}
