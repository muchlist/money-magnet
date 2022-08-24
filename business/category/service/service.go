package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/business/category/storer"
	"github.com/muchlist/moneymagnet/pkg/data"
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

func (s Core) CreateCategory(ctx context.Context, owner uuid.UUID, req model.NewCategory) (model.CategoryResp, error) {
	// TODO Validate user can create this category

	timeNow := time.Now()
	cat := model.Category{
		ID:           uuid.New(),
		Pocket:       req.PocketID,
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

func (s Core) EditCategory(ctx context.Context, owner uuid.UUID, newData model.UpdateCategory) (model.CategoryResp, error) {

	// Get existing Category
	CategoryExisting, err := s.repo.GetByID(ctx, newData.ID)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("get category by id: %w", err)
	}

	// TODO ensure this category can edited by owner

	// Modify data
	CategoryExisting.CategoryName = newData.CategoryName

	// Edit
	s.repo.Edit(ctx, &CategoryExisting)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("edit category: %w", err)
	}

	return CategoryExisting.ToCategoryResp(), nil
}

// FindAllCategory ...
func (s Core) FindAllCategory(ctx context.Context, pocketID uint64, filter data.Filters) ([]model.CategoryResp, data.Metadata, error) {

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
