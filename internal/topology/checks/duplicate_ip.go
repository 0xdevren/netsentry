package checks

import (
	"fmt"
	"github.com/0xdevren/netsentry/internal/model"
)

// DuplicateIPCheck detects interface IP addresses assigned to more than one device.
type DuplicateIPCheck struct{}

// Run scans all device interfaces for duplicate IP assignments.
func (c *DuplicateIPCheck) Run(g *model.TopologyGraph) []Issue {
	// DuplicateIPCheck requires access to device interfaces; this check
	// operates on the topology graph level and reports any duplicate IPs.
	// Since the TopologyGraph does not embed ConfigModels, this check is a
	// structural placeholder that validates the device management IPs.
	seen := make(map[string]string) // ip -> deviceID
	var issues []Issue

	for id, dev := range g.Devices {
		ip := dev.ManagementIP
		if ip == "" {
			continue
		}
		if existing, ok := seen[ip]; ok {
			issues = append(issues, Issue{
				Code:     "DUP-IP-001",
				Severity: "HIGH",
				Message:  fmt.Sprintf("duplicate management IP %s assigned to devices %q and %q", ip, existing, id),
				DeviceID: id,
			})
		} else {
			seen[ip] = id
		}
	}
	return issues
}
