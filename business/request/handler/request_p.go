package handler

import (
	"net/http"
	"strings"

	"github.com/muchlist/moneymagnet/business/request/model"
	"github.com/muchlist/moneymagnet/business/request/service"
	"github.com/muchlist/moneymagnet/business/zhelper"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
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

// @Summary      Create Join Request
// @Description  Create Join Request
// @Tags         Join
// @Accept       json
// @Produce      json
// @Param		 Body body model.NewRequestPocket true "Request Body"
// @Success      200  {object}  misc.ResponseSuccess{data=model.RequestPocket}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /request [post]
func (pt requestHandler) CreateRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-CreateRequest")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	var req model.NewRequestPocket
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

	result, err := pt.service.CreateRequest(ctx, claims, req.PocketID)
	if err != nil {
		pt.log.ErrorT(ctx, "error create request", err)
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

// @Summary      Action to Join Request
// @Description  Action to Join Request
// @Tags         Join
// @Accept       json
// @Produce      json
// @Param		 request_id path string true "request_id"
// @Param		 approve query bool false "approve"
// @Success      200  {object}  misc.ResponseMessage
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /request/{request_id}/action [post]
func (pt requestHandler) ApproveOrRejectRequest(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-ApproveOrRejectRequest")
	defer span.End()

	claims, err := mid.GetClaims(ctx)
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	id, err := web.ReadIDParam(r)
	if err != nil {
		pt.log.WarnT(ctx, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	isApprovedStr := web.ReadString(r.URL.Query(), "approve", "")
	isApprovedBool := false
	if isApprovedStr == "" {
		pt.log.WarnT(ctx, "bad request", err)
		web.ErrorResponse(w, http.StatusBadRequest, "?approve=<must be bool>")
		return
	}
	if strings.ToLower(isApprovedStr) == "true" {
		isApprovedBool = true
	}

	err = pt.service.ApproveRequest(ctx, claims, isApprovedBool, id)
	if err != nil {
		pt.log.ErrorT(ctx, "error change status request", err)
		statusCode, msg := zhelper.ParseError(err)
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

// @Summary      Get Request IN
// @Description  Get request you can approve
// @Tags         Join
// @Accept       json
// @Produce      json
// @Success      200  {object}  misc.ResponseSuccessList{data=[]model.RequestPocket}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /request/in [get]
func (pt requestHandler) FindRequestByApprover(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-FindRequestByApprover")
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

	result, metadata, err := pt.service.FindAllByApprover(ctx, claims, paging.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(ctx, "error find request", err)
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

// @Summary      Get Request OUT
// @Description  Get request created by you
// @Tags         Join
// @Accept       json
// @Produce      json
// @Success      200  {object}  misc.ResponseSuccessList{data=[]model.RequestPocket}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /request/out [get]
func (pt requestHandler) FindByRequester(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-FindByRequester")
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

	result, metadata, err := pt.service.FindAllByRequester(ctx, claims, paging.Filters{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
	})
	if err != nil {
		pt.log.ErrorT(ctx, "error find request", err)
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
