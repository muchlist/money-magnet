package handler

import (
	"net/http"

	"github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/business/user/service"
	"github.com/muchlist/moneymagnet/business/zhelper"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/observ/mmetric"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/validate"
	"github.com/muchlist/moneymagnet/pkg/xulid"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func NewUserHandler(log mlogger.Logger,
	validator validate.Validator,
	userService *service.Core) userHandler {
	return userHandler{
		log:       log,
		validator: validator,
		service:   userService,
	}
}

type userHandler struct {
	log       mlogger.Logger
	validator validate.Validator
	service   *service.Core
}

func (usr userHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-Register")
	defer span.End()

	var req model.UserRegisterReq
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(ctx, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := usr.validator.Struct(req)
	if err != nil {
		usr.log.WarnT(ctx, "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	message, err := usr.service.InsertUser(ctx, req)
	if err != nil {
		usr.log.ErrorT(ctx, "error insert user", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-Login")
	defer span.End()

	var req model.UserLoginReq
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(ctx, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := usr.validator.Struct(req)
	if err != nil {
		usr.log.WarnT(ctx, "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := usr.service.Login(ctx, req.Email, req.Password)
	if err != nil {

		// send metric
		mmetric.AddLoginFailedCounter(ctx)

		usr.log.ErrorT(ctx, "error login", err)
		statusCode, msg := zhelper.ParseError(err)
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
func (usr userHandler) EditSelfUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-EditSelfUser")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.UserUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(ctx, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Roles = nil

	// Not have validate, because no field required
	req.ID, err = xulid.Parse(claims.Identity)
	if err != nil {
		usr.log.ErrorT(ctx, "ulid from claims must be valid ulid", err)
		web.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := usr.service.PatchUser(ctx, req)
	if err != nil {
		usr.log.ErrorT(ctx, "error edit user", err)
		statusCode, msg := zhelper.ParseError(err)
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

// EditUser do edit user detail by admin
func (usr userHandler) EditUser(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-EditUser")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	id, err := web.ReadULIDParam(r)
	if err != nil {
		usr.log.WarnT(ctx, "error edit user", err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req model.UserUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(ctx, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Not have validate, because no field required
	req.ID = id

	result, err := usr.service.PatchUser(ctx, req)
	if err != nil {
		usr.log.ErrorT(ctx, "error edit user", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-UpdateFCM")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	fcm, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.WarnT(ctx, "fcm required", err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = usr.service.UpdateFCM(ctx, claims.Identity, fcm)
	if err != nil {
		usr.log.ErrorT(ctx, "error update fcm", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-DeleteUser")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// Get data from url path
	userIDToDelete, err := web.ReadULIDParam(r)
	if err != nil {
		usr.log.WarnT(ctx, err.Error(), err, mlogger.String("identity", claims.Identity))
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	claimsUUID, err := xulid.Parse(claims.Identity)
	if err != nil {
		usr.log.ErrorT(ctx, "uuid from claims must be uuid", err)
		web.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = usr.service.Delete(ctx, userIDToDelete, claimsUUID)
	if err != nil {
		usr.log.ErrorT(ctx, "error delete user", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-RefreshToken")
	defer span.End()

	type refresh struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	var req refresh
	err := web.ReadJSON(w, r, &req)
	if err != nil {
		usr.log.WarnT(ctx, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	errMap, err := usr.validator.Struct(req)
	if err != nil {
		usr.log.WarnT(ctx, "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := usr.service.Refresh(ctx, req.RefreshToken)
	if err != nil {
		usr.log.ErrorT(ctx, "error refresh token", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-Profile")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	result, err := usr.service.GetProfile(ctx, claims.Identity)
	if err != nil {
		usr.log.ErrorT(ctx, "error get profile", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-GetByID")
	defer span.End()

	// extract url path
	userID, err := web.ReadStrIDParam(r)
	if err != nil {
		usr.log.WarnT(ctx, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := usr.service.GetProfile(ctx, userID)
	if err != nil {
		usr.log.ErrorT(ctx, "error get user by id", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-FindByName")
	defer span.End()

	// extract url query
	name := web.ReadString(r.URL.Query(), "name", "")
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := usr.service.FindUserByName(ctx, name, paging.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		usr.log.ErrorT(ctx, "error get profile", err)
		statusCode, msg := zhelper.ParseError(err)
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
