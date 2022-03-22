package crypto

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func NewCrypto() crypto {
	return crypto{}
}

type crypto struct{}

// GenerateHash membuat hashpassword, hash password 1 dengan yang lainnya akan berbeda meskipun
// inputannya sama, sehingga untuk membandingkan hashpassword memerlukan method lain IsPWAndHashPWMatch
func (c crypto) GenerateHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("generate hash error: %w", err)
	}
	return string(passwordHash), nil
}

// IsPWAndHashPWMatch return true jika inputan password dan hashpassword sesuai
func (c crypto) IsPWAndHashPWMatch(password string, hashPass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	return err == nil
}
