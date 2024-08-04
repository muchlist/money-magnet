package mjwt

import (
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

type TokenType string

const (
	Access  TokenType = "Access"
	Refresh TokenType = "Refresh"
)

type CustomClaim struct {
	Identity string
	Name     string
	Exp      int64
	Type     TokenType
	Fresh    bool
	Roles    []string
}

func (c CustomClaim) GetULID() xulid.ULID {
	// Because we 100% sure Identity is ulid, and if not
	// maybe secret key got hacked. it's okay server panic
	return xulid.MustParse(c.Identity)
}
