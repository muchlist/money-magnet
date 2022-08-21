package storer

import (
	"context"
	"github.com/muchlist/moneymagnet/business/user/usermodel"
	"github.com/muchlist/moneymagnet/pkg/data"

	"github.com/google/uuid"
)

type UserStorer interface {
	UserSaver
	UserReader
}

type UserSaver interface {
	Insert(ctx context.Context, user *usermodel.User) error
	Edit(ctx context.Context, user *usermodel.User) error
	EditFCM(ctx context.Context, id uuid.UUID, fcm string) error
	Delete(ctx context.Context, id uuid.UUID) error
	ChangePassword(ctx context.Context, user *usermodel.User) error
}

type UserReader interface {
	GetByID(ctx context.Context, uuid uuid.UUID) (usermodel.User, error)
	GetByIDs(ctx context.Context, uuids []uuid.UUID) ([]usermodel.User, error)
	GetByEmail(ctx context.Context, email string) (usermodel.User, error)
	Find(ctx context.Context, name string, filter data.Filters) ([]usermodel.User, data.Metadata, error)
}
