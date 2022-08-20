package storer

import (
	"context"
	"github.com/muchlist/moneymagnet/bussines/user/usermodel"

	"github.com/google/uuid"
)

type UserReader interface {
	GetByID(ctx context.Context, uuids uuid.UUID) (usermodel.User, error)
	GetByIDs(ctx context.Context, uuids []uuid.UUID) ([]usermodel.User, error)
}
