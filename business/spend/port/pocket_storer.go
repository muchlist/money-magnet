package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type PocketStorer interface {
	GetByID(ctx context.Context, id xulid.ULID) (model.Pocket, error)
	Find(ctx context.Context, owner xulid.ULID, filter paging.Filters) ([]model.Pocket, paging.Metadata, error)
	FindUserPocketsByRelation(ctx context.Context, owner xulid.ULID, filter paging.Filters) ([]model.Pocket, paging.Metadata, error)

	UpdateBalance(ctx context.Context, pocketid xulid.ULID, balance int64, isSetOperaton bool) (int64, error)
}
