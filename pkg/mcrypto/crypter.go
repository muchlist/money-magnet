package mcrypto

type Crypter interface {
	GenerateHash(password string) ([]byte, error)
	IsPWAndHashPWMatch(password []byte, hashPass []byte) bool
}
