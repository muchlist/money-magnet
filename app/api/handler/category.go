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

func NewCatHandler(log mlogger.Logger, catService service.Core) catHandler {
	return catHandler{
		log:     log,
		service: catService,
	}
}

type catHandler struct {
	log     mlogger.Logger
	service service.Core
}

func (ch catHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.NewCategory
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		ch.log.WarnT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMessage, err := validate.Struct(req)
	if err != nil {
		ch.log.WarnT(traceID, "request not valid", err)
		web.ErrorResponse(w, http.StatusBadRequest, errMessage)
		return
	}

	userID, _ := uuid.Parse(claims.Identity)

	result, err := ch.service.CreateCategory(r.Context(), userID, req)
	if err != nil {
		ch.log.ErrorT(traceID, "error create pocket", err)
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
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.UpdateCategory
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		ch.log.WarnT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMessage, err := validate.Struct(req)
	if err != nil {
		ch.log.WarnT(traceID, "request not valid", err)
		web.ErrorResponse(w, http.StatusBadRequest, errMessage)
		return
	}

	userID, _ := uuid.Parse(claims.Identity)

	result, err := ch.service.EditCategory(r.Context(), userID, req)
	if err != nil {
		ch.log.ErrorT(traceID, "error rename category", err)
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
	traceID := web.ReadTraceID(r.Context())
	// claims, err := mid.GetClaims(r.Context())
	// if err != nil {
	// 	web.ServerErrorResponse(w, r, err)
	// 	return
	// }

	// extract url query
	pocketID, err := web.ReadUUIDParam(r)
	if err != nil {
		ch.log.WarnT(traceID, err.Error(), err)
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
		ch.log.ErrorT(traceID, "error find categories", err)
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
	traceID := web.ReadTraceID(r.Context())

	// extract url query
	categoryID, err := web.ReadUUIDParam(r)
	if err != nil {
		ch.log.WarnT(traceID, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ch.service.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		ch.log.ErrorT(traceID, "error delete categories", err)
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
