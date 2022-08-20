package mid

import (
	"context"
	"errors"
	"fmt"
	mjwt2 "github.com/muchlist/moneymagnet/pkg/mjwt"
	"net/http"
	"strings"

	"github.com/muchlist/moneymagnet/pkg/utils/slicer"
	"github.com/muchlist/moneymagnet/pkg/web"
)

const (
	headerKey = "Authorization"
	bearerKey = "Bearer"
)

// index
// 1. claims jwt validator
// 2. context jwt

// ==========================================================================
// claims jwt validator
// ==========================================================================

// RequiredRoles memerlukan salah satu role inputan agar diloloskan ke proses berikutnya
// token tidak perlu fresh
func RequiredRoles(rolesReq ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authStr := r.Header.Get(headerKey)

			// validate Authentication
			claims, err := validateAuthentication(authStr, false)
			if err != nil {
				web.ErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			}

			// validate roles
			if err := validateAuthorizationRole(claims.Roles, rolesReq); err != nil {
				web.ErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := setClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

// RequiredFreshRoles memerlukan salah satu role inputan agar diloloskan ke proses berikutnya
// token harus fresh (tidak hasil dari refresh token)
func RequiredFreshRoles(rolesReq ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authStr := r.Header.Get(headerKey)

			// validate Authentication
			claims, err := validateAuthentication(authStr, true)
			if err != nil {
				web.ErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			}

			// validate roles
			if err := validateAuthorizationRole(claims.Roles, rolesReq); err != nil {
				web.ErrorResponse(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := setClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func validateAuthentication(authHeader string, mustFresh bool) (mjwt2.CustomClaim, error) {
	if !strings.Contains(authHeader, bearerKey) {
		err := errors.New("expected authorization header format: bearer <token>")
		return mjwt2.CustomClaim{}, err
	}
	tokenString := strings.Split(authHeader, " ")
	if len(tokenString) != 2 {
		err := errors.New("expected authorization header format: bearer <token>")
		return mjwt2.CustomClaim{}, err
	}
	token, err := mjwt2.Glob.ValidateToken(tokenString[1])
	if err != nil {
		return mjwt2.CustomClaim{}, err
	}
	claims, err := mjwt2.Glob.ReadToken(token)
	if err != nil {
		return mjwt2.CustomClaim{}, err
	}
	if mustFresh {
		if !claims.Fresh {
			err := errors.New("expected fresh token")
			return mjwt2.CustomClaim{}, err
		}
	}
	return claims, nil
}

// validate roles
func validateAuthorizationRole(rolesHave []string, rolesAllowed []string) error {
	if (len(rolesAllowed) != 0) &&
		!(len(rolesAllowed) == 1 && rolesAllowed[0] == "") {
		for _, roleReq := range rolesAllowed {
			if !slicer.In(roleReq, rolesHave) {
				err := fmt.Errorf("expected role : %s", roleReq)
				return err
			}
		}
	}
	return nil
}

// ==========================================================================
// context jwt
// ==========================================================================

// ctxKey represents the type of value for the context key.
type ctxKey int

// key is used to store/retrieve a Claims value from a context.Context.
const key ctxKey = 1

// setClaims stores the claims in the context.
func setClaims(ctx context.Context, claims mjwt2.CustomClaim) context.Context {
	return context.WithValue(ctx, key, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) (mjwt2.CustomClaim, error) {
	v, ok := ctx.Value(key).(mjwt2.CustomClaim)
	if !ok {
		return mjwt2.CustomClaim{}, errors.New("claim value missing from context")
	}
	return v, nil
}
