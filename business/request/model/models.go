package model

import (
	"time"

	"github.com/google/uuid"
)

// RequestPocket unity used for domain and response because not have difference
type RequestPocket struct {
	ID          uint64     `json:"id" example:"2001"`
	RequesterID uuid.UUID  `json:"requester_id" example:"ba22d3c6-2cdd-40b4-a2aa-d68da8c88502"`
	PocketID    uuid.UUID  `json:"pocket_id" example:"968d4dfe-041a-4721-bd8a-4e60c507c671"`
	PocketName  string     `json:"pocket_name" example:"main pocket"`
	ApproverID  *uuid.UUID `json:"approver_id" example:"ba22d3c6-2cdd-40b4-a2aa-d68da8c88501"`
	IsApproved  bool       `json:"is_approved" example:"false"`
	IsRejected  bool       `json:"is_rejected" example:"false"`
	CreatedAt   time.Time  `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2022-09-10T17:03:15.091267+08:00"`
}

type NewRequestPocket struct {
	PocketID uuid.UUID `json:"pocket_id" example:"968d4dfe-041a-4721-bd8a-4e60c507c671"`
}

type FindBy struct {
	PocketIDs   []string
	ApproverID  string
	RequesterID string
}
