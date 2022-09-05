package model

import (
	"time"

	"github.com/google/uuid"
)

// RequestPocket unity used for domain and response because not have difference
type RequestPocket struct {
	ID          uint64     `json:"id"`
	RequesterID uuid.UUID  `json:"requester_id"`
	PocketID    uuid.UUID  `json:"pocket_id"`
	PocketName  string     `json:"pocket_name"`
	ApproverID  *uuid.UUID `json:"approver_id"`
	IsApproved  bool       `json:"is_approved"`
	IsRejected  bool       `json:"is_rejected"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// type RequestPocketNew struct {
// 	Pocket uint64
// }

type FindBy struct {
	PocketIDs   []string
	ApproverID  string
	RequesterID string
}
