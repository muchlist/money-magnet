package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mstore "github.com/muchlist/moneymagnet/business/pocket/mock_storer"
	"github.com/muchlist/moneymagnet/business/pocket/model"
	urmodel "github.com/muchlist/moneymagnet/business/user/model"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/stretchr/testify/assert"
)

var log = mlogger.New(mlogger.Options{
	Level:        "panic",
	Output:       "stdout",
	ContextField: nil,
})

func TestCreatePocketSuccess(t *testing.T) {
	// input output
	ctx := context.Background()
	ownerUUID := uuid.New()
	pocketUUID := uuid.New()
	payload := model.NewPocket{
		PocketName: "example pocket",
		EditorID:   []uuid.UUID{},
		WatcherID:  []uuid.UUID{},
		Icon:       1,
	}

	timeNow := time.Now()
	expect := model.PocketResp{
		ID:         pocketUUID,
		OwnerID:    ownerUUID,
		EditorID:   []uuid.UUID{ownerUUID},
		WatcherID:  []uuid.UUID{ownerUUID},
		PocketName: "example pocket",
		Icon:       1,
		Level:      1,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
		Version:    1,
	}

	// dependency
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock user
	userRepo := mstore.NewMockUserReader(ctrl)
	userRepo.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return(
		[]urmodel.User{
			{
				ID: ownerUUID,
			},
		}, nil,
	)

	// mock pocket
	pocketReplacePtr := model.Pocket{
		ID:         pocketUUID,
		OwnerID:    ownerUUID,
		EditorID:   []uuid.UUID{ownerUUID},
		WatcherID:  []uuid.UUID{ownerUUID},
		PocketName: "example pocket",
		Icon:       1,
		Level:      1,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
		Version:    1,
	}
	pocketRepo := mstore.NewMockPocketStorer(ctrl)
	pocketRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).SetArg(1, pocketReplacePtr).Return(nil)
	pocketRepo.EXPECT().InsertPocketUser(gomock.Any(), gomock.Any(), gomock.Eq(pocketReplacePtr.ID)).Return(nil)
	pocketRepo.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(x any, tFunc func(ctx context.Context) error) error {
		return tFunc(ctx)
	})

	// init service
	service := NewCore(log, pocketRepo, userRepo)

	// Assertion
	result, err := service.CreatePocket(ctx, ownerUUID, payload)
	assert.Nil(t, err)
	assert.Equal(t, expect, result)
}

func TestCreatePocketFailInsertUser(t *testing.T) {
	// input output
	ctx := context.Background()
	ownerUUID := uuid.New()
	pocketUUID := uuid.New()
	payload := model.NewPocket{
		PocketName: "example pocket",
		EditorID:   []uuid.UUID{},
		WatcherID:  []uuid.UUID{},
		Icon:       1,
	}

	timeNow := time.Now()

	// dependency
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock user
	userRepo := mstore.NewMockUserReader(ctrl)
	userRepo.EXPECT().GetByIDs(gomock.Any(), gomock.Any()).Return(
		[]urmodel.User{
			{
				ID: ownerUUID,
			},
		}, nil,
	)

	// mock pocket
	pocketReplacePtr := model.Pocket{
		ID:         pocketUUID,
		OwnerID:    ownerUUID,
		EditorID:   []uuid.UUID{ownerUUID},
		WatcherID:  []uuid.UUID{ownerUUID},
		PocketName: "example pocket",
		Icon:       1,
		Level:      1,
		CreatedAt:  timeNow,
		UpdatedAt:  timeNow,
		Version:    1,
	}
	pocketRepo := mstore.NewMockPocketStorer(ctrl)
	pocketRepo.EXPECT().Insert(gomock.Any(), gomock.Any()).SetArg(1, pocketReplacePtr).Return(errors.New("context deadline exceeded"))
	pocketRepo.EXPECT().WithinTransaction(gomock.Any(), gomock.Any()).DoAndReturn(func(x any, tFunc func(ctx context.Context) error) error {
		return tFunc(ctx)
	})

	// init service
	service := NewCore(log, pocketRepo, userRepo)

	// Assertion
	result, err := service.CreatePocket(ctx, ownerUUID, payload)
	assert.Equal(t, "insert pocket to db: context deadline exceeded", err.Error())
	assert.Equal(t, model.PocketResp{}, result)
}
