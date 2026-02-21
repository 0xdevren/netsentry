package model

// ACLAction specifies whether an ACL entry permits or denies traffic.
type ACLAction string

const (
	ACLActionPermit ACLAction = "permit"
	ACLActionDeny   ACLAction = "deny"
)

// ACLEntry is a single rule within an access control list.
type ACLEntry struct {
	// Sequence is the sequence number of the entry.
	Sequence int `json:"sequence,omitempty" yaml:"sequence,omitempty"`
	// Action is permit or deny.
	Action ACLAction `json:"action" yaml:"action"`
	// Protocol is the IP protocol (e.g. "tcp", "udp", "ip", "icmp").
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	// Source is the source address or network in CIDR or wildcard notation.
	Source string `json:"source,omitempty" yaml:"source,omitempty"`
	// Destination is the destination address or network.
	Destination string `json:"destination,omitempty" yaml:"destination,omitempty"`
	// SourcePort is an optional source port or range.
	SourcePort string `json:"source_port,omitempty" yaml:"source_port,omitempty"`
	// DestPort is an optional destination port or range.
	DestPort string `json:"dest_port,omitempty" yaml:"dest_port,omitempty"`
	// Log indicates whether matched traffic is logged.
	Log bool `json:"log,omitempty" yaml:"log,omitempty"`
	// Remark is a free-text comment associated with the entry.
	Remark string `json:"remark,omitempty" yaml:"remark,omitempty"`
}

// ACL represents an access control list with its entries.
type ACL struct {
	// Name is the ACL identifier.
	Name string `json:"name" yaml:"name"`
	// Type is "standard" or "extended".
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Entries is the ordered list of ACL entries.
	Entries []ACLEntry `json:"entries,omitempty" yaml:"entries,omitempty"`
}
