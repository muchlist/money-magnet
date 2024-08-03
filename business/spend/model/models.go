package model

import (
	"time"

	"github.com/google/uuid"
)

type Spend struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	UserName         string // Join
	PocketID         uuid.UUID
	PocketName       string // Join
	CategoryID       uuid.NullUUID
	CategoryName     string // Join
	CategoryIcon     int    // Join
	Name             string
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
		Name:             s.Name,
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

type SpendResp struct {
	ID               uuid.UUID     `json:"id" example:"f9339be2-6b05-4acb-a269-5309c39bae90"`
	UserID           uuid.UUID     `json:"user_id" example:"f9339be2-6b05-4acb-a269-5309c39bae89"`
	UserName         string        `json:"user_name" example:"Muchlis"`
	PocketID         uuid.UUID     `json:"pocket_id" example:"f9339be2-6b05-4acb-a269-5309c39bae91"`
	PocketName       string        `json:"pocket_name" example:"main pocket"`
	CategoryID       uuid.NullUUID `json:"category_id" example:"f9339be2-6b05-4acb-a269-5309c39bae92"`
	CategoryName     string        `json:"category_name" example:"food"`
	CategoryIcon     int           `json:"category_icon" example:"1"`
	Name             string        `json:"name" example:"Makan siang"`
	Price            int64         `json:"price" example:"50000"`
	BalanceSnapshoot int64         `json:"balance_snapshoot" example:"0"`
	IsIncome         bool          `json:"is_income" example:"false"`
	SpendType        int           `json:"type" example:"2"`
	Date             time.Time     `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
	CreatedAt        time.Time     `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt        time.Time     `json:"updated_at" example:"2022-09-10T17:03:15.091267+08:00"`
	Version          int           `json:"version" example:"1"`
}

type NewSpend struct {
	ID         uuid.NullUUID `json:"id" example:"f9339be2-6b05-4acb-a269-5309c39bae90"`
	PocketID   uuid.UUID     `json:"pocket_id" example:"f9339be2-6b05-4acb-a269-5309c39bae91"`
	CategoryID uuid.NullUUID `json:"category_id" example:"f9339be2-6b05-4acb-a269-5309c39bae92"`
	Name       string        `json:"name" example:"Makan siang"`
	Price      int64         `json:"price" example:"50000"`
	IsIncome   bool          `json:"is_income" example:"false"`
	SpendType  int           `json:"type" example:"2"`
	Date       time.Time     `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
}

type TransferSpend struct {
	PocketIDFrom uuid.UUID `json:"pocket_id_from" example:"f9339be2-6b05-4acb-a269-5309c39bae91"`
	PocketIDTo   uuid.UUID `json:"pocket_id_to" example:"f9339be2-6b05-4acb-a269-5309c39bae92"`
	Price        int64     `json:"price" example:"50000"`
	Date         time.Time `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
}

type UpdateSpend struct {
	ID         uuid.UUID     `json:"-"`
	CategoryID uuid.NullUUID `json:"category_id" example:"f9339be2-6b05-4acb-a269-5309c39bae92"`
	Name       *string       `json:"name" example:"Makan siang"`
	Price      *int64        `json:"price" example:"50000"`
	IsIncome   *bool         `json:"is_income" example:"false"`
	SpendType  *int          `json:"type" example:"2"`
	Date       *time.Time    `json:"date" example:"2022-09-10T17:03:15.091267+08:00"`
}
