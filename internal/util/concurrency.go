package util

import (
	"context"
	"sync"
)

// WorkerPool executes a set of jobs concurrently using a fixed number of workers.
// It is generic over the job and result types.
type WorkerPool[J any, R any] struct {
	concurrency int
	worker      func(ctx context.Context, job J) R
}

// NewWorkerPool constructs a WorkerPool with the given concurrency and worker function.
// Concurrency values <= 0 default to 1.
func NewWorkerPool[J any, R any](concurrency int, worker func(ctx context.Context, job J) R) *WorkerPool[J, R] {
	if concurrency <= 0 {
		concurrency = 1
	}
	return &WorkerPool[J, R]{concurrency: concurrency, worker: worker}
}

// Run submits all jobs to the pool and returns the collected results. The
// function blocks until all workers complete or the context is cancelled.
func (p *WorkerPool[J, R]) Run(ctx context.Context, jobs []J) []R {
	jobCh := make(chan J, len(jobs))
	resultCh := make(chan R, len(jobs))

	var wg sync.WaitGroup
	for i := 0; i < p.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-jobCh:
					if !ok {
						return
					}
					resultCh <- p.worker(ctx, j)
				}
			}
		}()
	}

	for _, j := range jobs {
		select {
		case <-ctx.Done():
			break
		case jobCh <- j:
		}
	}
	close(jobCh)

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []R
	for r := range resultCh {
		results = append(results, r)
	}
	return results
}
