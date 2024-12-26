package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type SpendStorer interface {
	SpendSaver
	SpendReader
}

type SpendSaver interface {
	Insert(ctx context.Context, spend *model.Spend) error
	Edit(ctx context.Context, spend *model.Spend) error
	Delete(ctx context.Context, id xulid.ULID) error
}

type SpendReader interface {
	GetByID(ctx context.Context, id xulid.ULID) (model.Spend, error)
	Find(ctx context.Context, spendFilter model.SpendFilter, filter paging.Filters) ([]model.Spend, paging.Metadata, error)
	FindWithCursor(ctx context.Context, spendFilter model.SpendFilter, filter paging.Cursor) ([]model.Spend, error)
	FindWithCursorMultiPockets(ctx context.Context, spendFilter model.SpendFilterMultiPocket, filter paging.Cursor) ([]model.Spend, error)
	CountAllPrice(ctx context.Context, pocketid xulid.ULID) (int64, error)
}

type Transactor interface {
	WithAtomic(ctx context.Context, tFunc func(ctx context.Context) error) error
}
