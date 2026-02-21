package validator

import (
	"context"
	"fmt"

	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/policy"
	"github.com/0xdevren/netsentry/internal/util"
)

// TopologyValidationRequest encapsulates a multi-device topology validation run.
type TopologyValidationRequest struct {
	// Graph is the topology graph to validate.
	Graph *model.TopologyGraph
	// Policy is the policy applied to each device in the graph.
	Policy *policy.Policy
	// Concurrency is the number of parallel validation workers.
	Concurrency int
}

// TopologyValidationResult is the aggregate result of a topology-wide validation.
type TopologyValidationResult struct {
	// DeviceReports maps device ID to its individual validation report.
	DeviceReports map[string]*policy.Report
	// Errors maps device ID to any error encountered during validation.
	Errors map[string]error
}

// TopologyValidator validates all devices in a topology graph in parallel.
type TopologyValidator struct {
	pool *util.WorkerPool[topologyJob, topologyResult]
}

type topologyJob struct {
	deviceID string
	config   *model.ConfigModel
	pol      *policy.Policy
}

type topologyResult struct {
	deviceID string
	report   *policy.Report
	err      error
}

// NewTopologyValidator constructs a TopologyValidator with the given concurrency.
func NewTopologyValidator(concurrency int) *TopologyValidator {
	if concurrency <= 0 {
		concurrency = 4
	}
	pool := util.NewWorkerPool[topologyJob, topologyResult](concurrency, func(ctx context.Context, j topologyJob) topologyResult {
		dv, err := NewDeviceValidator(DeviceValidatorOptions{Concurrency: 1})
		if err != nil {
			return topologyResult{deviceID: j.deviceID, err: fmt.Errorf("create validator: %w", err)}
		}
		report, err := dv.Validate(ctx, ValidationRequest{
			Config: j.config,
			Policy: j.pol,
		})
		return topologyResult{deviceID: j.deviceID, report: report, err: err}
	})
	return &TopologyValidator{pool: pool}
}

// Validate runs validation against every device in the topology graph.
func (tv *TopologyValidator) Validate(ctx context.Context, req TopologyValidationRequest) (*TopologyValidationResult, error) {
	if req.Graph == nil {
		return nil, fmt.Errorf("topology validator: TopologyGraph is required")
	}

	var jobs []topologyJob
	for id, dev := range req.Graph.Devices {
		cfg := &model.ConfigModel{Device: dev, GlobalSettings: make(map[string]string)}
		jobs = append(jobs, topologyJob{deviceID: id, config: cfg, pol: req.Policy})
	}

	rawResults := tv.pool.Run(ctx, jobs)

	result := &TopologyValidationResult{
		DeviceReports: make(map[string]*policy.Report),
		Errors:        make(map[string]error),
	}
	for _, r := range rawResults {
		if r.err != nil {
			result.Errors[r.deviceID] = r.err
		} else {
			result.DeviceReports[r.deviceID] = r.report
		}
	}

	return result, nil
}
