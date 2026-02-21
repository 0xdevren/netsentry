package plugins

import (
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/policy"
)

// RulePlugin is the interface that all external and built-in rule plugins must implement.
// Plugin authors implement this interface and register their plugin with the engine.
type RulePlugin interface {
	// Name returns the unique identifier of the plugin.
	Name() string
	// Validate evaluates the plugin's rule logic against the parsed configuration
	// and returns a slice of ValidationResult values.
	Validate(cfg *model.ConfigModel) []policy.ValidationResult
}

// Registry maintains the set of registered RulePlugin implementations.
type Registry struct {
	plugins map[string]RulePlugin
}

// NewRegistry constructs an empty plugin Registry.
func NewRegistry() *Registry {
	return &Registry{plugins: make(map[string]RulePlugin)}
}

// Register adds a plugin to the registry. An existing registration for the
// same name is silently replaced.
func (r *Registry) Register(p RulePlugin) {
	r.plugins[p.Name()] = p
}

// Get returns the plugin with the given name and a boolean indicating existence.
func (r *Registry) Get(name string) (RulePlugin, bool) {
	p, ok := r.plugins[name]
	return p, ok
}

// All returns all registered plugins.
func (r *Registry) All() []RulePlugin {
	out := make([]RulePlugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		out = append(out, p)
	}
	return out
}

// SNMPCommunityPlugin is a built-in example plugin that checks for insecure
// SNMP community strings.
type SNMPCommunityPlugin struct{}

// Name returns the plugin identifier.
func (s *SNMPCommunityPlugin) Name() string { return "snmp-community-check" }

// Validate checks that no line in the config declares a "public" SNMP community.
func (s *SNMPCommunityPlugin) Validate(cfg *model.ConfigModel) []policy.ValidationResult {
	var results []policy.ValidationResult
	for _, line := range cfg.Lines {
		if containsStr(line, "snmp-server community public") ||
			containsStr(line, "snmp-server community private") {
			results = append(results, policy.ValidationResult{
				RuleID:          "PLUGIN-SNMP-001",
				RuleDescription: "SNMP community string must not use default insecure values",
				Device:          cfg.Device,
				Status:          policy.StatusFail,
				Severity:        policy.SeverityHigh,
				Message:         "insecure SNMP community string detected: " + line,
				Remediation:     "Replace default SNMP community strings with unique, randomly generated values.",
			})
		}
	}
	if len(results) == 0 {
		results = append(results, policy.ValidationResult{
			RuleID:          "PLUGIN-SNMP-001",
			RuleDescription: "SNMP community string must not use default insecure values",
			Device:          cfg.Device,
			Status:          policy.StatusPass,
			Severity:        policy.SeverityHigh,
			Message:         "no insecure SNMP community strings detected",
		})
	}
	return results
}

func containsStr(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(s) < len(sub) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
