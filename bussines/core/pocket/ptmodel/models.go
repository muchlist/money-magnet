package ptmodel

import (
	"time"

	"github.com/google/uuid"
)

type Pocket struct {
	ID         uint64
	Owner      uuid.UUID
	Editor     []uuid.UUID
	Watcher    []uuid.UUID
	PocketName string
	Level      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int
}

type PocketNew struct {
	Owner      uuid.UUID   `json:"owner"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
	PocketName string      `json:"pocket_name"`
}

type PocketUpdate struct {
	ID         uint64      `json:"-"`
	Owner      *uuid.UUID  `json:"owner"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
	PocketName *string     `json:"pocket_name"`
	Version    *int        `json:"version"`
}
