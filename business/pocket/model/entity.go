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
	Users      []PocketUser // not included in get pocket by id
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

// GetOtherUsers returns a unique list of user IDs from
// EditorID, WatcherID, and OwnerID,
// excluding the given userID.
func (p *Pocket) GetOtherUsers(excludeUserID string) []string {
	userSet := ds.NewStringSet()

	if p.OwnerID.String() != excludeUserID {
		userSet.Add(p.OwnerID.String())
	}

	for _, id := range p.EditorID {
		if id != excludeUserID {
			userSet.Add(id)
		}
	}

	for _, id := range p.WatcherID {
		if id != excludeUserID {
			userSet.Add(id)
		}
	}

	return userSet.RevealNotEmptySorted()
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
