// Package dsl provides the parsing and validation layer for the NetSentry
// policy domain-specific language expressed as YAML.
package dsl

// grammar.go defines the structural grammar constants and helpers for
// validating policy file structure beyond basic YAML parsing.

// SupportedMatchKeys enumerates the valid keys within a match block.
var SupportedMatchKeys = []string{
	"contains",
	"not_contains",
	"regex",
	"required_block",
}

// SupportedActionKeys enumerates the valid keys within an action block.
var SupportedActionKeys = []string{
	"deny",
	"warn",
	"remediation",
}

// ValidSeverities lists valid severity values accepted in policy documents.
var ValidSeverities = []string{
	"CRITICAL",
	"HIGH",
	"MEDIUM",
	"LOW",
	"INFO",
}

// containsKey reports whether a string slice contains the given value.
func containsKey(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
