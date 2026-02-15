package concurrency

import (
	"context"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
)

const maxAutoWorkers = 16

// AutoWorkers resolves worker count with sane defaults.
// If requested <= 0, it uses min(2*NumCPU, 16).
func AutoWorkers(requested int) int {
	if requested > 0 {
		return requested
	}

	workers := runtime.NumCPU() * 2
	if workers < 1 {
		workers = 1
	}
	if workers > maxAutoWorkers {
		workers = maxAutoWorkers
	}
	return workers
}

// ClampWorkers ensures worker count is at least 1 and at most total when total > 0.
func ClampWorkers(requested, total int) int {
	workers := AutoWorkers(requested)
	if total > 0 && workers > total {
		workers = total
	}
	if workers < 1 {
		workers = 1
	}
	return workers
}

// MapOrdered processes items in parallel with bounded workers while preserving result order.
// If failFast is true, the first error cancels the pipeline and returns immediately.
// If failFast is false, processing continues and the first error is returned after all items finish.
func MapOrdered[T any, R any](
	ctx context.Context,
	items []T,
	workers int,
	failFast bool,
	fn func(context.Context, int, T) (R, error),
) ([]R, error) {
	results := make([]R, len(items))
	if len(items) == 0 {
		return results, nil
	}

	workers = ClampWorkers(workers, len(items))

	g, gctx := errgroup.WithContext(ctx)
	type job struct {
		index int
		item  T
	}
	jobs := make(chan job)

	var (
		firstErr   error
		firstErrMu sync.Mutex
	)
	setFirstErr := func(err error) {
		if err == nil {
			return
		}
		firstErrMu.Lock()
		defer firstErrMu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}

	for worker := 0; worker < workers; worker++ {
		g.Go(func() error {
			for {
				select {
				case <-gctx.Done():
					// In fail-fast mode, context cancellation is expected after first error.
					// Returning nil here prevents context.Canceled from masking the real error.
					if failFast {
						return nil
					}
					return gctx.Err()
				case currentJob, ok := <-jobs:
					if !ok {
						return nil
					}

					result, err := fn(gctx, currentJob.index, currentJob.item)
					if err != nil {
						if failFast {
							return err
						}
						setFirstErr(err)
						continue
					}
					results[currentJob.index] = result
				}
			}
		})
	}

	g.Go(func() error {
		defer close(jobs)
		for i, item := range items {
			select {
			case <-gctx.Done():
				if failFast {
					return nil
				}
				return gctx.Err()
			case jobs <- job{
				index: i,
				item:  item,
			}:
			}
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	// Preserve external cancellation signal even if fail-fast mode suppressed context errors.
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return results, nil
}
