package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

// Category is simple struct so can be unified with model
type Category struct {
	ID               xulid.ULID
	PocketID         xulid.ULID
	CategoryName     string
	CategoryIcon     int
	IsIncome         bool
	DefaultSpendType int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (c *Category) ToCategoryResp() CategoryResp {
	return CategoryResp{
		ID:               c.ID,
		PocketID:         c.PocketID,
		CategoryName:     c.CategoryName,
		CategoryIcon:     c.CategoryIcon,
		IsIncome:         c.IsIncome,
		DefaultSpendType: c.DefaultSpendType,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}
