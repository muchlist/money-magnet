package userrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
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

func (r Repo) Insert(ctx context.Context, user usermodel.User) (string, error) {
	ctxWT, cancel := context.WithTimeout(ctx, 3*time.Second)
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
		Suffix(db.Returning(keyEmail, keyName)).ToSql()

	if err != nil {
		return "", fmt.Errorf("build query Insert User: %w", err)
	}

	var email string
	var name string
	err = r.db.QueryRow(ctxWT, sqlStatement, args...).Scan(&email, &name)

	if err != nil {
		return "", db.ParseError(err)
	}

	return fmt.Sprintf("user %s with email %s successfuly created", name, email), nil
}
