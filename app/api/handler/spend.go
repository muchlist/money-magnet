package handler

import (
	"fmt"
	"net/http"
	"net/url"

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

// @Summary      Create Spend
// @Description  Create spend
// @Tags         Spend
// @Accept       json
// @Produce      json
// @Param		 Body body model.NewSpend true "Request Body"
// @Success      200  {object}  misc.ResponseSuccess{data=model.SpendResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /spends [post]
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

// @Summary      Update Spend
// @Description  Update spend
// @Tags         Spend
// @Accept       json
// @Produce      json
// @Param		 spend_id path string true "spend_id"
// @Param		 Body body model.UpdateSpend true "Request Body"
// @Success      200  {object}  misc.ResponseSuccess{data=model.SpendResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /spends/{spend_id} [patch]
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

// @Summary      Sync Spend Balance
// @Description  Sync spend to update pocket balance
// @Tags         Spend
// @Accept       json
// @Produce      json
// @Param		 spend_id path string true "spend_id"
// @Success      200  {object}  misc.ResponseMessage
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /spends/sync/{spend_id} [post]
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

// @Summary      Get Spend Detail
// @Description  Get spend detail by ID
// @Tags         Spend
// @Accept       json
// @Produce      json
// @Param 		 spend_id path string true "spend_id"
// @Success      200  {object}  misc.ResponseSuccessList{data=model.SpendResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /spends/{spend_id} [get]
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

func extractSpendFIlter(values url.Values) model.SpendFilter {
	rawFilter := model.SpendFilterRaw{
		User:      values.Get("user"),
		Category:  values.Get("category"),
		IsIncome:  values.Get("is_income"),
		Type:      values.Get("type"),
		DateStart: values.Get("date_start"),
		DateEnd:   values.Get("date_end"),
	}
	return rawFilter.ToModel()
}

// @Summary      Find Spend
// @Description  Find spend
// @Tags         Spend
// @Accept       json
// @Produce      json
// @Param 		 page query int false "page"
// @Param 		 page_size query int false "page-size"
// @Param 		 sort query string false "sort"
// @Param 		 user query string false "user"
// @Param 		 category query string false "category"
// @Param 		 is_income query bool false "is_income"
// @Param 		 type query string false "type"
// @Param 		 date_start query int false "date_start"
// @Param 		 date_end query int false "date_end"
// @Success      200  {object}  misc.ResponseSuccessList{data=[]model.SpendResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /spends [get]
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

	filter := extractSpendFIlter(r.URL.Query())
	filter.PocketID.UUID = pocketID

	result, metadata, err := pt.service.FindAllSpend(r.Context(), claims, filter, data.Filters{
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
