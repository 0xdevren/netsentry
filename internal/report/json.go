package report

import (
	"fmt"
	"io"
	"os"

	"encoding/json"

	"github.com/0xdevren/netsentry/internal/policy"
	"gopkg.in/yaml.v3"
)

// JSONReporter serialises the Report as indented JSON.
type JSONReporter struct {
	writer io.Writer
}

// NewJSONReporter constructs a JSONReporter writing to opts.Writer, opts.OutputPath, or stdout.
func NewJSONReporter(opts Options) *JSONReporter {
	return &JSONReporter{writer: resolveWriter(opts)}
}

// Generate encodes the report as indented JSON.
func (r *JSONReporter) Generate(report *policy.Report) error {
	enc := json.NewEncoder(r.writer)
	enc.SetIndent("", "  ")
	if err := enc.Encode(report); err != nil {
		return fmt.Errorf("json reporter: encode: %w", err)
	}
	return nil
}

// YAMLReporter serialises the Report as YAML.
type YAMLReporter struct {
	writer io.Writer
}

// NewYAMLReporter constructs a YAMLReporter.
func NewYAMLReporter(opts Options) *YAMLReporter {
	return &YAMLReporter{writer: resolveWriter(opts)}
}

// Generate encodes the report as YAML.
func (r *YAMLReporter) Generate(report *policy.Report) error {
	data, err := yaml.Marshal(report)
	if err != nil {
		return fmt.Errorf("yaml reporter: marshal: %w", err)
	}
	_, err = r.writer.Write(data)
	return err
}

// marshalYAML is used by the YAML reporter helper (stub kept for any future use).
func marshalYAML(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// resolveWriter returns the appropriate io.Writer from Options.
func resolveWriter(opts Options) io.Writer {
	if opts.Writer != nil {
		return opts.Writer
	}
	if opts.OutputPath != "" {
		f, err := os.Create(opts.OutputPath)
		if err == nil {
			return f
		}
	}
	return os.Stdout
}
