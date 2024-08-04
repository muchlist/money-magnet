package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type SpendResp struct {
	ID               xulid.ULID     `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FZZ"`
	UserID           xulid.ULID     `json:"user_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FNN"`
	UserName         string         `json:"user_name" example:"Muchlis"`
	PocketID         xulid.ULID     `json:"pocket_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FYY"`
	PocketName       string         `json:"pocket_name" example:"main pocket"`
	CategoryID       xulid.NullULID `json:"category_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FXX"`
	CategoryName     string         `json:"category_name" example:"food"`
	CategoryIcon     int            `json:"category_icon" example:"1"`
	Name             string         `json:"name" example:"Makan siang"`
	Price            int64          `json:"price" example:"50000"`
	BalanceSnapshoot int64          `json:"balance_snapshoot,omitempty" example:"0"`
	IsIncome         bool           `json:"is_income" example:"false"`
	SpendType        int            `json:"type" example:"2"`
	Date             time.Time      `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
	CreatedAt        time.Time      `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt        time.Time      `json:"updated_at" example:"2022-09-10T17:03:15.091267+08:00"`
	Version          int            `json:"version" example:"1"`
}

type NewSpend struct {
	ID         xulid.NullULID `json:"id" example:"01ARZ3NDEKTSV4RRFFQ69G5FZZ"`
	PocketID   xulid.ULID     `json:"pocket_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FYY"`
	CategoryID xulid.NullULID `json:"category_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FXX"`
	Name       string         `json:"name" example:"Makan siang"`
	Price      int64          `json:"price" example:"50000"`
	IsIncome   bool           `json:"is_income" example:"false"`
	SpendType  int            `json:"type" example:"2"`
	Date       time.Time      `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
}

type TransferSpend struct {
	PocketIDFrom xulid.ULID `json:"pocket_id_from" example:"01ARZ3NDEKTSV4RRFFQ69G5FYY"`
	PocketIDTo   xulid.ULID `json:"pocket_id_to" example:"01ARZ3NDEKTSV4RRFFQ69G5FXX"`
	Price        int64      `json:"price" example:"50000"`
	Date         time.Time  `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
}

type UpdateSpend struct {
	ID         xulid.ULID     `json:"-"`
	CategoryID xulid.NullULID `json:"category_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FXX"`
	Name       *string        `json:"name" example:"Makan siang"`
	Price      *int64         `json:"price" example:"50000"`
	IsIncome   *bool          `json:"is_income" example:"false"`
	SpendType  *int           `json:"type" example:"2"`
	Date       *time.Time     `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
}
