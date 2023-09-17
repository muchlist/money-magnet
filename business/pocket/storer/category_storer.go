package storer

import (
	"context"

	"github.com/muchlist/moneymagnet/business/category/model"
)

type CategorySaver interface {
	InsertMany(ctx context.Context, categories []model.Category) error
}
