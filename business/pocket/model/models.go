package model

import (
	"time"

	"github.com/google/uuid"
)

type Pocket struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	EditorID   []uuid.UUID
	WatcherID  []uuid.UUID
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
		OwnerID:    p.OwnerID,
		EditorID:   p.EditorID,
		WatcherID:  p.WatcherID,
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
	EditorID   []uuid.UUID `json:"editor_id"`
	WatcherID  []uuid.UUID `json:"watcher_id"`
	Icon       int         `json:"icon"`
}

type PocketUpdate struct {
	ID         uuid.UUID   `json:"-"`
	OwnerID    *uuid.UUID  `json:"owner_id"`
	EditorID   []uuid.UUID `json:"editor_id"`
	WatcherID  []uuid.UUID `json:"watcher_id"`
	PocketName *string     `json:"pocket_name"`
	Icon       *int        `json:"icon"`
	Version    *int        `json:"version"`
}

type PocketResp struct {
	ID         uuid.UUID   `json:"id"`
	OwnerID    uuid.UUID   `json:"owner_id"`
	EditorID   []uuid.UUID `json:"editor_id"`
	WatcherID  []uuid.UUID `json:"watcher_id"`
	PocketName string      `json:"pocket_name"`
	Icon       int         `json:"icon"`
	Level      int         `json:"level"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Version    int         `json:"version"`
}
