package dsl

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// RawRule is the raw YAML-decoded form of a policy rule before conversion.
type RawRule struct {
	ID          string                 `yaml:"id"`
	Description string                 `yaml:"description"`
	Severity    string                 `yaml:"severity"`
	Tags        []string               `yaml:"tags"`
	Match       map[string]interface{} `yaml:"match"`
	Action      map[string]interface{} `yaml:"action"`
	Enabled     *bool                  `yaml:"enabled"`
}

// RawPolicy is the raw YAML-decoded form of a policy file.
type RawPolicy struct {
	Name        string    `yaml:"name"`
	Version     string    `yaml:"version"`
	Description string    `yaml:"description"`
	Author      string    `yaml:"author"`
	Rules       []RawRule `yaml:"rules"`
}

// Parser parses raw YAML policy files into RawPolicy structures for
// subsequent validation and conversion.
type Parser struct{}

// NewParser constructs a new DSL Parser.
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile reads the file at path and returns a RawPolicy.
func (p *Parser) ParseFile(path string) (*RawPolicy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("dsl parser: read %q: %w", path, err)
	}
	return p.ParseBytes(data)
}

// ParseBytes parses YAML bytes into a RawPolicy.
func (p *Parser) ParseBytes(data []byte) (*RawPolicy, error) {
	var raw RawPolicy
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("dsl parser: yaml unmarshal: %w", err)
	}
	return &raw, nil
}
