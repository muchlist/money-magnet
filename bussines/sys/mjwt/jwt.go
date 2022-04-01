package mjwt

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/muchlist/moneymagnet/foundation/tools/slicer"
)

const (
	CLAIMS         = "claims"
	identityKey    = "identity"
	nameKey        = "name"
	rolesKey       = "roles"
	pocketRolesKey = "pocket_roles"
	tokenTypeKey   = "type"
	expKey         = "exp"
	freshKey       = "fresh"
)

var (
	ErrCastingClaims = errors.New("fail to type casting")
	ErrInvalidToken  = errors.New("token not valid")
)

func New(secretKey string) *core {
	if secretKey == "" {
		log.Fatal("secret key cannot be empty")
	}
	return &core{
		secretKey: []byte(secretKey),
	}
}

type core struct {
	secretKey []byte
}

// GenerateToken membuat token jwt untuk login header, untuk menguji nilai payloadnya
// dapat menggunakan situs jwt.io
func (j *core) GenerateToken(claims CustomClaim) (string, error) {

	jwtClaim := jwt.MapClaims{
		identityKey:    claims.Identity,
		nameKey:        claims.Name,
		rolesKey:       claims.Roles,
		pocketRolesKey: claims.PocketRoles,
		expKey:         claims.Exp,
		tokenTypeKey:   claims.Type,
		freshKey:       claims.Fresh,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaim)

	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to signed token: %w", err)
	}

	return signedToken, nil
}

// ReadToken membaca inputan token dan menghasilkan pointer struct CustomClaim
// struct CustomClaim digunakan untuk nilai passing antar middleware
func (j *core) ReadToken(token *jwt.Token) (*CustomClaim, error) {
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
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err.Error(), ErrCastingClaims)
	}

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
func (j *core) ValidateToken(tokenString string) (*jwt.Token, error) {
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
