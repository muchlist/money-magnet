package ptservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptmodel"
	"github.com/muchlist/moneymagnet/bussines/sys/errr"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
	"github.com/muchlist/moneymagnet/foundation/utils/slicer"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("data not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Service manages the set of APIs for user access.
type Service struct {
	log  mlogger.Logger
	repo PocketStorer
}

// NewService constructs a core for user api access.
func NewService(
	log mlogger.Logger,
	repo PocketStorer,
) Service {
	return Service{
		log:  log,
		repo: repo,
	}
}

func (s Service) CreatePocket(ctx context.Context, owner uuid.UUID, req ptmodel.PocketNew) (ptmodel.PocketResp, error) {
	if req.Editor == nil {
		req.Editor = []uuid.UUID{owner}
	}
	if req.Watcher == nil {
		req.Watcher = []uuid.UUID{owner}
	}

	timeNow := time.Now()
	pocket := ptmodel.Pocket{
		Owner:      owner,
		Editor:     req.Editor,
		Watcher:    req.Watcher,
		PocketName: req.PocketName,
		Level:      1,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
		Version:    1,
	}

	err := s.repo.Insert(ctx, &pocket)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("insert pocket to db: %w", err)
	}

	return pocket.ToPocketResp(), nil
}

type AddPersonData struct {
	Owner      uuid.UUID
	Person     uuid.UUID
	PocketID   uint64
	IsReadOnly bool
}

func (s Service) AddPerson(ctx context.Context, data AddPersonData) (ptmodel.PocketResp, error) {

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, data.PocketID)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Check if owner not have access to pocket
	if !(pocketExisting.Owner == data.Owner ||
		slicer.In(data.Owner, pocketExisting.Editor)) {
		return ptmodel.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	if data.IsReadOnly {
		// add to wathcer
		pocketExisting.Watcher = append(pocketExisting.Watcher, data.Person)
	} else {
		// add to editor
		pocketExisting.Editor = append(pocketExisting.Editor, data.Person)
	}

	// Edit
	s.repo.Edit(ctx, &pocketExisting)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("edit pocket: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
}

type RemovePersonData struct {
	Owner    uuid.UUID
	Person   uuid.UUID
	PocketID uint64
}

// RemovePerson will remove person from both editor and watcher
func (s Service) RemovePerson(ctx context.Context, data RemovePersonData) (ptmodel.PocketResp, error) {

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, data.PocketID)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Check if owner not have access to pocket
	if !(pocketExisting.Owner == data.Owner ||
		slicer.In(data.Owner, pocketExisting.Editor)) {
		return ptmodel.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	pocketExisting.Editor = slicer.RemoveFrom(data.Person, pocketExisting.Editor)
	pocketExisting.Watcher = slicer.RemoveFrom(data.Person, pocketExisting.Watcher)

	// Edit
	s.repo.Edit(ctx, &pocketExisting)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("edit pocket: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
}
