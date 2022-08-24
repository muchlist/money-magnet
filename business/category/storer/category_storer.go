package storer

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/pkg/data"
)

type CategoryStorer interface {
	CategorySaver
	CategoryReader
}

type CategorySaver interface {
	Insert(ctx context.Context, category *model.Category) error
	Edit(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CategoryReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Category, error)
	Find(ctx context.Context, pocketID uint64, filter data.Filters) ([]model.Category, data.Metadata, error)
}
