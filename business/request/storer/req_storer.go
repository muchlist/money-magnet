package storer

import (
	"context"

	"github.com/muchlist/moneymagnet/business/request/model"
	"github.com/muchlist/moneymagnet/pkg/data"
)

type RequestStorer interface {
	RequestSaver
	RequestReader
	Transactor
}

type RequestSaver interface {
	Insert(ctx context.Context, pocket *model.RequestPocket) error
	UpdateStatus(ctx context.Context, pocket *model.RequestPocket) error
}

type RequestReader interface {
	GetByID(ctx context.Context, id uint64) (model.RequestPocket, error)
	Find(ctx context.Context, findBy model.FindBy, filter data.Filters) ([]model.RequestPocket, data.Metadata, error)
}

type Transactor interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}
