package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type UserStorer interface {
	UserSaver
	UserReader
}

type UserSaver interface {
	Insert(ctx context.Context, user *model.User) error
	Edit(ctx context.Context, user *model.User) error
	EditFCM(ctx context.Context, id xulid.ULID, fcm string) error
	Delete(ctx context.Context, id xulid.ULID) error
	ChangePassword(ctx context.Context, user *model.User) error
}

type UserReader interface {
	GetByID(ctx context.Context, ulid xulid.ULID) (model.User, error)
	GetByIDs(ctx context.Context, ulids []string) ([]model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	Find(ctx context.Context, name string, filter paging.Filters) ([]model.User, paging.Metadata, error)
}
