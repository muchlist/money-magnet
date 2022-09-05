package storer

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/pocket/model"
)

type PocketStorer interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Pocket, error)
	Edit(ctx context.Context, Pocket *model.Pocket) error
	InsertPocketUser(ctx context.Context, userIDs []uuid.UUID, pocketID uuid.UUID) error
}
