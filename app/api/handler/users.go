package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/business/user/service"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/validate"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func NewUserHandler(log mlogger.Logger,
	validator validate.Validator,
	userService service.Core) userHandler {
	return userHandler{
		log:       log,
		validator: validator,
		service:   userService,
	}
}

type userHandler struct {
	log       mlogger.Logger
	validator validate.Validator
	service   service.Core
}

func (usr userHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req model.UserRegisterReq
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := usr.validator.Struct(req)
	if err != nil {
		usr.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	message, err := usr.service.InsertUser(r.Context(), req)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error insert user", err)
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

	var req model.UserLoginReq
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := usr.validator.Struct(req)
	if err != nil {
		usr.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := usr.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error login", err)
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

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.UserUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Not have validate, because no field required
	req.ID, err = uuid.Parse(claims.Identity)
	if err != nil {
		usr.log.ErrorT(r.Context(), "uuid from claims must be uuid", err)
		web.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := usr.service.FetchUser(r.Context(), req)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error edit user", err)
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

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	id, err := web.ReadUUIDParam(r)
	if err != nil {
		usr.log.WarnT(r.Context(), "error edit user", err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req model.UserUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Not have validate, because no field required
	req.ID = id

	result, err := usr.service.FetchUser(r.Context(), req)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error edit user", err)
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

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	fcm, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.WarnT(r.Context(), "fcm required", err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = usr.service.UpdateFCM(r.Context(), claims.Identity, fcm)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error update fcm", err)
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

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	userIDToDelete, err := web.ReadUUIDParam(r)
	if err != nil {
		usr.log.WarnT(r.Context(), err.Error(), err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	claimsUUID, err := uuid.Parse(claims.Identity)
	if err != nil {
		usr.log.ErrorT(r.Context(), "uuid from claims must be uuid", err)
		web.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = usr.service.Delete(r.Context(), userIDToDelete, claimsUUID)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error delete user", err)
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

	type refresh struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	var req refresh
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	errMap, err := usr.validator.Struct(req)
	if err != nil {
		usr.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := usr.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error refresh token", err)
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

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	result, err := usr.service.GetProfile(r.Context(), claims.Identity)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error get profile", err)
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

	// extract url path
	userID, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := usr.service.GetProfile(r.Context(), userID)
	if err != nil {
		usr.log.ErrorT(r.Context(), "error get user by id", err)
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
		usr.log.ErrorT(r.Context(), "error get profile", err)
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
