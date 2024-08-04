package model

import (
	"time"

	"github.com/muchlist/moneymagnet/pkg/ds"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type Pocket struct {
	ID         xulid.ULID
	OwnerID    xulid.ULID
	EditorID   []string
	WatcherID  []string
	Users      []PocketUser
	PocketName string
	Balance    int64
	Currency   string
	Icon       int
	Level      int
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Version    int
}

type PocketUser struct {
	ID   xulid.ULID `json:"id" example:"01J4EXF94QDMR5XT9KN527XEP6"`
	Role string     `json:"role" example:"owner"`
	Name string     `json:"name" example:"muchlis"`
}

func (p *Pocket) Sanitize() {
	editorSet := ds.NewStringSet()
	editorSet.AddAll(p.EditorID)

	watcherSet := ds.NewStringSet()
	watcherSet.AddAll(p.WatcherID)

	p.EditorID = editorSet.RevealSorted()
	p.WatcherID = watcherSet.RevealSorted()
}

func (p *Pocket) ToPocketResp() PocketResp {
	return PocketResp{
		ID:         p.ID,
		OwnerID:    p.OwnerID,
		EditorID:   p.EditorID,
		WatcherID:  p.WatcherID,
		Users:      p.Users,
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
