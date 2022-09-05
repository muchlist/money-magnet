package storer

import (
	"context"

	"github.com/muchlist/moneymagnet/business/user/model"

	"github.com/google/uuid"
)

type UserReader interface {
	GetByID(ctx context.Context, uuids uuid.UUID) (model.User, error)
	GetByIDs(ctx context.Context, uuids []uuid.UUID) ([]model.User, error)
}
