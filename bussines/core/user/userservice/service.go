package userservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/user/storer"
	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
	"github.com/muchlist/moneymagnet/bussines/sys/data"
	"github.com/muchlist/moneymagnet/bussines/sys/errr"
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
	repo   storer.UserStorer
	crypto mcrypto.Crypter
	jwt    mjwt.TokenHandler
}

// NewService constructs a core for user api access.
func NewService(
	log mlogger.Logger,
	repo storer.UserStorer,
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

// Login return detail user with access token and refresh token
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

// InsertUser used for register user
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

// FetchUser do edit user with ignoring nil field
// ID is required
func (s Service) FetchUser(ctx context.Context, req usermodel.UserUpdate) (usermodel.UserResp, error) {
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return usermodel.UserResp{}, ErrInvalidID
	}

	userExisting, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("get user: %w", err)
	}

	if req.Email != nil {
		userExisting.Email = *req.Email
	}
	if req.Name != nil {
		userExisting.Name = *req.Name
	}
	if req.Roles != nil {
		userExisting.Roles = req.Roles
	}
	if req.PocketRoles != nil {
		userExisting.PocketRoles = req.PocketRoles
	}
	if req.Password != nil {
		hashPassword, err := s.crypto.GenerateHash(*req.Password)
		if err != nil {
			return usermodel.UserResp{}, fmt.Errorf("generate hashpw when edit user: %w", err)
		}
		userExisting.Password = hashPassword
	}
	if req.Fcm != nil {
		userExisting.Fcm = *req.Fcm
	}

	if err := s.repo.Edit(ctx, &userExisting); err != nil {
		return usermodel.UserResp{}, fmt.Errorf("edit user: %w", err)
	}

	return userExisting.ToUserResp(), nil
}

// UpdateFCM do save fcm to database
func (s Service) UpdateFCM(ctx context.Context, id string, fcm string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	if err := s.repo.EditFCM(ctx, userID, fcm); err != nil {
		return fmt.Errorf("edit fcm: %w", err)
	}
	return nil
}

// Delete ...
func (s Service) Delete(ctx context.Context, userIDToDelete string, userIDExecutor string) error {
	userID, err := uuid.Parse(userIDToDelete)
	if err != nil {
		return ErrInvalidID
	}
	if userIDExecutor == userIDToDelete {
		return errr.New("cannot delete self profile", 400)
	}
	return s.repo.Delete(ctx, userID)
}

// Refresh do refresh token,
// access token in reslt is new but tagged as not fresh
func (s Service) Refresh(ctx context.Context, refreshToken string) (usermodel.UserResp, error) {
	// validate token, signature and exp etc...
	token, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return usermodel.UserResp{}, err
	}
	claims, err := s.jwt.ReadToken(token)
	if err != nil {
		return usermodel.UserResp{}, err
	}

	// cek claims type token
	if claims.Type != mjwt.Refresh {
		return usermodel.UserResp{}, mjwt.ErrInvalidToken
	}

	userID, _ := uuid.Parse(claims.Identity)
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("%v: %w", err, ErrInvalidEmailOrPass)
	}

	expired := time.Now().Add(time.Minute * expiredJWTToken).Unix()

	AccessClaims := mjwt.CustomClaim{
		Identity:    user.ID.String(),
		Name:        user.Name,
		Exp:         expired,
		Type:        mjwt.Access,
		Fresh:       false,
		Roles:       user.Roles,
		PocketRoles: user.PocketRoles,
	}

	accessToken, err := s.jwt.GenerateToken(AccessClaims)
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

// GetProfile do load user by id
func (s Service) GetProfile(ctx context.Context, id string) (usermodel.UserResp, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return usermodel.UserResp{}, ErrInvalidID
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return usermodel.UserResp{}, fmt.Errorf("get by id: %w", err)
	}
	return user.ToUserResp(), nil
}

// FindUserByName do find user filter by *name*
func (s Service) FindUserByName(ctx context.Context, name string, filter data.Filters) ([]usermodel.UserResp, data.Metadata, error) {
	users, metadata, err := s.repo.Find(ctx, name, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find user: %w", err)
	}
	usersResult := make([]usermodel.UserResp, len(users))
	for i := range users {
		usersResult[i] = users[i].ToUserResp()
	}
	return usersResult, metadata, nil
}
