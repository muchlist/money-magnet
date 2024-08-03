package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/user/model"

	"github.com/google/uuid"
)

//go:generate mockgen -source user_storer.go -destination mockport/mock_user_storer.go -package mockport
type UserReader interface {
	GetByID(ctx context.Context, uuids uuid.UUID) (model.User, error)
	GetByIDs(ctx context.Context, uuids []uuid.UUID) ([]model.User, error)
}
