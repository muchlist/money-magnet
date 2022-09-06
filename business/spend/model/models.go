package model

import (
	"time"

	"github.com/google/uuid"
)

type Spend struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	UserName      string // Join
	PocketID      uuid.UUID
	PocketName    string // Join
	CategoryID    uuid.UUID
	CategoryName  string // Join
	CategoryID2   uuid.UUID
	CategoryName2 string // Join
	Name          string
	Price         int64
	Balance       int64
	IsIncome      bool
	SpendType     int
	Date          time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Version       int
}

func (s *Spend) ToResp() SpendResp {
	return SpendResp{
		ID:            s.ID,
		UserID:        s.UserID,
		UserName:      s.UserName,
		PocketID:      s.PocketID,
		PocketName:    s.PocketName,
		CategoryID:    s.CategoryID,
		CategoryName:  s.CategoryName,
		CategoryID2:   s.CategoryID2,
		CategoryName2: s.CategoryName2,
		Name:          s.Name,
		Price:         s.Price,
		Balance:       s.Balance,
		IsIncome:      s.IsIncome,
		SpendType:     s.SpendType,
		Date:          s.Date,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
		Version:       s.Version,
	}
}

type SpendResp struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	UserName      string    `json:"user_name"`
	PocketID      uuid.UUID `json:"pocket_id"`
	PocketName    string    `json:"pocket_name"`
	CategoryID    uuid.UUID `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	CategoryID2   uuid.UUID `json:"category_id_2"`
	CategoryName2 string    `json:"category_name_2"`
	Name          string    `json:"name"`
	Price         int64     `json:"price"`
	Balance       int64     `json:"balance"`
	IsIncome      bool      `json:"is_income"`
	SpendType     int       `json:"type"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Version       int       `json:"version"`
}

type NewSpend struct {
	PocketID    uuid.UUID `json:"pocket_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	CategoryID2 uuid.UUID `json:"category_id_2"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"`
	IsIncome    bool      `json:"is_income"`
	SpendType   int       `json:"type"`
	Date        time.Time `json:"date"`
}

type UpdateSpend struct {
	ID          uuid.UUID     `json:"-"`
	CategoryID2 uuid.NullUUID `json:"category_id_2"`
	CategoryID  uuid.NullUUID `json:"category_id"`
	Name        *string       `json:"name"`
	Price       *int64        `json:"price"`
	IsIncome    *bool         `json:"is_income"`
	SpendType   *int          `json:"type"`
	Date        *time.Time    `json:"date"`
}
