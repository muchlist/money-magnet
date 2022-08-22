package ptservice

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mstore "github.com/muchlist/moneymagnet/business/pocket/mock_storer"
	"github.com/muchlist/moneymagnet/business/pocket/ptmodel"
	"github.com/muchlist/moneymagnet/business/user/usermodel"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/stretchr/testify/assert"
)

var log = mlogger.New("panic", "stdout")

func TestCreatePocketSuccess(t *testing.T) {
	// input output
	ctx := context.Background()
	ownerUUID := uuid.New()
	payload := ptmodel.PocketNew{
		PocketName: "example pocket",
		Editor:     []uuid.UUID{},
		Watcher:    []uuid.UUID{},
		Icon:       1,
	}

	timeNow := time.Now()
	expect := ptmodel.PocketResp{
		ID:         1,
		Owner:      ownerUUID,
		Editor:     []uuid.UUID{ownerUUID},
		Watcher:    []uuid.UUID{ownerUUID},
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
		[]usermodel.User{
			{
				ID: ownerUUID,
			},
		}, nil,
	)

	// mock pocket
	pocketReplacePtr := ptmodel.Pocket{
		ID:         1,
		Owner:      ownerUUID,
		Editor:     []uuid.UUID{ownerUUID},
		Watcher:    []uuid.UUID{ownerUUID},
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

	// init service
	service := NewService(log, pocketRepo, userRepo)

	// Assertion
	result, err := service.CreatePocket(ctx, ownerUUID, payload)
	assert.Nil(t, err)
	assert.Equal(t, expect, result)
}
