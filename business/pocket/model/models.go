package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/pkg/utils/ds"
)

type Pocket struct {
	ID         uuid.UUID
	OwnerID    uuid.UUID
	EditorID   []uuid.UUID
	WatcherID  []uuid.UUID
	PocketName string
	Balance    int64
	Currency   string
	Icon       int
	Level      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int
}

func (p *Pocket) Sanitize() {
	editorSet := ds.NewUUIDSet()
	editorSet.AddAll(p.EditorID)

	watcherSet := ds.NewUUIDSet()
	watcherSet.AddAll(p.WatcherID)

	p.EditorID = editorSet.Reveal()
	p.WatcherID = watcherSet.Reveal()
}

func (p *Pocket) ToPocketResp() PocketResp {
	return PocketResp{
		ID:         p.ID,
		OwnerID:    p.OwnerID,
		EditorID:   p.EditorID,
		WatcherID:  p.WatcherID,
		PocketName: p.PocketName,
		Balance:    p.Balance,
		Currency:   p.Currency,
		Icon:       p.Icon,
		Level:      p.Level,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		Version:    p.Version,
	}
}

type NewPocket struct {
	PocketName string      `json:"pocket_name" validate:"required" example:"dompet utama"`
	Currency   string      `json:"currency" example:"RP."`
	EditorID   []uuid.UUID `json:"editor_id" example:"ba22d3c6-2cdd-40b4-a2aa-d68da8c88502"`
	WatcherID  []uuid.UUID `json:"watcher_id" example:"ba22d3c6-2cdd-40b4-a2aa-d68da8c88502"`
	Icon       int         `json:"icon" example:"1"`
}

type PocketUpdate struct {
	ID         uuid.UUID `json:"-"`
	PocketName *string   `json:"pocket_name" example:"dompet utama"`
	Currency   *string   `json:"currency" example:"RP."`
	Icon       *int      `json:"icon" example:"1"`
}

type PocketResp struct {
	ID         uuid.UUID   `json:"id" example:"968d4dfe-041a-4721-bd8a-4e60c507c671"`
	OwnerID    uuid.UUID   `json:"owner_id" example:"ba22d3c6-2cdd-40b4-a2aa-d68da8c88502"`
	EditorID   []uuid.UUID `json:"editor_id" example:"ba22d3c6-2cdd-40b4-a2aa-d68da8c88502"`
	WatcherID  []uuid.UUID `json:"watcher_id" example:"ba22d3c6-2cdd-40b4-a2aa-d68da8c88502"`
	PocketName string      `json:"pocket_name" example:"dompet utama"`
	Balance    int64       `json:"balance" example:"50000"`
	Currency   string      `json:"currency" example:"RP."`
	Icon       int         `json:"icon" example:"1"`
	Level      int         `json:"level" example:"1"`
	CreatedAt  time.Time   `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt  time.Time   `json:"updated_at" example:"2022-09-10T17:03:15.091267+08:00"`
	Version    int         `json:"version" example:"2"`
}
