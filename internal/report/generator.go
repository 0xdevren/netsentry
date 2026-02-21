// Package report implements the report generation subsystem for NetSentry.
// All reporters implement the Reporter interface.
package report

import (
	"fmt"
	"io"

	"github.com/0xdevren/netsentry/internal/policy"
)

// Reporter is the interface implemented by all output format generators.
type Reporter interface {
	// Generate writes the report to the configured destination.
	Generate(report *policy.Report) error
}

// Format enumerates the supported output formats.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
	FormatHTML  Format = "html"
	FormatPDF   Format = "pdf"
)

// Options configures a reporter instance.
type Options struct {
	// Format is the output format.
	Format Format
	// Writer is the destination for text-based reporters. Defaults to os.Stdout.
	Writer io.Writer
	// OutputPath is the destination file path for file-based reporters (HTML, PDF).
	OutputPath string
	// NoColor disables ANSI colour codes in table output.
	NoColor bool
}

// New constructs the Reporter appropriate for the given Options.
func New(opts Options) (Reporter, error) {
	switch opts.Format {
	case FormatTable, "":
		return NewTableReporter(opts), nil
	case FormatJSON:
		return NewJSONReporter(opts), nil
	case FormatYAML:
		return NewYAMLReporter(opts), nil
	case FormatHTML:
		return NewHTMLReporter(opts)
	case FormatPDF:
		return NewPDFReporter(opts)
	default:
		return nil, fmt.Errorf("report: unsupported format %q", opts.Format)
	}
}
