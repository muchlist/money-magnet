package port

import "context"

type ETagStorer interface {
	GetTagByPocketID(ctx context.Context, pocketID string) (int64, error)
	SetTagByPocketID(ctx context.Context, pocketID string, value int64) error
}
