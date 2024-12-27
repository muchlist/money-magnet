package port

import (
	"context"

	"github.com/muchlist/moneymagnet/business/notification/model"
)

type NotificationSender interface {
	SendNotificationToUser(ctx context.Context, payload model.SendMessage) error
}
