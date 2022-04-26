package storer

import (
	"context"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
)

type UserReader interface {
	GetByID(ctx context.Context, uuids uuid.UUID) (usermodel.User, error)
	GetByIDs(ctx context.Context, uuids []uuid.UUID) ([]usermodel.User, error)
}
