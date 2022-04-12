package ptservice

import (
	"context"

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
}

type PocketReader interface {
	GetByID(ctx context.Context, id uint64) (ptmodel.Pocket, error)
	Find(ctx context.Context, name string, filter data.Filters) ([]ptmodel.Pocket, data.Metadata, error)
}
