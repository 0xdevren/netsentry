package topology

import (
	"github.com/0xdevren/netsentry/internal/model"
)

// Builder constructs a Graph from a collection of ConfigModels by inferring
// device adjacencies from routing protocol neighbours.
type Builder struct{}

// NewBuilder constructs a Builder.
func NewBuilder() *Builder { return &Builder{} }

// Build constructs a topology Graph from the supplied list of device configs.
// Adjacencies are inferred from BGP and OSPF neighbour configuration.
func (b *Builder) Build(configs []*model.ConfigModel) *Graph {
	g := NewGraph()

	for _, cfg := range configs {
		g.AddDevice(cfg.Device)
	}

	// Infer BGP adjacencies.
	for _, cfg := range configs {
		if cfg.BGPConfig == nil {
			continue
		}
		srcID := cfg.Device.ID
		if srcID == "" {
			srcID = cfg.Device.Hostname
		}
		for _, neighbor := range cfg.BGPConfig.Neighbors {
			targetID := b.findDeviceByIP(configs, neighbor.Address)
			if targetID == "" {
				continue
			}
			g.AddLink(model.TopologyLink{
				SourceDevice: srcID,
				TargetDevice: targetID,
				Protocol:     "bgp",
			})
		}
	}

	// Infer OSPF adjacencies from shared network prefixes.
	// Two devices in the same OSPF area with overlapping networks are considered adjacent.
	for i, cfgA := range configs {
		if cfgA.OSPFConfig == nil {
			continue
		}
		idA := cfgA.Device.ID
		if idA == "" {
			idA = cfgA.Device.Hostname
		}
		for j, cfgB := range configs {
			if i >= j || cfgB.OSPFConfig == nil {
				continue
			}
			idB := cfgB.Device.ID
			if idB == "" {
				idB = cfgB.Device.Hostname
			}
			if b.sharedOSPFArea(cfgA.OSPFConfig, cfgB.OSPFConfig) {
				g.AddLink(model.TopologyLink{
					SourceDevice: idA,
					TargetDevice: idB,
					Protocol:     "ospf",
				})
			}
		}
	}

	return g
}

// findDeviceByIP returns the device ID whose management IP matches addr.
func (b *Builder) findDeviceByIP(configs []*model.ConfigModel, addr string) string {
	for _, cfg := range configs {
		if cfg.Device.ManagementIP == addr {
			id := cfg.Device.ID
			if id == "" {
				return cfg.Device.Hostname
			}
			return id
		}
		for _, iface := range cfg.Interfaces {
			if iface.IPAddress == addr {
				id := cfg.Device.ID
				if id == "" {
					return cfg.Device.Hostname
				}
				return id
			}
		}
	}
	return ""
}

// sharedOSPFArea reports whether two OSPF configs share at least one area ID.
func (b *Builder) sharedOSPFArea(a, c *model.OSPFConfig) bool {
	areaSet := make(map[string]struct{})
	for _, area := range a.Areas {
		areaSet[area.ID] = struct{}{}
	}
	for _, area := range c.Areas {
		if _, ok := areaSet[area.ID]; ok {
			return true
		}
	}
	return false
}
