package userrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
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
	sb squirrel.StatementBuilderType
}

// NewRepo constructs a data for api access..
func NewRepo(sqlDB *pgxpool.Pool) UserRepoAssumer {
	return Repo{
		db: sqlDB,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// =========================================================================
// MANIPULATOR

// Insert implements UserRepoAssumer
func (r Repo) Insert(ctx context.Context, user *usermodel.User) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	timeNow := time.Now()

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
			timeNow,
			timeNow).
		Suffix(db.Returning(keyID)).ToSql()

	if err != nil {
		return fmt.Errorf("build query insert user: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&user.ID)
	if err != nil {
		db.ParseError(err)
	}

	return nil
}

// Edit implements UserRepoAssumer
func (r Repo) Edit(ctx context.Context, user *usermodel.User) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(squirrel.Eq{
			keyName:        user.Name,
			keyRoles:       user.Roles,
			keyPocketRoles: user.PocketRoles,
			keyUpdatedAt:   time.Now(),
			keyVersion:     user.Version + 1,
		}).
		Where(squirrel.Eq{keyID: user.ID}).
		Suffix(db.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query edit user: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&user.Version)
	if err != nil {
		db.ParseError(err)
	}

	return nil
}

// Delete implements UserRepoAssumer
func (r Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Delete(keyTable).
		Where(squirrel.And{
			squirrel.Eq{keyID: id},
		}).ToSql()
	if err != nil {
		return fmt.Errorf("build query delete user: %w", err)
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

// ChangePassword implements UserRepoAssumer
func (r Repo) ChangePassword(ctx context.Context, user *usermodel.User) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlStatement, args, err := r.sb.Update(keyTable).
		SetMap(squirrel.Eq{
			keyName:      user.Password,
			keyUpdatedAt: time.Now(),
			keyVersion:   user.Version + 1,
		}).
		Where(squirrel.Eq{keyID: user.ID}).
		Suffix(db.Returning(keyVersion)).
		ToSql()

	if err != nil {
		return fmt.Errorf("build query change password user: %w", err)
	}

	err = r.db.QueryRow(ctx, sqlStatement, args...).Scan(&user.Version)
	if err != nil {
		db.ParseError(err)
	}

	return nil
}

// =========================================================================
// GETTER

// GetByID implements UserRepoAssumer
func (r Repo) GetByID(ctx context.Context, id int) (usermodel.User, error) {
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
	).From(keyTable).Where(squirrel.Eq{keyID: id}).ToSql()

	if err != nil {
		return usermodel.User{}, fmt.Errorf("build query get user by id: %w", err)
	}

	var user usermodel.User
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
		db.ParseError(err)
	}

	return user, nil
}

// GetByEmail implements UserRepoAssumer
func (r Repo) GetByEmail(ctx context.Context, email string) (usermodel.User, error) {
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
	).From(keyTable).Where(squirrel.Eq{keyEmail: email}).ToSql()

	if err != nil {
		return usermodel.User{}, fmt.Errorf("build query get user by email: %w", err)
	}

	var user usermodel.User
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
		db.ParseError(err)
	}

	return user, nil
}

// Find implements UserRepoAssumer
func (r Repo) Find(ctx context.Context, name string, filter db.Filters) ([]usermodel.User, error) {

	// Validation filter
	filter.SortSafelist = []string{"name", "-name", "updated_at", "-updated_at"}
	if err := filter.Validate(); err != nil {
		return []usermodel.User{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	sqlFrom := r.sb.Select(
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
		sqlFrom = sqlFrom.Where(squirrel.ILike{keyName: fmt.Sprint("%", name, "%")})
	}

	sqlStatement, args, err := sqlFrom.OrderBy(filter.SortColumnDirection(), keyCreatedAt+" ASC").
		Limit(uint64(filter.Limit())).
		Offset(uint64(filter.Offset())).
		ToSql()

	if err != nil {
		return []usermodel.User{}, fmt.Errorf("build query find user: %w", err)
	}

	rows, err := r.db.Query(ctx, sqlStatement, args...)
	if err != nil {
		return []usermodel.User{}, db.ParseError(err)
	}
	defer rows.Close()

	users := make([]usermodel.User, 0)
	for rows.Next() {
		var user usermodel.User
		err := rows.Scan(
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
			return []usermodel.User{}, db.ParseError(err)
		}
		users = append(users, user)
	}

	if len(users) == 0 {
		return []usermodel.User{}, db.ErrDBNotFound
	}

	return users, nil
}
