package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/muchlist/moneymagnet/business/pocket/model"
	"github.com/muchlist/moneymagnet/business/pocket/storer"
	"github.com/muchlist/moneymagnet/business/shared"
	"github.com/muchlist/moneymagnet/pkg/data"
	"github.com/muchlist/moneymagnet/pkg/db"
	"github.com/muchlist/moneymagnet/pkg/errr"
	"github.com/muchlist/moneymagnet/pkg/mjwt"

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

func (s Core) CreatePocket(ctx context.Context, claims mjwt.CustomClaim, req model.NewPocket) (model.PocketResp, error) {
	// Validate editor and watcher uuids
	combineUserUUIDs := append(req.EditorID, req.WatcherID...)
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

	if req.EditorID == nil || len(req.EditorID) == 0 {
		req.EditorID = []uuid.UUID{claims.GetUUID()}
	}
	if req.WatcherID == nil || len(req.WatcherID) == 0 {
		req.WatcherID = []uuid.UUID{claims.GetUUID()}
	}

	timeNow := time.Now()
	pocket := model.Pocket{
		OwnerID:    claims.GetUUID(),
		EditorID:   req.EditorID,
		WatcherID:  req.WatcherID,
		PocketName: req.PocketName,
		Level:      1,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
		Version:    1,
	}

	// Run IN Transaction
	transErr := s.repo.WithinTransaction(
		ctx, func(ctx context.Context) error {

			// insert pocket
			err = s.repo.Insert(ctx, &pocket)
			if err != nil {
				return fmt.Errorf("insert pocket to db: %w", err)
			}

			// insert relation
			uuidUserSet := ds.NewUUIDSet()
			uuidUserSet.Add(claims.GetUUID())
			uuidUserSet.AddAll(combineUserUUIDs)
			uniqueUsers := uuidUserSet.Reveal()

			err = s.repo.InsertPocketUser(ctx, uniqueUsers, pocket.ID)
			if err != nil {
				return fmt.Errorf("loop insert pocket_user to db: %w", err)
			}
			return nil
		},
	)

	if transErr != nil {
		return model.PocketResp{}, transErr
	}

	return pocket.ToPocketResp(), nil
}

func (s Core) UpdatePocket(ctx context.Context, claims mjwt.CustomClaim, newData model.PocketUpdate) (model.PocketResp, error) {
	// Validate Pocket Roles
	canEdit, _ := shared.IsCanEditOrWatch(newData.ID, claims.PocketRoles)
	if !canEdit {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, newData.ID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	// Modify data
	if newData.PocketName != nil {
		pocketExisting.PocketName = *newData.PocketName
	}
	if newData.Currency != nil {
		pocketExisting.Currency = *newData.Currency
	}
	if newData.Icon != nil {
		pocketExisting.Icon = *newData.Icon
	}

	// Edit
	err = s.repo.Edit(ctx, &pocketExisting)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("edit pocket: %w", err)
	}

	return pocketExisting.ToPocketResp(), nil
}

func (s Core) AddPerson(ctx context.Context, claims mjwt.CustomClaim, data AddPersonData) (model.PocketResp, error) {

	// Validate Pocket Roles
	canEdit, _ := shared.IsCanEditOrWatch(data.PocketID, claims.PocketRoles)
	if !canEdit {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, data.PocketID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
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
		pocketExisting.WatcherID = append(pocketExisting.WatcherID, data.Person)
	} else {
		// add to editor
		pocketExisting.EditorID = append(pocketExisting.EditorID, data.Person)
	}

	// Run IN Transaction
	transErr := s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		// Edit
		err = s.repo.Edit(ctx, &pocketExisting)
		if err != nil {
			return fmt.Errorf("edit pocket: %w", err)
		}

		// insert to related table
		err = s.repo.InsertPocketUser(ctx, []uuid.UUID{data.Person}, pocketExisting.ID)
		if err != nil {
			return fmt.Errorf("insert pocket_user to db: %w", err)
		}

		return nil
	})
	if transErr != nil {
		return model.PocketResp{}, transErr
	}

	return pocketExisting.ToPocketResp(), nil
}

// RemovePerson will remove person from both editor and watcher
func (s Core) RemovePerson(ctx context.Context, claims mjwt.CustomClaim, data RemovePersonData) (model.PocketResp, error) {

	// Validate Pocket Roles
	canEdit, _ := shared.IsCanEditOrWatch(data.PocketID, claims.PocketRoles)
	if !canEdit {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	// Get existing Pocket
	pocketExisting, err := s.repo.GetByID(ctx, data.PocketID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket by id: %w", err)
	}

	pocketExisting.EditorID = slicer.RemoveFrom(data.Person, pocketExisting.EditorID)
	pocketExisting.WatcherID = slicer.RemoveFrom(data.Person, pocketExisting.WatcherID)

	// Run IN Transaction
	transErr := s.repo.WithinTransaction(ctx, func(ctx context.Context) error {
		// Edit
		err = s.repo.Edit(ctx, &pocketExisting)
		if err != nil {
			return fmt.Errorf("edit pocket: %w", err)
		}

		// delete from related table
		err = s.repo.DeletePocketUser(ctx, data.Person, pocketExisting.ID)
		if err != nil {
			return fmt.Errorf("delete pocket_user from db: %w", err)
		}
		return nil
	})
	if transErr != nil {
		return model.PocketResp{}, transErr
	}

	return pocketExisting.ToPocketResp(), nil
}

// GetDetail ...
func (s Core) GetDetail(ctx context.Context, claims mjwt.CustomClaim, pocketID uuid.UUID) (model.PocketResp, error) {

	// Validate Pocket Roles
	canEdit, canWatch := shared.IsCanEditOrWatch(pocketID, claims.PocketRoles)
	if !(canEdit || canWatch) {
		return model.PocketResp{}, errr.New("not have access to this pocket", 400)
	}

	// Get existing Pocket
	pocketDetail, err := s.repo.GetByID(ctx, pocketID)
	if err != nil {
		return model.PocketResp{}, fmt.Errorf("get pocket detail by id: %w", err)
	}

	return pocketDetail.ToPocketResp(), nil
}

// FindAllPocket ...
func (s Core) FindAllPocket(ctx context.Context, claims mjwt.CustomClaim, filter data.Filters) ([]model.PocketResp, data.Metadata, error) {

	// Get existing Pocket
	pockets, metadata, err := s.repo.FindUserPockets(ctx, claims.GetUUID(), filter)
	if err != nil {
		return nil, data.Metadata{}, fmt.Errorf("find pocket user: %w", err)
	}

	pocketResult := make([]model.PocketResp, len(pockets))
	for i := range pockets {
		pocketResult[i] = pockets[i].ToPocketResp()
	}

	return pocketResult, metadata, nil
}
