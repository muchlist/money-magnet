package model

import (
	"time"

	"github.com/google/uuid"
)

// Category is simple struct so can be unified with model
type Category struct {
	ID           uuid.UUID
	PocketID     uuid.UUID
	CategoryName string
	IsIncome     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *Category) ToCategoryResp() CategoryResp {
	return CategoryResp{
		ID:           c.ID,
		PocketID:     c.PocketID,
		CategoryName: c.CategoryName,
		IsIncome:     c.IsIncome,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

type NewCategory struct {
	PocketID     uuid.UUID `json:"pocket_id" validate:"required" example:"f9339be2-6b05-4acb-a269-5309c39bae91"`
	CategoryName string    `json:"category_name" validate:"required" example:"gaji"`
	IsIncome     bool      `json:"is_income" example:"true"`
}

type UpdateCategory struct {
	ID           uuid.UUID `json:"-" validate:"required"`
	CategoryName string    `json:"category_name" validate:"required" example:"gaji_2"`
}

type CategoryResp struct {
	ID           uuid.UUID `json:"id" example:"bead2cb0-692e-41d2-a623-c44d1e19f2a0"`
	PocketID     uuid.UUID `json:"pocket_id" example:"f9339be2-6b05-4acb-a269-5309c39bae91"`
	CategoryName string    `json:"category_name" example:"gaji"`
	IsIncome     bool      `json:"is_income" example:"true"`
	CreatedAt    time.Time `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt    time.Time `json:"update_at" example:"2022-09-10T17:03:15.091267+08:00"`
}
