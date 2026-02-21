package report

import (
	"fmt"
	"html/template"
	"os"

	"github.com/0xdevren/netsentry/internal/policy"
)

// HTMLReporter renders the report as a self-contained HTML document.
type HTMLReporter struct {
	opts Options
}

// NewHTMLReporter constructs an HTMLReporter.
func NewHTMLReporter(opts Options) (*HTMLReporter, error) {
	if opts.OutputPath == "" && opts.Writer == nil {
		return nil, fmt.Errorf("html reporter: OutputPath or Writer is required")
	}
	return &HTMLReporter{opts: opts}, nil
}

// Generate renders the report as HTML.
func (r *HTMLReporter) Generate(report *policy.Report) error {
	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"statusClass": func(s policy.ValidationStatus) string {
			switch s {
			case policy.StatusPass:
				return "pass"
			case policy.StatusFail:
				return "fail"
			case policy.StatusWarn:
				return "warn"
			default:
				return "skip"
			}
		},
		"severityClass": func(s policy.Severity) string {
			switch s {
			case policy.SeverityCritical:
				return "critical"
			case policy.SeverityHigh:
				return "high"
			case policy.SeverityMedium:
				return "medium"
			case policy.SeverityLow:
				return "low"
			default:
				return "info"
			}
		},
		"scoreColor": func(score float64) string {
			if score >= 90 {
				return "#28a745"
			}
			if score >= 70 {
				return "#ffc107"
			}
			return "#dc3545"
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("html reporter: parse template: %w", err)
	}

	w := resolveWriter(r.opts)
	if closer, ok := w.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	if err := tmpl.Execute(w, report); err != nil {
		return fmt.Errorf("html reporter: execute template: %w", err)
	}
	return nil
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>NetSentry Report - {{.Device.Hostname}}</title>
<style>
  body { font-family: 'Segoe UI', sans-serif; background: #0f1117; color: #e0e0e0; margin: 0; padding: 2rem; }
  h1 { color: #00bfff; font-size: 1.8rem; margin-bottom: 0.25rem; }
  .meta { color: #888; font-size: 0.9rem; margin-bottom: 2rem; }
  .summary { display: flex; gap: 1rem; flex-wrap: wrap; margin-bottom: 2rem; }
  .stat { background: #1a1d27; border-radius: 8px; padding: 1rem 1.5rem; min-width: 100px; text-align: center; }
  .stat .label { font-size: 0.75rem; color: #888; text-transform: uppercase; }
  .stat .value { font-size: 1.8rem; font-weight: bold; margin-top: 0.25rem; }
  .score-value { color: {{scoreColor .Summary.Score}}; }
  table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
  th { background: #1a1d27; padding: 0.75rem 1rem; text-align: left; font-size: 0.75rem;
       text-transform: uppercase; color: #888; border-bottom: 1px solid #2a2d3a; }
  td { padding: 0.75rem 1rem; border-bottom: 1px solid #1a1d27; font-size: 0.875rem; }
  tr:hover td { background: #1a1d27; }
  .pass { color: #28a745; font-weight: 600; }
  .fail { color: #dc3545; font-weight: 600; }
  .warn { color: #ffc107; font-weight: 600; }
  .skip { color: #6c757d; font-weight: 600; }
  .critical { color: #dc3545; font-weight: 700; }
  .high { color: #fd7e14; font-weight: 600; }
  .medium { color: #ffc107; }
  .low { color: #20c997; }
  .info { color: #6c757d; }
  .badge { display: inline-block; padding: 0.2rem 0.5rem; border-radius: 4px; font-size: 0.75rem; }
</style>
</head>
<body>
<h1>NetSentry Compliance Report</h1>
<div class="meta">Device: <strong>{{.Device.Hostname}}</strong> &nbsp;|&nbsp;
  Policy: <strong>{{.Policy}}</strong>{{if .PolicyVersion}} v{{.PolicyVersion}}{{end}}</div>
<div class="summary">
  <div class="stat"><div class="label">Score</div>
    <div class="value score-value">{{printf "%.0f" .Summary.Score}}%</div></div>
  <div class="stat"><div class="label">Passed</div>
    <div class="value pass">{{.Summary.Passed}}</div></div>
  <div class="stat"><div class="label">Failed</div>
    <div class="value fail">{{.Summary.Failed}}</div></div>
  <div class="stat"><div class="label">Warnings</div>
    <div class="value warn">{{.Summary.Warnings}}</div></div>
  <div class="stat"><div class="label">Skipped</div>
    <div class="value skip">{{.Summary.Skipped}}</div></div>
  <div class="stat"><div class="label">Total</div>
    <div class="value">{{.Summary.Total}}</div></div>
</div>
<table>
  <thead>
    <tr><th>Rule ID</th><th>Status</th><th>Severity</th><th>Message</th><th>Remediation</th></tr>
  </thead>
  <tbody>
    {{range .Results}}
    <tr>
      <td>{{.RuleID}}</td>
      <td><span class="{{statusClass .Status}}">{{.Status}}</span></td>
      <td><span class="{{severityClass .Severity}}">{{.Severity}}</span></td>
      <td>{{.Message}}</td>
      <td>{{.Remediation}}</td>
    </tr>
    {{end}}
  </tbody>
</table>
</body>
</html>`

// PDFReporter writes a PDF report. In this implementation it produces an HTML
// file with a .pdf extension. Full PDF rendering requires a headless browser
// or wkhtmltopdf, which is documented in DEVELOPER_GUIDE.md.
type PDFReporter struct {
	opts Options
}

// NewPDFReporter constructs a PDFReporter.
func NewPDFReporter(opts Options) (*PDFReporter, error) {
	if opts.OutputPath == "" {
		opts.OutputPath = "report.pdf.html"
	}
	return &PDFReporter{opts: opts}, nil
}

// Generate writes an HTML document suitable for headless print-to-PDF.
func (r *PDFReporter) Generate(report *policy.Report) error {
	// Delegate to HTMLReporter with print-optimised path.
	htmlOpts := r.opts
	htmlOpts.Format = FormatHTML
	if htmlOpts.OutputPath != "" {
		htmlOpts.Writer, _ = os.Create(htmlOpts.OutputPath)
	}
	hr, err := NewHTMLReporter(htmlOpts)
	if err != nil {
		return fmt.Errorf("pdf reporter: %w", err)
	}
	return hr.Generate(report)
}
