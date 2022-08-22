package ctmodel

import (
	"time"

	"github.com/google/uuid"
)

// Category is simple struct so can be unified with model
type Category struct {
	ID           uuid.UUID
	Pocket       uint64
	CategoryName string
	IsIncome     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *Category) ToCategoryResp() CategoryResp {
	return CategoryResp{
		ID:           c.ID,
		Pocket:       c.Pocket,
		CategoryName: c.CategoryName,
		IsIncome:     c.IsIncome,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

type NewCategory struct {
	PocketID     uint64 `json:"pocket_id" validate:"required"`
	CategoryName string `json:"category_name" validate:"required"`
	IsIncome     bool   `json:"is_income"`
}

type UpdateCategory struct {
	ID           uuid.UUID
	PocketID     *uint64 `json:"pocket_id" validate:"required"`
	CategoryName *string `json:"category_name" validate:"required"`
	IsIncome     *bool   `json:"is_income"`
}

type CategoryResp struct {
	ID           uuid.UUID `json:"id"`
	Pocket       uint64    `json:"pocket_id"`
	CategoryName string    `json:"category_name"`
	IsIncome     bool      `json:"is_income"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"update_at"`
}
