package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type UserStorer interface {
	GetByIDs(ctx context.Context, ulids []string) ([]model.User, error)

	// also we need to delete fcm when error send message
	EditFCM(ctx context.Context, id xulid.ULID, fcms []string) error
	AppendFCM(ctx context.Context, id xulid.ULID, fcm string) error
	RemoveFCM(ctx context.Context, id xulid.ULID, fcm string) error
}
