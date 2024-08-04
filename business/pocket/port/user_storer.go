package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

//go:generate mockgen -source user_storer.go -destination mockport/mock_user_storer.go -package mockport
type UserReader interface {
	GetByID(ctx context.Context, ulid xulid.ULID) (model.User, error)
	GetByIDs(ctx context.Context, ulids []string) ([]model.User, error)
}
