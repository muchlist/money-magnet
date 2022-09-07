package handler

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/request/service"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/validate"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func NewRequestHandler(log mlogger.Logger,
	validator validate.Validator,
	requestService service.Core) requestHandler {
	return requestHandler{
		log:       log,
		validator: validator,
		service:   requestService,
	}
}

type requestHandler struct {
	log       mlogger.Logger
	validator validate.Validator
	service   service.Core
}

// CreateRequest ...
func (pt requestHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req struct {
		PocketID uuid.UUID `json:"pocket_id" validator:"required,uuid"`
	}
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		pt.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := pt.validator.Struct(req)
	if err != nil {
		pt.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := pt.service.CreateRequest(r.Context(), claims.GetUUID(), req.PocketID)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error create request", err)
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

// ApproveOrRejectRequest ...
func (pt requestHandler) ApproveOrRejectRequest(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	id, err := web.ReadIDParam(r)
	if err != nil {
		pt.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	isApprovedStr := web.ReadString(r.URL.Query(), "approve", "")
	isApprovedBool := false
	if isApprovedStr == "" {
		pt.log.WarnT(r.Context(), "bad request", err)
		web.ErrorResponse(w, http.StatusBadRequest, "?approve=<must be bool>")
		return
	}
	if strings.ToLower(isApprovedStr) == "true" {
		isApprovedBool = true
	}

	err = pt.service.ApproveRequest(r.Context(), claims.GetUUID(), isApprovedBool, id)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error change status request", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": "status request changed",
	}
	err = web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}

// FindRequestByApprover...
func (pt requestHandler) FindRequestByApprover(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url query
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := pt.service.FindAllByApprover(r.Context(), claims.GetUUID(), data.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(r.Context(), "error find request", err)
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

// FindByRequester ...
func (pt requestHandler) FindByRequester(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url query
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := pt.service.FindAllByRequester(r.Context(), claims.GetUUID(), data.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(r.Context(), "error find request", err)
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
