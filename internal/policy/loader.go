package policy

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Loader reads and deserialises policy definitions from YAML files.
type Loader struct{}

// NewLoader constructs a new Loader instance.
func NewLoader() *Loader {
	return &Loader{}
}

// LoadFile reads the policy at the given filesystem path and returns the parsed Policy.
func (l *Loader) LoadFile(path string) (*Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("policy loader: read file %q: %w", path, err)
	}
	return l.LoadBytes(data)
}

// LoadBytes parses YAML-encoded policy bytes and returns the parsed Policy.
func (l *Loader) LoadBytes(data []byte) (*Policy, error) {
	var p Policy
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("policy loader: yaml unmarshal: %w", err)
	}
	if err := l.validate(&p); err != nil {
		return nil, fmt.Errorf("policy loader: validation: %w", err)
	}
	return &p, nil
}

// validate checks the structural integrity of a parsed Policy.
func (l *Loader) validate(p *Policy) error {
	if p.Name == "" {
		return fmt.Errorf("policy name is required")
	}
	seen := make(map[string]struct{}, len(p.Rules))
	for i, r := range p.Rules {
		if r.ID == "" {
			return fmt.Errorf("rule at index %d has no id", i)
		}
		if _, dup := seen[r.ID]; dup {
			return fmt.Errorf("duplicate rule id %q at index %d", r.ID, i)
		}
		seen[r.ID] = struct{}{}
		if !r.Severity.IsValid() {
			return fmt.Errorf("rule %q has invalid severity %q", r.ID, r.Severity)
		}
	}
	return nil
}
