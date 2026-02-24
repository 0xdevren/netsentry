package policy

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/0xdevren/netsentry/internal/model"
)

// Matcher evaluates a MatchSpec against a ConfigModel.
type Matcher struct {
	// cache maps compiled regular expressions keyed on the pattern string.
	cache map[string]*regexp.Regexp
}

// NewMatcher constructs a new Matcher with an initialised regex cache.
func NewMatcher() *Matcher {
	return &Matcher{cache: make(map[string]*regexp.Regexp)}
}

// Match evaluates the provided MatchSpec against the given ConfigModel and
// returns true when the match condition is satisfied.
func (m *Matcher) Match(spec MatchSpec, cfg *model.ConfigModel) (bool, error) {
	switch {
	case spec.Contains != "":
		return m.matchContains(spec.Contains, cfg), nil
	case spec.NotContains != "":
		return m.matchNotContains(spec.NotContains, cfg), nil
	case spec.Regex != "":
		return m.matchRegex(spec.Regex, cfg)
	case spec.RequiredBlock != "":
		return m.matchRequiredBlock(spec.RequiredBlock, cfg), nil
	default:
		return false, fmt.Errorf("matcher: MatchSpec has no defined condition")
	}
}

// matchContains returns true if any configuration line contains the substring.
func (m *Matcher) matchContains(text string, cfg *model.ConfigModel) bool {
	for _, line := range cfg.Lines {
		if strings.Contains(line, text) {
			return true
		}
	}
	return false
}

// matchNotContains returns true if no configuration line contains the substring.
func (m *Matcher) matchNotContains(text string, cfg *model.ConfigModel) bool {
	for _, line := range cfg.Lines {
		if strings.Contains(line, text) {
			return false
		}
	}
	return true
}

// matchRegex returns true if any configuration line matches the regular expression.
func (m *Matcher) matchRegex(pattern string, cfg *model.ConfigModel) (bool, error) {
	re, err := m.compileRegex(pattern)
	if err != nil {
		return false, err
	}
	for _, line := range cfg.Lines {
		if re.MatchString(line) {
			return true, nil
		}
	}
	return false, nil
}

// matchRequiredBlock returns true if any configuration line begins with the block prefix.
func (m *Matcher) matchRequiredBlock(prefix string, cfg *model.ConfigModel) bool {
	for _, line := range cfg.Lines {
		if len(line) >= len(prefix) && line[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// compileRegex returns a compiled regex from cache, compiling and caching it on first use.
func (m *Matcher) compileRegex(pattern string) (*regexp.Regexp, error) {
	if re, ok := m.cache[pattern]; ok {
		return re, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("matcher: invalid regex %q: %w", pattern, err)
	}
	m.cache[pattern] = re
	return re, nil
}
