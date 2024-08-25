package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/request/model"
	"github.com/muchlist/moneymagnet/pkg/paging"
)

type RequestStorer interface {
	RequestSaver
	RequestReader
}

type RequestSaver interface {
	Insert(ctx context.Context, pocket *model.RequestPocket) error
	UpdateStatus(ctx context.Context, pocket *model.RequestPocket) error
}

type RequestReader interface {
	GetByID(ctx context.Context, id uint64) (model.RequestPocket, error)
	Find(ctx context.Context, findBy model.FindBy, filter paging.Filters) ([]model.RequestPocket, paging.Metadata, error)
}

type Transactor interface {
	WithAtomic(ctx context.Context, tFunc func(ctx context.Context) error) error
}
