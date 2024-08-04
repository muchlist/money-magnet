package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/category/model"
)

//go:generate mockgen -source category_storer.go -destination mockport/mock_category_storer.go -package mockport
type CategorySaver interface {
	InsertMany(ctx context.Context, categories []model.Category) error
}
