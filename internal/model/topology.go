package model

// TopologyLink represents a directed connection between two devices.
type TopologyLink struct {
	// SourceDevice is the ID of the originating device.
	SourceDevice string `json:"source_device" yaml:"source_device"`
	// SourceInterface is the interface on the source device.
	SourceInterface string `json:"source_interface,omitempty" yaml:"source_interface,omitempty"`
	// TargetDevice is the ID of the destination device.
	TargetDevice string `json:"target_device" yaml:"target_device"`
	// TargetInterface is the interface on the target device.
	TargetInterface string `json:"target_interface,omitempty" yaml:"target_interface,omitempty"`
	// Protocol is the protocol that established this adjacency (e.g. "ospf", "bgp", "cdp").
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

// TopologyGraph is a collection of devices and their inter-device links.
type TopologyGraph struct {
	// Devices is the map of device ID to Device.
	Devices map[string]Device `json:"devices,omitempty" yaml:"devices,omitempty"`
	// Links is the list of directed topology links.
	Links []TopologyLink `json:"links,omitempty" yaml:"links,omitempty"`
}
