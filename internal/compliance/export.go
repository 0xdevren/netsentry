package compliance

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ExportFormat defines the serialisation format for compliance export.
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportYAML ExportFormat = "yaml"
)

// Export writes a BaselineEntry slice to the given file path in the specified format.
func Export(entries []BaselineEntry, path string, format ExportFormat) error {
	var data []byte
	var err error

	switch format {
	case ExportJSON:
		data, err = json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return fmt.Errorf("compliance export: json: %w", err)
		}
	case ExportYAML:
		data, err = yaml.Marshal(entries)
		if err != nil {
			return fmt.Errorf("compliance export: yaml: %w", err)
		}
	default:
		return fmt.Errorf("compliance export: unsupported format %q", format)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("compliance export: write %q: %w", path, err)
	}
	return nil
}
