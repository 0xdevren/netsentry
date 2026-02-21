package model

// BGPNeighbor represents a single BGP peer relationship.
type BGPNeighbor struct {
	// Address is the IP address of the peer.
	Address string `json:"address" yaml:"address"`
	// RemoteAS is the autonomous system number of the peer.
	RemoteAS int `json:"remote_as" yaml:"remote_as"`
	// Description is an optional description label.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	// Password indicates whether MD5 authentication is configured (value hidden).
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	// UpdateSource is the interface used as the BGP source.
	UpdateSource string `json:"update_source,omitempty" yaml:"update_source,omitempty"`
	// Shutdown indicates whether the neighbor is administratively disabled.
	Shutdown bool `json:"shutdown,omitempty" yaml:"shutdown,omitempty"`
	// NextHopSelf overrides next-hop to self for advertised prefixes.
	NextHopSelf bool `json:"next_hop_self,omitempty" yaml:"next_hop_self,omitempty"`
	// RouteMapIn is the inbound route-map applied to this neighbor.
	RouteMapIn string `json:"route_map_in,omitempty" yaml:"route_map_in,omitempty"`
	// RouteMapOut is the outbound route-map applied to this neighbor.
	RouteMapOut string `json:"route_map_out,omitempty" yaml:"route_map_out,omitempty"`
	// PrefixListIn is the inbound prefix-list name.
	PrefixListIn string `json:"prefix_list_in,omitempty" yaml:"prefix_list_in,omitempty"`
	// PrefixListOut is the outbound prefix-list name.
	PrefixListOut string `json:"prefix_list_out,omitempty" yaml:"prefix_list_out,omitempty"`
}

// BGPNetwork is a network advertised via BGP.
type BGPNetwork struct {
	// Prefix is the network prefix in CIDR notation.
	Prefix string `json:"prefix" yaml:"prefix"`
	// Mask is the subnet mask for the prefix (IPv4 only).
	Mask string `json:"mask,omitempty" yaml:"mask,omitempty"`
}

// BGPConfig holds the BGP protocol configuration for a device.
type BGPConfig struct {
	// LocalAS is the autonomous system number of the local device.
	LocalAS int `json:"local_as" yaml:"local_as"`
	// RouterID is the BGP router identifier.
	RouterID string `json:"router_id,omitempty" yaml:"router_id,omitempty"`
	// Neighbors is the list of configured BGP peers.
	Neighbors []BGPNeighbor `json:"neighbors,omitempty" yaml:"neighbors,omitempty"`
	// Networks is the list of networks advertised.
	Networks []BGPNetwork `json:"networks,omitempty" yaml:"networks,omitempty"`
}
