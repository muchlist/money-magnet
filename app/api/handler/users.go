package handler

import (
	"errors"
	"github.com/muchlist/moneymagnet/bussines/user/usermodel"
	"github.com/muchlist/moneymagnet/bussines/user/userservice"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/validate"
	"net/http"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
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
		usr.log.ErrorT(traceID, "request not valid", err)
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
		usr.log.ErrorT(traceID, "request not valid", err)
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

// EditSelfUser
// TODO : remove edit roles and pocket roles by user input
func (usr userHandler) EditSelfUser(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req usermodel.UserUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Not have validate, because no field required

	req.ID = claims.Identity

	result, err := usr.service.FetchUser(r.Context(), req)
	if err != nil {
		usr.log.ErrorT(traceID, "error edit user", err)
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

func (usr userHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	id, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.ErrorT(traceID, "error edit user", err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req usermodel.UserUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Not have validate, because no field required
	req.ID = id

	result, err := usr.service.FetchUser(r.Context(), req)
	if err != nil {
		usr.log.ErrorT(traceID, "error edit user", err)
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

func (usr userHandler) UpdateFCM(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	fcm, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.ErrorT(traceID, "fcm required", err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = usr.service.UpdateFCM(r.Context(), claims.Identity, fcm)
	if err != nil {
		usr.log.ErrorT(traceID, "error update fcm", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": "success",
	}
	err = web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

func (usr userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	userIDToDelete, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.ErrorT(traceID, err.Error(), err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = usr.service.Delete(r.Context(), userIDToDelete, claims.Identity)
	if err != nil {
		usr.log.ErrorT(traceID, "error delete user", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": "success",
	}
	err = web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

func (usr userHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())

	type refresh struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	var req refresh
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	errMessage, err := validate.Struct(req)
	if err != nil {
		usr.log.ErrorT(traceID, "request not valid", err)
		web.ErrorResponse(w, http.StatusBadRequest, errMessage)
		return
	}

	result, err := usr.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		usr.log.ErrorT(traceID, "error refresh token", err)
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

// ==================================================GET
func (usr userHandler) Profile(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	result, err := usr.service.GetProfile(r.Context(), claims.Identity)
	if err != nil {
		usr.log.ErrorT(traceID, "error get profile", err)
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

// GetByID...
func (usr userHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())

	// extract url path
	userID, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.ErrorT(traceID, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := usr.service.GetProfile(r.Context(), userID)
	if err != nil {
		usr.log.ErrorT(traceID, "error get user by id", err)
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

func (usr userHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())

	// extract url query
	name := web.ReadString(r.URL.Query(), "name", "")
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := usr.service.FindUserByName(r.Context(), name, data.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		usr.log.ErrorT(traceID, "error get profile", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"metadata": metadata,
		"data":     result,
	}
	err = web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

// =================================================FUNC

func parseError(err error) (int, string) {
	switch err := err.(type) {
	case errr.StatusCodeError:
		return err.StatusCode, err.Error()
	default:
		if errors.Is(err, db.ErrDBDuplicatedEntry) ||
			errors.Is(err, db.ErrDBNotFound) ||
			errors.Is(err, db.ErrDBParentNotFound) ||
			errors.Is(err, userservice.ErrInvalidID) ||
			errors.Is(err, db.ErrDBSortFilter) {
			return http.StatusBadRequest, err.Error()
		}

		if errors.Is(err, mjwt.ErrInvalidToken) {
			return http.StatusUnauthorized, err.Error()
		}

		if errors.Is(err, userservice.ErrInvalidEmailOrPass) {
			return http.StatusBadRequest, "invalid email or password"
		}

		return http.StatusInternalServerError, err.Error()
	}
}
