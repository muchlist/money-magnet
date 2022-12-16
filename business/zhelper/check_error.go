package zhelper

import (
	"errors"
	"net/http"

	"github.com/muchlist/moneymagnet/business/user/service"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
)

func ParseError(err error) (int, string) {
	switch err := err.(type) {
	case errr.StatusCodeError:
		return err.StatusCode, err.Error()
	default:
		if errors.Is(err, db.ErrDBDuplicatedEntry) ||
			errors.Is(err, db.ErrDBNotFound) ||
			errors.Is(err, db.ErrDBRelationNotFound) ||
			errors.Is(err, service.ErrInvalidID) ||
			errors.Is(err, db.ErrDBSortFilter) {
			return http.StatusBadRequest, err.Error()
		}

		if errors.Is(err, mjwt.ErrInvalidToken) {
			return http.StatusUnauthorized, err.Error()
		}

		if errors.Is(err, service.ErrInvalidEmailOrPass) {
			return http.StatusBadRequest, "invalid email or password"
		}

		return http.StatusInternalServerError, err.Error()
	}
}
