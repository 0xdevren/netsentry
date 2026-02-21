package model

// ConfigModel is the canonical, vendor-neutral structured representation
// of a parsed network device configuration. It is the primary input to
// the validation pipeline.
type ConfigModel struct {
	// Device contains identity and metadata for the device.
	Device Device `json:"device" yaml:"device"`
	// RawText is the original configuration text before parsing.
	RawText string `json:"raw_text,omitempty" yaml:"raw_text,omitempty"`
	// Interfaces is the list of logical and physical interfaces.
	Interfaces []Interface `json:"interfaces,omitempty" yaml:"interfaces,omitempty"`
	// ACLs is the list of access control lists defined on the device.
	ACLs []ACL `json:"acls,omitempty" yaml:"acls,omitempty"`
	// BGPConfig holds BGP protocol configuration, if present.
	BGPConfig *BGPConfig `json:"bgp,omitempty" yaml:"bgp,omitempty"`
	// OSPFConfig holds OSPF protocol configuration, if present.
	OSPFConfig *OSPFConfig `json:"ospf,omitempty" yaml:"ospf,omitempty"`
	// StaticRoutes is the list of static routing entries.
	StaticRoutes []StaticRoute `json:"static_routes,omitempty" yaml:"static_routes,omitempty"`
	// VLANs is the list of VLANs configured on the device.
	VLANs []VLAN `json:"vlans,omitempty" yaml:"vlans,omitempty"`
	// GlobalSettings holds key-value pairs for global configuration items
	// that do not map to a structured sub-model (e.g. hostname, logging servers).
	GlobalSettings map[string]string `json:"global_settings,omitempty" yaml:"global_settings,omitempty"`
	// Lines holds the raw configuration lines for regex/contains matching.
	Lines []string `json:"lines,omitempty" yaml:"lines,omitempty"`
}

// HasLine reports whether the configuration contains the given exact line.
func (c *ConfigModel) HasLine(line string) bool {
	for _, l := range c.Lines {
		if l == line {
			return true
		}
	}
	return false
}

// ContainsText reports whether any line in the configuration contains the given substring.
func (c *ConfigModel) ContainsText(text string) bool {
	for _, l := range c.Lines {
		if contains(l, text) {
			return true
		}
	}
	return false
}

// contains is an inlined string containment check to avoid importing strings
// in the hot validation path.
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
