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

func (usr userHandler) Register(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	var req usermodel.UserRegisterReq
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMessage, err := validate.Struct(req)
	if err != nil {
		usr.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, errMessage)
		return
	}

	message, err := usr.service.InsertUser(r.Context(), req)
	if err != nil {
		usr.log.ErrorT(traceID, "error insert user", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": message,
	}
	err = web.WriteJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

func (usr userHandler) Login(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	var req usermodel.UserLoginReq
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMessage, err := validate.Struct(req)
	if err != nil {
		usr.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, errMessage)
		return
	}

	result, err := usr.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		usr.log.ErrorT(traceID, "error login", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": result,
	}
	err = web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

func parseError(err error) (int, string) {
	if errors.Is(err, db.ErrDBDuplicatedEntry) ||
		errors.Is(err, db.ErrDBNotFound) ||
		errors.Is(err, db.ErrDBParentNotFound) {
		return http.StatusBadRequest, err.Error()
	}

	if errors.Is(err, userservice.ErrInvalidEmailOrPass) {
		return http.StatusBadRequest, "invalid email or password"
	}

	return http.StatusInternalServerError, err.Error()
}
