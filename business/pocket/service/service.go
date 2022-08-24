package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/business/pocket/storer"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/errr"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/utils/ds"
	"github.com/muchlist/moneymagnet/pkg/utils/slicer"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("data not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Core manages the set of APIs for user access.
type Core struct {
	log      mlogger.Logger
	repo     storer.PocketStorer
	userRepo storer.UserReader
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	repo storer.PocketStorer,
	userRepo storer.UserReader,
) Core {
	return Core{
		log:      log,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s Core) CreatePocket(ctx context.Context, owner uuid.UUID, req model.PocketNew) (model.PocketResp, error) {
	// Validate editor and watcher uuids
	combineUserUUIDs := append(req.Editor, req.Watcher...)
	users, err := s.userRepo.GetByIDs(ctx, combineUserUUIDs)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get users: %w", err)
	}
	for _, id := range combineUserUUIDs {
		found := false
		for _, user := range users {
			if user.ID == id {
				found = true
				break
			}
		}
		if !found {
			return model.PocketResp{}, errr.New(fmt.Sprintf("uuid %s is not have valid user", id), 400)
		}
	}

	if req.Editor == nil || len(req.Editor) == 0 {
		req.Editor = []uuid.UUID{owner}
	}
	if req.Watcher == nil || len(req.Watcher) == 0 {
		req.Watcher = []uuid.UUID{owner}
	}

	timeNow := time.Now()
	pocket := model.Pocket{
		Owner:      owner,
		Editor:     req.Editor,
		Watcher:    req.Watcher,
		PocketName: req.PocketName,
		Level:      1,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
		Version:    1,
	}

	err = s.repo.Insert(ctx, &pocket)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("insert pocket to db: %w", err)
	}

	// insert relation
	uuidUserSet := ds.NewUUIDSet()
	uuidUserSet.Add(owner)
	uuidUserSet.AddAll(combineUserUUIDs)
	uniqueUsers := uuidUserSet.Reveal()

	err = s.repo.InsertPocketUser(ctx, uniqueUsers, pocket.ID)
	if err != nil {
		return pocket.ToPocketResp(), fmt.Errorf("loop insert pocket_user to db: %w", err)
	}

	return pocket.ToPocketResp(), nil
}

func (s Core) RenamePocket(ctx context.Context, owner uuid.UUID, pocketID uint64, newName string) (model.PocketResp, error) {

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, pocketID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Check if owner not have access to pocket
	if !(pocketExisting.Owner == owner ||
		slicer.In(owner, pocketExisting.Editor)) {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	// Modify data
	pocketExisting.PocketName = newName

	// Edit
	s.repo.Edit(ctx, &pocketExisting)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("edit pocket: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
}

func (s Core) AddPerson(ctx context.Context, data AddPersonData) (model.PocketResp, error) {

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, data.PocketID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Check if owner not have access to pocket
	if !(pocketExisting.Owner == data.Owner ||
		slicer.In(data.Owner, pocketExisting.Editor)) {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	// Check if person to add is exist
	_, err = s.userRepo.GetByID(ctx, data.Person)
	if err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return model.PocketResp{}, errr.New("account is not exist", 400)
		}
		return model.PocketResp{}, fmt.Errorf("get user by id : %w", err)
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
		return model.PocketResp{}, fmt.Errorf("edit pocket: %w", err)
	}

	// insert to related table
	err = s.repo.InsertPocketUser(ctx, []uuid.UUID{data.Person}, pocketExisting.ID)
	if err != nil {
		return pocketExisting.ToPocketResp(), fmt.Errorf("insert pocket_user to db: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
}

// RemovePerson will remove person from both editor and watcher
func (s Core) RemovePerson(ctx context.Context, data RemovePersonData) (model.PocketResp, error) {
	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, data.PocketID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Check if owner not have access to pocket
	if !(pocketExisting.Owner == data.Owner ||
		slicer.In(data.Owner, pocketExisting.Editor)) {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	pocketExisting.Editor = slicer.RemoveFrom(data.Person, pocketExisting.Editor)
	pocketExisting.Watcher = slicer.RemoveFrom(data.Person, pocketExisting.Watcher)

	// Edit
	s.repo.Edit(ctx, &pocketExisting)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("edit pocket: %w", err)
	}

	// delete from related table
	err = s.repo.DeletePocketUser(ctx, data.Person, pocketExisting.ID)
	if err != nil {
		return pocketExisting.ToPocketResp(), fmt.Errorf("delete pocket_user from db: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
}

// GetDetail ...
func (s Core) GetDetail(ctx context.Context, userID string, pocketID uint64) (model.PocketResp, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return model.PocketResp{}, ErrInvalidID
	}

	// Get existing Pocket
	pocketDetail, err := s.repo.GetByID(ctx, pocketID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket detail by id: %w", err)
	}

	// Check if owner not have access to pocket
	if !(pocketDetail.Owner == userUUID ||
		slicer.In(userUUID, pocketDetail.Watcher)) {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	return pocketDetail.ToPocketResp(), nil
}

// FindAllPocket ...
func (s Core) FindAllPocket(ctx context.Context, userID string, filter data.Filters) ([]model.PocketResp, data.Metadata, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, data.Metadata{}, ErrInvalidID
	}

	// Get existing Pocket
	pockets, metadata, err := s.repo.FindUserPockets(ctx, userUUID, filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find pocket user: %w", err)
	}

	pocketResult := make([]model.PocketResp, len(pockets))
	for i := range pockets {
		pocketResult[i] = pockets[i].ToPocketResp()
	}

	return pocketResult, metadata, nil
}
