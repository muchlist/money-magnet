package mjwt

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
