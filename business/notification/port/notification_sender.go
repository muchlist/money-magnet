package port

import (
	"context"

	"github.com/muchlist/moneymagnet/pkg/mfirebase"
)

type FCMSender interface {
	SendMessage(ctx context.Context, payload mfirebase.Payload) ([]string /*invalid token*/, error)
}
