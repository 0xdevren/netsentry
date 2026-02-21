package model

// StaticRoute represents a single static routing entry.
type StaticRoute struct {
	// Destination is the destination network prefix in CIDR notation.
	Destination string `json:"destination" yaml:"destination"`
	// NextHop is the next-hop IP address or interface.
	NextHop string `json:"next_hop" yaml:"next_hop"`
	// AdminDistance is the administrative distance of the route (0-255).
	AdminDistance int `json:"admin_distance,omitempty" yaml:"admin_distance,omitempty"`
	// Tag is an optional route tag for policy routing.
	Tag int `json:"tag,omitempty" yaml:"tag,omitempty"`
	// Name is an optional description for the route.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// Permanent indicates the route is not removed when the next-hop is unreachable.
	Permanent bool `json:"permanent,omitempty" yaml:"permanent,omitempty"`
}
