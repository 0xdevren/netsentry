package report

import "fmt"

// PDFReporter is defined in html.go. This file is intentionally empty to
// satisfy the directory placeholder. The PDF implementation delegates to
// HTMLReporter for print-to-PDF workflows.
var _ = fmt.Sprintf // ensure package is importable
