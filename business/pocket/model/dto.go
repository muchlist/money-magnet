package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type NewPocket struct {
	PocketName string   `json:"pocket_name" validate:"required" example:"dompet utama"`
	Currency   string   `json:"currency" example:"RP."`
	EditorID   []string `json:"editor_id" example:"01J4EXF94QDMR5XT9KN527XEP6"`
	WatcherID  []string `json:"watcher_id" example:"01J4EXF94QDMR5XT9KN527XEP6"`
	Icon       int      `json:"icon" example:"1"`
}

type PocketUpdate struct {
	ID         xulid.ULID `json:"-"`
	PocketName *string    `json:"pocket_name" example:"dompet utama"`
	Currency   *string    `json:"currency" example:"RP."`
	Icon       *int       `json:"icon" example:"1"`
}

type PocketResp struct {
	ID         xulid.ULID   `json:"id" example:"01J4EXF94QDMR5XT9KN527XEP8"`
	OwnerID    xulid.ULID   `json:"owner_id" example:"01J4EXF94QDMR5XT9KN527XEP6"`
	EditorID   []string     `json:"editor_id" example:"01J4EXF94QDMR5XT9KN527XEP6"`
	WatcherID  []string     `json:"watcher_id" example:"01J4EXF94QDMR5XT9KN527XEP6"`
	Users      []PocketUser `json:"users"`
	PocketName string       `json:"pocket_name" example:"dompet utama"`
	Balance    int64        `json:"balance" example:"50000"`
	Currency   string       `json:"currency" example:"RP."`
	Icon       int          `json:"icon" example:"1"`
	Level      int          `json:"level" example:"1"`
	CreatedAt  time.Time    `json:"created_at" example:"2022-09-10T17:03:15.091267+08:00"`
	UpdatedAt  time.Time    `json:"updated_at" example:"2022-09-10T17:03:15.091267+08:00"`
	Version    int          `json:"version" example:"2"`
}
