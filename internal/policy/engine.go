package policy

import (
	"context"
	"fmt"
	"sync"

	"github.com/0xdevren/netsentry/internal/model"
)

// EngineOptions configures the behaviour of the policy Engine.
type EngineOptions struct {
	// Concurrency is the number of parallel rule evaluation goroutines.
	// A value of 0 or less defaults to 4.
	Concurrency int
}

// Engine orchestrates concurrent rule evaluation across all rules in a policy.
type Engine struct {
	evaluator *Evaluator
	opts      EngineOptions
}

// NewEngine constructs an Engine with the given options.
func NewEngine(opts EngineOptions) *Engine {
	if opts.Concurrency <= 0 {
		opts.Concurrency = 4
	}
	return &Engine{
		evaluator: NewEvaluator(NewMatcher()),
		opts:      opts,
	}
}

// job is an internal unit of work for the worker pool.
type job struct {
	rule Rule
	cfg  *model.ConfigModel
}

// Run evaluates all enabled rules in the policy against the given configuration
// using a worker pool. It respects context cancellation.
func (e *Engine) Run(ctx context.Context, p *Policy, cfg *model.ConfigModel) ([]ValidationResult, error) {
	if p == nil {
		return nil, fmt.Errorf("engine: nil policy")
	}
	if cfg == nil {
		return nil, fmt.Errorf("engine: nil config model")
	}

	jobs := make(chan job, len(p.Rules))
	results := make(chan ValidationResult, len(p.Rules))

	var wg sync.WaitGroup
	for i := 0; i < e.opts.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}
					results <- e.evaluator.Evaluate(j.rule, j.cfg)
				}
			}
		}()
	}

	// Enqueue jobs.
	for _, rule := range p.Rules {
		select {
		case <-ctx.Done():
			close(jobs)
			wg.Wait()
			close(results)
			return nil, fmt.Errorf("engine: context cancelled before all jobs enqueued: %w", ctx.Err())
		case jobs <- job{rule: rule, cfg: cfg}:
		}
	}
	close(jobs)

	// Collect results.
	go func() {
		wg.Wait()
		close(results)
	}()

	var out []ValidationResult
	for r := range results {
		out = append(out, r)
	}

	if ctx.Err() != nil {
		return out, fmt.Errorf("engine: context error: %w", ctx.Err())
	}

	return out, nil
}
