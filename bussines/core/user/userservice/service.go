package userservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
	"github.com/muchlist/moneymagnet/bussines/sys/mjwt"
	"github.com/muchlist/moneymagnet/foundation/mcrypto"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound           = errors.New("data not found")
	ErrInvalidID          = errors.New("ID is not in its proper form")
	ErrInvalidEmailOrPass = errors.New("email or password not valid")
)

const (
	expiredJWTToken        = 60 * 1       // 1 Hour
	expiredJWTRefreshToken = 15 * 24 * 10 // 15 days
)

// Service manages the set of APIs for user access.
type Service struct {
	log    mlogger.Logger
	repo   UserStorer
	crypto mcrypto.Crypter
	jwt    mjwt.TokenHandler
}

// NewService constructs a core for user api access.
func NewService(
	log mlogger.Logger,
	repo UserStorer,
	crypto mcrypto.Crypter,
	jwt mjwt.TokenHandler,
) Service {
	return Service{
		log:    log,
		repo:   repo,
		crypto: crypto,
		jwt:    jwt,
	}
}

// Login ...
func (s Service) Login(ctx context.Context, email, password string) (usermodel.UserResp, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("%v: %w", err, ErrInvalidEmailOrPass)
	}

	if !s.crypto.IsPWAndHashPWMatch([]byte(password), user.Password) {
		return usermodel.UserResp{}, fmt.Errorf("%v: %w", err, ErrInvalidEmailOrPass)
	}

	expired := time.Now().Add(time.Minute * expiredJWTToken).Unix()

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID.String(),
		Name:        user.Name,
		Exp:         expired,
		Type:        mjwt.Access,
		Fresh:       true,
		Roles:       user.Roles,
		PocketRoles: user.PocketRoles,
	}

	expired = time.Now().Add(time.Minute * expiredJWTRefreshToken).Unix()
	RefreshClaims := mjwt.CustomClaim{
		Identity:    user.ID.String(),
		Name:        user.Name,
		Exp:         expired,
		Type:        mjwt.Refresh,
		Fresh:       false,
		Roles:       user.Roles,
		PocketRoles: user.PocketRoles,
	}

	accessToken, err := s.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("fail to generate token when login: %w", err)
	}
	refreshToken, err := s.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("fail to generate token when login: %w", err)
	}

	response := usermodel.UserResp{
		ID:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		Roles:        user.Roles,
		PocketRoles:  user.PocketRoles,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Version:      user.Version,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return response, nil
}

// InsertUser
func (s Service) InsertUser(ctx context.Context, req usermodel.UserRegisterReq) (usermodel.UserResp, error) {

	hashPassword, err := s.crypto.GenerateHash(req.Password)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("generate hashpw when insert user: %w", err)
	}

	if req.Roles == nil {
		req.Roles = []string{}
	}
	if req.PocketRoles == nil {
		req.PocketRoles = []string{}
	}

	timeNow := time.Now()
	user := usermodel.User{
		ID:          uuid.New(),
		Email:       req.Email,
		Name:        req.Name,
		Password:    hashPassword,
		Roles:       req.Roles,
		PocketRoles: req.PocketRoles,
		Fcm:         "",
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
		Version:     1,
	}

	err = s.repo.Insert(ctx, &user)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("insert user to db: %w", err)
	}

	return usermodel.UserResp{
		ID:          user.ID,
		Email:       user.Email,
		Name:        user.Email,
		Roles:       user.Roles,
		PocketRoles: user.PocketRoles,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Version:     user.Version,
	}, nil
}

func (s Service) GetProfile(ctx context.Context, id string) (usermodel.UserResp, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("get by id: %w", err)
	}
	return user.ToUserResp(), nil
}
