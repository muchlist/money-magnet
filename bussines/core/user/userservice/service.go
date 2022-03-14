package userservice

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
	"github.com/muchlist/moneymagnet/bussines/core/user/userrepo"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Service manages the set of APIs for user access.
type Service struct {
	log  mlogger.Logger
	repo userrepo.UserRepoAssumer
}

// NewService constructs a core for user api access.
func NewService(log mlogger.Logger, repo userrepo.UserRepoAssumer) Service {
	return Service{
		log:  log,
		repo: repo,
	}
}

// InsertUser melakukan register user
func (s Service) InsertUser(ctx context.Context, user usermodel.UserReq) (string, error) {

	userInput := usermodel.User{
		ID:          uuid.New(),
		Email:       "whois.muchlas@gmail.com",
		Name:        "Muchlis",
		Password:    "secret",
		Roles:       []string{"admin"},
		PocketRoles: []string{"asdsadadasd:read"},
		Fcm:         "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// return s.repo.Insert(ctx, userInput)
	msg, err := s.repo.Insert(ctx, userInput)
	if err != nil {
		return "", err
	}
	return msg, nil

	// // cek ketersediaan id
	// _, err := s.dao.CheckIDAvailable(ctx, user.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// // END cek ketersediaan id

	// hashPassword, err := u.crypto.GenerateHash(user.Password)
	// if err != nil {
	// 	return nil, err
	// }

	// user.Password = hashPassword
	// user.Timestamp = time.Now().Unix()

	// insertedID, err := u.dao.InsertUser(ctx, user)
	// if err != nil {
	// 	return nil, err
	// }
	// return insertedID, nil
}
