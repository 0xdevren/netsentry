package model

// OSPFArea represents a single OSPF area with its type and networks.
type OSPFArea struct {
	// ID is the area identifier (e.g. "0", "0.0.0.0", "10").
	ID string `json:"id" yaml:"id"`
	// Type is the area type: "backbone", "stub", "nssa", or "normal".
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// Networks are the subnets participating in this area.
	Networks []string `json:"networks,omitempty" yaml:"networks,omitempty"`
}

// OSPFRedistribution defines a routing protocol redistributed into OSPF.
type OSPFRedistribution struct {
	// Source is the source protocol (e.g. "bgp", "static", "connected").
	Source string `json:"source" yaml:"source"`
	// Metric is the redistributed metric value.
	Metric int `json:"metric,omitempty" yaml:"metric,omitempty"`
	// MetricType is the OSPF metric type (1 or 2).
	MetricType int `json:"metric_type,omitempty" yaml:"metric_type,omitempty"`
	// RouteMap optionally filters redistributed routes.
	RouteMap string `json:"route_map,omitempty" yaml:"route_map,omitempty"`
}

// OSPFConfig holds the OSPF protocol configuration for a device.
type OSPFConfig struct {
	// ProcessID is the OSPF process identifier.
	ProcessID int `json:"process_id" yaml:"process_id"`
	// RouterID is the OSPF router identifier.
	RouterID string `json:"router_id,omitempty" yaml:"router_id,omitempty"`
	// Areas are the configured OSPF areas.
	Areas []OSPFArea `json:"areas,omitempty" yaml:"areas,omitempty"`
	// Redistributions lists protocols redistributed into OSPF.
	Redistributions []OSPFRedistribution `json:"redistributions,omitempty" yaml:"redistributions,omitempty"`
	// PassiveInterfaces lists interfaces that do not send OSPF hellos.
	PassiveInterfaces []string `json:"passive_interfaces,omitempty" yaml:"passive_interfaces,omitempty"`
	// DefaultPassive indicates that all interfaces are passive by default.
	DefaultPassive bool `json:"default_passive,omitempty" yaml:"default_passive,omitempty"`
}
