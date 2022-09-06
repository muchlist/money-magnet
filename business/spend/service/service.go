package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/business/spend/storer"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/errr"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/utils/slicer"
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
	pocketRepo storer.PocketReader
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	repo storer.SpendStorer,
	pocketRepo storer.PocketReader,
) Core {
	return Core{
		log:        log,
		repo:       repo,
		pocketRepo: pocketRepo,
	}
}

func (s Core) CreateSpend(ctx context.Context, userID uuid.UUID, req model.NewSpend) (model.SpendResp, error) {
	// Get Pocket to validate user can be write
	pocket, err := s.pocketRepo.GetByID(ctx, req.PocketID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	isUserIsOwner := pocket.OwnerID == userID
	isUserCanEdit := slicer.In(userID, pocket.EditorID)

	// not owner or editor
	if !(isUserIsOwner || isUserCanEdit) {
		return model.SpendResp{}, errr.New("user doesn't have access to write this pocket", 400)
	}

	timeNow := time.Now()
	spend := model.Spend{
		ID:          uuid.New(),
		UserID:      userID,
		PocketID:    req.PocketID,
		CategoryID:  req.CategoryID,
		CategoryID2: req.CategoryID2,
		Name:        req.Name,
		Price:       req.Price,
		Balance:     0, // TODO : how to get this
		IsIncome:    req.IsIncome,
		SpendType:   req.SpendType,
		Date:        req.Date,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
		Version:     1,
	}

	err = s.repo.Insert(ctx, &spend)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("insert spend to db: %w", err)
	}

	return spend.ToResp(), nil
}

func (s Core) UpdatePartialSpend(ctx context.Context, userID uuid.UUID, req model.UpdateSpend) (model.SpendResp, error) {

	// Get existing Spend
	spendExisting, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get spend by id: %w", err)
	}

	// validate id creator
	if spendExisting.UserID != userID {
		return model.SpendResp{}, errr.New("user cannot edit this transaction", 400)
	}

	// Get Pocket to validate user can be write
	pocket, err := s.pocketRepo.GetByID(ctx, spendExisting.PocketID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	isUserIsOwner := pocket.OwnerID == userID
	isUserCanEdit := slicer.In(userID, pocket.EditorID)

	// not owner or editor
	if !(isUserIsOwner || isUserCanEdit) {
		return model.SpendResp{}, errr.New("user doesn't have access to write this pocket", 400)
	}

	// Modify data
	if req.CategoryID.Valid {
		spendExisting.CategoryID = req.CategoryID.UUID
	}
	if req.CategoryID.Valid {
		spendExisting.CategoryID2 = req.CategoryID2.UUID
	}
	if req.Name != nil {
		spendExisting.Name = *req.Name
	}
	if req.Price != nil {
		spendExisting.Price = *req.Price
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

	// Edit
	s.repo.Edit(ctx, &spendExisting)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("edit spend: %w", err)
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
func (s Core) FindAllSpend(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID, filter data.Filters) ([]model.SpendResp, data.Metadata, error) {
	// Get Pocket to validate user can be write
	pocket, err := s.pocketRepo.GetByID(ctx, pocketID)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// validate
	isUserIsOwner := pocket.OwnerID == userID
	isUserCanEdit := slicer.In(userID, pocket.EditorID)
	isUserCanWatch := slicer.In(userID, pocket.WatcherID)

	// --- not owner or editor
	if !(isUserIsOwner || isUserCanEdit || isUserCanWatch) {
		return nil, data.Metadata{}, errr.New("user doesn't have access to read this pocket", 400)
	}

	spends, metadata, err := s.repo.Find(ctx, pocketID, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find spend by pocketID: %w", err)
	}

	spendResult := make([]model.SpendResp, len(spends))
	for i := range spends {
		spendResult[i] = spends[i].ToResp()
	}

	return spendResult, metadata, nil
}
