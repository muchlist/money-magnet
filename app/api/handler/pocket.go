package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptmodel"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptservice"
	"github.com/muchlist/moneymagnet/bussines/sys/data"
	"github.com/muchlist/moneymagnet/bussines/sys/mid"
	"github.com/muchlist/moneymagnet/bussines/sys/validate"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
	"github.com/muchlist/moneymagnet/foundation/web"
)

func NewPocketHandler(log mlogger.Logger, pocketService ptservice.Service) pocketHandler {
	return pocketHandler{
		log:     log,
		service: pocketService,
	}
}

type pocketHandler struct {
	log     mlogger.Logger
	service ptservice.Service
}

func (pt pocketHandler) CreatePocket(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req ptmodel.PocketNew
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		pt.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMessage, err := validate.Struct(req)
	if err != nil {
		pt.log.ErrorT(traceID, "request not valid", err)
		web.ErrorResponse(w, http.StatusBadRequest, errMessage)
		return
	}

	userID, _ := uuid.Parse(claims.Identity)

	result, err := pt.service.CreatePocket(r.Context(), userID, req)
	if err != nil {
		pt.log.ErrorT(traceID, "error create pocket", err)
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

func (pt pocketHandler) RenamePocket(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	type Req struct {
		ID         uint64 `json:"id" validate:"required"`
		PocketName string `json:"pocket_name" validate:"required"`
	}

	var req Req
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		pt.log.ErrorT(traceID, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMessage, err := validate.Struct(req)
	if err != nil {
		pt.log.ErrorT(traceID, "request not valid", err)
		web.ErrorResponse(w, http.StatusBadRequest, errMessage)
		return
	}

	userID, _ := uuid.Parse(claims.Identity)

	result, err := pt.service.RenamePocket(r.Context(), userID, req.ID, req.PocketName)
	if err != nil {
		pt.log.ErrorT(traceID, "error rename pocket", err)
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

// GetByID...
func (pt pocketHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url path
	pocketID, err := web.ReadIDParam(r)
	if err != nil {
		pt.log.ErrorT(traceID, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := pt.service.GetDetail(r.Context(), claims.Identity, pocketID)
	if err != nil {
		pt.log.ErrorT(traceID, "error get pocket by id", err)
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

func (pt pocketHandler) FindUserPocket(w http.ResponseWriter, r *http.Request) {
	traceID := web.ReadTraceID(r.Context())
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url query
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := pt.service.FindAllPocket(r.Context(), claims.Identity, data.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(traceID, "error find pocket", err)
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
