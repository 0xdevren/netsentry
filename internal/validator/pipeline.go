package validator

import (
	"context"
	"fmt"

	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/policy"
)

// PipelineStep is a single processing step in the validation pipeline.
type PipelineStep func(ctx context.Context, state *PipelineState) error

// PipelineState carries all mutable state through the validation pipeline.
type PipelineState struct {
	// RawConfig is the raw configuration bytes loaded from the source.
	RawConfig []byte
	// Config is the parsed and structured configuration model.
	Config *model.ConfigModel
	// Policy is the loaded and validated policy definition.
	Policy *policy.Policy
	// Report is the assembled validation report (populated by the last step).
	Report *policy.Report
	// Strict indicates whether warnings are treated as failures.
	Strict bool
}

// Pipeline chains multiple PipelineSteps sequentially.
type Pipeline struct {
	steps []PipelineStep
}

// NewPipeline constructs an empty Pipeline.
func NewPipeline() *Pipeline { return &Pipeline{} }

// AddStep appends a step to the pipeline.
func (p *Pipeline) AddStep(step PipelineStep) *Pipeline {
	p.steps = append(p.steps, step)
	return p
}

// Run executes all steps in order, passing the shared PipelineState. Execution
// stops immediately if any step returns a non-nil error.
func (p *Pipeline) Run(ctx context.Context, state *PipelineState) error {
	for i, step := range p.steps {
		if ctx.Err() != nil {
			return fmt.Errorf("pipeline: context cancelled before step %d: %w", i, ctx.Err())
		}
		if err := step(ctx, state); err != nil {
			return fmt.Errorf("pipeline: step %d: %w", i, err)
		}
	}
	return nil
}
