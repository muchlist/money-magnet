package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/muchlist/moneymagnet/business/user/service"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/lrucache"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
)

func parseError(err error) (int, string) {
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

const idempotencyKey = "Idempotency"

func idempotencyInjector(r *http.Request, cc lrucache.CacheStorer, data lrucache.Payload) {
	if r.Header.Get(idempotencyKey) != "" {
		cc.Set(fmt.Sprintf("%s-%s", r.URL.Path, r.Header.Get(idempotencyKey)), data)
	}
}

func idempotencyExtract(r *http.Request, cc lrucache.CacheStorer) (lrucache.Payload, bool) {
	if r.Header.Get(idempotencyKey) == "" {
		return lrucache.Payload{}, false
	}
	return cc.Get(fmt.Sprintf("%s-%s", r.URL.Path, r.Header.Get(idempotencyKey)))
}
