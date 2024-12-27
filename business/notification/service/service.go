package service

import (
	"context"
	"fmt"

	"github.com/muchlist/moneymagnet/business/notification/model"
	"github.com/muchlist/moneymagnet/business/notification/port"
	"github.com/muchlist/moneymagnet/pkg/mfirebase"
	"github.com/muchlist/moneymagnet/pkg/mlogger"
	"github.com/muchlist/moneymagnet/pkg/observ"
)

// Core manages the set of APIs for user access.
type Core struct {
	log       mlogger.Logger
	fcmSender port.FCMSender
	userRepo  port.UserStorer
}

// NewCore constructs a core for user api access.
func NewCore(
	log mlogger.Logger,
	fcmSender port.FCMSender,
	userRepo port.UserStorer,
) *Core {
	return &Core{
		log:       log,
		fcmSender: fcmSender,
		userRepo:  userRepo,
	}
}

func (s *Core) SendNotificationToUser(ctx context.Context, payload model.SendMessage) error {
	ctx, span := observ.GetTracer().Start(ctx, "service-SendNotificationToUser")
	defer span.End()

	// get all token from users
	users, err := s.userRepo.GetByIDs(ctx, payload.UserIds)
	if err != nil {
		return fmt.Errorf("get users by ids: %w", err)
	}

	allTokens := make([]string, 0)
	for _, user := range users {
		allTokens = append(allTokens, user.Fcm...)
	}

	failedTokens, err := s.fcmSender.SendMessage(ctx, mfirebase.Payload{
		Title:          payload.Title,
		Message:        payload.Message,
		ReceiverTokens: allTokens,
	})

	if err != nil {
		return fmt.Errorf("send message failed: %w", err)
	}

	if len(failedTokens) != 0 {
		failedTokenMap := make(map[string]struct{})
		for _, token := range failedTokens {
			failedTokenMap[token] = struct{}{}
		}

		for _, user := range users {
			var newTokens []string
			for _, token := range user.Fcm {
				if _, found := failedTokenMap[token]; !found {
					newTokens = append(newTokens, token)
				}
			}

			// if any token changes, update fcm user
			if len(newTokens) != len(user.Fcm) {
				err := s.userRepo.EditFCM(ctx, user.ID, newTokens)
				if err != nil {
					return fmt.Errorf("failed modify fcm after invalid tokens detected: %w", err)
				}
			}
		}
	}

	return nil
}
