package service

import (
	"context"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/category/model"
	"github.com/muchlist/moneymagnet/business/category/port"
	pocketPort "github.com/muchlist/moneymagnet/business/pocket/port"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
	"github.com/muchlist/moneymagnet/pkg/slicer"
	"github.com/muchlist/moneymagnet/pkg/xulid"
)

// Core manages the set of APIs for category access.
type Core struct {
	log          mlogger.Logger
	repo         port.CategoryStorer
	pockerReader pocketPort.PocketReader
}

// NewCore constructs a core for category api access.
func NewCore(
	log mlogger.Logger,
	repo port.CategoryStorer,
	pockerReader pocketPort.PocketReader,
) *Core {
	return &Core{
		log:          log,
		repo:         repo,
		pockerReader: pockerReader,
	}
}

func (s *Core) CreateCategory(ctx context.Context, claims mjwt.CustomClaim, req model.NewCategory) (model.CategoryResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "category-service-CreateCategory")
	defer span.End()

	// Get existing Pocket
	pocketExisting, err := s.pockerReader.GetByID(ctx, req.PocketID)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) {
		return model.CategoryResp{}, errr.New("not have access to this pocket", 400)
	}

	timeNow := time.Now()
	cat := model.Category{
		ID:               xulid.Instance().NewULID(),
		PocketID:         req.PocketID,
		CategoryName:     req.CategoryName,
		CategoryIcon:     req.CategoryIcon,
		IsIncome:         req.IsIncome,
		DefaultSpendType: req.DefaultSpendType,
		CreatedAt:        timeNow,
		UpdatedAt:        timeNow,
	}

	if err := s.repo.Insert(ctx, &cat); err != nil {
		return model.CategoryResp{}, fmt.Errorf("insert category to db: %w", err)
	}

	return cat.ToCategoryResp(), nil
}

func (s *Core) EditCategory(ctx context.Context, claims mjwt.CustomClaim, newData model.UpdateCategory) (model.CategoryResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "category-service-EditCategory")
	defer span.End()

	// Get existing Category
	categoryExisting, err := s.repo.GetByID(ctx, newData.ID.String())
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("get category by id: %w", err)
	}

	// Get existing Pocket
	pocketExisting, err := s.pockerReader.GetByID(ctx, categoryExisting.PocketID)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) {
		return model.CategoryResp{}, errr.New("not have access to this pocket", 400)
	}

	// Modify data
	categoryExisting.CategoryName = newData.CategoryName
	categoryExisting.CategoryIcon = newData.CategoryIcon
	categoryExisting.DefaultSpendType = newData.DefaultSpendType

	// Edit
	err = s.repo.Edit(ctx, &categoryExisting)
	if err != nil {
		return model.CategoryResp{}, fmt.Errorf("edit category: %w", err)
	}

	return categoryExisting.ToCategoryResp(), nil
}

// FindAllCategory ...
func (s *Core) FindAllCategory(ctx context.Context, pocketID xulid.ULID, filter paging.Filters) ([]model.CategoryResp, paging.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "category-service-FindAllCategory")
	defer span.End()

	// Get category
	cats, metadata, err := s.repo.Find(ctx, pocketID.String(), filter)
	if err != nil {
		return nil, paging.Metadata{}, fmt.Errorf("find category: %w", err)
	}

	catResults := make([]model.CategoryResp, len(cats))
	for i := range cats {
		catResults[i] = cats[i].ToCategoryResp()
	}

	return catResults, metadata, nil
}

// DeleteCategory ...
func (s *Core) DeleteCategory(ctx context.Context, categoryID xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "category-service-DeleteCategory")
	defer span.End()

	err := s.repo.Delete(ctx, categoryID.String())
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}
