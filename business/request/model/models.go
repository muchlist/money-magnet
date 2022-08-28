package model

import (
	"time"

	"github.com/google/uuid"
)

// RequestPocket unity used for domain and response because not have difference
type RequestPocket struct {
	ID         uint64     `json:"id"`
	Requester  uuid.UUID  `json:"requester"`
	Pocket     uint64     `json:"pocket"`
	PocketName string     `json:"pocket_name"`
	Approver   *uuid.UUID `json:"approver"`
	IsApproved *bool      `json:"is_approved"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// type RequestPocketNew struct {
// 	Pocket uint64
// }
