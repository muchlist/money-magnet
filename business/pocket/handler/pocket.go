package handler

import (
	"net/http"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/business/pocket/service"
	"github.com/muchlist/moneymagnet/business/zhelper"
	"github.com/muchlist/moneymagnet/pkg/lrucache"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/validate"

	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/web"
)

func NewPocketHandler(log mlogger.Logger,
	validator validate.Validator,
	cache lrucache.CacheStorer,
	pocketService service.Core) pocketHandler {
	return pocketHandler{
		log:       log,
		validator: validator,
		cache:     cache,
		service:   pocketService,
	}
}

type pocketHandler struct {
	log       mlogger.Logger
	validator validate.Validator
	cache     lrucache.CacheStorer
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-CreatePocket")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.NewPocket
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		pt.log.WarnT(ctx, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	errMap, err := pt.validator.Struct(req)
	if err != nil {
		pt.log.WarnT(ctx, "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := pt.service.CreatePocket(ctx, claims, req)
	if err != nil {
		pt.log.ErrorT(ctx, "error create pocket", err)
		statusCode, msg := zhelper.ParseError(err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-UpdatePocket")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url path
	pocketID, err := web.ReadULIDParam(r)
	if err != nil {
		pt.log.WarnT(ctx, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req model.PocketUpdate
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		pt.log.WarnT(ctx, "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = pocketID

	errMap, err := pt.validator.Struct(req)
	if err != nil {
		pt.log.WarnT(ctx, "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := pt.service.UpdatePocket(ctx, claims, req)
	if err != nil {
		pt.log.ErrorT(ctx, "error update pocket", err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-GetByID")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url path
	pocketID, err := web.ReadULIDParam(r)
	if err != nil {
		pt.log.WarnT(ctx, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := pt.service.GetDetail(ctx, claims, pocketID)
	if err != nil {
		pt.log.ErrorT(ctx, "error get pocket by id", err)
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
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-GetByID")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	// extract url query
	sort := web.ReadString(r.URL.Query(), "sort", "")
	page := web.ReadInt(r.URL.Query(), "page", 0)
	pageSize := web.ReadInt(r.URL.Query(), "page_size", 0)

	result, metadata, err := pt.service.FindAllPocket(ctx, claims, paging.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(ctx, "error find pocket", err)
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
