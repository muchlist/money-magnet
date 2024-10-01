package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

//go:generate mockgen -source pocket_storer.go -destination mockport/mock_pocket_storer.go -package mockport
type PocketStorer interface {
	PocketSaver
	PocketReader
}

type PocketSaver interface {
	Insert(ctx context.Context, Pocket *model.Pocket) error
	Edit(ctx context.Context, Pocket *model.Pocket) error
	Delete(ctx context.Context, id xulid.ULID) error
	UpdateBalance(ctx context.Context, pocketID xulid.ULID, balance int64, isSetOperaton bool) (int64, error)

	// many to many relation
	InsertPocketUser(ctx context.Context, userIDs []string, pocketID xulid.ULID) error
	DeletePocketUser(ctx context.Context, userID xulid.ULID, pocketID xulid.ULID) error
}

type PocketReader interface {
	GetByID(ctx context.Context, id xulid.ULID) (model.Pocket, error)
	GetFirst(ctx context.Context, ownerID string) (model.Pocket, error)
	Find(ctx context.Context, owner xulid.ULID, filter paging.Filters) ([]model.Pocket, paging.Metadata, error)
	FindUserPocketsByRelation(ctx context.Context, owner xulid.ULID, filter paging.Filters) ([]model.Pocket, paging.Metadata, error)
}

type Transactor interface {
	WithAtomic(ctx context.Context, tFunc func(ctx context.Context) error) error
}
