package storer

import (
	"context"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/data"

	"github.com/google/uuid"
)

type PocketStorer interface {
	PocketSaver
	PocketReader
}

type PocketSaver interface {
	Insert(ctx context.Context, Pocket *model.Pocket) error
	Edit(ctx context.Context, Pocket *model.Pocket) error
	Delete(ctx context.Context, id uuid.UUID) error

	// many to many relation
	InsertPocketUser(ctx context.Context, userIDs []uuid.UUID, pocketID uuid.UUID) error
	DeletePocketUser(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) error
}

type PocketReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Pocket, error)
	Find(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error)
	FindUserPockets(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error)
}
