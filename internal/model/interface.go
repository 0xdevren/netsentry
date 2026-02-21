package model

// Interface represents a single network interface on a device.
type Interface struct {
	// Name is the interface identifier (e.g. "GigabitEthernet0/0", "eth0").
	Name string `json:"name" yaml:"name"`
	// Description is the operator-configured description string.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// IPAddress is the primary IPv4 address in CIDR notation.
	IPAddress string `json:"ip_address,omitempty" yaml:"ip_address,omitempty"`
	// IPv6Address is the primary IPv6 address in CIDR notation.
	IPv6Address string `json:"ipv6_address,omitempty" yaml:"ipv6_address,omitempty"`
	// SubnetMask is the subnet mask for the IP address (IPv4 only).
	SubnetMask string `json:"subnet_mask,omitempty" yaml:"subnet_mask,omitempty"`
	// Shutdown indicates whether the interface is administratively shut down.
	Shutdown bool `json:"shutdown" yaml:"shutdown"`
	// VLANMode is the switchport mode: "access", "trunk", or "routed".
	VLANMode string `json:"vlan_mode,omitempty" yaml:"vlan_mode,omitempty"`
	// AccessVLAN is the VLAN ID when the interface is in access mode.
	AccessVLAN int `json:"access_vlan,omitempty" yaml:"access_vlan,omitempty"`
	// TrunkAllowedVLANs lists allowed VLANs in trunk mode.
	TrunkAllowedVLANs []int `json:"trunk_allowed_vlans,omitempty" yaml:"trunk_allowed_vlans,omitempty"`
	// SpanningTreePortFast indicates whether portfast is enabled.
	SpanningTreePortFast bool `json:"stp_portfast,omitempty" yaml:"stp_portfast,omitempty"`
	// MTU is the configured maximum transmission unit.
	MTU int `json:"mtu,omitempty" yaml:"mtu,omitempty"`
	// Bandwidth is the configured interface bandwidth in kbps.
	Bandwidth int `json:"bandwidth,omitempty" yaml:"bandwidth,omitempty"`
	// InboundACL is the name of the ACL applied inbound on this interface.
	InboundACL string `json:"inbound_acl,omitempty" yaml:"inbound_acl,omitempty"`
	// OutboundACL is the name of the ACL applied outbound on this interface.
	OutboundACL string `json:"outbound_acl,omitempty" yaml:"outbound_acl,omitempty"`
	// Attributes holds additional vendor-specific key-value interface attributes.
	Attributes map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`
}
