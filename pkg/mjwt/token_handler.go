package mjwt

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenHandler interface {
	GenerateToken(claims CustomClaim) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	ReadToken(token *jwt.Token) (CustomClaim, error)
}
