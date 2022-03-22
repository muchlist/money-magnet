package crypto

type Crypto interface {
	GenerateHash(password string) (string, error)
	IsPWAndHashPWMatch(password string, hashPass string) bool
}
