package storer

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/category/ctmodel"
	"github.com/muchlist/moneymagnet/pkg/data"
)

type CategoryStorer interface {
	CategorySaver
	CategoryReader
}

type CategorySaver interface {
	Insert(ctx context.Context, category *ctmodel.Category) error
	Edit(ctx context.Context, category *ctmodel.Category) error
	Delete(ctx context.Context, id uint64) error
}

type CategoryReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (ctmodel.Category, error)
	Find(ctx context.Context, pocketID uint64, filter data.Filters) ([]ctmodel.Category, data.Metadata, error)
}
