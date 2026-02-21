// Package juniper provides parsers for Juniper JunOS device configurations.
package juniper

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/0xdevren/netsentry/internal/model"
)

// JunOSParser parses Juniper JunOS device configurations.
// JunOS uses a hierarchical curly-brace format distinct from Cisco IOS.
type JunOSParser struct{}

// NewJunOSParser constructs a JunOSParser.
func NewJunOSParser() *JunOSParser { return &JunOSParser{} }

// DeviceType returns the platform this parser handles.
func (p *JunOSParser) DeviceType() model.DeviceType { return model.DeviceTypeJuniperOS }

// Parse converts JunOS configuration (both set-format and hierarchical) into
// a ConfigModel. The parser handles both "set system host-name R1" flat stanzas
// and hierarchical block formats.
func (p *JunOSParser) Parse(_ context.Context, data []byte, device model.Device) (*model.ConfigModel, error) {
	lines := splitLines(data)
	cfg := &model.ConfigModel{
		Device:         device,
		RawText:        string(data),
		Lines:          lines,
		GlobalSettings: make(map[string]string),
	}

	// Detect format: set-based or hierarchical.
	isSetFormat := false
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "set ") {
			isSetFormat = true
			break
		}
	}

	if isSetFormat {
		return p.parseSetFormat(cfg, lines)
	}
	return p.parseHierarchical(cfg, lines)
}

// parseSetFormat handles "set path value" style JunOS configuration.
func (p *JunOSParser) parseSetFormat(cfg *model.ConfigModel, lines []string) (*model.ConfigModel, error) {
	interfaceMap := make(map[string]*model.Interface)
	bgp := &model.BGPConfig{}
	hasBGP := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "set ") {
			continue
		}
		stmt := strings.TrimPrefix(line, "set ")
		parts := strings.Fields(stmt)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "system":
			if len(parts) >= 3 && parts[1] == "host-name" {
				cfg.Device.Hostname = parts[2]
				cfg.GlobalSettings["hostname"] = parts[2]
			}
			if len(parts) >= 3 && parts[1] == "ntp" && parts[2] == "server" && len(parts) >= 4 {
				cfg.GlobalSettings["ntp_server"] = parts[3]
			}

		case "interfaces":
			if len(parts) >= 2 {
				ifName := parts[1]
				if _, ok := interfaceMap[ifName]; !ok {
					interfaceMap[ifName] = &model.Interface{Name: ifName, Attributes: make(map[string]string)}
				}
				iface := interfaceMap[ifName]
				stmt := strings.Join(parts[2:], " ")
				if strings.HasPrefix(stmt, "description ") {
					iface.Description = strings.TrimPrefix(stmt, "description ")
				}
				if strings.Contains(stmt, "address ") {
					addr := stmt[strings.Index(stmt, "address ")+len("address "):]
					iface.IPAddress = strings.Fields(addr)[0]
				}
				if strings.Contains(stmt, "disable") {
					iface.Shutdown = true
				}
			}

		case "protocols":
			if len(parts) >= 2 && parts[1] == "bgp" {
				hasBGP = true
				stmt := strings.Join(parts[2:], " ")
				if strings.HasPrefix(stmt, "group ") {
					// Extract neighbor within group.
					subparts := strings.Fields(strings.TrimPrefix(stmt, "group "))
					if len(subparts) >= 3 && subparts[1] == "neighbor" {
						addr := subparts[2]
						found := false
						for i, n := range bgp.Neighbors {
							if n.Address == addr {
								_ = i
								found = true
							}
						}
						if !found {
							bgp.Neighbors = append(bgp.Neighbors, model.BGPNeighbor{Address: addr})
						}
					}
				}
				if strings.HasPrefix(stmt, "local-as ") {
					if as, err := strconv.Atoi(strings.TrimPrefix(stmt, "local-as ")); err == nil {
						bgp.LocalAS = as
					}
				}
			}

		case "routing-options":
			if len(parts) >= 3 && parts[1] == "static" && parts[2] == "route" && len(parts) >= 5 {
				dest := parts[3]
				nh := ""
				for i, part := range parts {
					if part == "next-hop" && i+1 < len(parts) {
						nh = parts[i+1]
						break
					}
				}
				cfg.StaticRoutes = append(cfg.StaticRoutes, model.StaticRoute{Destination: dest, NextHop: nh})
			}
			if len(parts) >= 3 && parts[1] == "autonomous-system" {
				if as, err := strconv.Atoi(parts[2]); err == nil && !hasBGP {
					bgp.LocalAS = as
				}
			}
		}
	}

	for _, iface := range interfaceMap {
		cfg.Interfaces = append(cfg.Interfaces, *iface)
	}
	if hasBGP {
		cfg.BGPConfig = bgp
	}

	return cfg, nil
}

// parseHierarchical handles JunOS hierarchical brace-format configuration.
func (p *JunOSParser) parseHierarchical(cfg *model.ConfigModel, lines []string) (*model.ConfigModel, error) {
	// Simple block-walking parser for hierarchical format.
	type frame struct{ block string }
	stack := []frame{}

	interfaceMap := make(map[string]*model.Interface)
	var currentIface *model.Interface

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "##" {
			continue
		}

		if strings.HasSuffix(trimmed, "{") {
			block := strings.TrimSuffix(trimmed, "{")
			block = strings.TrimSpace(block)
			stack = append(stack, frame{block: block})

			// Track current interface context.
			if len(stack) == 2 && stack[0].block == "interfaces" {
				ifName := block
				if _, ok := interfaceMap[ifName]; !ok {
					interfaceMap[ifName] = &model.Interface{Name: ifName, Attributes: make(map[string]string)}
				}
				currentIface = interfaceMap[ifName]
			}
			continue
		}

		if trimmed == "}" {
			if len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if len(stack) == 1 && stack[0].block == "interfaces" {
					_ = top
				} else if len(stack) == 0 {
					currentIface = nil
				}
			}
			continue
		}

		// Attribute assignment inside a block.
		stmt := strings.TrimSuffix(trimmed, ";")

		// System block.
		if len(stack) >= 1 && stack[0].block == "system" {
			if strings.HasPrefix(stmt, "host-name ") {
				cfg.Device.Hostname = strings.TrimPrefix(stmt, "host-name ")
				cfg.GlobalSettings["hostname"] = cfg.Device.Hostname
			}
		}

		// Interface attributes.
		if currentIface != nil && len(stack) >= 2 {
			if strings.HasPrefix(stmt, "description ") {
				currentIface.Description = strings.TrimPrefix(stmt, "description ")
			}
			if strings.HasPrefix(stmt, "address ") {
				currentIface.IPAddress = strings.TrimPrefix(stmt, "address ")
			}
			if stmt == "disable" {
				currentIface.Shutdown = true
			}
		}

		// Static routes.
		if len(stack) >= 2 && stack[0].block == "routing-options" && stack[1].block == "static" {
			if strings.HasPrefix(stmt, "next-hop ") {
				nh := strings.TrimPrefix(stmt, "next-hop ")
				if len(cfg.StaticRoutes) > 0 {
					cfg.StaticRoutes[len(cfg.StaticRoutes)-1].NextHop = nh
				}
			}
		}
		if len(stack) >= 2 && stack[0].block == "routing-options" && stack[1].block == "static" {
			if strings.HasPrefix(stmt, "route ") {
				dest := strings.TrimPrefix(stmt, "route ")
				cfg.StaticRoutes = append(cfg.StaticRoutes, model.StaticRoute{Destination: strings.Fields(dest)[0]})
			}
		}
		_ = fmt.Sprintf // avoid unused import
	}

	for _, iface := range interfaceMap {
		cfg.Interfaces = append(cfg.Interfaces, *iface)
	}

	return cfg, nil
}

// splitLines splits raw bytes on newlines.
func splitLines(data []byte) []string {
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimRight(line, "\r")
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	return lines
}
