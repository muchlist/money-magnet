package mcrypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func New() core {
	return core{}
}

type core struct{}

// GenerateHash membuat hashpassword, hash password 1 dengan yang lainnya akan berbeda meskipun
// inputannya sama, sehingga untuk membandingkan hashpassword memerlukan method lain IsPWAndHashPWMatch
func (c core) GenerateHash(password string) ([]byte, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return nil, fmt.Errorf("generate hash error: %w", err)
	}
	return passwordHash, nil
}

// IsPWAndHashPWMatch return true jika inputan password dan hashpassword sesuai
func (c core) IsPWAndHashPWMatch(password []byte, hashPass []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashPass, password)
	return err == nil
}
