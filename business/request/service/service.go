package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/business/request/model"
	"github.com/muchlist/moneymagnet/business/request/port"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
)

// Core manages the set of APIs for request access.
type Core struct {
	log        mlogger.Logger
	repo       port.RequestStorer
	pocketRepo port.PocketStorer
	txManager  port.Transactor
}

// NewCore constructs a core for request api access.
func NewCore(
	log mlogger.Logger,
	repo port.RequestStorer,
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

func (s Core) CreateRequest(ctx context.Context, claims mjwt.CustomClaim, pocketID uuid.UUID) (model.RequestPocket, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-CreateRequest")
	defer span.End()

	timeNow := time.Now()

	// GET Pocket BY ID
	pocket, err := s.pocketRepo.GetByID(ctx, pocketID)
	if err != nil {
		return model.RequestPocket{}, fmt.Errorf("get pocket by id: %w", err)
	}

	req := model.RequestPocket{
		RequesterID: claims.GetUUID(),
		PocketID:    pocketID,
		PocketName:  pocket.PocketName,
		ApproverID:  &pocket.OwnerID,
		CreatedAt:   timeNow,
		UpdatedAt:   timeNow,
	}
	if err = s.repo.Insert(ctx, &req); err != nil {
		return model.RequestPocket{}, fmt.Errorf("insert request: %w", err)
	}

	// TODO SEND TO FCM

	return req, nil
}

func (s Core) ApproveRequest(ctx context.Context, claims mjwt.CustomClaim, IsApproved bool, requestID uint64) error {
	ctx, span := observ.GetTracer().Start(ctx, "service-ApproveRequest")
	defer span.End()

	// GET Request by ID
	req, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return fmt.Errorf("get request by id: %w", err)
	}

	if *req.ApproverID != claims.GetUUID() {
		return errr.New("the user does not have access rights to approve this request", 400)
	}

	// check if either have true value
	if req.IsApproved || req.IsRejected {
		return errr.New("This request has been processed before", 400)
	}

	if IsApproved {
		req.IsApproved = true
	} else {
		req.IsRejected = true
	}

	if err = s.repo.UpdateStatus(ctx, &req); err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	// IF req.IsApproved update in pocket editor and watcher
	// else return
	if req.IsRejected {
		return nil
	}

	// TODO SEND TO FCM

	// Get existing Pocket
	pocketExisting, err := s.pocketRepo.GetByID(ctx, req.PocketID)
	if err != nil {
		return fmt.Errorf("get pocket by id: %w", err)
	}

	// add to wathcer
	pocketExisting.WatcherID = append(pocketExisting.WatcherID, req.RequesterID)
	// add to editor
	pocketExisting.EditorID = append(pocketExisting.EditorID, req.RequesterID)

	// TRANSACTION
	err = s.txManager.WithAtomic(ctx, func(ctx context.Context) error {
		// Edit
		err = s.pocketRepo.Edit(ctx, &pocketExisting)
		if err != nil {
			return fmt.Errorf("edit pocket: %w", err)
		}
		// insert to related table
		err = s.pocketRepo.InsertPocketUser(ctx, []uuid.UUID{req.RequesterID}, pocketExisting.ID)
		if err != nil {
			return fmt.Errorf("insert pocket_user to db: %w", err)
		}

		return nil
	})

	return err
}

// FindAllByRequester ...
func (s Core) FindAllByRequester(ctx context.Context, claims mjwt.CustomClaim, filter data.Filters) ([]model.RequestPocket, data.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindAllByRequester")
	defer span.End()

	// Get All Request
	findBy := model.FindBy{
		RequesterID: claims.Identity,
	}

	reqs, metadata, err := s.repo.Find(ctx, findBy, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find request: %w", err)
	}

	return reqs, metadata, nil
}

// FindAllByApprover ...
func (s Core) FindAllByApprover(ctx context.Context, claims mjwt.CustomClaim, filter data.Filters) ([]model.RequestPocket, data.Metadata, error) {
	ctx, span := observ.GetTracer().Start(ctx, "service-FindAllByApprover")
	defer span.End()

	// Get All Request
	findBy := model.FindBy{
		ApproverID: claims.Identity,
	}

	reqs, metadata, err := s.repo.Find(ctx, findBy, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find request: %w", err)
	}

	return reqs, metadata, nil
}
