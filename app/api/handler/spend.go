package handler

import (
	"fmt"
	"net/http"

	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/business/spend/service"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/validate"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func NewSpendHandler(log mlogger.Logger,
	validator validate.Validator,
	spendService service.Core) spendHandler {
	return spendHandler{
		log:       log,
		validator: validator,
		service:   spendService,
	}
}

type spendHandler struct {
	log       mlogger.Logger
	validator validate.Validator
	service   service.Core
}

func (pt spendHandler) CreateSpend(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.NewSpend
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

	result, err := pt.service.CreateSpend(r.Context(), claims, req)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error create spend", err)
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

func (pt spendHandler) EditSpend(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url path
	spendID, err := web.ReadUUIDParam(r)
	if err != nil {
		pt.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req model.UpdateSpend
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		pt.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	req.ID = spendID

	errMap, err := pt.validator.Struct(req)
	if err != nil {
		pt.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := pt.service.UpdatePartialSpend(r.Context(), claims, req)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error update spend", err)
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
func (pt spendHandler) GetByID(w http.ResponseWriter, r *http.Request) {

	// extract url path
	spendID, err := web.ReadUUIDParam(r)
	if err != nil {
		pt.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := pt.service.GetDetail(r.Context(), spendID)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error get spend by id", err)
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

func (pt spendHandler) FindSpend(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	pocketID, err := web.ReadUUIDParam(r)
	if err != nil {
		pt.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// extract url query
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := pt.service.FindAllSpend(r.Context(), claims, pocketID, data.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(r.Context(), "error find spend", err)
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

// SyncBalance...
func (pt spendHandler) SyncBalance(w http.ResponseWriter, r *http.Request) {
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url path
	pocketID, err := web.ReadUUIDParam(r)
	if err != nil {
		pt.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newBalance, err := pt.service.SyncBalance(r.Context(), claims, pocketID)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error sync balance", err)
		statusCode, msg := parseError(err)
		web.ErrorResponse(w, statusCode, msg)
		return
	}
	env := web.Envelope{
		"data": fmt.Sprintf("new balance for pocket_id %s has set to %d", pocketID, newBalance),
	}
	err = web.WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}
}
