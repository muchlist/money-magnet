package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type PocketStorer interface {
	GetByID(ctx context.Context, id xulid.ULID) (model.Pocket, error)
	Find(ctx context.Context, owner xulid.ULID, filter data.Filters) ([]model.Pocket, data.Metadata, error)
	FindUserPocketsByRelation(ctx context.Context, owner xulid.ULID, filter data.Filters) ([]model.Pocket, data.Metadata, error)

	UpdateBalance(ctx context.Context, pocketid xulid.ULID, balance int64, isSetOperaton bool) (int64, error)
}
