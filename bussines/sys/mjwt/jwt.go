package mjwt

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/muchlist/moneymagnet/foundation/tools/slicer"
)

const (
	CLAIMS         = "claims"
	identityKey    = "identity"
	nameKey        = "name"
	rolesKey       = "roles"
	pocketRolesKey = "roles"
	tokenTypeKey   = "type"
	expKey         = "exp"
	freshKey       = "fresh"
)

var (
	ErrCastingClaims = errors.New("fail to type casting")
	ErrInvalidToken  = errors.New("token not valid")
)

func NewJwt(secretKey string) *jwtUtils {
	if secretKey == "" {
		log.Fatal("secret key cannot be empty")
	}
	return &jwtUtils{
		secretKey: []byte(secretKey),
	}
}

type jwtUtils struct {
	secretKey []byte
}

// GenerateToken membuat token jwt untuk login header, untuk menguji nilai payloadnya
// dapat menggunakan situs jwt.io
func (j *jwtUtils) GenerateToken(claims CustomClaim) (string, error) {
	expired := time.Now().Add(time.Minute * claims.ExtraMinute).Unix()

	jwtClaim := jwt.MapClaims{}
	jwtClaim[identityKey] = claims.Identity
	jwtClaim[nameKey] = claims.Name
	jwtClaim[rolesKey] = claims.Roles
	jwtClaim[pocketRolesKey] = claims.PocketRoles
	jwtClaim[expKey] = expired
	jwtClaim[tokenTypeKey] = claims.Type
	jwtClaim[freshKey] = claims.Fresh

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return signedToken, nil
}

// ReadToken membaca inputan token dan menghasilkan pointer struct CustomClaim
// struct CustomClaim digunakan untuk nilai passing antar middleware
func (j *jwtUtils) ReadToken(token *jwt.Token) (*CustomClaim, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, ErrCastingClaims
	}

	identity, ok := claims[identityKey].(string)
	if !ok {
		return nil, ErrCastingClaims
	}
	name, ok := claims[nameKey].(string)
	if !ok {
		return nil, ErrCastingClaims
	}
	exp, ok := claims[expKey].(float64)
	if !ok {
		return nil, ErrCastingClaims
	}
	tokenType, ok := claims[tokenTypeKey].(string)
	if !ok {
		return nil, ErrCastingClaims
	}
	fresh, ok := claims[freshKey].(bool)
	if !ok {
		return nil, ErrCastingClaims
	}
	roles, err := slicer.ToStringSlice(claims[rolesKey])
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err.Error(), ErrCastingClaims)
	}
	pocketRoles, err := slicer.ToStringSlice(claims[pocketRolesKey])

	customClaim := CustomClaim{
		Identity:    identity,
		Name:        name,
		Exp:         int64(exp),
		Roles:       roles,
		PocketRoles: pocketRoles,
		Type:        TokenType(tokenType),
		Fresh:       fresh,
	}

	return &customClaim, nil
}

// ValidateToken memvalidasi apakah token string masukan valid, termasuk memvalidasi apabila field exp nya kadaluarsa
func (j *jwtUtils) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	// Jika expired akan muncul disini asalkan ada claims exp
	if err != nil {
		return nil, ErrInvalidToken
	}

	return token, nil
}
