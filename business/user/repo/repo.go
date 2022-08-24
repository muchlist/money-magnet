package repo

import (
	"context"
	"fmt"
	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/pkg/data"
	db2 "github.com/muchlist/moneymagnet/pkg/db"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	keyTable       = "users"
	keyID          = "id"
	keyEmail       = "email"
	keyName        = "name"
	keyPassword    = "password"
	keyRoles       = "roles"
	keyPocketRoles = "pocket_roles"
	keyFCM         = "fcm"
	keyCreatedAt   = "created_at"
	keyUpdatedAt   = "updated_at"
	keyVersion     = "version"
)

// Repo manages the set of APIs for user access.
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
func (r Repo) Insert(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Insert(keyTable).
		Columns(
			keyID,
			keyName,
			keyEmail,
			keyPassword,
			keyRoles,
			keyPocketRoles,
			keyFCM,
			keyCreatedAt,
			keyUpdatedAt).
		Values(
			user.ID,
			user.Name,
			user.Email,
			user.Password,
			user.Roles,
			user.PocketRoles,
			user.Fcm,
			user.CreatedAt,
			user.UpdatedAt).
		Suffix(db2.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query insert user: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&user.ID)
	if err != nil {
		return db2.ParseError(err)
	}

	return nil
}

// Edit ...
func (r Repo) Edit(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyName:        user.Name,
			keyEmail:       user.Email,
			keyRoles:       user.Roles,
			keyPocketRoles: user.PocketRoles,
			keyFCM:         user.Fcm,
			keyUpdatedAt:   time.Now(),
			keyVersion:     user.Version + 1,
		}).
		Where(sq.Eq{keyID: user.ID}).
		Suffix(db2.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit user: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&user.Version)
	if err != nil {
		return db2.ParseError(err)
	}

	return nil
}

func (r Repo) EditFCM(ctx context.Context, id uuid.UUID, fcm string) error {
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

	_, err = r.db.Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db2.ParseError(err)
	}

	return nil
}

// Delete ...
func (r Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(sq.And{
			sq.Eq{keyID: id},
		}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete user: %w", err)
	}

	res, err := r.db.Exec(ctx, sqlStatement, args...)
	if err != nil {
		return db2.ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return db2.ErrDBNotFound
	}

	return nil
}

// ChangePassword ...
func (r Repo) ChangePassword(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(sq.Eq{
			keyName:      user.Password,
			keyUpdatedAt: time.Now(),
			keyVersion:   user.Version + 1,
		}).
		Where(sq.Eq{keyID: user.ID}).
		Suffix(db2.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query change password user: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&user.Version)
	if err != nil {
		return db2.ParseError(err)
	}

	return nil
}

// =========================================================================
// GETTER

// GetByID get one user by uuid
func (r Repo) GetByID(ctx context.Context, uuid uuid.UUID) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyName,
		keyEmail,
		keyPassword,
		keyRoles,
		keyPocketRoles,
		keyFCM,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyID: uuid}).ToSql()

	if err != nil {
		return model.User{}, fmt.Errorf("build query get user by id: %w", err)
	}

	var user model.User
	err = r.db.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Roles,
			&user.PocketRoles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
	if err != nil {
		return model.User{}, db2.ParseError(err)
	}

	return user, nil
}

// GetByIDs get many user by []uuid
func (r Repo) GetByIDs(ctx context.Context, uuids []uuid.UUID) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyName,
		keyEmail,
		keyRoles,
		keyPocketRoles,
		keyFCM,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyID: uuids}).ToSql()

	if err != nil {
		return nil, fmt.Errorf("build query get user by ids: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, db2.ParseError(err)
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
			&user.PocketRoles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
		if err != nil {
			return nil, db2.ParseError(err)
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
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Select(
		keyID,
		keyName,
		keyEmail,
		keyPassword,
		keyRoles,
		keyPocketRoles,
		keyFCM,
		keyCreatedAt,
		keyUpdatedAt,
		keyVersion,
	).From(keyTable).Where(sq.Eq{keyEmail: email}).ToSql()

	if err != nil {
		return model.User{}, fmt.Errorf("build query get user by email: %w", err)
	}

	var user model.User
	err = r.db.QueryRow(ctx, sqlStatement, args...).
		Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Roles,
			&user.PocketRoles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
	if err != nil {
		return model.User{}, db2.ParseError(err)
	}

	return user, nil
}

// Find get all user
func (r Repo) Find(ctx context.Context, name string, filter data.Filters) ([]model.User, data.Metadata, error) {

	// Validation filter
	filter.SortSafelist = []string{"name", "-name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return nil, data.Metadata{}, db2.ErrDBSortFilter
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
		keyPocketRoles,
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
		return nil, data.Metadata{}, fmt.Errorf("build query find user: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStatement, args...)
	if err != nil {
		return nil, data.Metadata{}, db2.ParseError(err)
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
			&user.PocketRoles,
			&user.Fcm,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version)
		if err != nil {
			return nil, data.Metadata{}, db2.ParseError(err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, data.Metadata{}, err
	}

	metadata := data.CalculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return users, metadata, nil
}
