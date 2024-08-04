package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/business/spend/port"
	"github.com/muchlist/moneymagnet/constant"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/slicer"
	"github.com/muchlist/moneymagnet/pkg/xulid"

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
	repo       port.SpendStorer
	pocketRepo port.PocketStorer
	txManager  port.Transactor
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	repo port.SpendStorer,
	pocketRepo port.PocketStorer,
	txManager port.Transactor,
) Core {
	return Core{
		log:        log,
		repo:       repo,
		pocketRepo: pocketRepo,
		txManager:  txManager,
	}
}

func (s Core) CreateSpend(ctx context.Context, claims mjwt.CustomClaim, req model.NewSpend) (model.SpendResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-CreateSpend")
	defer span.End()

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, req.PocketID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) {
		return model.SpendResp{}, errr.New("not have access to this pocket", 400)
	}

	timeNow := time.Now()
	spendID := xulid.Instance().NewULID()
	if req.ID.Valid {
		spendID = req.ID.ULID
	}
	spend := model.Spend{
		ID:               spendID,
		UserID:           claims.GetULID(),
		PocketID:         req.PocketID,
		CategoryID:       req.CategoryID,
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

	transErr := s.txManager.WithAtomic(ctx, func(ctx context.Context) error {
		err = s.repo.Insert(ctx, &spend)
		if err != nil {
			return fmt.Errorf("insert spend to db: %w", err)
		}

		newBalance, err := s.pocketRepo.UpdateBalance(ctx, spend.PocketID, spend.Price, false)
		if err != nil {
			return fmt.Errorf("fail to change balance: %w", err)
		}
		spend.BalanceSnapshoot = newBalance

		return nil
	})

	if transErr != nil {
		return model.SpendResp{}, transErr
	}

	return spend.ToResp(), nil
}

func (s Core) TransferToPocketAsSpend(ctx context.Context, claims mjwt.CustomClaim, req model.TransferSpend) error {
	ctx, span := observ.GetTracer().Start(ctx, "service-TransferToPocketAsSpend")
	defer span.End()

	// Get existing Pocket (from)
	fromPocket, err := s.pocketRepo.GetByID(ctx, req.PocketIDFrom)
	if err != nil {
		return fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), fromPocket.EditorID) {
		return errr.New("not have access to this pocket", 400)
	}

	// Get existing Pocket (to)
	toPocket, err := s.pocketRepo.GetByID(ctx, req.PocketIDTo)
	if err != nil {
		return fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), toPocket.EditorID) {
		return errr.New("not have access to this pocket", 400)
	}

	// check balance and price value
	if req.Price <= 0 {
		return errr.New("the transfer value must be more than zero", 400)
	}
	if fromPocket.Balance < req.Price {
		return errr.New("balance must be more than the transfer value", 400)
	}

	transErr := s.txManager.WithAtomic(ctx, func(ctx context.Context) error {

		timeNow := time.Now()

		// spend for pocket-from
		spendID := xulid.Instance().NewULID()
		spend := model.Spend{
			ID:       spendID,
			UserID:   claims.GetULID(),
			PocketID: req.PocketIDFrom,
			CategoryID: xulid.NullULID{
				ULID:  xulid.MustParse(constant.CAT_TRANSFER_OUT_ID),
				Valid: true,
			},
			Name:      fmt.Sprintf("Transfer To %s", toPocket.PocketName),
			Price:     -req.Price,
			IsIncome:  false,
			SpendType: 0,
			Date:      req.Date,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Version:   1,
		}

		// spend for pocket-to
		spendIDTo := xulid.Instance().NewULID()
		spendTo := model.Spend{
			ID:       spendIDTo,
			UserID:   claims.GetULID(),
			PocketID: req.PocketIDTo,
			CategoryID: xulid.NullULID{
				ULID:  xulid.MustParse(constant.CAT_TRANSFER_IN_ID),
				Valid: true,
			},
			Name:      fmt.Sprintf("Transfer From %s", fromPocket.PocketName),
			Price:     req.Price,
			IsIncome:  true,
			SpendType: 0,
			Date:      req.Date,
			CreatedAt: timeNow,
			UpdatedAt: timeNow,
			Version:   1,
		}

		// prevent deadlock we must order execution based on consistency value
		// in this case order by uuid
		spends := []model.Spend{spend, spendTo}
		// Sort spends by UUID
		sort.Slice(spends, func(i, j int) bool {
			return spends[i].ID.String() < spends[j].ID.String()
		})

		for _, ss := range spends {
			err = s.repo.Insert(ctx, &ss)
			if err != nil {
				return fmt.Errorf("insert spend to db - %s: %w", ss.PocketName, err)
			}

			_, err := s.pocketRepo.UpdateBalance(ctx, ss.PocketID, ss.Price, false)
			if err != nil {
				return fmt.Errorf("fail to change balance - %s: %w", ss.PocketName, err)
			}
		}

		return nil
	})

	if transErr != nil {
		return transErr
	}

	return nil
}

func (s Core) UpdatePartialSpend(ctx context.Context, claims mjwt.CustomClaim, req model.UpdateSpend) (model.SpendResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-UpdatePartialSpend")
	defer span.End()

	// Get existing Spend
	spendExisting, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get spend by id: %w", err)
	}

	// validate id creator
	if spendExisting.UserID != claims.GetULID() {
		return model.SpendResp{}, errr.New("user cannot edit this transaction", 400)
	}

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, spendExisting.PocketID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) {
		return model.SpendResp{}, errr.New("not have access to this pocket", 400)
	}

	// Modify data
	if req.CategoryID.Valid {
		spendExisting.CategoryID = req.CategoryID
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
func (s Core) GetDetail(ctx context.Context, spendID xulid.ULID) (model.SpendResp, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-GetDetail")
	defer span.End()

	// Get existing Spend
	spendDetail, err := s.repo.GetByID(ctx, spendID)
	if err != nil {
		return model.SpendResp{}, fmt.Errorf("get spend detail by id: %w", err)
	}

	return spendDetail.ToResp(), nil
}

// FindAllSpend ...
func (s Core) FindAllSpend(ctx context.Context, claims mjwt.CustomClaim, spendFilter model.SpendFilter, filter data.Filters) ([]model.SpendResp, data.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindAllSpend")
	defer span.End()

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, spendFilter.PocketID.ULID)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) &&
		!slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.WatcherID) {
		return nil, data.Metadata{}, errr.New("not have access to this pocket", 400)
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
func (s Core) SyncBalance(ctx context.Context, claims mjwt.CustomClaim, pocketID xulid.ULID) (int64, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-SyncBalance")
	defer span.End()

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, pocketID)
	if err != nil {
		return 0, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) &&
		!slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.WatcherID) {
		return 0, errr.New("not have access to this pocket", 400)
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
