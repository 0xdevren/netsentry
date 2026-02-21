package checks

import (
	"fmt"
	"github.com/0xdevren/netsentry/internal/model"
)

// LoopCheck detects routing loops in the topology link graph using DFS cycle detection.
type LoopCheck struct{}

// Run detects directed cycles in the topology.
func (c *LoopCheck) Run(g *model.TopologyGraph) []Issue {
	// Build adjacency list.
	adj := make(map[string][]string)
	for _, link := range g.Links {
		adj[link.SourceDevice] = append(adj[link.SourceDevice], link.TargetDevice)
	}

	visited := make(map[string]int) // 0=unvisited, 1=in-stack, 2=done
	var issues []Issue

	var dfs func(node string, path []string)
	dfs = func(node string, path []string) {
		visited[node] = 1
		for _, neighbor := range adj[node] {
			if visited[neighbor] == 1 {
				issues = append(issues, Issue{
					Code:     "LOOP-001",
					Severity: "CRITICAL",
					Message:  fmt.Sprintf("routing loop detected: %v -> %s", path, neighbor),
					DeviceID: node,
				})
				return
			}
			if visited[neighbor] == 0 {
				dfs(neighbor, append(path, neighbor))
			}
		}
		visited[node] = 2
	}

	for id := range g.Devices {
		if visited[id] == 0 {
			dfs(id, []string{id})
		}
	}

	return issues
}

// SubnetOverlapCheck detects overlapping management IP subnets.
type SubnetOverlapCheck struct{}

// Run checks for subnet overlaps across devices.
func (s *SubnetOverlapCheck) Run(g *model.TopologyGraph) []Issue {
	// Placeholder: full CIDR overlap check requires interface-level data.
	// At topology graph level, we flag devices with identical /24 management prefixes.
	prefixMap := make(map[string]string)
	var issues []Issue

	for id, dev := range g.Devices {
		ip := dev.ManagementIP
		if ip == "" {
			continue
		}
		prefix := subnetPrefix(ip, 24)
		if prefix == "" {
			continue
		}
		if existing, ok := prefixMap[prefix]; ok {
			issues = append(issues, Issue{
				Code:     "SUBNET-OVERLAP-001",
				Severity: "MEDIUM",
				Message:  fmt.Sprintf("devices %q and %q share /24 prefix %s", existing, id, prefix),
				DeviceID: id,
			})
		} else {
			prefixMap[prefix] = id
		}
	}
	return issues
}

// subnetPrefix extracts the first prefixLen-bit network prefix from an IP string.
func subnetPrefix(ip string, prefixLen int) string {
	parts := splitIP(ip)
	if len(parts) != 4 {
		return ""
	}
	octets := prefixLen / 8
	joined := ""
	for i := 0; i < octets && i < 4; i++ {
		if i > 0 {
			joined += "."
		}
		joined += parts[i]
	}
	return joined
}

// splitIP splits an IPv4 address string on dots.
func splitIP(ip string) []string {
	var parts []string
	current := ""
	for _, c := range ip {
		if c == '.' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
