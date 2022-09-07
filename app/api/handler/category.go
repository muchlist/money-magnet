package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/business/category/service"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/validate"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func NewCatHandler(log mlogger.Logger,
	validator validate.Validator,
	catService service.Core) catHandler {
	return catHandler{
		log:       log,
		validator: validator,
		service:   catService,
	}
}

type catHandler struct {
	log       mlogger.Logger
	validator validate.Validator
	service   service.Core
}

func (ch catHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.NewCategory
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		ch.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := ch.validator.Struct(req)
	if err != nil {
		ch.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	userID, _ := uuid.Parse(claims.Identity)

	result, err := ch.service.CreateCategory(r.Context(), userID, req)
	if err != nil {
		ch.log.ErrorT(r.Context(), "error create pocket", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": result,
	}
	err = web.WriteJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

func (ch catHandler) EditCategory(w http.ResponseWriter, r *http.Request) {
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.UpdateCategory
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		ch.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := ch.validator.Struct(req)
	if err != nil {
		ch.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	userID, _ := uuid.Parse(claims.Identity)

	result, err := ch.service.EditCategory(r.Context(), userID, req)
	if err != nil {
		ch.log.ErrorT(r.Context(), "error rename category", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": result,
	}
	err = web.WriteJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

func (ch catHandler) FindPocketCategory(w http.ResponseWriter, r *http.Request) {
	// extract url query
	pocketID, err := web.ReadUUIDParam(r)
	if err != nil {
		ch.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := ch.service.FindAllCategory(r.Context(), pocketID, data.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		ch.log.ErrorT(r.Context(), "error find categories", err)
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

func (ch catHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// extract url query
	categoryID, err := web.ReadUUIDParam(r)
	if err != nil {
		ch.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ch.service.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		ch.log.ErrorT(r.Context(), "error delete categories", err)
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
