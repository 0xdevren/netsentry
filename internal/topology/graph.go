// Package topology provides multi-device network topology graph construction
// and analysis.
package topology

import (
	"fmt"

	"github.com/0xdevren/netsentry/internal/model"
)

// Graph is the in-memory representation of the network topology.
type Graph struct {
	nodes map[string]*model.Device
	edges []model.TopologyLink
}

// NewGraph constructs an empty Graph.
func NewGraph() *Graph {
	return &Graph{nodes: make(map[string]*model.Device)}
}

// AddDevice adds or replaces a device node in the graph.
func (g *Graph) AddDevice(d model.Device) {
	id := d.ID
	if id == "" {
		id = d.Hostname
	}
	g.nodes[id] = &d
}

// AddLink adds a directed link between two devices.
func (g *Graph) AddLink(link model.TopologyLink) {
	g.edges = append(g.edges, link)
}

// Devices returns a copy of all device nodes.
func (g *Graph) Devices() []model.Device {
	out := make([]model.Device, 0, len(g.nodes))
	for _, d := range g.nodes {
		out = append(out, *d)
	}
	return out
}

// Links returns all directed links.
func (g *Graph) Links() []model.TopologyLink {
	return g.edges
}

// Neighbors returns device IDs directly connected to the given device.
func (g *Graph) Neighbors(deviceID string) []string {
	var neighbors []string
	seen := make(map[string]struct{})
	for _, e := range g.edges {
		if e.SourceDevice == deviceID {
			if _, ok := seen[e.TargetDevice]; !ok {
				neighbors = append(neighbors, e.TargetDevice)
				seen[e.TargetDevice] = struct{}{}
			}
		}
	}
	return neighbors
}

// Validate ensures all referenced device IDs in links exist as nodes.
func (g *Graph) Validate() error {
	for _, e := range g.edges {
		if _, ok := g.nodes[e.SourceDevice]; !ok {
			return fmt.Errorf("topology: link references unknown source device %q", e.SourceDevice)
		}
		if _, ok := g.nodes[e.TargetDevice]; !ok {
			return fmt.Errorf("topology: link references unknown target device %q", e.TargetDevice)
		}
	}
	return nil
}

// ToModel converts the Graph to a model.TopologyGraph for serialisation.
func (g *Graph) ToModel() *model.TopologyGraph {
	tg := &model.TopologyGraph{
		Devices: make(map[string]model.Device, len(g.nodes)),
		Links:   g.edges,
	}
	for id, d := range g.nodes {
		tg.Devices[id] = *d
	}
	return tg
}
