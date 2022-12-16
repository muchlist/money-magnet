package handler

import (
	"net/http"

	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/business/category/service"
	"github.com/muchlist/moneymagnet/business/zhelper"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/mid"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
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

// @Summary      Create Category
// @Description  Create Category for Spend
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param		 Body body model.NewCategory true "Request Body"
// @Success      200  {object}  misc.ResponseSuccess{data=model.CategoryResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /categories [post]
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

	result, err := ch.service.CreateCategory(r.Context(), claims, req)
	if err != nil {
		ch.log.ErrorT(r.Context(), "error create pocket", err)
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

// @Summary      Edit Category
// @Description  Edit category name
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param		 category_id path string true "category_id"
// @Param		 Body body model.UpdateCategory true "Request Body"
// @Success      200  {object}  misc.ResponseSuccess{data=model.CategoryResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /categories/{category_id} [put]
func (ch catHandler) EditCategory(w http.ResponseWriter, r *http.Request) {
	claims, err := mid.GetClaims(r.Context())
	if err != nil {
		web.ServerErrorResponse(w, r, err)
		return
	}

	categoryID, err := web.ReadUUIDParam(r)
	if err != nil {
		ch.log.WarnT(r.Context(), err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var req model.UpdateCategory
	err = web.ReadJSON(w, r, &req)
	if err != nil {
		ch.log.WarnT(r.Context(), "bad json", err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	req.ID = categoryID

	errMap, err := ch.validator.Struct(req)
	if err != nil {
		ch.log.WarnT(r.Context(), "request not valid", err)
		web.ErrorPayloadResponse(w, err.Error(), errMap)
		return
	}

	result, err := ch.service.EditCategory(r.Context(), claims, req)
	if err != nil {
		ch.log.ErrorT(r.Context(), "error rename category", err)
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

// @Summary      Find Category
// @Description  Find category for 1 pocket
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param 		 pocket_id path string true "pocket_id"
// @Success      200  {object}  misc.ResponseSuccessList{data=[]model.CategoryResp}
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /categories/from-pocket/{pocket_id} [get]
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

// @Summary      Delete Category
// @Description  Delete category by id
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param 		 category_id path string true "category_id"
// @Success      200  {object}  misc.ResponseMessage
// @Failure      400  {object}  misc.ResponseErr
// @Failure      500  {object}  misc.Response500Err
// @Router       /categories/{category_id} [delete]
func (ch catHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := observ.GetTracer().Start(r.Context(), "handler-DeleteCategory")
	defer span.End()

	// extract url query
	categoryID, err := web.ReadUUIDParam(r)
	if err != nil {
		ch.log.WarnT(ctx, err.Error(), err)
		web.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ch.service.DeleteCategory(ctx, categoryID)
	if err != nil {
		ch.log.ErrorT(ctx, "error delete categories", err)
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
