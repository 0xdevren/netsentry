package sdk

// Report mirrors the policy.Report struct for SDK consumers who do not import internal packages.
type Report struct {
	Device        Device          `json:"device"`
	Policy        string          `json:"policy"`
	PolicyVersion string          `json:"policy_version,omitempty"`
	Results       []Result        `json:"results"`
	Summary       ReportSummary   `json:"summary"`
}

// Device is the SDK representation of a network device.
type Device struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	Type     string `json:"type"`
}

// Result is the SDK representation of a single rule evaluation outcome.
type Result struct {
	RuleID      string `json:"rule_id"`
	Status      string `json:"status"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Remediation string `json:"remediation,omitempty"`
}

// ReportSummary is the SDK representation of compliance summary metrics.
type ReportSummary struct {
	Total    int     `json:"total"`
	Passed   int     `json:"passed"`
	Failed   int     `json:"failed"`
	Warnings int     `json:"warnings"`
	Skipped  int     `json:"skipped"`
	Score    float64 `json:"score"`
}
