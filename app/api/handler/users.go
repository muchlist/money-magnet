package handler

import (
	"errors"
	"net/http"

	"github.com/muchlist/moneymagnet/bussines/core/user/usermodel"
	"github.com/muchlist/moneymagnet/bussines/core/user/userservice"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
	"github.com/muchlist/moneymagnet/bussines/sys/validate"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
	"github.com/muchlist/moneymagnet/foundation/web"
)

func NewUserHandler(log mlogger.Logger, userService userservice.Service) userHandler {
	return userHandler{
		log:     log,
		service: userService,
	}
}

type userHandler struct {
	log     mlogger.Logger
	service userservice.Service
}

func (usr userHandler) Get(w http.ResponseWriter, r *http.Request) {

	traceID := web.ReadTraceID(r.Context())

	user := usermodel.UserReq{
		Name:     "aaa",
		Email:    "aaa@muchlis.dev",
		Password: "123131312",
	}

	errMap, err := validate.Struct(user)
	if err != nil {
		web.ErrorResponse(w, http.StatusBadRequest, errMap)
		return
	}

	message, err := usr.service.InsertUser(r.Context(), usermodel.UserReq{})
	if err != nil {
		usr.log.ErrorT(traceID, "insert user", err)
		if errors.Is(err, db.ErrDBDuplicatedEntry) ||
			errors.Is(err, db.ErrDBNotFound) ||
			errors.Is(err, db.ErrDBParentNotFound) {
			web.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		web.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	env := web.Envelope{
		"data": message,
	}
	err = web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}
