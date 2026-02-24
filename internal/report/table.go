package report

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/0xdevren/netsentry/internal/policy"
)

// TableReporter renders the report as a formatted ASCII table.
type TableReporter struct {
	writer  io.Writer
	noColor bool
}

// NewTableReporter constructs a TableReporter.
func NewTableReporter(opts Options) *TableReporter {
	return &TableReporter{writer: resolveWriter(opts), noColor: opts.NoColor}
}

// Generate writes the formatted table report.
func (r *TableReporter) Generate(report *policy.Report) error {
	w := r.writer

	// Use buffered writer for better I/O performance.
	buf := bufio.NewWriter(w)
	defer buf.Flush()

	// Header.
	fmt.Fprintf(buf, "\nDEVICE : %s\n", report.Device.String())
	fmt.Fprintf(buf, "POLICY : %s", report.Policy)
	if report.PolicyVersion != "" {
		fmt.Fprintf(buf, " (v%s)", report.PolicyVersion)
	}
	fmt.Fprintln(buf)
	fmt.Fprintln(buf, strings.Repeat("-", 72))

	// Table.
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"RULE-ID", "STATUS", "SEVERITY", "MESSAGE"})
	table.SetBorder(false)
	table.SetHeaderLine(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator(" ")
	table.SetAutoWrapText(true)
	table.SetColWidth(50)

	for _, res := range report.Results {
		statusStr := r.formatStatus(res.Status)
		severityStr := r.formatSeverity(res.Severity)
		msg := res.Message
		if len(msg) > 60 {
			msg = msg[:57] + "..."
		}
		table.Append([]string{res.RuleID, statusStr, severityStr, msg})
	}
	table.Render()

	// Summary.
	fmt.Fprintln(buf, strings.Repeat("-", 72))
	fmt.Fprintln(buf, "SUMMARY:")
	s := report.Summary
	fmt.Fprintf(buf, "  Passed   : %d\n", s.Passed)
	fmt.Fprintf(buf, "  Failed   : %d\n", s.Failed)
	fmt.Fprintf(buf, "  Warnings : %d\n", s.Warnings)
	fmt.Fprintf(buf, "  Skipped  : %d\n", s.Skipped)
	fmt.Fprintf(buf, "  Score    : %.0f%%\n", s.Score)
	fmt.Fprintln(buf)

	return nil
}

func (r *TableReporter) formatStatus(s policy.ValidationStatus) string {
	if r.noColor || os.Getenv("NO_COLOR") != "" {
		return string(s)
	}
	switch s {
	case policy.StatusPass:
		return color.GreenString(string(s))
	case policy.StatusFail:
		return color.RedString(string(s))
	case policy.StatusWarn:
		return color.YellowString(string(s))
	case policy.StatusSkip:
		return color.CyanString(string(s))
	default:
		return color.MagentaString(string(s))
	}
}

func (r *TableReporter) formatSeverity(s policy.Severity) string {
	if r.noColor || os.Getenv("NO_COLOR") != "" {
		return string(s)
	}
	switch s {
	case policy.SeverityCritical:
		return color.New(color.FgRed, color.Bold).Sprint(string(s))
	case policy.SeverityHigh:
		return color.RedString(string(s))
	case policy.SeverityMedium:
		return color.YellowString(string(s))
	case policy.SeverityLow:
		return color.CyanString(string(s))
	default:
		return color.WhiteString(string(s))
	}
}
