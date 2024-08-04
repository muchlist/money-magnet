package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type User struct {
	ID        xulid.ULID
	Email     string
	Name      string
	Password  []byte
	Roles     []string
	Fcm       string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

func (u *User) ToUserResp() UserResp {
	return UserResp{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Roles:        u.Roles,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Version:      u.Version,
		AccessToken:  "",
		RefreshToken: "",
	}
}
