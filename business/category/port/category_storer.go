package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/pkg/paging"
)

type CategoryStorer interface {
	CategorySaver
	CategoryReader
}

type CategorySaver interface {
	Insert(ctx context.Context, category *model.Category) error
	InsertMany(ctx context.Context, categories []model.Category) error
	Edit(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
}

type CategoryReader interface {
	GetByID(ctx context.Context, id string) (model.Category, error)
	Find(ctx context.Context, pocketID string, filter paging.Filters) ([]model.Category, paging.Metadata, error)
}
