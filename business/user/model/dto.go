package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type UserLoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserRegisterReq struct {
	Name     string   `json:"name" validate:"required"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required"`
	Roles    []string `json:"roles"`
}

type UserResp struct {
	ID                  xulid.ULID `json:"id"`
	Email               string     `json:"email"`
	Name                string     `json:"name"`
	Roles               []string   `json:"roles"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	Version             int        `json:"version"`
	AccessToken         string     `json:"access_token,omitempty"`
	RefreshToken        string     `json:"refresh_token,omitempty"`
	AccessTokenExpired  int64      `json:"access_token_expired,omitempty"`
	RefreshTokenExpired int64      `json:"refresh_token_expired,omitempty"`
}

type UserUpdate struct {
	ID       xulid.ULID `json:"-"`
	Email    *string    `json:"email"`
	Name     *string    `json:"name"`
	Password *string    `json:"password"`
	Roles    []string   `json:"roles"`
	Fcm      *string    `json:"fcm"`
}
