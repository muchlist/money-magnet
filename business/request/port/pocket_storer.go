package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type PocketStorer interface {
	GetByID(ctx context.Context, id xulid.ULID) (model.Pocket, error)
	Edit(ctx context.Context, Pocket *model.Pocket) error
	InsertPocketUser(ctx context.Context, userIDs []string, pocketid xulid.ULID) error
}
