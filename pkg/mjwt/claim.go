package mjwt

import "github.com/google/uuid"

type TokenType string

const (
	Access  TokenType = "Access"
	Refresh TokenType = "Refresh"
)

type CustomClaim struct {
	Identity    string
	Name        string
	Exp         int64
	Type        TokenType
	Fresh       bool
	Roles       []string
	PocketRoles []string
}

func (c CustomClaim) GetUUID() uuid.UUID {
	// Because we 100% sure Identity is uuid, and if not
	// maybe secret key got hacked. it's okay server panic
	return uuid.MustParse(c.Identity)
}
