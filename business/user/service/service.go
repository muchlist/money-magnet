package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/business/user/port"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/xulid"

	"github.com/muchlist/moneymagnet/pkg/mcrypto"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound           = errors.New("data not found")
	ErrInvalidID          = errors.New("ID is not in its proper form")
	ErrInvalidEmailOrPass = errors.New("email or password not valid")
)

const (
	expiredJWTToken        = 60 * 1       // 1 Hour
	expiredJWTRefreshToken = 60 * 24 * 10 // 60 days
)

// Core manages the set of APIs for user access.
type Core struct {
	log    mlogger.Logger
	repo   port.UserStorer
	crypto mcrypto.Crypter
	jwt    mjwt.TokenHandler
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	repo port.UserStorer,
	crypto mcrypto.Crypter,
	jwt mjwt.TokenHandler,
) Core {
	return Core{
		log:    log,
		repo:   repo,
		crypto: crypto,
		jwt:    jwt,
	}
}

// Login return detail user with access token and refresh token
func (s Core) Login(ctx context.Context, email, password string) (model.UserResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-Login")
	defer span.End()

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("%v: %w", err, ErrInvalidEmailOrPass)
	}

	if !s.crypto.IsPWAndHashPWMatch([]byte(password), user.Password) {
		return model.UserResp{}, fmt.Errorf("%v: %w", err, ErrInvalidEmailOrPass)
	}

	expired := time.Now().Add(time.Minute * expiredJWTToken).Unix()

	AccessClaims := mjwt.CustomClaim{
		Identity: user.ID.String(),
		Name:     user.Name,
		Exp:      expired,
		Type:     mjwt.Access,
		Fresh:    true,
		Roles:    user.Roles,
	}

	expiredRefresh := time.Now().Add(time.Minute * expiredJWTRefreshToken).Unix()
	RefreshClaims := mjwt.CustomClaim{
		Identity: user.ID.String(),
		Name:     user.Name,
		Exp:      expiredRefresh,
		Type:     mjwt.Refresh,
		Fresh:    false,
		Roles:    user.Roles,
	}

	accessToken, err := s.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("fail to generate token when login: %w", err)
	}
	refreshToken, err := s.jwt.GenerateToken(RefreshClaims)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("fail to generate token when login: %w", err)
	}

	response := model.UserResp{
		ID:                  user.ID,
		Email:               user.Email,
		Name:                user.Name,
		Roles:               user.Roles,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
		Version:             user.Version,
		AccessToken:         accessToken,
		AccessTokenExpired:  expired,
		RefreshToken:        refreshToken,
		RefreshTokenExpired: expiredRefresh,
	}

	return response, nil
}

// InsertUser used for register user
func (s Core) InsertUser(ctx context.Context, req model.UserRegisterReq) (model.UserResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-InsertUser")
	defer span.End()

	hashPassword, err := s.crypto.GenerateHash(req.Password)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("generate hashpw when insert user: %w", err)
	}

	if req.Roles == nil {
		req.Roles = []string{}
	}

	timeNow := time.Now()
	user := model.User{
		ID:        xulid.Instance().NewULID(),
		Email:     req.Email,
		Name:      req.Name,
		Password:  hashPassword,
		Roles:     req.Roles,
		Fcm:       "",
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		Version:   1,
	}

	err = s.repo.Insert(ctx, &user)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("insert user to db: %w", err)
	}

	return model.UserResp{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Email,
		Roles:     user.Roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Version:   user.Version,
	}, nil
}

// PatchUser do edit user with ignoring nil field
// ID is required
func (s Core) PatchUser(ctx context.Context, req model.UserUpdate) (model.UserResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-PatchUser")
	defer span.End()

	userExisting, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("get user: %w", err)
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
	if req.Password != nil {
		hashPassword, err := s.crypto.GenerateHash(*req.Password)
		if err != nil {
			return model.UserResp{}, fmt.Errorf("generate hashpw when edit user: %w", err)
		}
		userExisting.Password = hashPassword
	}
	if req.Fcm != nil {
		userExisting.Fcm = *req.Fcm
	}

	if err := s.repo.Edit(ctx, &userExisting); err != nil {
		return model.UserResp{}, fmt.Errorf("edit user: %w", err)
	}

	return userExisting.ToUserResp(), nil
}

// UpdateFCM do save fcm to database
func (s Core) UpdateFCM(ctx context.Context, id string, fcm string) error {
	ctx, span := observ.GetTracer().Start(ctx, "service-UpdateFCM")
	defer span.End()

	userID, err := xulid.Parse(id)
	if err != nil {
		return ErrInvalidID
	}
	if err := s.repo.EditFCM(ctx, userID, fcm); err != nil {
		return fmt.Errorf("edit fcm: %w", err)
	}
	return nil
}

// Delete ...
func (s Core) Delete(ctx context.Context, userIDToDelete xulid.ULID, userIDExecutor xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "service-Delete")
	defer span.End()

	if userIDExecutor == userIDToDelete {
		return errr.New("cannot delete self profile", 400)
	}
	return s.repo.Delete(ctx, userIDToDelete)
}

// Refresh do refresh token,
// access token in reslt is new but tagged as not fresh
func (s Core) Refresh(ctx context.Context, refreshToken string) (model.UserResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-Refresh")
	defer span.End()

	// validate token, signature and exp etc...
	token, err := s.jwt.ValidateToken(refreshToken)
	if err != nil {
		return model.UserResp{}, err
	}
	claims, err := s.jwt.ReadToken(token)
	if err != nil {
		return model.UserResp{}, err
	}

	// cek claims type token
	if claims.Type != mjwt.Refresh {
		return model.UserResp{}, mjwt.ErrInvalidToken
	}

	userID, _ := xulid.Parse(claims.Identity)
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("%v: %w", err, ErrInvalidEmailOrPass)
	}

	expired := time.Now().Add(time.Minute * expiredJWTToken).Unix()

	AccessClaims := mjwt.CustomClaim{
		Identity: user.ID.String(),
		Name:     user.Name,
		Exp:      expired,
		Type:     mjwt.Access,
		Fresh:    false,
		Roles:    user.Roles,
	}

	accessToken, err := s.jwt.GenerateToken(AccessClaims)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("fail to generate token when login: %w", err)
	}

	response := model.UserResp{
		ID:                 user.ID,
		Email:              user.Email,
		Name:               user.Name,
		Roles:              user.Roles,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
		Version:            user.Version,
		AccessToken:        accessToken,
		AccessTokenExpired: expired,
		RefreshToken:       refreshToken,
	}

	return response, nil
}

// GetProfile do load user by id
func (s Core) GetProfile(ctx context.Context, id string) (model.UserResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-GetProfile")
	defer span.End()

	userID, err := xulid.Parse(id)
	if err != nil {
		return model.UserResp{}, ErrInvalidID
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return model.UserResp{}, fmt.Errorf("get by id: %w", err)
	}
	return user.ToUserResp(), nil
}

// FindUserByName do find user filter by *name*
func (s Core) FindUserByName(ctx context.Context, name string, filter data.Filters) ([]model.UserResp, data.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindUserByName")
	defer span.End()

	users, metadata, err := s.repo.Find(ctx, name, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find user: %w", err)
	}
	usersResult := make([]model.UserResp, len(users))
	for i := range users {
		usersResult[i] = users[i].ToUserResp()
	}
	return usersResult, metadata, nil
}
