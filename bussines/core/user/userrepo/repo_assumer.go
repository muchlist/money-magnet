package userrepo

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
)

type UserRepoAssumer interface {
	UserSaver
	UserReader
}

type UserSaver interface {
	Insert(ctx context.Context, user *usermodel.User) error
	Edit(ctx context.Context, user *usermodel.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, user *usermodel.User) error
}

type UserReader interface {
	GetByID(ctx context.Context, id int) (usermodel.User, error)
	GetByEmail(ctx context.Context, email string) (usermodel.User, error)
	Find(ctx context.Context, name string, filter db.Filters) ([]usermodel.User, error)
}
