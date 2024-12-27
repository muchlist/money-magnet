// package bg is used for background processes management.
// It provides functionality to run background jobs, handle panic recovery,
// and ensure graceful shutdown by waiting for all background processes to finish.
package bg

import (
	"context"
	"log"
)

type BackgroundJob struct {
	JobTitle string
	Execute  func(ctx context.Context)
}

func RunSafeBackground(ctx context.Context, job BackgroundJob) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recover panic when run %s : %v", job.JobTitle, err)
			}
		}()

		detachedCtx := NewDetachContext(ctx)
		job.Execute(detachedCtx)
	}()
}
