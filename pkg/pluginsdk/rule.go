package pluginsdk

import (
	"github.com/0xdevren/netsentry/internal/model"
	"github.com/0xdevren/netsentry/internal/policy"
)

// RulePlugin is implemented by external rule plugins that add custom validation
// logic beyond what the YAML policy engine supports. Register with the plugin
// registry using plugins.NewRegistry().Register(plugin).
type RulePlugin interface {
	// Name returns the unique plugin identifier.
	Name() string
	// Validate evaluates the plugin's rule set against the parsed config and
	// returns any validation results.
	Validate(cfg *model.ConfigModel) []policy.ValidationResult
}
