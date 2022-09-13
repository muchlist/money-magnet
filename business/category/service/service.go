package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/business/category/storer"
	"github.com/muchlist/moneymagnet/business/shared"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
)

// Core manages the set of APIs for category access.
type Core struct {
	log  mlogger.Logger
	repo storer.CategoryStorer
}

// NewCore constructs a core for category api access.
func NewCore(
	log mlogger.Logger,
	repo storer.CategoryStorer,
) Core {
	return Core{
		log:  log,
		repo: repo,
	}
}

func (s Core) CreateCategory(ctx context.Context, claims mjwt.CustomClaim, req model.NewCategory) (model.CategoryResp, error) {
	// if cannot edit return error
	canEdit, _ := shared.IsCanEditOrWatch(req.PocketID, claims.PocketRoles)
	if !canEdit {
		return model.CategoryResp{}, errr.New("user doesn't have access to write this resource", 400)
	}

	timeNow := time.Now()
	cat := model.Category{
		ID:           uuid.New(),
		PocketID:     req.PocketID,
		CategoryName: req.CategoryName,
		IsIncome:     req.IsIncome,
		CreatedAt:    timeNow,
		UpdatedAt:    timeNow,
	}

	if err := s.repo.Insert(ctx, &cat); err != nil {
		return model.CategoryResp{}, fmt.Errorf("insert category to db: %w", err)
	}

	return cat.ToCategoryResp(), nil
}

func (s Core) EditCategory(ctx context.Context, claims mjwt.CustomClaim, newData model.UpdateCategory) (model.CategoryResp, error) {

	// Get existing Category
	categoryExisting, err := s.repo.GetByID(ctx, newData.ID)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("get category by id: %w", err)
	}

	// if cannot edit return error
	canEdit, _ := shared.IsCanEditOrWatch(categoryExisting.PocketID, claims.PocketRoles)
	if !canEdit {
		return model.CategoryResp{}, errr.New("user doesn't have access to write this resource", 400)
	}

	// Modify data
	categoryExisting.CategoryName = newData.CategoryName

	// Edit
	err = s.repo.Edit(ctx, &categoryExisting)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("edit category: %w", err)
	}

	return categoryExisting.ToCategoryResp(), nil
}

// FindAllCategory ...
func (s Core) FindAllCategory(ctx context.Context, pocketID uuid.UUID, filter data.Filters) ([]model.CategoryResp, data.Metadata, error) {

	// Get category
	cats, metadata, err := s.repo.Find(ctx, pocketID, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find category: %w", err)
	}

	catResults := make([]model.CategoryResp, len(cats))
	for i := range cats {
		catResults[i] = cats[i].ToCategoryResp()
	}

	return catResults, metadata, nil
}

// DeleteCategory ...
func (s Core) DeleteCategory(ctx context.Context, categoryID uuid.UUID) error {
	err := s.repo.Delete(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}
