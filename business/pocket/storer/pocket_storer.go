package storer

import (
	"context"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/data"

	"github.com/google/uuid"
)

/*
mockgen -source=business/pocket/storer/pocket_storer.go -destination=business/pocket/mock_storer/pocket_storer.go
*/

type PocketStorer interface {
	PocketSaver
	PocketReader
	Transactor
}

type PocketSaver interface {
	Insert(ctx context.Context, Pocket *model.Pocket) error
	Edit(ctx context.Context, Pocket *model.Pocket) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateBalance(ctx context.Context, pocketID uuid.UUID, balance int64, isSetOperaton bool) (int64, error)

	// many to many relation
	InsertPocketUser(ctx context.Context, userIDs []uuid.UUID, pocketID uuid.UUID) error
	DeletePocketUser(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) error
}

type PocketReader interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Pocket, error)
	Find(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error)
	FindUserPocketsByRelation(ctx context.Context, owner uuid.UUID, filter data.Filters) ([]model.Pocket, data.Metadata, error)
}

type Transactor interface {
	WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error
}
