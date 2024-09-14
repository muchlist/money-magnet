package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/ctype"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type Spend struct {
	ID               xulid.ULID
	UserID           xulid.ULID
	UserName         string // Join
	PocketID         xulid.ULID
	PocketName       string // Join
	CategoryID       xulid.NullULID
	CategoryName     string // Join
	CategoryIcon     int    // Join
	Name             ctype.UppercaseString
	Price            int64
	BalanceSnapshoot int64
	IsIncome         bool
	SpendType        int
	Date             time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Version          int
}

func (s *Spend) ToResp() SpendResp {
	return SpendResp{
		ID:               s.ID,
		UserID:           s.UserID,
		UserName:         s.UserName,
		PocketID:         s.PocketID,
		PocketName:       s.PocketName,
		CategoryID:       s.CategoryID,
		CategoryName:     s.CategoryName,
		CategoryIcon:     s.CategoryIcon,
		Name:             ctype.FromUppercaseString(s.Name),
		Price:            s.Price,
		BalanceSnapshoot: s.BalanceSnapshoot,
		IsIncome:         s.IsIncome,
		SpendType:        s.SpendType,
		Date:             s.Date,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
		Version:          s.Version,
	}
}
