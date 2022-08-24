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

// Service manages the set of APIs for category access.
type Service struct {
	log  mlogger.Logger
	repo storer.CategoryStorer
}

// NewService constructs a core for category api access.
func NewService(
	log mlogger.Logger,
	repo storer.CategoryStorer,
) Service {
	return Service{
		log:  log,
		repo: repo,
	}
}

func (s Service) CreateCategory(ctx context.Context, owner uuid.UUID, req model.NewCategory) (model.CategoryResp, error) {
	// TODO Validate user can create this category

	timeNow := time.Now()
	cat := model.Category{
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

func (s Service) RenameCategory(ctx context.Context, owner uuid.UUID, CategoryID uuid.UUID, newName string) (model.CategoryResp, error) {

	// Get existing Category
	CategoryExisting, err := s.repo.GetByID(ctx, CategoryID)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("get category by id: %w", err)
	}

	// TODO ensure this category can edited by owner

	// Modify data
	CategoryExisting.CategoryName = newName

	// Edit
	s.repo.Edit(ctx, &CategoryExisting)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("edit category: %w", err)
	}

	return CategoryExisting.ToCategoryResp(), nil
}

// FindAllCategory ...
func (s Service) FindAllCategory(ctx context.Context, pocketID uint64, filter data.Filters) ([]model.CategoryResp, data.Metadata, error) {

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
