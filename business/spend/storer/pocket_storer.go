package storer

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/data"
)

type PocketReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Pocket, error)
	Find(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error)
	FindUserPockets(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error)
}
