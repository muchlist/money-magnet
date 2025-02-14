package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	notifModel "github.com/muchlist/moneymagnet/business/notification/model"
	"github.com/muchlist/moneymagnet/business/spend/model"
	"github.com/muchlist/moneymagnet/business/spend/port"
	"github.com/muchlist/moneymagnet/constant"
	"github.com/muchlist/moneymagnet/pkg/bg"
	"github.com/muchlist/moneymagnet/pkg/ctype"
	"github.com/muchlist/moneymagnet/pkg/daterange"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/observ"
	"github.com/muchlist/moneymagnet/pkg/paging"
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
	log                mlogger.Logger
	repo               port.SpendStorer
	pocketRepo         port.PocketStorer
	eTagRepo           port.ETagStorer
	notificationSender port.NotificationSender
	txManager          port.Transactor
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	repo port.SpendStorer,
	pocketRepo port.PocketStorer,
	eTagRepo port.ETagStorer,
	notificationSender port.NotificationSender,
	txManager port.Transactor,
) *Core {
	return &Core{
		log:                log,
		repo:               repo,
		pocketRepo:         pocketRepo,
		eTagRepo:           eTagRepo,
		notificationSender: notificationSender,
		txManager:          txManager,
	}
}

func (s *Core) CreateSpend(ctx context.Context, claims mjwt.CustomClaim, req model.NewSpend) (model.SpendResp, error) {
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

	isIncome := req.Price > 0

	spend := model.Spend{
		ID:               spendID,
		UserID:           claims.GetULID(),
		PocketID:         req.PocketID,
		CategoryID:       req.CategoryID,
		Name:             ctype.ToUppercaseString(req.Name),
		Price:            req.Price,
		BalanceSnapshoot: 0,
		IsIncome:         isIncome,
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

	// updating eTag
	bg.RunSafeBackground(ctx, bg.BackgroundJob{
		JobTitle: "set etag for pocket",
		Execute: func(ctx context.Context) {
			err := s.eTagRepo.SetTagByPocketID(ctx, pocketExisting.ID.String(), time.Now().UnixMilli())
			if err != nil {
				s.log.ErrorT(ctx, fmt.Sprintf("error set eTag for pocket %s", pocketExisting.ID.String()), err)
			}
		},
	})

	// send notification to other user if any
	otherUsers := pocketExisting.GetOtherUsers(claims.Identity)
	if len(otherUsers) != 0 {
		bg.RunSafeBackground(ctx, bg.BackgroundJob{
			JobTitle: "Send Notification Create Spend",
			Execute: func(ctx context.Context) {
				err := s.notificationSender.SendNotificationToUser(ctx, notifModel.SendMessage{
					Title:   fmt.Sprintf("Penambahan record pada %s oleh %s", pocketExisting.PocketName, claims.Name),
					Message: fmt.Sprintf("%s %d", req.Name, req.Price),
					UserIds: otherUsers,
				})
				if err != nil {
					s.log.ErrorT(ctx, "error send notification to user", err)
				}
			},
		})
	}

	return spend.ToResp(), nil
}

func (s *Core) TransferToPocketAsSpend(ctx context.Context, claims mjwt.CustomClaim, req model.TransferSpend) error {
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
			Name:      ctype.ToUppercaseString(fmt.Sprintf("Transfer To %s", toPocket.PocketName)),
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
			Name:      ctype.ToUppercaseString(fmt.Sprintf("Transfer From %s", fromPocket.PocketName)),
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

	// updating eTag
	bg.RunSafeBackground(ctx, bg.BackgroundJob{
		JobTitle: "set etag for pockets",
		Execute: func(ctx context.Context) {
			err := s.eTagRepo.SetTagByPocketID(ctx, req.PocketIDFrom.String(), time.Now().UnixMilli())
			if err != nil {
				s.log.ErrorT(ctx, fmt.Sprintf("error set eTag for pocket %s", req.PocketIDFrom.String()), err)
			}

			err = s.eTagRepo.SetTagByPocketID(ctx, req.PocketIDTo.String(), time.Now().UnixMilli())
			if err != nil {
				s.log.ErrorT(ctx, fmt.Sprintf("error set eTag for pocket %s", req.PocketIDTo.String()), err)
			}
		},
	})

	return nil
}

func (s *Core) UpdatePartialSpend(ctx context.Context, claims mjwt.CustomClaim, req model.UpdateSpend) (model.SpendResp, error) {
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
		name := ctype.ToUppercaseString(*req.Name)
		spendExisting.Name = name
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

		if *req.Price > 0 {
			spendExisting.IsIncome = true
		}
	}

	// Edit
	transErr := s.txManager.WithAtomic(ctx, func(ctx context.Context) error {
		err := s.repo.Edit(ctx, &spendExisting)
		if err != nil {
			return fmt.Errorf("edit spend: %w", err)
		}

		if diff != 0 {
			newBalance, err := s.pocketRepo.UpdateBalance(ctx, spendExisting.PocketID, diff, false)
			if err != nil {
				return fmt.Errorf("fail to change balance: %w", err)
			}
			spendExisting.BalanceSnapshoot = newBalance
		}

		return nil
	})
	if transErr != nil {
		return model.SpendResp{}, transErr
	}

	// updating eTag
	bg.RunSafeBackground(ctx, bg.BackgroundJob{
		JobTitle: "set etag for pocket",
		Execute: func(ctx context.Context) {
			err := s.eTagRepo.SetTagByPocketID(ctx, spendExisting.PocketID.String(), time.Now().UnixMilli())
			if err != nil {
				s.log.ErrorT(ctx, fmt.Sprintf("error set eTag for pocket %s", spendExisting.PocketID.String()), err)
			}
		},
	})

	// send notification to other user if any
	otherUsers := pocketExisting.GetOtherUsers(claims.Identity)
	if len(otherUsers) != 0 {
		bg.RunSafeBackground(ctx, bg.BackgroundJob{
			JobTitle: "Send Notification Update Spend",
			Execute: func(ctx context.Context) {
				err := s.notificationSender.SendNotificationToUser(ctx, notifModel.SendMessage{
					Title:   fmt.Sprintf("Perubahan record pada %s oleh %s", pocketExisting.PocketName, claims.Name),
					Message: fmt.Sprintf("%s %d", spendExisting.Name, spendExisting.Price),
					UserIds: otherUsers,
				})
				if err != nil {
					s.log.ErrorT(ctx, "error send notification to user", err)
				}
			},
		})
	}

	return spendExisting.ToResp(), nil
}

func (s *Core) DeleteSpend(ctx context.Context, claims mjwt.CustomClaim, spendID xulid.ULID) error {
	ctx, span := observ.GetTracer().Start(ctx, "service-UpdatePartialSpend")
	defer span.End()

	// Get existing Spend
	spendExisting, err := s.repo.GetByID(ctx, spendID)
	if err != nil {
		return fmt.Errorf("get spend by id: %w", err)
	}

	// validate id creator
	if spendExisting.UserID != claims.GetULID() {
		return errr.New("user cannot delete this transaction", 400)
	}

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, spendExisting.PocketID)
	if err != nil {
		return fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) {
		return errr.New("not have access to this pocket", 400)
	}

	reverseExistingPriceToDelete := -spendExisting.Price

	// Edit
	transErr := s.txManager.WithAtomic(ctx, func(ctx context.Context) error {
		err := s.repo.Delete(ctx, spendID)
		if err != nil {
			return fmt.Errorf("delete spend: %w", err)
		}

		_, err = s.pocketRepo.UpdateBalance(ctx, spendExisting.PocketID, reverseExistingPriceToDelete, false)
		if err != nil {
			return fmt.Errorf("fail to updating balance: %w", err)
		}

		return nil
	})
	if transErr != nil {
		return transErr
	}

	// updating eTag
	bg.RunSafeBackground(ctx, bg.BackgroundJob{
		JobTitle: "set etag for pocket",
		Execute: func(ctx context.Context) {
			err := s.eTagRepo.SetTagByPocketID(ctx, spendExisting.PocketID.String(), time.Now().UnixMilli())
			if err != nil {
				s.log.ErrorT(ctx, fmt.Sprintf("error set eTag for pocket %s", spendExisting.PocketID.String()), err)
			}
		},
	})

	// send notification to other user if any
	otherUsers := pocketExisting.GetOtherUsers(claims.Identity)
	if len(otherUsers) != 0 {
		bg.RunSafeBackground(ctx, bg.BackgroundJob{
			JobTitle: "Send Notification Delete Spend",
			Execute: func(ctx context.Context) {
				err := s.notificationSender.SendNotificationToUser(ctx, notifModel.SendMessage{
					Title:   fmt.Sprintf("Penghapusan record pada %s oleh %s", pocketExisting.PocketName, claims.Name),
					Message: fmt.Sprintf("%s %d", spendExisting.Name, spendExisting.Price),
					UserIds: otherUsers,
				})
				if err != nil {
					s.log.ErrorT(ctx, "error send notification to user", err)
				}
			},
		})
	}

	return nil
}

// GetDetail ...
func (s *Core) GetDetail(ctx context.Context, spendID xulid.ULID) (model.SpendResp, error) {
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
func (s *Core) FindAllSpend(ctx context.Context, claims mjwt.CustomClaim, spendFilter model.SpendFilter, filter paging.Filters) ([]model.SpendResp, paging.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindAllSpend")
	defer span.End()

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, spendFilter.PocketID.ULID)
	if err != nil {
		return nil, paging.Metadata{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) &&
		!slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.WatcherID) {
		return nil, paging.Metadata{}, errr.New("not have access to this pocket", 400)
	}

	spends, metadata, err := s.repo.Find(ctx, spendFilter, filter)
	if err != nil {
		return nil, paging.Metadata{}, fmt.Errorf("find spend by pocketID: %w", err)
	}

	spendResult := make([]model.SpendResp, len(spends))
	for i := range spends {
		spendResult[i] = spends[i].ToResp()
	}

	return spendResult, metadata, nil
}

// FindAllSpendByCursor ...
func (s *Core) FindAllSpendByCursor(ctx context.Context, claims mjwt.CustomClaim, spendFilter model.SpendFilter, filter paging.Cursor) ([]model.SpendResp, paging.CursorMetadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindAllSpend")
	defer span.End()

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, spendFilter.PocketID.ULID)
	if err != nil {
		return nil, paging.CursorMetadata{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) &&
		!slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.WatcherID) {
		return nil, paging.CursorMetadata{}, errr.New("not have access to this pocket", 400)
	}

	spends, err := s.repo.FindWithCursor(ctx, spendFilter, filter)
	if err != nil {
		return nil, paging.CursorMetadata{}, fmt.Errorf("find all spend by pocketID with cursor: %w", err)
	}

	// Menentukan cursor selanjutnya
	var reverseCursor string
	var nextCursor string
	if len(spends) > 0 {
		reverseCursor = spends[0].Date.Format(time.RFC3339) // Set prev cursor
	}
	if len(spends) > int(filter.GetPageSize()) {
		nextCursor = spends[filter.GetPageSize()-1].Date.Format(time.RFC3339) // Set next cursor apabila ditemukan data lebih dari limit
		spends = spends[:filter.GetPageSize()]                                // Hapus data yang kelebihan
	}

	spendResult := make([]model.SpendResp, len(spends))
	for i := range spends {
		spendResult[i] = spends[i].ToResp()
	}

	return spendResult, paging.CursorMetadata{
		CurrentCursor: filter.GetCursor(),
		CursorType:    filter.GetCursorType(),
		PageSize:      filter.GetPageSize(),
		NextCursor:    nextCursor,
		NextPage:      "",
		ReverseCursor: reverseCursor,
		ReversePage:   "",
	}, nil
}

type AutoDateRangeParams struct {
	PocketID  xulid.ULID
	Claims    mjwt.CustomClaim
	Filter    paging.Cursor
	RangeType string
	TimeZone  string
	ETag      string
}

func (s *Core) FindAllSpendByCursorAutoDateRange(ctx context.Context, params AutoDateRangeParams) ([]model.SpendResp, paging.CursorMetadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindAllSpendByCursorAutoDateRange")
	defer span.End()

	// only support eTag for rangeType week right now
	islast7daysType := params.RangeType == "last-7-days"
	isPageOne := params.Filter.GetCursor() == ""
	var tagSaved int64 = 0

	if islast7daysType && isPageOne {
		// always check eTag from redis
		tag, err := s.eTagRepo.GetTagByPocketID(ctx, params.PocketID.String())
		if err != nil {
			return nil,
				paging.CursorMetadata{},
				errr.New(fmt.Sprintf("error get tag: %s", err.Error()), 500)
		}

		if tag != 0 {
			tagSaved = tag
			tagInput, _ := strconv.ParseInt(params.ETag, 10, 64)
			if tag == tagInput {
				return nil,
					paging.CursorMetadata{},
					errr.New(
						"data is uptodate",
						http.StatusNotModified,
					)
			}
		}
	}

	dateRange, err := daterange.ParseDateRange(params.RangeType, params.TimeZone)
	if err != nil {
		return nil, paging.CursorMetadata{}, errr.New(err.Error(), 400)
	}

	spendFilter := model.SpendFilter{
		DateStart: &dateRange.StartDate,
		DateEnd:   &dateRange.EndDate,
	}
	spendFilter.PocketID.ULID = params.PocketID

	// use exisitng logic to get spend by cursor
	results, metaPaging, err := s.FindAllSpendByCursor(ctx, params.Claims, spendFilter, params.Filter)
	if err != nil {
		return nil, paging.CursorMetadata{}, err
	}

	// updating eTag if empty
	if islast7daysType && (tagSaved == 0) && isPageOne {
		tagSaved = time.Now().UnixMilli() // updating tag coz existing is empty
		bg.RunSafeBackground(ctx, bg.BackgroundJob{
			JobTitle: "set etag for pocket",
			Execute: func(ctx context.Context) {
				err := s.eTagRepo.SetTagByPocketID(ctx, params.PocketID.String(), tagSaved)
				if err != nil {
					s.log.ErrorT(ctx, fmt.Sprintf("error set eTag for pocket %s", params.PocketID.String()), err)
				}
			},
		})
	}

	metaPaging.ETag = fmt.Sprintf("%d", tagSaved)

	return results, metaPaging, err
}

// FindAllSpendMultiPocketByCursor ...
func (s *Core) FindAllSpendMultiPocketByCursor(ctx context.Context, claims mjwt.CustomClaim, spendFilter model.SpendFilterMultiPocket, filter paging.Cursor) ([]model.SpendResp, paging.CursorMetadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindAllSpendMultiPocketByCursor")
	defer span.End()

	// Must Have PocketID
	if len(spendFilter.Pockets) == 0 {
		return nil, paging.CursorMetadata{}, errr.New("pocket id is required", http.StatusBadRequest)
	}

	// Get existing Pocket
	// TECHDEBT MVP : Validate just first pocket rather than all pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, spendFilter.Pockets[0])
	if err != nil {
		return nil, paging.CursorMetadata{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Validate Pocket Roles Editor
	if !slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.EditorID) &&
		!slicer.In(xulid.MustParse(claims.Identity).String(), pocketExisting.WatcherID) {
		return nil, paging.CursorMetadata{}, errr.New("not have access to this pocket", 400)
	}

	spends, err := s.repo.FindWithCursorMultiPockets(ctx, spendFilter, filter)
	if err != nil {
		return nil, paging.CursorMetadata{}, fmt.Errorf("find all spend by multi filter with cursor: %w", err)
	}

	// Menentukan cursor selanjutnya
	var reverseCursor string
	var nextCursor string
	if len(spends) > 0 {
		reverseCursor = spends[0].Date.Format(time.RFC3339) // Set prev cursor
	}
	if len(spends) > int(filter.GetPageSize()) {
		nextCursor = spends[filter.GetPageSize()-1].Date.Format(time.RFC3339) // Set next cursor apabila ditemukan data lebih dari limit
		spends = spends[:filter.GetPageSize()]                                // Hapus data yang kelebihan
	}

	spendResult := make([]model.SpendResp, len(spends))
	for i := range spends {
		spendResult[i] = spends[i].ToResp()
	}

	return spendResult, paging.CursorMetadata{
		CurrentCursor: filter.GetCursor(),
		CursorType:    filter.GetCursorType(),
		PageSize:      filter.GetPageSize(),
		NextCursor:    nextCursor,
		NextPage:      "",
		ReverseCursor: reverseCursor,
		ReversePage:   "",
	}, nil
}

// SyncBalance ...
func (s *Core) SyncBalance(ctx context.Context, claims mjwt.CustomClaim, pocketID xulid.ULID) (int64, error) {
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
