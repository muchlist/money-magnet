package handler

import (
	"net/http"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/business/pocket/service"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/validate"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func NewPocketHandler(log mlogger.Logger,
	validator validate.Validator,
	pocketService service.Core) pocketHandler {
	return pocketHandler{
		log:       log,
		validator: validator,
		service:   pocketService,
	}
}

type pocketHandler struct {
	log       mlogger.Logger
	validator validate.Validator
	service   service.Core
}

// @Summary      Create Pocket
// @Description  Create Pocket
// @Tags         Pocket
// @Accept       json
// @Produce      json
// @Param		 Body body model.NewPocket true "Request Body"
// @Success      200  {object}  misc.ResponseSuccess{data=model.PocketResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /pockets [post]
func (pt pocketHandler) CreatePocket(w http.ResponseWriter, r *http.Request) {
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.NewPocket
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

	result, err := pt.service.CreatePocket(r.Context(), claims, req)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error create pocket", err)
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

// @Summary      Update Pocket
// @Description  Update Pocket
// @Tags         Pocket
// @Accept       json
// @Produce      json
// @Param		 pocket_id path string true "pocket_id"
// @Param		 Body body model.PocketUpdate true "Request Body"
// @Success      200  {object}  misc.ResponseSuccess{data=model.PocketResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /pockets/{pocket_id} [patch]
func (pt pocketHandler) UpdatePocket(w http.ResponseWriter, r *http.Request) {

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

	var req model.PocketUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		pt.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = pocketID

	errMap, err := pt.validator.Struct(req)
	if err != nil {
		pt.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := pt.service.UpdatePocket(r.Context(), claims, req)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error update pocket", err)
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

// @Summary      Get Pocket Detail
// @Description  Get Pocket Detail by ID
// @Tags         Pocket
// @Accept       json
// @Produce      json
// @Param 		 pocket_id path string true "pocket_id"
// @Success      200  {object}  misc.ResponseSuccessList{data=model.PocketResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /pockets/{pocket_id} [get]
func (pt pocketHandler) GetByID(w http.ResponseWriter, r *http.Request) {

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

	result, err := pt.service.GetDetail(r.Context(), claims, pocketID)
	if err != nil {
		pt.log.ErrorT(r.Context(), "error get pocket by id", err)
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

// @Summary      Find Pocket
// @Description  Find pocket
// @Tags         Pocket
// @Accept       json
// @Produce      json
// @Param 		 page query int false "page"
// @Param 		 page_size query int false "page-size"
// @Param 		 sort query string false "sort"
// @Success      200  {object}  misc.ResponseSuccessList{data=[]model.PocketResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /pockets [get]
func (pt pocketHandler) FindUserPocket(w http.ResponseWriter, r *http.Request) {

	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url query
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := pt.service.FindAllPocket(r.Context(), claims, data.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(r.Context(), "error find pocket", err)
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
