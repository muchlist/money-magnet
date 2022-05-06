package ptservice

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/ptmodel"
	"github.com/muchlist/moneymagnet/bussines/core/pocket/storer"
	"github.com/muchlist/moneymagnet/bussines/sys/db"
	"github.com/muchlist/moneymagnet/bussines/sys/errr"
	"github.com/muchlist/moneymagnet/foundation/mlogger"
	"github.com/muchlist/moneymagnet/foundation/utils/ds"
	"github.com/muchlist/moneymagnet/foundation/utils/slicer"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound  = errors.New("data not found")
	ErrInvalidID = errors.New("ID is not in its proper form")
)

// Service manages the set of APIs for user access.
type Service struct {
	log      mlogger.Logger
	repo     storer.PocketStorer
	userRepo storer.UserReader
}

// NewService constructs a core for user api access.
func NewService(
	log mlogger.Logger,
	repo storer.PocketStorer,
	userRepo storer.UserReader,
) Service {
	return Service{
		log:      log,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s Service) CreatePocket(ctx context.Context, owner uuid.UUID, req ptmodel.PocketNew) (ptmodel.PocketResp, error) {
	// Validate editor and watcher uuids
	combineUserUUIDs := append(req.Editor, req.Watcher...)
	users, err := s.userRepo.GetByIDs(ctx, combineUserUUIDs)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("get users: %w", err)
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
			return ptmodel.PocketResp{}, errr.New(fmt.Sprintf("uuid %s is not have valid user", id), 400)
		}
	}

	if req.Editor == nil || len(req.Editor) == 0 {
		req.Editor = []uuid.UUID{owner}
	}
	if req.Watcher == nil || len(req.Watcher) == 0 {
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

	err = s.repo.Insert(ctx, &pocket)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("insert pocket to db: %w", err)
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

func (s Service) RenamePocket(ctx context.Context, owner uuid.UUID, pocketID uint64, newName string) (ptmodel.PocketResp, error) {

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, pocketID)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Check if owner not have access to pocket
	if !(pocketExisting.Owner == owner ||
		slicer.In(owner, pocketExisting.Editor)) {
		return ptmodel.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	// Modify data
	pocketExisting.PocketName = newName

	// Edit
	s.repo.Edit(ctx, &pocketExisting)
	if err != nil {
		return ptmodel.PocketResp{}, fmt.Errorf("edit pocket: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
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

	// Check if person to add is exist
	_, err = s.userRepo.GetByID(ctx, data.Person)
	if err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return ptmodel.PocketResp{}, errr.New("account is not exist", 400)
		}
		return ptmodel.PocketResp{}, fmt.Errorf("get user by id : %w", err)
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

	// insert to related table
	err = s.repo.InsertPocketUser(ctx, []uuid.UUID{data.Person}, pocketExisting.ID)
	if err != nil {
		return pocketExisting.ToPocketResp(), fmt.Errorf("insert pocket_user to db: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
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

	// delete from related table
	err = s.repo.DeletePocketUser(ctx, data.Person, pocketExisting.ID)
	if err != nil {
		return pocketExisting.ToPocketResp(), fmt.Errorf("delete pocket_user from db: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
}
