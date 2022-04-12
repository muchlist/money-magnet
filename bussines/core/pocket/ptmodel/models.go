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

func (p *Pocket) ToPocketResp() PocketResp {
	return PocketResp{
		ID:         p.ID,
		Owner:      p.Owner,
		Editor:     p.Editor,
		Watcher:    p.Watcher,
		PocketName: p.PocketName,
		Level:      p.Level,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		Version:    p.Version,
	}
}

type PocketNew struct {
	PocketName string      `json:"pocket_name"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
}

type PocketUpdate struct {
	ID         uint64      `json:"-"`
	Owner      *uuid.UUID  `json:"owner"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
	PocketName *string     `json:"pocket_name"`
	Version    *int        `json:"version"`
}

type PocketResp struct {
	ID         uint64      `json:"id"`
	Owner      uuid.UUID   `json:"owner"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
	PocketName string      `json:"pocket_name"`
	Level      int         `json:"level"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Version    int         `json:"version"`
}
