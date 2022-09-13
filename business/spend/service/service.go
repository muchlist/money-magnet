package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/shared"
	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/business/spend/storer"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("data not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for user access.
type Core struct {
	log        mlogger.Logger
	repo       storer.SpendStorer
	pocketRepo storer.PocketStorer
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	repo storer.SpendStorer,
	pocketRepo storer.PocketStorer,
) Core {
	return Core{
		log:        log,
		repo:       repo,
		pocketRepo: pocketRepo,
	}
}

func (s Core) CreateSpend(ctx context.Context, claims mjwt.CustomClaim, req model.NewSpend) (model.SpendResp, error) {

	canEdit, _ := shared.IsCanEditOrWatch(req.PocketID, claims.PocketRoles)
	if !canEdit {
		return model.SpendResp{}, errr.New("user doesn't have access to write this pocket", 400)
	}

	timeNow := time.Now()
	spendID := uuid.New()
	if req.ID.Valid {
		spendID = req.ID.UUID
	}
	spend := model.Spend{
		ID:               spendID,
		UserID:           claims.GetUUID(),
		PocketID:         req.PocketID,
		CategoryID:       req.CategoryID,
		CategoryID2:      req.CategoryID2,
		Name:             req.Name,
		Price:            req.Price,
		BalanceSnapshoot: 0,
		IsIncome:         req.IsIncome,
		SpendType:        req.SpendType,
		Date:             req.Date,
		CreatedAt:        timeNow,
		UpdatedAt:        timeNow,
		Version:          1,
	}

	err := s.repo.Insert(ctx, &spend)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("insert spend to db: %w", err)
	}

	newBalance, err := s.pocketRepo.UpdateBalance(ctx, spend.PocketID, spend.Price, false)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("fail to change balance: %w", err)
	}
	spend.BalanceSnapshoot = newBalance

	return spend.ToResp(), nil
}

func (s Core) UpdatePartialSpend(ctx context.Context, claims mjwt.CustomClaim, req model.UpdateSpend) (model.SpendResp, error) {

	// Get existing Spend
	spendExisting, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get spend by id: %w", err)
	}

	// validate id creator
	if spendExisting.UserID != claims.GetUUID() {
		return model.SpendResp{}, errr.New("user cannot edit this transaction", 400)
	}

	// validate pocket roles
	canEdit, _ := shared.IsCanEditOrWatch(spendExisting.PocketID, claims.PocketRoles)
	if !canEdit {
		return model.SpendResp{}, errr.New("user doesn't have access to write this pocket", 400)
	}

	// Modify data
	if req.CategoryID.Valid {
		spendExisting.CategoryID = req.CategoryID
	}
	if req.CategoryID2.Valid {
		spendExisting.CategoryID2 = req.CategoryID2
	}
	if req.Name != nil {
		spendExisting.Name = *req.Name
	}
	if req.IsIncome != nil {
		spendExisting.IsIncome = *req.IsIncome
	}
	if req.SpendType != nil {
		spendExisting.SpendType = *req.SpendType
	}
	if req.Date != nil {
		spendExisting.Date = *req.Date
	}

	// more logic if price change
	var diff int64
	if req.Price != nil {
		diff = *req.Price - spendExisting.Price
		spendExisting.Price = *req.Price
	}

	// Edit
	err = s.repo.Edit(ctx, &spendExisting)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("edit spend: %w", err)
	}

	if diff != 0 {
		newBalance, err := s.pocketRepo.UpdateBalance(ctx, spendExisting.PocketID, diff, false)
		if err != nil {
			return model.SpendResp{}, fmt.Errorf("fail to change balance: %w", err)
		}
		spendExisting.BalanceSnapshoot = newBalance
	}

	return spendExisting.ToResp(), nil
}

// GetDetail ...
func (s Core) GetDetail(ctx context.Context, spendID uuid.UUID) (model.SpendResp, error) {
	// Get existing Spend
	spendDetail, err := s.repo.GetByID(ctx, spendID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get spend detail by id: %w", err)
	}

	return spendDetail.ToResp(), nil
}

// FindAllSpend ...
func (s Core) FindAllSpend(ctx context.Context, claims mjwt.CustomClaim, spendFilter model.SpendFilter, filter data.Filters) ([]model.SpendResp, data.Metadata, error) {

	// if cannot edit and cannot watch, return error
	canEdit, canWatch := shared.IsCanEditOrWatch(spendFilter.PocketID.UUID, claims.PocketRoles)
	if !(canEdit || canWatch) {
		return nil, data.Metadata{}, errr.New("user doesn't have access to read this resource", 400)
	}

	spends, metadata, err := s.repo.Find(ctx, spendFilter, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find spend by pocketID: %w", err)
	}

	spendResult := make([]model.SpendResp, len(spends))
	for i := range spends {
		spendResult[i] = spends[i].ToResp()
	}

	return spendResult, metadata, nil
}

// SyncBalance ...
func (s Core) SyncBalance(ctx context.Context, claims mjwt.CustomClaim, pocketID uuid.UUID) (int64, error) {

	// if cannot edit and cannot watch, return error
	canEdit, canWatch := shared.IsCanEditOrWatch(pocketID, claims.PocketRoles)
	if !(canEdit || canWatch) {
		return 0, errr.New("user doesn't have access to read this resource", 400)
	}

	balance, err := s.repo.CountAllPrice(ctx, pocketID)
	if err != nil {
		return 0, fmt.Errorf("aggregate all price on pocket: %w", err)
	}

	newBalance, err := s.pocketRepo.UpdateBalance(ctx, pocketID, balance, true)
	if err != nil {
		return 0, fmt.Errorf("fail update balance: %w", err)
	}

	return newBalance, nil
}
