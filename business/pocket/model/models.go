package model

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
	Icon       int
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
		Icon:       p.Icon,
		Level:      p.Level,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		Version:    p.Version,
	}
}

type PocketNew struct {
	PocketName string      `json:"pocket_name" validate:"required"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
	Icon       int         `json:"icon"`
}

type PocketUpdate struct {
	ID         uint64      `json:"-"`
	Owner      *uuid.UUID  `json:"owner"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
	PocketName *string     `json:"pocket_name"`
	Icon       *int        `json:"icon"`
	Version    *int        `json:"version"`
}

type PocketResp struct {
	ID         uint64      `json:"id"`
	Owner      uuid.UUID   `json:"owner"`
	Editor     []uuid.UUID `json:"editor"`
	Watcher    []uuid.UUID `json:"watcher"`
	PocketName string      `json:"pocket_name"`
	Icon       int         `json:"icon"`
	Level      int         `json:"level"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Version    int         `json:"version"`
}
