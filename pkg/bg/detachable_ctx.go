package bg

import (
	"context"
	"time"
)

// NewDetachContext create new context with no cancelation
// but stil have values from parent context
func NewDetachContext(parent context.Context) context.Context {
	return Detach{Ctx: parent}
}

// Detach implement context.Context but with no cancelation
type Detach struct {
	Ctx context.Context
}

func (d Detach) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (d Detach) Done() <-chan struct{} {
	return nil
}

func (d Detach) Err() error {
	return nil
}

func (d Detach) Value(key any) any {
	return d.Ctx.Value(key)
}
