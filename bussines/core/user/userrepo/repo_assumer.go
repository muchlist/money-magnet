package userrepo

import (
	"context"

	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
)

type UserRepoAssumer interface {
	UserSaver
	UserReader
}

type UserSaver interface {
	Insert(ctx context.Context, user usermodel.User) (string, error)
	// Edit(ctx context.Context, userInput dto.UserEditModel) (*dto.UserModel, rest_err.APIError)
	// Delete(ctx context.Context, id int, filterMerchant int) rest_err.APIError
	// ChangePassword(ctx context.Context, input dto.UserModel) rest_err.APIError
}

type UserReader interface {
	// GetByID(ctx context.Context, id int) (*dto.UserModel, rest_err.APIError)
	// GetByEmail(ctx context.Context, email string) (*dto.UserModel, rest_err.APIError)
	// FindWithPagination(ctx context.Context, opt FindPaginationParams) ([]dto.UserModel, rest_err.APIError)
}
