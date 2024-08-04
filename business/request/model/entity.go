package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

// RequestPocket unity used for domain and response because not have difference
type RequestPocket struct {
	ID          uint64      `json:"id" example:"2001"`
	RequesterID xulid.ULID  `json:"requester_id" example:"01J4EXF94QDMR5XT9KN527XEP6"`
	PocketID    xulid.ULID  `json:"pocket_id" example:"01J4EXF94QDMR5XT9KN527XEP8"`
	PocketName  string      `json:"pocket_name" example:"main pocket"`
	ApproverID  *xulid.ULID `json:"approver_id" example:"01ARZ3NDEKTSV4RRFFQ69G5FXX"`
	IsApproved  bool        `json:"is_approved" example:"false"`
	IsRejected  bool        `json:"is_rejected" example:"false"`
	CreatedAt   time.Time   `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt   time.Time   `json:"updated_at" example:"2022-09-10T17:03:15.091267+08:00"`
}

type NewRequestPocket struct {
	PocketID xulid.ULID `json:"pocket_id" example:"01J4EXF94QDMR5XT9KN527XEP8"`
}

type FindBy struct {
	PocketIDs   []string
	ApproverID  string
	RequesterID string
}
