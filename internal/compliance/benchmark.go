package compliance

import (
	"fmt"
	"github.com/0xdevren/netsentry/internal/policy"
)

// BenchmarkLevel defines a CIS benchmark level.
type BenchmarkLevel int

const (
	// LevelOne is the basic CIS benchmark profile (minimal operational impact).
	LevelOne BenchmarkLevel = 1
	// LevelTwo is the strict CIS benchmark profile (maximum security).
	LevelTwo BenchmarkLevel = 2
)

// BenchmarkResult represents compliance against a named benchmark.
type BenchmarkResult struct {
	// DeviceID is the assessed device.
	DeviceID string `json:"device_id"`
	// Benchmark is the benchmark name (e.g. "CIS-IOS-L1").
	Benchmark string `json:"benchmark"`
	// Level is the benchmark level.
	Level BenchmarkLevel `json:"level"`
	// Score is the compliance percentage.
	Score float64 `json:"score"`
	// PassedControls is the number of passing controls.
	PassedControls int `json:"passed_controls"`
	// TotalControls is the total number of evaluated controls.
	TotalControls int `json:"total_controls"`
}

// BenchmarkEvaluator assesses compliance against a named CIS benchmark.
type BenchmarkEvaluator struct {
	// Name is the benchmark identifier.
	Name string
	// Level is the benchmark level.
	Level BenchmarkLevel
}

// NewCISIOSLevel1 returns a CIS IOS Level 1 benchmark evaluator.
func NewCISIOSLevel1() *BenchmarkEvaluator {
	return &BenchmarkEvaluator{Name: "CIS-IOS-L1", Level: LevelOne}
}

// NewCISIOSLevel2 returns a CIS IOS Level 2 benchmark evaluator.
func NewCISIOSLevel2() *BenchmarkEvaluator {
	return &BenchmarkEvaluator{Name: "CIS-IOS-L2", Level: LevelTwo}
}

// Evaluate assesses a validation report against this benchmark profile.
func (b *BenchmarkEvaluator) Evaluate(report *policy.Report) (*BenchmarkResult, error) {
	if report == nil {
		return nil, fmt.Errorf("benchmark: nil report")
	}
	return &BenchmarkResult{
		DeviceID:       report.Device.ID,
		Benchmark:      b.Name,
		Level:          b.Level,
		Score:          report.Summary.Score,
		PassedControls: report.Summary.Passed,
		TotalControls:  report.Summary.Total,
	}, nil
}
