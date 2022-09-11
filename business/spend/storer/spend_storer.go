package storer

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/pkg/data"
)

type SpendStorer interface {
	SpendSaver
	SpendReader
	Transactor
}

type SpendSaver interface {
	Insert(ctx context.Context, spend *model.Spend) error
	Edit(ctx context.Context, spend *model.Spend) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SpendReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Spend, error)
	Find(ctx context.Context, spendFilter model.SpendFilter, filter data.Filters) ([]model.Spend, data.Metadata, error)
	CountAllPrice(ctx context.Context, pocketID uuid.UUID) (int64, error)
}

type Transactor interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}
