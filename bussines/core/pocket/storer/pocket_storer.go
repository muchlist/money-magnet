package storer

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptmodel"
	"github.com/muchlist/moneymagnet/bussines/sys/data"
)

type PocketStorer interface {
	PocketSaver
	PocketReader
}

type PocketSaver interface {
	Insert(ctx context.Context, Pocket *ptmodel.Pocket) error
	Edit(ctx context.Context, Pocket *ptmodel.Pocket) error
	Delete(ctx context.Context, id uint64) error

	// many to many relation
	InsertPocketUser(ctx context.Context, userIDs []uuid.UUID, pocketID uint64) error
	DeletePocketUser(ctx context.Context, userID uuid.UUID, pocketID uint64) error
}

type PocketReader interface {
	GetByID(ctx context.Context, id uint64) (ptmodel.Pocket, error)
	Find(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error)
	FindUserPockets(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error)
}
