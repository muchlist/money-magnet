package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type NewCategory struct {
	PocketID     xulid.ULID `json:"pocket_id" validate:"required" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	CategoryName string     `json:"category_name" validate:"required" example:"gaji"`
	CategoryIcon int        `json:"category_icon" example:"0"`
	IsIncome     bool       `json:"is_income" example:"true"`
}

type UpdateCategory struct {
	ID           xulid.ULID `json:"-" validate:"required"`
	CategoryName string     `json:"category_name" validate:"required" example:"gaji_2"`
	CategoryIcon int        `json:"category_icon" example:"0"`
}

type CategoryResp struct {
	ID           xulid.ULID `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FZZ"`
	PocketID     xulid.ULID `json:"pocket_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FAV"`
	CategoryName string     `json:"category_name" example:"gaji"`
	CategoryIcon int        `json:"category_icon" example:"0"`
	IsIncome     bool       `json:"is_income" example:"true"`
	CreatedAt    time.Time  `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt    time.Time  `json:"update_at" example:"2022-09-10T17:03:15.091267+08:00"`
}
